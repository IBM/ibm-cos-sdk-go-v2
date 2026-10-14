module github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources

go 1.26.0

require github.com/IBM/ibm-cos-sdk-go-v2 v1.1.0

require (
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.68 // indirect
	github.com/aws/smithy-go v1.27.10 // indirect
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../credentials
