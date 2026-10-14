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

const opListBucketReplicationFailures = "ListBucketReplicationFailures"

// ListBucketReplicationFailures lists long-term replication failures for a given bucket.
// Specifically, this lists failures from the per-bucket failure queue and does not list
// short-term failures from the global replication queue that are pending once-per-hour retries.
func (c *Client) ListBucketReplicationFailures(ctx context.Context, params *ListBucketReplicationFailuresInput, optFns ...func(*Options)) (*ListBucketReplicationFailuresOutput, error) {
	if params == nil {
		params = &ListBucketReplicationFailuresInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, opListBucketReplicationFailures, params, optFns, c.addOperationListBucketReplicationFailuresMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*ListBucketReplicationFailuresOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type ListBucketReplicationFailuresInput struct {
	// The bucket name.
	//
	// Bucket is a required field
	Bucket *string `location:"uri" locationName:"Bucket" type:"string" required:"true"`

	// Obfuscated string token for pagination. This is returned on incomplete listing responses.
	ContinuationToken *string `location:"querystring" locationName:"continuation-token" type:"string"`

	// Encoding type to use for Key. Only valid value is url.
	EncodingType types.EncodingType `location:"querystring" locationName:"encoding-type" type:"string"`

	// ISO-8601 / epoch ms time from which to start the listing (exclusive).
	FirstSyncAttemptedBefore *string `location:"querystring" locationName:"first-sync-attempted-before" type:"string"`

	// Maximum number of entries to return (default 1000).
	MaxKeys *int32 `location:"querystring" locationName:"max-keys" type:"integer"`

	noSmithyDocumentSerde
}

func (in *ListBucketReplicationFailuresInput) bindEndpointParams(p *EndpointParameters) {
	p.Bucket = in.Bucket
}

type ListBucketReplicationFailuresOutput struct {
	// Name of bucket being listed.
	Name *string `type:"string"`

	// Value provided via the first-sync-attempted-before request param.
	FirstSyncAttemptedBefore *string `type:"string"`

	// Maximum accepted value is 1000.
	MaxKeys *int32 `type:"integer"`

	// Whether or not the results are truncated.
	IsTruncated *bool `type:"boolean"`

	// Encoding type, if specified in the request.
	EncodingType types.EncodingType `type:"string"`

	// Number of keys returned with this request.
	KeyCount *int32 `type:"integer"`

	// Continuation token for pagination.
	ContinuationToken *string `type:"string"`

	// Next continuation token. Present if this result was truncated.
	NextContinuationToken *string `type:"string"`

	// Metadata about each object returned.
	Contents []types.ObjectSyncAttempted `type:"list"`

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationListBucketReplicationFailuresMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpListBucketReplicationFailures{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpListBucketReplicationFailures{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, opListBucketReplicationFailures); err != nil {
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
	if err = addOpListBucketReplicationFailuresValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opListBucketReplicationFailures(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addListBucketReplicationFailuresUpdateEndpoint(stack, options); err != nil {
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

// ListBucketReplicationFailuresPaginatorOptions is the paginator options for ListBucketReplicationFailures
type ListBucketReplicationFailuresPaginatorOptions struct {
	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool

	// Limit is the maximum number of items to return.
	Limit int32
}

// ListBucketReplicationFailuresPaginator is a paginator for ListBucketReplicationFailures
type ListBucketReplicationFailuresPaginator struct {
	options   ListBucketReplicationFailuresPaginatorOptions
	client    ListBucketReplicationFailuresAPIClient
	params    *ListBucketReplicationFailuresInput
	nextToken *string
	firstPage bool
}

// NewListBucketReplicationFailuresPaginator returns a new ListBucketReplicationFailuresPaginator
func NewListBucketReplicationFailuresPaginator(client ListBucketReplicationFailuresAPIClient, params *ListBucketReplicationFailuresInput, optFns ...func(options *ListBucketReplicationFailuresPaginatorOptions)) *ListBucketReplicationFailuresPaginator {
	if params == nil {
		params = &ListBucketReplicationFailuresInput{}
	}

	options := ListBucketReplicationFailuresPaginatorOptions{}
	if params.MaxKeys != nil {
		options.Limit = *params.MaxKeys
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListBucketReplicationFailuresPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.ContinuationToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListBucketReplicationFailuresPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) > 0)
}

// NextPage retrieves the next ListBucketReplicationFailures page.
func (p *ListBucketReplicationFailuresPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListBucketReplicationFailuresOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.ContinuationToken = p.nextToken

	if p.options.Limit > 0 {
		params.MaxKeys = &p.options.Limit
	}

	optFns = append(optFns,
		addIsPaginatorUserAgent,
	)
	result, err := p.client.ListBucketReplicationFailures(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextContinuationToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}

// ListBucketReplicationFailuresAPIClient is a client that may be implemented for a ListBucketReplicationFailures
type ListBucketReplicationFailuresAPIClient interface {
	ListBucketReplicationFailures(context.Context, *ListBucketReplicationFailuresInput, ...func(*Options)) (*ListBucketReplicationFailuresOutput, error)
}

var _ ListBucketReplicationFailuresAPIClient = (*Client)(nil)

func (v *ListBucketReplicationFailuresInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opListBucketReplicationFailures(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: opListBucketReplicationFailures,
	}
}

func getListBucketReplicationFailuresBucketMember(input interface{}) (*string, bool) {
	in := input.(*ListBucketReplicationFailuresInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}

func addListBucketReplicationFailuresUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getListBucketReplicationFailuresBucketMember,
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
