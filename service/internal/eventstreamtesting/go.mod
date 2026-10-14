module github.com/IBM/ibm-cos-sdk-go-v2/service/internal/eventstreamtesting

go 1.26.0

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v1.1.0
	github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream v1.6.10
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.68
)

require github.com/aws/smithy-go v1.27.10 // indirect

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/
