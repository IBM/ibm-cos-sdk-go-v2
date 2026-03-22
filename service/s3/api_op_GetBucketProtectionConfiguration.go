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

const opGetBucketProtectionConfiguration = "GetBucketProtectionConfiguration"

func (c *Client) GetBucketProtectionConfiguration(ctx context.Context, params *GetBucketProtectionConfigurationInput, optFns ...func(*Options)) (*GetBucketProtectionConfigurationOutput, error) {
	if params == nil {
		params = &GetBucketProtectionConfigurationInput{}
	}
	result, metadata, err := c.invokeOperation(ctx, opGetBucketProtectionConfiguration, params, optFns, c.addOperationGetBucketProtectionConfigurationMiddlewares)
	if err != nil {
		return nil, err
	}
	// Expecting the result interface to be of type *GetBucketProtectionConfigurationOutput
	// extract and convert value hold by result interface to type *GetBucketProtectionConfigurationOutput
	out := result.(*GetBucketProtectionConfigurationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetBucketProtectionConfigurationInput struct {
	// Bucket is a required field
	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`
	noSmithyDocumentSerde
}

type GetBucketProtectionConfigurationOutput struct {
	IbmProtectionManagementState *string

	// Bucket protection configuration
	ProtectionConfiguration *types.ProtectionConfiguration

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata
}

func (in *GetBucketProtectionConfigurationInput) bindEndpointParams(p *EndpointParameters) {
	p.Bucket = in.Bucket
}

func (c *Client) addOperationGetBucketProtectionConfigurationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	// Add Serialize Custom Middleware
	err = stack.Serialize.Add(&awsRestxml_serializeOpGetBucketProtectionConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	// Add Deserialize Custom Middleware
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpGetBucketProtectionConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	// Add Operation
	if err = addProtocolFinalizerMiddlewares(stack, options, "GetBucketProtectionConfiguration"); err != nil {
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
	// here I have to add Validation for GetBucketProtection
	//if err = addOpGetBucketProtectionConfigurationValidationMiddleware(stack); err != nil {
	//	return err
	//}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetBucketProtectionConfiguration(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addGetBucketProtectionConfigurationUpdateEndpoint(stack, options); err != nil {
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

func (v *GetBucketProtectionConfigurationInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opGetBucketProtectionConfiguration(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: opGetBucketProtectionConfiguration,
	}
}

func getGetBucketProtectionConfigurationBucketMember(input interface{}) (*string, bool) {
	in := input.(*GetBucketProtectionConfigurationInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}

func addGetBucketProtectionConfigurationUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getGetBucketProtectionConfigurationBucketMember,
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
