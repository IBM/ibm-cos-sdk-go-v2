# IBM Cloud Object Storage - Go SDK v2 Beta Release

Welcome to the Go SDK v2 Beta! This release is an early version of our SDK and is intended for testing and feedback purposes.

This package allows Go developers to write software that interacts with [IBM
Cloud Object Storage](https://www.ibm.com/cloud/object-storage).  It is a fork of the [``AWS SDK for Go v2``](https://github.com/aws/aws-sdk-go-v2) library and can stand as a drop-in replacement if the application needs to connect to object storage using an S3-like API and does not make use of other AWS services.

***

## Important Notes

This is a **beta release** — APIs may change before the GA release and use in production environments is **not recommended** at this time.

***

## Feedback & Issue Reporting

We value your feedback! Please help us improve by:

Please report any bugs or issues in the [GitHub Issues](https://github.com/ibm/ibm-cos-sdk-go-v2/issues/new) section. We also welcome suggestions for enhancements or reports of unexpected behavior.

## Notice

IBM has added a [Language Support Policy](#language-support-policy). Language versions will be deprecated on the published schedule without additional notice.

## Documentation

* [Core documentation for IBM COS](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-getting-started-cloud-object-storage)
* [Go API reference documentation](https://ibm.github.io/ibm-cos-sdk-go-v2/)
* [REST API reference documentation](https://cloud.ibm.com/docs/services/cloud-object-storage?topic=cloud-object-storage-compatibility-api)

For release notes, see the [CHANGELOG](CHANGELOG.md).

* [Getting the SDK](#getting-the-sdk)
* [Example code](#example-code)
* [Getting help](#getting-help)

## Quick start

You'll need:

* An instance of COS.
* An API key from [IBM Cloud Identity and Access Management](https://cloud.ibm.com/docs/account?topic=account-userapikey&interface=ui) with at least `Writer` permissions.
* The ID of the instance of COS that you are working with.
* Token acquisition endpoint
* Service endpoint

These values can be found in the IBM Cloud Console by [generating a 'service credential'](https://cloud.ibm.com/docs/services/cloud-object-storage/iam?topic=cloud-object-storage-service-credentials#service-credentials).

## Archive Tier Support

You can automatically archive objects after a specified length of time or after a specified date. Once archived, a temporary copy of an object can be restored for access as needed. Restore time may take up to 15 hours.

An archive policy is set at the bucket level by calling the ``PutBucketLifecycleConfiguration`` method on a client instance. A newly added or modified archive policy applies to new objects uploaded and does not affect existing objects. For more detail, see the [documentation](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-using-go).

## Immutable Object Storage

Users can configure buckets with an Immutable Object Storage policy to prevent objects from being modified or deleted for a defined period of time. The retention period can be specified on a per-object basis, or objects can inherit a default retention period set on the bucket. It is also possible to set open-ended and permanent retention periods. Immutable Object Storage meets the rules set forth by the SEC governing record retention, and IBM Cloud administrators are unable to bypass these restrictions. For more detail, see the [IBM Cloud documentation](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-using-go).

Note: Immutable Object Storage does not support Aspera transfers via the SDK to upload objects or directories at this stage.

## Accelerated Archive

Users can set an archive rule that would allow data restore from an archive in 2 hours or 12 hours.

## Getting the SDK

To begin using the SDK, initialize your project with Go modules and install the required dependencies using go get. The following example demonstrates how to create a bucket in IBM Cloud Object Storage using the v2 SDK and S3-compatible APIs.  The SDK requires a minimum version of Go 1.23 or newer.

### Initialize Project

```sh
mkdir ~/ibmcos
cd ~/ibmcos
go mod init ibmcos
```

### Add SDK Dependencies

```sh
go get github.com/IBM/ibm-cos-sdk-go-v2/config
go get github.com/IBM/ibm-cos-sdk-go-v2/service/s3
go get github.com/IBM/ibm-cos-sdk-go-v2/credentials
```

### Example code

In your preferred editor, copy the code below into `main.go`, and modify the API key `API_KEY` and instance ID `RESOURCE_INSTANCE_ID` fields with your credentials.

```go
package main

import ( "context"
 "fmt"
 "log"
 "strconv"
 "time"

 "github.com/IBM/ibm-cos-sdk-go-v2/aws"
 "github.com/IBM/ibm-cos-sdk-go-v2/config"
 "github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
 "github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

// -------------------------------------------------------------------
// IBM COS Example - Create Bucket using API KEY with Go v2 SDK
// -------------------------------------------------------------------

const (
 region            = "us-south"
 apiKey            = "API_KEY"
 serviceInstanceID = "RESOURCE_INSTANCE_ID"
 authEndpoint      = "https://iam.cloud.ibm.com/identity/token"
 endpoint          = "https://s3.us-south.cloud-object-storage.appdomain.cloud"
)

func main() {
 
 // Load configuration using IAM Auth (API key)
 cfg, err := config.LoadDefaultConfig(context.TODO(),
  config.WithCredentialsProvider(ibmiam.NewStaticCredentials(authEndpoint, apiKey, serviceInstanceID)),
  config.WithEndpoint(endpoint),
  config.WithRegion(region),
 )
 if err != nil {
  log.Fatalf("Failed to load configuration: %v", err)
 }

 client := s3.NewFromConfig(cfg)

 // Create unique bucket name
 ms := time.Now().UnixMilli()
 bucketName := "v2-go-bucket-" + strconv.FormatInt(ms, 10)

 // Create bucket input
 createBucketInput := &s3.CreateBucketInput{
  Bucket: aws.String(bucketName),
 }

 // Execute bucket creation
 _, err = client.CreateBucket(context.Background(), createBucketInput)
 if err != nil {
  log.Fatalf("Failed to create bucket: %v", err)
 }
 fmt.Println("Bucket Created Successfully!")
 fmt.Printf("Bucket Name: %s\n", bucketName)
}
```

***

### Compile and Execute

```sh
go run main.go
```

More examples can be found [here](./examples).

***

## Getting Help

Feel free to use GitHub issues for tracking bugs and feature requests, but for help please use one of the following resources:

* Read a quick start guide in [IBM Cloud Docs](https://cloud.ibm.com/docs/cloud-object-storage?topic=cloud-object-storage-using-go-v2).
* Ask a question on [Stack Overflow](https://stackoverflow.com/questions/tagged/object-storage+ibm) and tag it with `ibm` and `object-storage`.
* Open a support ticket with [IBM Cloud Support](https://cloud.ibm.com/unifiedsupport/supportcenter/)
* If it turns out that you may have found a bug, please [open an issue](https://github.com/ibm/ibm-cos-sdk-go-v2/issues/new).

## Language Support Policy

IBM supports [current public releases](https://golang.org/doc/devel/release.html). IBM will deprecate language versions 90 days after a version reaches end-of-life. All clients will need to upgrade to a supported version before the end of the grace period.

## License

This SDK is distributed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
see LICENSE.txt and NOTICE.txt for more information.
