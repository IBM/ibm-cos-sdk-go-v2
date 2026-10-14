package manager_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
	"golang.org/x/sys/unix"
)

type downloadPipeClient struct {
	data   []byte
	ranges []string
	mu     sync.Mutex
}

func (c *downloadPipeClient) HeadObject(_ context.Context, _ *s3.HeadObjectInput, _ ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{
		ContentLength: aws.Int64(int64(len(c.data))),
		PartsCount:    aws.Int32(3),
	}, nil
}

func (c *downloadPipeClient) GetObject(_ context.Context, input *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	c.mu.Lock()
	c.ranges = append(c.ranges, aws.ToString(input.Range))
	c.mu.Unlock()

	var start, end int
	if _, err := fmt.Sscanf(aws.ToString(input.Range), "bytes=%d-%d", &start, &end); err != nil {
		return nil, err
	}
	if end >= len(c.data) {
		end = len(c.data) - 1
	}
	part := append([]byte(nil), c.data[start:end+1]...)

	return &s3.GetObjectOutput{
		Body:          io.NopCloser(bytes.NewReader(part)),
		ContentLength: aws.Int64(int64(len(part))),
	}, nil
}

func TestDownloadWithPipe(t *testing.T) {
	pipePath := filepath.Join(t.TempDir(), "download.pipe")
	if err := unix.Mkfifo(pipePath, 0600); err != nil {
		t.Fatalf("mkfifo failed: %v", err)
	}

	payload := []byte("abcdefghijklmnopqrstuvwxyz")
	client := &downloadPipeClient{data: payload}
	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.Concurrency = 2
	})

	readCh := make(chan []byte, 1)
	readErrCh := make(chan error, 1)
	go func() {
		reader, err := os.OpenFile(pipePath, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			readErrCh <- err
			return
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			readErrCh <- err
			return
		}
		readCh <- data
	}()

	var progress []int64
	var progressMu sync.Mutex

	n, err := downloader.DownloadWithPipe(context.Background(), &manager.DownloadPipeInput{
		PipePath: aws.String(pipePath),
		GetObjectInput: &s3.GetObjectInput{
			Bucket: aws.String("bucket"),
			Key:    aws.String("key"),
		},
		ProgressFn: func(bytesDownloaded, totalBytes int64) {
			progressMu.Lock()
			progress = append(progress, bytesDownloaded)
			progressMu.Unlock()
			if totalBytes != int64(len(payload)) {
				readErrCh <- fmt.Errorf("expected total bytes %d, got %d", len(payload), totalBytes)
			}
		},
	})
	if err != nil {
		t.Fatalf("DownloadWithPipe failed: %v", err)
	}
	if n != int64(len(payload)) {
		t.Fatalf("expected %d downloaded bytes, got %d", len(payload), n)
	}

	select {
	case err := <-readErrCh:
		t.Fatalf("pipe reader failed: %v", err)
	case got := <-readCh:
		if !bytes.Equal(got, payload) {
			t.Fatalf("expected payload %q, got %q", payload, got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for pipe reader")
	}

	progressMu.Lock()
	if len(progress) == 0 {
		progressMu.Unlock()
		t.Fatal("expected progress callback to be invoked")
	}
	finalProgress := progress[len(progress)-1]
	progressMu.Unlock()
	if finalProgress != int64(len(payload)) {
		t.Fatalf("expected final progress %d, got %d", len(payload), finalProgress)
	}

	client.mu.Lock()
	gotRanges := append([]string(nil), client.ranges...)
	client.mu.Unlock()
	sort.Strings(gotRanges)
	expectedRanges := []string{"bytes=0-8", "bytes=9-17", "bytes=18-25"}
	sort.Strings(expectedRanges)
	if fmt.Sprint(gotRanges) != fmt.Sprint(expectedRanges) {
		t.Fatalf("expected ranges %v, got %v", expectedRanges, gotRanges)
	}
}
