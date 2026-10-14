package manager

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/logging"
)

// UploadPipeInput provides the input parameters for uploading from a named pipe.
type UploadPipeInput struct {
	// Required: Path to the named pipe (FIFO)
	PipePath *string

	// Required: Standard S3 PutObjectInput
	PutObjectInput *s3.PutObjectInput

	// Optional: Timeout for pipe operations in seconds. Reserved for compatibility.
	PipeTimeout *int64

	// Optional: Progress callback
	ProgressFn func(bytesUploaded int64)
}

// pipeReader wraps a pipe and tracks bytes read.
type pipeReader struct {
	r         io.Reader
	bytesRead int64
	mu        sync.Mutex
	callback  func(int64)
}

func (pr *pipeReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)

	pr.mu.Lock()
	pr.bytesRead += int64(n)
	bytesRead := pr.bytesRead
	pr.mu.Unlock()

	if pr.callback != nil && n > 0 {
		pr.callback(bytesRead)
	}

	return n, err
}

func (pr *pipeReader) BytesRead() int64 {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	return pr.bytesRead
}

func uploadWithPipe(ctx context.Context, u *Uploader, input *UploadPipeInput, options ...func(*Uploader)) (*UploadOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("invalid parameter: input is required")
	}
	if input.PipePath == nil {
		return nil, fmt.Errorf("invalid parameter: PipePath is required")
	}
	if input.PutObjectInput == nil {
		return nil, fmt.Errorf("invalid parameter: PutObjectInput is required")
	}

	// Ensures we don't need nil checks later on
	logger := logging.WithContext(ctx, u.Logger)

	logger.Logf(logging.Debug, "pipe upload starting, bucket %s, key %s, pipe %s",
		aws.ToString(input.PutObjectInput.Bucket),
		aws.ToString(input.PutObjectInput.Key),
		aws.ToString(input.PipePath))

	pipe, err := os.OpenFile(*input.PipePath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		logger.Logf(logging.Debug, "failed to open pipe %s, %v", aws.ToString(input.PipePath), err)
		return nil, fmt.Errorf("failed to open pipe: %w", err)
	}
	defer pipe.Close()

	logger.Logf(logging.Debug, "pipe opened for reading, %s", aws.ToString(input.PipePath))

	reader := &pipeReader{
		r:        pipe,
		callback: input.ProgressFn,
	}

	putObjectInput := *input.PutObjectInput
	putObjectInput.Body = reader

	result, err := u.Upload(ctx, &putObjectInput, options...)
	if err != nil {
		logger.Logf(logging.Debug, "pipe upload failed after %d bytes read from %s, %v",
			reader.BytesRead(), aws.ToString(input.PipePath), err)
		return nil, err
	}

	logger.Logf(logging.Debug, "pipe upload complete, %d bytes read from %s",
		reader.BytesRead(), aws.ToString(input.PipePath))

	return result, nil
}
