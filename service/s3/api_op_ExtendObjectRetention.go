package s3

import (
	"context"
	"fmt"
	"time"

	awsmiddleware "github.com/IBM/ibm-cos-sdk-go-v2/aws/middleware"
	v4 "github.com/IBM/ibm-cos-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/IBM/ibm-cos-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const opExtendObjectRetention = "ExtendObjectRetention"

func (c *Client) ExtendObjectRetention(ctx context.Context, params *ExtendObjectRetentionInput, optFns ...func(*Options)) (*ExtendObjectRetentionOutput, error) {
	if params == nil {
		params = &ExtendObjectRetentionInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, opExtendObjectRetention, params, optFns, c.addOperationExtendObjectRetentionMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ExtendObjectRetentionOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ExtendObjectRetentionInput struct {

	// Additional time, in seconds, to add to the existing retention period for
	// the object. If this field and New-Retention-Time and/or New-Retention-Expiration-Date
	// are specified, a 400 error will be returned. If none of the Request Headers
	// are specified, a 400 error will be returned to the user. The retention period
	// of an object may be extended up to bucket maximum retention period from the
	// time of the request.
	AdditionalRetentionPeriod *int64 `location:"header" locationName:"Additional-Retention-Period" type:"integer"`

	// Bucket is a required field
	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

	// Retention Period in seconds for the object. The Retention will be enforced
	// from the current time until current time + the value in this header. This
	// value has to be within the ranges defined for the bucket.
	ExtendRetentionFromCurrentTime *int64 `location:"header" locationName:"Extend-Retention-From-Current-Time" type:"integer"`

	// Key is a required field
	Key *string `location:"uri" locationName:"Key" min:"1" type:"string" required:"true"`

	NewRetentionExpirationDate *time.Time `location:"header" locationName:"New-Retention-Expiration-Date" type:"timestamp" timestampFormat:"iso8601"`

	NewRetentionPeriod *int64 `location:"header" locationName:"New-Retention-Period" type:"integer"`

	noSmithyDocumentSerde
}

func (in *ExtendObjectRetentionInput) bindEndpointParams(p *EndpointParameters) {

	p.Bucket = in.Bucket
	p.Key = in.Key

}

type ExtendObjectRetentionOutput struct {

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationExtendObjectRetentionMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpExtendObjectRetention{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpExtendObjectRetention{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, opExtendObjectRetention); err != nil {
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
	if err = addRequestChecksumMetricsTracking(stack, options); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	// change validation here
	//if err = addOpExtendObjectRetentionValidationMiddleware(stack); err != nil {
	//	return err
	//}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opExtendObjectRetention(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	//if err = addExtendObjectRetentionInputChecksumMiddlewares(stack, options); err != nil {
	//	return err
	//}
	if err = addExtendObjectRetentionUpdateEndpoint(stack, options); err != nil {
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

func (v *ExtendObjectRetentionInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opExtendObjectRetention(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: opExtendObjectRetention,
	}
}

// getExtendObjectRetentionRequestAlgorithmMember gets the request checksum algorithm value
// provided as input.
//func getExtendObjectRetentionRequestAlgorithmMember(input interface{}) (string, bool) {
//	in := input.(*ExtendObjectRetentionInput)
//	if len(in.ChecksumAlgorithm) == 0 {
//		return "", false
//	}
//	return string(in.ChecksumAlgorithm), true
//}

//func addExtendObjectRetentionInputChecksumMiddlewares(stack *middleware.Stack, options Options) error {
//	return addInputChecksumMiddleware(stack, internalChecksum.InputMiddlewareOptions{
//		GetAlgorithm:                     getExtendObjectRetentionRequestAlgorithmMember,
//		RequireChecksum:                  true,
//		RequestChecksumCalculation:       options.RequestChecksumCalculation,
//		EnableTrailingChecksum:           false,
//		EnableComputeSHA256PayloadHash:   true,
//		EnableDecodedContentLengthHeader: true,
//	})
//}

// getExtendObjectRetentionBucketMember returns a pointer to string denoting a provided
// bucket member value and a boolean indicating if the input has a modeled bucket
// name,
func getExtendObjectRetentionBucketMember(input interface{}) (*string, bool) {
	in := input.(*ExtendObjectRetentionInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addExtendObjectRetentionUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getExtendObjectRetentionBucketMember,
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
