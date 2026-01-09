package s3

import (
	"context"
	"fmt"

	awsmiddleware "github.com/IBM/ibm-cos-sdk-go-v2/aws/middleware"
	v4 "github.com/IBM/ibm-cos-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/IBM/ibm-cos-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const opAddLegalHold = "AddLegalHold"

func (c *Client) AddLegalHold(ctx context.Context, params *AddLegalHoldInput, optFns ...func(*Options)) (*AddLegalHoldOutput, error) {
	if params == nil {
		params = &AddLegalHoldInput{}
	}
	result, metadata, err := c.invokeOperation(ctx, opAddLegalHold, params, optFns, c.addOperationAddLegalHoldMiddlewares)
	if err != nil {
		return nil, err
	}
	out := result.(*AddLegalHoldOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type AddLegalHoldInput struct {

	// Bucket is a required field
	Bucket *string

	// Key is a required field
	Key *string

	// RetentionLegalHoldId is a required field
	RetentionLegalHoldId *string

	noSmithyDocumentSerde
}

type AddLegalHoldOutput struct {
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (in *AddLegalHoldInput) bindEndpointParams(p *EndpointParameters) {
	p.Bucket = in.Bucket
}

func (c *Client) addOperationAddLegalHoldMiddlewares(stack *middleware.Stack, options Options) error {

	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	// Add serializer
	err := stack.Serialize.Add(&awsRestxml_serializeOpAddLegalHold{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpAddLegalHold{}, middleware.After)
	if err != nil {
		return err
	}
	// Add Middleware
	if err = addProtocolFinalizerMiddlewares(stack, options, "AddLegalHold"); err != nil {
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
	// Add Validation Middleware
	if err = addOpAddLegalHoldValidationMiddleware(stack); err != nil {
		return err
	}
	// Add Initializer
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opAddLegalHold(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addAddLegalHoldUpdateEndpoint(stack, options); err != nil {
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

func (v *AddLegalHoldInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opAddLegalHold(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: opAddLegalHold,
	}
}

func getAddLegalHoldBucketMember(input interface{}) (*string, bool) {
	in := input.(*AddLegalHoldInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}

func addAddLegalHoldUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getAddLegalHoldBucketMember,
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
