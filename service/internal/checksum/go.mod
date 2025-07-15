module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/checksum

go 1.23.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.36.3
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url v1.12.15
	github.com/aws/smithy-go v1.22.2
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
