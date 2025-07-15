module github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2

go 1.23.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.36.3
	github.com/aws/smithy-go v1.22.2
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../
