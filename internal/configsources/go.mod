module github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources

go 1.25.0

require github.com/IBM/ibm-cos-sdk-go-v2 v1.0.1

require (
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/aws/smithy-go v1.27.2 // indirect
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../credentials
