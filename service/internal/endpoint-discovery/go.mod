module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/endpoint-discovery

go 1.24.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.0.0
	github.com/aws/smithy-go v1.24.0
)

require github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67 // indirect

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../
