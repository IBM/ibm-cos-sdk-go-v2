package manager

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

const (
	invalidParameterError = "InvalidParameter"
	pipeOpenError         = "PipeOpenError"
)

var ErrPipeWriteTimeout = errors.New("pipe write timeout")

// DownloadPipeInput provides the input parameters for downloading to a named pipe.
type DownloadPipeInput struct {
	PipePath        *string
	GetObjectInput  *s3.GetObjectInput
	PipeTimeout     *int64
	MaxBufferMemory *int64
	MaxRetries      *int
	ProgressFn      func(bytesDownloaded, totalBytes int64)
}

type memoryManager struct {
	maxMemory       int64
	availableMemory int64
	mu              sync.Mutex
	cond            *sync.Cond
}

func newMemoryManager(maxMemory int64) *memoryManager {
	mm := &memoryManager{maxMemory: maxMemory, availableMemory: maxMemory}
	mm.cond = sync.NewCond(&mm.mu)
	return mm
}

func (mm *memoryManager) Reserve(size int64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	for mm.availableMemory < size {
		mm.cond.Wait()
	}

	mm.availableMemory -= size
}

func (mm *memoryManager) Release(size int64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.availableMemory += size
	mm.cond.Broadcast()
}

type orderedWriter struct {
	writer       io.Writer
	nextIndex    int64
	buffer       map[int64][]byte
	mu           sync.Mutex
	flushMu      sync.Mutex
	chunkSize    int
	writeTimeout time.Duration
	memoryMgr    *memoryManager
}

func newOrderedWriter(writer io.Writer, chunkSize int, writeTimeout time.Duration, memoryMgr *memoryManager) *orderedWriter {
	if chunkSize == 0 {
		chunkSize = 64 * 1024
	}
	if writeTimeout == 0 {
		writeTimeout = 30 * time.Second
	}

	return &orderedWriter{
		writer:       writer,
		nextIndex:    1,
		buffer:       make(map[int64][]byte),
		chunkSize:    chunkSize,
		writeTimeout: writeTimeout,
		memoryMgr:    memoryMgr,
	}
}

func (ow *orderedWriter) WritePart(ctx context.Context, partNum int64, data []byte) error {
	ow.mu.Lock()
	ow.buffer[partNum] = data
	ow.mu.Unlock()

	return ow.flushSequential(ctx)
}

func (ow *orderedWriter) flushSequential(ctx context.Context) error {
	ow.flushMu.Lock()
	defer ow.flushMu.Unlock()

	for {
		ow.mu.Lock()
		data, exists := ow.buffer[ow.nextIndex]
		if !exists {
			ow.mu.Unlock()
			return nil
		}
		dataSize := int64(len(data))
		delete(ow.buffer, ow.nextIndex)
		ow.nextIndex++
		ow.mu.Unlock()

		if err := ow.writeChunked(ctx, data); err != nil {
			return err
		}

		ow.memoryMgr.Release(dataSize)
	}
}

func (ow *orderedWriter) writeChunked(ctx context.Context, data []byte) error {
	for i := 0; i < len(data); i += ow.chunkSize {
		end := min(i+ow.chunkSize, len(data))

		writeCtx, cancel := context.WithTimeout(ctx, ow.writeTimeout)
		errChan := make(chan error, 1)
		chunk := data[i:end]
		go func() {
			_, err := ow.writer.Write(chunk)
			errChan <- err
		}()

		select {
		case <-writeCtx.Done():
			cancel()
			return ErrPipeWriteTimeout
		case err := <-errChan:
			cancel()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ow *orderedWriter) Flush(context.Context) error {
	ow.mu.Lock()
	defer ow.mu.Unlock()

	if len(ow.buffer) > 0 {
		return fmt.Errorf("parts remain unwritten in buffer")
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ObjectMetadata contains metadata about the S3 object to be downloaded.
type ObjectMetadata struct {
	ContentLength int64
	PartsCount    int64
	PartSize      int64
}

func (d *Downloader) GetObjectMetadata(ctx context.Context, input *s3.GetObjectInput) (*ObjectMetadata, error) {
	headOutput, err := d.S3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket:    input.Bucket,
		Key:       input.Key,
		VersionId: input.VersionId,
	}, d.ClientOptions...)
	if err != nil {
		return nil, err
	}

	metadata := &ObjectMetadata{
		ContentLength: aws.ToInt64(headOutput.ContentLength),
		PartsCount:    int64(aws.ToInt32(headOutput.PartsCount)),
	}
	if metadata.PartsCount == 0 {
		metadata.PartsCount = 1
	}

	return metadata, nil
}

func (m *ObjectMetadata) computePartSize() {
	if m.ContentLength == 0 || m.PartsCount == 0 {
		m.PartSize = 0
		return
	}

	m.PartSize = m.ContentLength / m.PartsCount
	if m.ContentLength%m.PartsCount != 0 {
		m.PartSize++
	}
}

func getMaxBufferMemory(d *Downloader, metadata *ObjectMetadata) int64 {
	concurrency := int64(d.Concurrency)
	if concurrency == 0 {
		concurrency = DefaultDownloadConcurrency
	}

	partSize := metadata.PartSize
	if partSize == 0 {
		partSize = d.PartSize
		if partSize == 0 {
			partSize = DefaultDownloadPartSize
		}
	}

	return concurrency * partSize * 2
}

func validateDownloadPipeInput(input *DownloadPipeInput) error {
	if input == nil {
		return fmt.Errorf("%s: input is required", invalidParameterError)
	}
	if input.PipePath == nil {
		return fmt.Errorf("%s: PipePath is required", invalidParameterError)
	}
	if input.GetObjectInput == nil {
		return fmt.Errorf("%s: GetObjectInput is required", invalidParameterError)
	}
	if input.MaxBufferMemory != nil && *input.MaxBufferMemory <= 0 {
		return fmt.Errorf("%s: MaxBufferMemory must be greater than 0", invalidParameterError)
	}

	return nil
}

func downloadWithPipe(ctx context.Context, d *Downloader, input *DownloadPipeInput, options ...func(*Downloader)) (n int64, err error) {
	if err := validateDownloadPipeInput(input); err != nil {
		return 0, err
	}

	dl := *d
	for _, option := range options {
		option(&dl)
	}

	// Ensures we don't need nil checks later on
	dl.Logger = logging.WithContext(ctx, dl.Logger)

	dl.Logger.Logf(logging.Debug, "pipe download starting, bucket %s, key %s, pipe %s",
		aws.ToString(input.GetObjectInput.Bucket),
		aws.ToString(input.GetObjectInput.Key),
		aws.ToString(input.PipePath))

	pipe, err := os.OpenFile(*input.PipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		dl.Logger.Logf(logging.Debug, "failed to open pipe %s, %v",
			aws.ToString(input.PipePath), err)
		return 0, fmt.Errorf("%s: failed to open pipe: %w", pipeOpenError, err)
	}
	defer pipe.Close()

	dl.Logger.Logf(logging.Debug, "pipe opened for writing, %s", aws.ToString(input.PipePath))

	metadata, err := dl.GetObjectMetadata(ctx, input.GetObjectInput)
	if err != nil {
		dl.Logger.Logf(logging.Debug, "failed to get object metadata, %v", err)
		return 0, err
	}
	metadata.computePartSize()

	maxBufferMemory := getMaxBufferMemory(&dl, metadata)
	if input.MaxBufferMemory != nil {
		maxBufferMemory = *input.MaxBufferMemory
	}

	bufferSize := getPipeBufferSize(pipe)
	memoryMgr := newMemoryManager(maxBufferMemory)

	timeout := 30 * time.Second
	if input.PipeTimeout != nil {
		timeout = time.Duration(*input.PipeTimeout) * time.Second
	}

	orderedWriter := newOrderedWriter(pipe, bufferSize, timeout, memoryMgr)
	totalSize := metadata.ContentLength
	partsCount := metadata.PartsCount
	chunkSize := metadata.PartSize

	workerCtx, cancelWorkers := context.WithCancel(ctx)
	defer cancelWorkers()

	var wg sync.WaitGroup
	partsChan := make(chan int64, dl.Concurrency)
	errChan := make(chan error, dl.Concurrency)

	var downloadedBytes int64
	var mu sync.Mutex

	maxRetries := 3
	if input.MaxRetries != nil {
		maxRetries = *input.MaxRetries
	}

	worker := func() {
		defer wg.Done()

		for pn := range partsChan {
			select {
			case <-workerCtx.Done():
				errChan <- workerCtx.Err()
				return
			default:
			}

			if err := func() error {
				start := (pn - 1) * chunkSize
				end := start + chunkSize - 1
				if end >= totalSize {
					end = totalSize - 1
				}
				expectedSize := end - start + 1
				byteRange := fmt.Sprintf("bytes=%d-%d", start, end)

				memoryMgr.Reserve(expectedSize)
				memoryTransferred := false
				defer func() {
					if !memoryTransferred {
						memoryMgr.Release(expectedSize)
					}
				}()

				var getOutput *s3.GetObjectOutput
				var buf *bytes.Buffer
				var lastErr error

				for attempt := 1; attempt <= maxRetries; attempt++ {
					select {
					case <-workerCtx.Done():
						return workerCtx.Err()
					default:
					}

					getOutput, lastErr = dl.S3.GetObject(workerCtx, &s3.GetObjectInput{
						Bucket:     input.GetObjectInput.Bucket,
						Key:        input.GetObjectInput.Key,
						VersionId:  input.GetObjectInput.VersionId,
						PartNumber: aws.Int32(int32(pn)),
						Range:      aws.String(byteRange),
					}, dl.ClientOptions...)
					if lastErr != nil {
						if attempt < maxRetries {
							dl.Logger.Logf(logging.Debug,
								"failed to download part %d on attempt %d of %d, retrying, %v",
								pn, attempt, maxRetries, lastErr)
							time.Sleep(time.Duration(attempt) * time.Second)
							continue
						}
						dl.Logger.Logf(logging.Debug, "failed to download part %d after %d attempts, %v",
							pn, maxRetries, lastErr)
						return fmt.Errorf("failed to download part %d after %d attempts: %w", pn, maxRetries, lastErr)
					}

					buf = new(bytes.Buffer)
					_, lastErr = io.Copy(buf, getOutput.Body)
					getOutput.Body.Close()
					if lastErr != nil {
						if attempt < maxRetries {
							dl.Logger.Logf(logging.Debug,
								"failed to read part %d on attempt %d of %d, retrying, %v",
								pn, attempt, maxRetries, lastErr)
							time.Sleep(time.Duration(attempt) * time.Second)
							continue
						}
						dl.Logger.Logf(logging.Debug, "failed to read part %d after %d attempts, %v",
							pn, maxRetries, lastErr)
						return fmt.Errorf("failed to read part %d after %d attempts: %w", pn, maxRetries, lastErr)
					}

					break
				}

				actualSize := int64(buf.Len())
				if actualSize != expectedSize {
					memoryMgr.Release(expectedSize)
					memoryMgr.Reserve(actualSize)
					expectedSize = actualSize
				}

				if err := orderedWriter.WritePart(workerCtx, pn, buf.Bytes()); err != nil {
					dl.Logger.Logf(logging.Debug, "failed to write part %d to pipe, %v", pn, err)
					return fmt.Errorf("failed to write part %d: %w", pn, err)
				}

				memoryTransferred = true

				mu.Lock()
				downloadedBytes += actualSize
				progress := downloadedBytes
				if input.ProgressFn != nil {
					input.ProgressFn(downloadedBytes, totalSize)
				}
				mu.Unlock()

				dl.Logger.Logf(logging.Debug, "wrote part %d of %d to pipe, %d bytes, %d of %d total",
					pn, partsCount, actualSize, progress, totalSize)

				return nil
			}(); err != nil {
				errChan <- err
				cancelWorkers()
				return
			}
		}
	}

	if dl.Concurrency == 0 {
		dl.Concurrency = DefaultDownloadConcurrency
	}

	dl.Logger.Logf(logging.Debug,
		"pipe download plan, contentLength %d, parts %d, partSize %d, concurrency %d, "+
			"maxBufferMemory %d, pipeBufferSize %d, writeTimeout %s, maxRetries %d",
		totalSize, partsCount, chunkSize, dl.Concurrency,
		maxBufferMemory, bufferSize, timeout, maxRetries)

	for i := 0; i < dl.Concurrency; i++ {
		wg.Add(1)
		go worker()
	}

	producerDone := make(chan struct{})
	go func() {
		defer close(producerDone)
		defer close(partsChan)

		for partNum := int64(1); partNum <= partsCount; partNum++ {
			select {
			case <-workerCtx.Done():
				return
			case partsChan <- partNum:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			dl.Logger.Logf(logging.Debug, "pipe download failed, %v", err)
			<-producerDone
			return 0, err
		}
	}

	if err := orderedWriter.Flush(context.Background()); err != nil {
		dl.Logger.Logf(logging.Debug, "pipe download flush failed, %v", err)
		return 0, err
	}

	dl.Logger.Logf(logging.Debug, "pipe download complete, %d bytes written to %s",
		downloadedBytes, aws.ToString(input.PipePath))

	return downloadedBytes, nil
}
