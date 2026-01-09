module github.com/IBM/ibm-cos-sdk-go-v2/feature/s3/transfermanager

go 1.24.0

toolchain go1.24.4

require (
	github.com/IBM/ibm-cos-sdk-go-v2 v0.0.1
	github.com/IBM/ibm-cos-sdk-go-v2/config v1.29.14
	github.com/IBM/ibm-cos-sdk-go-v2/service/s3 v1.79.3
	github.com/aws/smithy-go v1.22.2
)

require (
	github.com/IBM/go-sdk-core/v5 v5.20.1 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/credentials v1.17.67 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/internal/v4a v1.3.34 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/checksum v1.7.1 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared v1.18.15 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.26.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	go.mongodb.org/mongo-driver v1.17.2 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

replace github.com/IBM/ibm-cos-sdk-go-v2 => ../../../

replace github.com/IBM/ibm-cos-sdk-go-v2/aws => ../../../aws/

replace github.com/IBM/ibm-cos-sdk-go-v2/aws/protocol/eventstream => ../../../aws/protocol/eventstream/

replace github.com/IBM/ibm-cos-sdk-go-v2/config => ../../../config/

replace github.com/IBM/ibm-cos-sdk-go-v2/credentials => ../../../credentials/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/configsources => ../../../internal/configsources/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/endpoints/v2 => ../../../internal/endpoints/v2/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/ini => ../../../internal/ini/

replace github.com/IBM/ibm-cos-sdk-go-v2/internal/v4a => ../../../internal/v4a/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/accept-encoding => ../../../service/internal/accept-encoding/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/checksum => ../../../service/internal/checksum/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/presigned-url => ../../../service/internal/presigned-url/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/internal/s3shared => ../../../service/internal/s3shared/

replace github.com/IBM/ibm-cos-sdk-go-v2/service/s3 => ../../../service/s3/
