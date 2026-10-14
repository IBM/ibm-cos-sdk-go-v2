package manager_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
	"golang.org/x/sys/unix"
)

type uploadPipeClient struct {
	body []byte
	mu   sync.Mutex
}

func (c *uploadPipeClient) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	payload, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.body = append([]byte(nil), payload...)
	c.mu.Unlock()

	return &s3.PutObjectOutput{
		ETag: aws.String("etag"),
	}, nil
}

func (*uploadPipeClient) UploadPart(context.Context, *s3.UploadPartInput, ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	panic("unexpected UploadPart")
}

func (*uploadPipeClient) CreateMultipartUpload(context.Context, *s3.CreateMultipartUploadInput, ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	panic("unexpected CreateMultipartUpload")
}

func (*uploadPipeClient) CompleteMultipartUpload(context.Context, *s3.CompleteMultipartUploadInput, ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	panic("unexpected CompleteMultipartUpload")
}

func (*uploadPipeClient) AbortMultipartUpload(context.Context, *s3.AbortMultipartUploadInput, ...func(*s3.Options)) (*s3.AbortMultipartUploadOutput, error) {
	panic("unexpected AbortMultipartUpload")
}

func TestUploadWithPipe(t *testing.T) {
	pipePath := filepath.Join(t.TempDir(), "upload.pipe")
	if err := unix.Mkfifo(pipePath, 0600); err != nil {
		t.Fatalf("mkfifo failed: %v", err)
	}

	client := &uploadPipeClient{}
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.Concurrency = 1
	})

	payload := []byte("pipe upload payload")
	var progress []int64
	var progressMu sync.Mutex

	errCh := make(chan error, 1)
	go func() {
		_, err := uploader.UploadWithPipe(context.Background(), &manager.UploadPipeInput{
			PipePath: aws.String(pipePath),
			PutObjectInput: &s3.PutObjectInput{
				Bucket: aws.String("bucket"),
				Key:    aws.String("key"),
			},
			ProgressFn: func(bytesUploaded int64) {
				progressMu.Lock()
				progress = append(progress, bytesUploaded)
				progressMu.Unlock()
			},
		})
		errCh <- err
	}()

	writer, err := os.OpenFile(pipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		t.Fatalf("open writer failed: %v", err)
	}

	if _, err := writer.Write(payload); err != nil {
		t.Fatalf("write payload failed: %v", err)
	}
	writer.Close()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("UploadWithPipe failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("UploadWithPipe timed out")
	}

	client.mu.Lock()
	got := append([]byte(nil), client.body...)
	client.mu.Unlock()

	if !bytes.Equal(got, payload) {
		t.Fatalf("expected payload %q, got %q", payload, got)
	}

	progressMu.Lock()
	defer progressMu.Unlock()
	if len(progress) == 0 {
		t.Fatal("expected progress callback to be invoked")
	}
	if progress[len(progress)-1] != int64(len(payload)) {
		t.Fatalf("expected final progress %d, got %d", len(payload), progress[len(progress)-1])
	}
}
