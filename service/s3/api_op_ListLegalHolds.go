package s3

import (
	"context"
	"fmt"
	"time"

	awsmiddleware "github.com/IBM/ibm-cos-sdk-go-v2/aws/middleware"
	v4 "github.com/IBM/ibm-cos-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/IBM/ibm-cos-sdk-go-v2/service/s3/internal/customizations"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const opListLegalHolds = "ListLegalHolds"

func (c *Client) ListLegalHolds(ctx context.Context, params *ListLegalHoldsInput, optFns ...func(*Options)) (*ListLegalHoldsOutput, error) {
	if params == nil {
		params = &ListLegalHoldsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, opListLegalHolds, params, optFns, c.addOperationListLegalHoldsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListLegalHoldsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListLegalHoldsInput struct {

	// Bucket is a required field
	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

	// Key is a required field
	Key *string `location:"uri" locationName:"Key" min:"1" type:"string" required:"true"`

	noSmithyDocumentSerde
}

func (in *ListLegalHoldsInput) bindEndpointParams(p *EndpointParameters) {
	p.Bucket = in.Bucket
	p.Key = in.Key
}

type ListLegalHoldsOutput struct {
	CreateTime *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	LegalHolds []types.LegalHold `type:"list"`

	// Retention period to store on the object in seconds. The object can be neither
	// overwritten nor deleted until the amount of time specified in the retention
	// period has elapsed. If this field and Retention-Expiration-Date are specified
	// a 400 error is returned. If neither is specified the bucket's DefaultRetention
	// period will be used. 0 is a legal value assuming the bucket's minimum retention
	// period is also 0.
	RetentionPeriod *int64 `type:"integer"`

	RetentionPeriodExpirationDate *time.Time `type:"timestamp" timestampFormat:"iso8601"`

	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListLegalHoldsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpListLegalHolds{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpListLegalHolds{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, opListLegalHolds); err != nil {
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
	//if err = addOpListLegalHoldsValidationMiddleware(stack); err != nil {
	//	return err
	//}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListLegalHolds(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addListLegalHoldsUpdateEndpoint(stack, options); err != nil {
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

func (v *ListLegalHoldsInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opListLegalHolds(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: opListLegalHolds,
	}
}

// getListLegalHoldsBucketMember returns a pointer to string denoting a provided
// bucket member valueand a boolean indicating if the input has a modeled bucket
// name,
func getListLegalHoldsBucketMember(input interface{}) (*string, bool) {
	in := input.(*ListLegalHoldsInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addListLegalHoldsUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getListLegalHoldsBucketMember,
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
