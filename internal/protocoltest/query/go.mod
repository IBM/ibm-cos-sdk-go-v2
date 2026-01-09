module github.com/IBM/ibm-cos-sdk-go-v2/internal/protocoltest/query

go 1.24.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v0.0.1
	github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources v1.3.34
	github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 v2.6.34
	github.com/aws/smithy-go v1.22.2
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/
