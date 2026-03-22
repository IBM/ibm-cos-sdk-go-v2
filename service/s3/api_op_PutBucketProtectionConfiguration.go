package s3

import (
	"context"
	"fmt"

	awsmiddleware "github.com/IBM/ibm-cos-sdk-go-v2/aws/middleware"
	v4 "github.com/IBM/ibm-cos-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/IBM/ibm-cos-sdk-go-v2/service/s3/internal/customizations"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const opPutBucketProtectionConfiguration = "PutBucketProtectionConfiguration"

func (c *Client) PutBucketProtectionConfiguration(ctx context.Context, params *PutBucketProtectionConfigurationInput, optFns ...func(*Options)) (*PutBucketProtectionConfigurationOutput, error) {
	if params == nil {
		params = &PutBucketProtectionConfigurationInput{}
	}
	result, metadata, err := c.invokeOperation(ctx, opPutBucketProtectionConfiguration, params, optFns, c.addOperationPutBucketProtectionConfigurationMiddlewares)
	if err != nil {
		return nil, err
	}
	out := result.(*PutBucketProtectionConfigurationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type PutBucketProtectionConfigurationInput struct {
	//CompleteMultipartUploadInput struct

	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

	ProtectionConfiguration *types.ProtectionConfiguration `type:"structure"`

	noSmithyDocumentSerde
}

type PutBucketProtectionConfigurationOutput struct {
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (in *PutBucketProtectionConfigurationInput) bindEndpointParams(p *EndpointParameters) {
	p.Bucket = in.Bucket
}

func (c *Client) addOperationPutBucketProtectionConfigurationMiddlewares(stack *middleware.Stack, options Options) error {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err := stack.Serialize.Add(&awsRestxml_serializeOpPutBucketProtectionConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpPutBucketProtectionConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addProtocolFinalizerMiddlewares(stack, options, "PutBucketProtectionConfiguration"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}
	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addPutBucketContextMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addIsExpressUserAgent(stack); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	if err = addOpPutBucketProtectionConfigurationValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutBucketProtectionConfiguration(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addPutBucketProtectionConfigurationUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSerializeImmutableHostnameBucketMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func (v *PutBucketProtectionConfigurationInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opPutBucketProtectionConfiguration(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: opPutBucketProtectionConfiguration,
	}
}

func getPutBucketProtectionConfigurationBucketMember(input interface{}) (*string, bool) {
	in := input.(*PutBucketProtectionConfigurationInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}

func addPutBucketProtectionConfigurationUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getPutBucketProtectionConfigurationBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}
