# V2 Migration Guide

## Table of Contents

1. [Introduction](#1-introduction)
    - [1.1 Purpose of the Migration Guide](#11-purpose-of-the-migration-guide)
    - [1.2 Scope](#12-scope)
    - [1.3 Prerequisites](#13-prerequisites)
2. [Changes in V2](#2-changes-in-v2)
    - [2.1 Client Construction](#21-client-construction)
    - [2.2 Invoking API Operations](#22-invoking-api-operations)
    - [2.3 Service Data Types](#23-service-data-types)
    - [2.4 Enumeration Values](#24-enumeration-values)
    - [2.5 Pointer Parameters](#25-pointer-parameters)
    - [2.6 Error Types](#26-error-types)
    - [2.7 Paginators](#27-paginators)
    - [2.8 Waiters](#28-waiters)
    - [2.9 Presigned Requests](#29-presigned-requests)
    - [2.10 HTTP Request/Response](#210-http-requestresponse)
    - [2.11 S3 Transfer Manager](#211-s3-transfer-manager)
3. [Migration Overview](#3-migration-overview)
    - [3.1 Current System Architecture (v1)](#31-current-system-architecture-v1)
    - [3.2 Target System Architecture (v2)](#32-target-system-architecture-v2)
    - [3.3 Key Differences at a Glance](#33-key-differences-at-a-glance)
    - [3.4 Migration Examples](#34-migration-examples)
4. [Step-by-Step Migration Process](#4-step-by-step-migration-process)
    - [4.1 Update Dependencies (Go Modules)](#41-update-dependencies-go-modules)
    - [4.2 Refactor S3 Client Initialization](#42-refactor-s3-client-initialization---configuration-loading)
    - [4.3 Update API Calls](#43-update-api-calls)
    - [4.4 Handle Context and Error Changes](#44-handle-context-and-error-changes)
5. [Handling Errors for Go V2](#5-handling-errors-for-go-v2)
    - [5.1 Logging Errors](#51-logging-errors)
    - [5.2 Client Errors](#52-client-errors)
    - [5.3 API Error Responses](#53-api-error-responses)
    - [5.4 Retrieving Request Identifiers](#54-retrieving-request-identifiers)
6. [Comparison Table](#6-comparison-table)

---

### Examples can be referred in the [examples](examples) directory

## 1. Introduction

IBM COS SDK for Go v2 provides APIs and utilities that developers can use to build Go applications that use COS services. The SDK removes the complexity of coding directly against a web service interface. It hides a lot of the lower-level plumbing, such as authentication, request retries, and error handling. This guide helps you migrate projects from IBM COS SDK for Go v1 to v2. The v2 SDK is modular, faster to work with, and aligns with modern Go practices.

### 1.1 Purpose of the Migration Guide

Developers familiar with AWS SDK v2 patterns will find the new IBM COS v2 structure intuitive. Applications written using v1 may need updates in client construction, request models, error handling, and credentials.

### 1.2 Scope

The migration guide includes but is not limited to:

- Refactoring code to use new client initialization patterns, configurations, and API operation structures.
- Updating module dependencies from v1 to v2.
- Understanding key behavioural and structural changes introduced in v2.
- Examples for common COS operations.
- Guidance on testing and error handling in v2.

### 1.3 Prerequisites

- **Go 1.24+** installed and configured.
- Working knowledge of Go modules.
- An active IBM COS account.

---

## 2. Changes in V2

### 2.1 Client Construction

You can construct clients in the IBM COS SDK for Go using either the `New` or `NewFromConfig` constructor functions in the client's package. When migrating from the IBM COS SDK for Go v1, we recommend that you use the `NewFromConfig` variant, which will return a new service client using values from an `aws.Config`. The `aws.Config` value will have been created while loading the SDK shared configuration using `config.LoadDefaultConfig`.

#### 2.1.1 Example 1

#### V1 Code `(using s3.New)`

```go
import "github.com/IBM/ibm-cos-sdk-go/aws/session"
import "github.com/IBM/ibm-cos-sdk-go/service/s3"

// ...

sess, err := session.NewSession()
if err != nil {
    // handle error
}

client := s3.New(sess)
````

#### V2 Code `(using s3.NewFromConfig)`

```go
import "context"
import "github.com/IBM/ibm-cos-sdk-go-v2/config"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
    // handle error
}

client := s3.NewFromConfig(cfg)
```

#### 2.1.2 Example 2

#### V1 Code

```go
import "github.com/IBM/ibm-cos-sdk-go/aws"
import "github.com/IBM/ibm-cos-sdk-go/aws/session"
import "github.com/IBM/ibm-cos-sdk-go/service/s3"

// ...

sess, err := session.NewSession()
if err != nil {
    // handle error
}

client := s3.New(sess, &aws.Config{
    Region: aws.String("us-south"),
})
```

#### V2 Code

```go
import "context"
import "github.com/IBM/ibm-cos-sdk-go-v2/config"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
    // handle error
}

client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.Region = "us-south"
})
```

### 2.2 Invoking API Operations

The number of S3 client operation methods has been reduced significantly. The `<OperationName>Request`, `<OperationName>WithContext`, and `<OperationName>` methods have all been consolidated into a single operation method, `<OperationName>`.

#### Example

#### V1 Code

```go
import "context"
import "github.com/IBM/ibm-cos-sdk-go/service/s3"

// ...

client := s3.New(sess)

// Pattern 1
output, err := client.PutObject(&s3.PutObjectInput{
    // input parameters
})

// Pattern 2
output, err := client.PutObjectWithContext(context.TODO(), &s3.PutObjectInput{
    // input parameters
})

// Pattern 3
req, output := client.PutObjectRequest(context.TODO(), &s3.PutObjectInput{
    // input parameters
})
err := req.Send()
```

#### V2 Code

```go
import "context"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"

// ...

client := s3.NewFromConfig(cfg)

output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
    // input parameters
})
```

### 2.3 Service Data Types

The top-level input and output types of an operation are found in the S3 client package. The input and output type for a given operation follow the pattern of `<OperationName>Input` and `<OperationName>Output`, where `OperationName` is the name of the operation you are invoking. For example, the input and output shapes for the `PutObject` operation are `PutObjectInput` and `PutObjectOutput` respectively.

### 2.4 Enumeration Values

The SDK now provides a typed experience for all API enumeration fields. Rather than using a string literal value copied from the service API reference documentation, you can now use one of the concrete types found in the S3 client's `types` package.

For example, you can provide the `PutObjectInput` operation with an ACL to be applied on an object. In the COS SDK for Go v1, this parameter was a **string** type. In the COS SDK for Go v2, this parameter is now a `types.ObjectCannedACL`. The `types` package provides generated constants for the valid enumeration values that can be assigned to this field. For example `types.ObjectCannedACLPrivate` is the constant for the *"private"* canned ACL value. This value can be used in place of managing string constants within your application.

### 2.5 Pointer Parameters

The IBM COS SDK for Go v1 required pointer references to be passed for all input parameters. The COS SDK for Go v2 has simplified the experience by removing the need to pass input values as pointers where possible. This change means that many operations no longer require your application to pass pointer references for the following types: `uint8`, `uint16`, `uint32`, `int8`, `int16`, `int32`, `float32`, `float64`, `bool`. Similarly, slice and map element types have been updated accordingly to reflect whether their elements must be passed as pointer references.

The `aws` package contains helper functions for creating pointers for the Go built-in types; these helpers should be used to more easily handle creating pointer types for these Go types. Similarly, helper methods are provided for safely de-referencing pointer values for these types. For example, `aws.String` converts from `string` ⇒ `*string`. Inversely, `aws.ToString` converts from `*string` ⇒ `string`. When upgrading your application from IBM COS SDK for Go v1 to v2, you must migrate usage of the helpers for converting from the pointer types to the non-pointer variants. For example, `aws.StringValue` must be updated to `aws.ToString`.

### 2.6 Error Types

The COS SDK for Go takes full advantage of the error wrapping functionality introduced in Go 1.13. For example, the `GetObject` operation can return a `NoSuchKey` error if attempting to retrieve an object key that doesn't exist. You can use `errors.As` to test whether the returned operation error is a `types.NoSuchKey` error. In the event a specific type for an error is not modeled, you can utilize the `smithy.APIError` interface type for inspecting the returned error code and message from the service. This functionality replaces `awserr.Error` and the other `awserr` functionality from the COS SDK for Go v1.

#### V1 Code

```go
import "github.com/IBM/ibm-cos-sdk-go/aws/awserr"
import "github.com/IBM/ibm-cos-sdk-go/service/s3"

// ...

client := s3.New(sess)

output, err := s3.GetObject(&s3.GetObjectInput{
    // input parameters
})
if err != nil {
    if awsErr, ok := err.(awserr.Error); ok {
        if awsErr.Code() == "NoSuchKey" {
            // handle NoSuchKey
        } else {
            // handle other codes
        }
        return
    }
    // handle a error
}
```

#### V2 Code

```go
import "context"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3/types"
import "github.com/aws/smithy-go"

// ...

client := s3.NewFromConfig(cfg)

output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
    // input parameters
})
if err != nil {
    var nsk *types.NoSuchKey
    if errors.As(err, &nsk) {
        // handle NoSuchKey error
        return
    }
    var apiErr smithy.APIError
    if errors.As(err, &apiErr) {
        code := apiErr.ErrorCode()
        message := apiErr.ErrorMessage()
        // handle error code
        return
    }
    // handle error
    return
}
```

### 2.7 Paginators

Paginators are no longer invoked as methods on the S3 client. To use a paginator for an operation you must construct a paginator for an operation using one of the paginator constructor methods. For example, to paginate over the `ListObjectsV2` operation you must construct its paginator using the `s3.NewListObjectsV2Paginator`. This constructor returns a `ListObjectsV2Paginator` which provides the methods `HasMorePages` and `NextPage` for determining whether there are more pages to retrieve and invoking the operation to retrieve the next page respectively.

#### V1 Code

```go
import "fmt"
import "github.com/IBM/ibm-cos-sdk-go/service/s3"

// ...

client := s3.New(sess)

params := &s3.ListObjectsV2Input{
    // input parameters
}

totalObjects := 0
err := client.ListObjectsV2Pages(params, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
    totalObjects += len(output.Contents)
    return !lastPage
})
if err != nil {
    // handle error
}
fmt.Println("total objects:", totalObjects)
```

#### V2 Code

```go
import "context"
import "fmt"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"

// ...

client := s3.NewFromConfig(cfg)

params := &s3.ListObjectsV2Input{
    // input parameters
}

totalObjects := 0
paginator := s3.NewListObjectsV2Paginator(client, params)
for paginator.HasMorePages() {
    output, err := paginator.NextPage(context.TODO())
    if err != nil {
        // handle error
    }
    totalObjects += len(output.Contents)
}
fmt.Println("total objects:", totalObjects)
```

### 2.8 Waiters

Waiters are no longer invoked as methods on the S3 client. To use a waiter you first construct the desired waiter type, and then invoke the wait method. For example, to wait for a S3 Bucket to exist, you must construct a `BucketExists` waiter. Use the `s3.NewBucketExistsWaiter` constructor to create a `s3.BucketExistsWaiter`. The `s3.BucketExistsWaiter` provides a `Wait` method which can be used to wait for a bucket to become available.

### 2.9 Presigned Requests

COS SDK for Go exposes specific `PresignClient` implementations in the S3 package with specific APIs for supported presignable operations. Uses of `Presign` and `PresignRequest` must be converted to use the presigning client.

#### V1 Code

```go
import (
    "fmt"
    "time"
    "github.com/IBM/ibm-cos-sdk-go/aws"
    "github.com/IBM/ibm-cos-sdk-go/aws/session"
    "github.com/IBM/ibm-cos-sdk-go/service/s3"
)

func main() {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    svc := s3.New(sess)
    req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
        Bucket: aws.String("demo-bucket"),
        Key:    aws.String("key"),
    })

    // pattern 1
    url1, err := req.Presign(20 * time.Minute)
    if err != nil {
        panic(err)
    }
    fmt.Println(url1)

    // pattern 2
    url2, header, err := req.PresignRequest(20 * time.Minute)
    if err != nil {
        panic(err)
    }
    fmt.Println(url2, header)
}
```

#### V2 Code

```go
import (
    "context"
    "fmt"
    "time"
    "github.com/IBM/ibm-cos-sdk-go-v2/aws"
    "github.com/IBM/ibm-cos-sdk-go-v2/config"
    "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

func main() {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        panic(err)
    }

    svc := s3.NewPresignClient(s3.NewFromConfig(cfg))
    req, err := svc.PresignGetObject(context.Background(), &s3.GetObjectInput{
        Bucket: aws.String("demo-bucket"),
        Key:    aws.String("key"),
    }, func(o *s3.PresignOptions) {
        o.Expires = 20 * time.Minute
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(req.Method, req.URL, req.SignedHeader)
}
```

### 2.10 HTTP Request/Response

The `HTTPRequest` and `HTTPResponse` fields from `Request` are now exposed in specific middleware phases. Since middleware is transport-agnostic, you must perform a type assertion on the middleware input or output to reveal the underlying HTTP request or response. Request handlers which reference `Request.HTTPRequest` and `Request.HTTPResponse` must be migrated to middleware.

#### V1 Code

```go
import (
    "github.com/IBM/ibm-cos-sdk-go/aws/request"
    "github.com/IBM/ibm-cos-sdk-go/aws/session"
)

func withHeader(header, val string) request.Option {
    return func(r *request.Request) {
        request.HTTPRequest.Header.Set(header, val)
    }
}

func main() {
    sess := session.Must(session.NewSession())
    sess.Handlers.Build.PushBack(withHeader("x-user-header", "..."))
    svc := s3.New(sess)
    // ...
}
```

#### V2 Code

```go
import (
    "context"
    "fmt"
    "github.com/IBM/ibm-cos-sdk-go-v2/config"
    "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
    "github.com/aws/smithy-go/middleware"
    smithyhttp "github.com/aws/smithy-go/transport/http"
)

type withHeader struct {
    header, val string
}

// implements middleware.BuildMiddleware, which runs AFTER a request has been
// serialized and can operate on the transport request
var _ middleware.BuildMiddleware = (*withHeader)(nil)

func (*withHeader) ID() string {
    return "withHeader"
}

func (m *withHeader) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (
    out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
    req, ok := in.Request.(*smithyhttp.Request)
    if !ok {
        return out, metadata, fmt.Errorf("unrecognized transport type %T", in.Request)
    }

    req.Header.Set(m.header, m.val)
    return next.HandleBuild(ctx, in)
}

func WithHeader(header, val string) func (*s3.Options) {
    return func(o *s3.Options) {
        o.APIOptions = append(o.APIOptions, func (s *middleware.Stack) error {
            return s.Build.Add(&withHeader{
                header: header,
                val: val,
            }, middleware.After)
        })
    }
}

func main() {
    cfg, err := config.LoadDefaultConfig(context.Background())
    if err != nil {
        // ...
    }

    svc := s3.NewFromConfig(cfg, WithHeader("x-user-header", "..."))
    // ...
}
```

### 2.11 S3 Transfer Manager

The S3 transfer manager is available for managing uploads and downloads of objects concurrently. This package is located in a Go module outside the S3 client import path. This module can be retrieved by using `go get github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager`.

- `s3.NewUploader` and `s3.NewUploaderWithClient` have been replaced with the constructor method `manager.NewUploader` for creating an Upload manager client.
- `s3.NewDownloader` and `s3.NewDownloaderWithClient` have been replaced with a single constructor method `manager.NewDownloader` for creating a Download manager client.

---

## 3\. Migration Overview

### 3.1 Current System Architecture (v1)

In **COS SDK for Go v1**, the S3 client and related utilities are built around the **session-based architecture**. Key characteristics include:

- **Session-based Initialization:** The `session.NewSession` function creates a shared session object that holds configuration and credentials.

  ```go
  sess, err := session.NewSession(&aws.Config{
  Region: aws.String("us-west-2"),
  })
  if err != nil {
  panic(err)
  }
  svc := s3.New(sess)
  ```

- **S3 operations:**
  - Uploads/Downloads with `s3manager.Uploader` and `s3manager.Downloader`:
  
      ```go
      uploader := s3manager.NewUploader(sess)
      _, err := uploader.Upload(&s3manager.UploadInput{
      Bucket: aws.String("my-bucket"),
      Key:    aws.String("my-key"),
      Body:   file,
      })
      ```

  - Listing with `ListObjectsV2Pages` for pagination:
  
      ```go
      err := svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
      Bucket: aws.String("my-bucket"),
      }, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
      for _, obj := range page.Contents {
      fmt.Println(*obj.Key)
      }
      return !lastPage
      })
      ```

- **Error handling** with `awserr.Error`:

  ```go
  if aerr, ok := err.(awserr.Error); ok {
  fmt.Println(aerr.Code(), aerr.Message())
  }
  ```

### 3.2 Target System Architecture (v2)

The **IBM COS SDK for Go v2** introduces a **modular architecture** and modern Go practices:

- **Modular config** via `config.LoadDefaultConfig` → `s3.NewFromConfig`. Instead of sessions, v2 uses `config.LoadDefaultConfig` to load credentials and region settings, returning a Config object used to create service clients.

  ```go
  cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
  if err != nil {
  panic(err)
  }
  client := s3.NewFromConfig(cfg)
  ```

- **S3 operations**:

  - Uploads/Downloads via `feature/s3/manager`:
  
      ```go
      uploader := manager.NewUploader(client)
      _, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
      Bucket: aws.String("my-bucket"),
      Key:    aws.String("my-key"),
      Body:   file,
      })
      ```

  - Listing objects via typed paginator structs:
  
      ```go
      paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
      Bucket: aws.String("my-bucket"),
      })
      for paginator.HasMorePages() {
      page, err := paginator.NextPage(context.TODO())
      if err != nil {
      panic(err)
      }
      for _, obj := range page.Contents {
      fmt.Println(*obj.Key)
      }
      }
      ```

  - Presigning via a dedicated `PresignClient`:
  
      ```go
      presignClient := s3.NewPresignClient(client)
      url, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
      Bucket: aws.String("my-bucket"),
      Key:    aws.String("my-key"),
      })
      ```

- **Error handling** via `smithy.OperationError` and `smithy.APIError`:

  ```go
  var apiErr smithy.APIError
  if errors.As(err, &apiErr) {
  fmt.Println(apiErr.ErrorCode(), apiErr.ErrorMessage())
  }
  ```

### 3.3 Key Differences at a Glance

- **Mandatory `context.Context`:** All v2 API calls require `context.Context` as the first parameter, enabling cancellation and timeouts.
- **Credentials & providers:** v2 uses the `credentials` and `config` packages with provider interfaces. There are provider constructors such as `credentials.NewStaticCredentialsProvider(accessKey, secretKey, token)` you can pass to `config.LoadDefaultConfig` via `config.WithCredentialsProvider(...)`. The default provider chain behaviour (env vars → shared credentials file → IAM) is preserved but implemented differently.
- **Managers / Utilities:** Upload/download manager moved to `feature/s3/manager`. Multipart upload / download and concurrency controls are supported by the manager package.
- **New Paginator Types:** v2 replaces callback-based pagination (`ListObjectsV2Pages`) with typed paginator structs (`NewListObjectsV2Paginator`).
- **Smithy-based Error Handling:** v2 introduces smithy-go error interfaces (`OperationError`, `APIError`) and modeled service errors for better granularity.
- **Presign Client:** v2 uses a dedicated `PresignClient` for generating presigned URLs, improving clarity and separation of concerns. Create `s3.NewPresignClient(s3Client)` and call `PresignGetObject`/`PresignPutObject` with options such as `s3.WithPresignExpires(...)`. This replaces v1 `req.Presign(...)`.

### 3.4 Migration Examples

> Note: these snippets are minimal, focused on showing the migration pattern. Replace error handling/logging to match your app.

#### 1. Initialize configuration and S3 client (v1 → v2)

##### v1 (common pattern)

```go
// v1
sess, err := session.NewSession(&aws.Config{
    Region: aws.String("us-west-2"),
})
svc := s3.New(sess)
```

##### v2 (recommended)

```go
// v2
import (
    "context"
    "github.com/IBM/ibm-cos-sdk-go-v2/config"
    "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

ctx := context.Background()
cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
if err != nil {
    // handle error
}
s3Client := s3.NewFromConfig(cfg)
```

##### 2. Static credentials (v1 → v2)

##### v1

```go
sess, _ := session.NewSession(&aws.Config{
    Credentials: credentials.NewStaticCredentials("AKID", "SECRET", ""),
})
s3Client := s3.New(sess)
```

##### v2

```go
import (
    "context"
    "github.com/IBM/ibm-cos-sdk-go-v2/config"
    "github.com/IBM/ibm-cos-sdk-go-v2/credentials"
)

// In v2, pass the provider to LoadDefaultConfig:
ctx := context.Background()
cfg, err := config.LoadDefaultConfig(ctx,
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKID", "SECRET", "")),
    config.WithRegion("us-west-2"),
)
if err != nil {
    // handle error
}
s3Client := s3.NewFromConfig(cfg)
```

##### 3. PutObject (upload small object)

##### v1 `code`

```go
_, err := svc.PutObject(&s3.PutObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("path/file.txt"),
    Body:   strings.NewReader("hello"),
})
```

##### v2 `code`

```go
import (
    "github.com/IBM/ibm-cos-sdk-go-v2/aws"
    "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("path/file.txt"),
    Body:   strings.NewReader("hello"),
})
```

Note: v2 methods take an explicit `context.Context` as the first parameter.

##### 4. Upload large file — use manager (uploader)

#### v1 (s3manager upload)

```go
uploader := s3manager.NewUploader(sess)
f, _ := os.Open("bigfile.bin")
result, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("bigfile.bin"),
    Body:   f,
})
```

##### v2 (feature/s3/manager)

```go
import (
    "os"
    "github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager"
)

f, _ := os.Open("bigfile.bin")
uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
    // optional: set concurrency, part size, etc.
    u.PartSize = 5 * 1024 * 1024 // 5MB
})
result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("bigfile.bin"),
    Body:   f,
})
```

##### 5. Download large file — manager downloader

#### v1 (s3manager download)

```go
downloader := s3manager.NewDownloaderWithClient(s3Client)
result, err := downloader.Download(file, &s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("bigfile.bin"),
})
```

##### v2 (feature/s3/manager downloader)

```go
import "github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/manager"

file, _ := os.Create("bigfile.bin")
downloader := manager.NewDownloader(s3Client)
n, err := downloader.Download(context.TODO(), file, &s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("bigfile.bin"),
})
```

##### 6. Presigned GetObject URL

##### v1 `Presign code`

```go
req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("file.txt"),
})
urlStr, err := req.Presign(15 * time.Minute)
```

##### v2 `Presign code`

```go
import (
    "time"
    "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

presignClient := s3.NewPresignClient(s3Client)
presignedReq, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String("file.txt"),
}, s3.WithPresignExpires(15*time.Minute))
if err != nil {
    // handle
}
fmt.Println(presignedReq.URL)
```

##### 7. List buckets, List objects

v2 examples are basically the same call patterns but with `context` and v2 types:

```go
// list buckets
out, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

// list objects (v2 recommends ListObjectsV2)
objs, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
    Bucket: aws.String("my-bucket"),
    Prefix: aws.String("logs/"),
})
```

---

## 4\. Step-by-Step Migration Process

### 4.1 Update Dependencies (Go Modules)

- Replace v1 modules with v2 counterparts in `go.mod`.
- Run `go get` and `go mod tidy`.

The following list are some Go modules provided by the SDK:

| Module | Description |
| :--- | :--- |
| `github.com/IBM/ibm-cos-sdk-go-v2` | The SDK core |
| `github.com/IBM/ibm-cos-sdk-go-v2/config` | Shared Configuration Loading |
| `github.com/IBM/ibm-cos-sdk-go-v2/credentials` | Credential Providers |

The SDK's service clients and higher level utilities modules are nested under the following import paths:

| Import Root | Description |
| :--- | :--- |
| `github.com/IBM/ibm-cos-sdk-go-v2/service/` | Service Client Modules |
| `github.com/IBM/ibm-cos-sdk-go-v2/feature/` | High-Level utilities for S3 services like Transfer Manager |

### 4.2 Refactor S3 Client Initialization - Configuration Loading

- **v1: session-based**

  ```go
  // V1 using NewSession
  import "github.com/IBM/ibm-cos-sdk-go/aws/session"
  // ...
  sess, err := session.NewSession()
  if err != nil {
      // handle error
  }
  ```

- **v2: config-based with `config.LoadDefaultConfig` then `s3.NewFromConfig`**

  ```go
  // V2 using LoadDefaultConfig
  import "context"
  import "github.com/IBM/ibm-cos-sdk-go-v2/config"
  // ...
  cfg, err := config.LoadDefaultConfig(context.TODO())
  if err != nil {
      // handle error
  }
  ```

### 4.3 Update API Calls

- PutObject, GetObject, ListObjectsV2 via v2 methods.
- Use `NewListObjectsV2Paginator` for pagination.
- Use `feature/s3/manager` for multipart transfers.
- Use `s3.NewPresignClient` for presigned URLs.

### 4.4 Handle Context and Error Changes

- All v2 operations take `context.Context`.
- Use `errors.As(...)` to inspect smithy or modeled service errors.

---

## 5\. Handling Errors for Go V2

The IBM COS SDK for Go returns errors that satisfy the Go `error` interface type. Use the `Error()` method to get a formatted string of the SDK error message without special handling. Errors returned by the SDK may implement an `Unwrap` method. The SDK uses `Unwrap` to provide additional context to errors while exposing the underlying error or error chain. Use `Unwrap` with `errors.As` to handle unwrapping error chains.

Always check whether an error occurred after invoking a function or method that can return an `error` interface type. The most basic form of error handling looks like this:

```go
if err != nil {
    // Handle error
    return
}
```

### 5.1 Logging Errors

The simplest form of error handling is traditionally to log or print the error message before returning or exiting from the application.

```go
import "log"

// ...

if err != nil {
    log.Printf("error: %s", err.Error())
    return
}
```

### 5.2 Client Errors

The SDK wraps all errors returned by service clients. `OperationError` provides contextual information about the service name and operation that is associated with an underlying error. This information can be useful for applications that perform batches of operations to one or more services, with a centralized error handling mechanism. Your application can use `errors.As` to access this `OperationError` metadata.

```go
import "log"
import "github.com/aws/smithy-go"

// ...

if err != nil {
    var oe *smithy.OperationError
    if errors.As(err, &oe) {
        log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
    }
    return
}   
```

### 5.3 API Error Responses

All service API response errors implement the `smithy.APIError` interface type. This interface can be used to handle both modeled or un-modeled service error responses. This type provides access to the error code and message returned by the service. Additionally, this type provides indication of whether the fault of the error was due to the client or server if known.

```go
import "log"
import "github.com/aws/smithy-go"

// ...

if err != nil {
    var ae smithy.APIError
    if errors.As(err, &ae) {
        log.Printf("code: %s, message: %s, fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
    }
    return
}
```

### 5.4 Retrieving Request Identifiers

S3 requests contain identifiers that can be used to assist AWS Support with troubleshooting your request. You can use `s3.ResponseError` and call `ServiceRequestID()` and `ServiceHostID()` to retrieve the request ID and host ID.

```go
import "log"
import "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"

// ...

if err != nil {
    var re s3.ResponseError
    if errors.As(err, &re) {
        log.Printf("requestID: %s, hostID: %s request failure", re.ServiceRequestID(), re.ServiceHostID());
    }
    return
}    
```

---

## 6\. Comparison Table

| Area | v1 (ibm-cos-sdk-go) | v2 (ibm-cos-sdk-go-v2) | Notes |
| :--- | :--- | :--- | :--- |
| **Package layout** | `github.com/aws/ibm-cos-sdk-go/service/s3` + `s3/s3manager` | `github.com/IBM/ibm-cos-sdk-go-v2/service/s3` + `feature/s3/manager` | v2 is modular; manager lives under `feature/s3/manager` |
| **Config / Session** | `session.NewSession` / `aws.Config` | `config.LoadDefaultConfig` → `s3.NewFromConfig` | Session replaced by config in v2 |
| **Context usage** | Optional `WithContext` methods | `context.Context` as first parameter | Built-in cancellation and timeouts |
| **Paginators** | `ListObjectsV2Pages` helper | `NewListObjectsV2Paginator` | Typed paginator structs in v2 |
| **Upload/Download manager** | `s3manager.Uploader`/`Downloader` | `manager.Uploader`/`Downloader` | Tune PartSize and Concurrency |
| **Error handling** | `awserr.Error` | `smithy.OperationError` / `smithy.APIError` | Use `errors.As` to unwrap modeled errors |
| **Presigning** | s3 Presign helpers | `s3.NewPresignClient` | Separate presign client |
| **Credentials** | Env/shared config or explicit | Default chain via config; override provider | Prefer env/shared config |
| **Minimum Go version** | Older Go versions supported | Modern Go required | Ensure toolchain meets v2 requirements |
