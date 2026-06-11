module github.com/IBM/ibm-cos-sdk-go-v2/internal/protocoltest/smithyrpcv2cbor

go 1.25.0

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.0.1
	github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources v1.3.34
	github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 v2.6.34
	github.com/aws/smithy-go v1.27.2
)

require github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67 // indirect

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/
