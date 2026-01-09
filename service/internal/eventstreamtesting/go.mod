module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/eventstreamtesting

go 1.24.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v0.0.1
	github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream v1.6.10
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67
)

require github.com/aws/smithy-go v1.22.2 // indirect

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
