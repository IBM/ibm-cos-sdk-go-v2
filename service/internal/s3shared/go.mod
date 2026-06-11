module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared

go 1.25.0

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.0.1
	github.com/aws/smithy-go v1.27.2
)

require github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67 // indirect

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../
