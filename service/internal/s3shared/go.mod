module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared

go 1.24.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v0.0.1
	github.com/aws/smithy-go v1.22.2
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../
