package s3

//
//func (c *Client) ListBucketsExtended(ctx context.Context, params *ListBucketsExtendedInput, optFns ...func(*Options)) (*ListBucketsExtendedOutput, error) {
//	return &ListBucketsExtendedOutput{}, nil
//}
//
//type BucketExtended struct {
//	CreationDate *time.Time `type:"timestamp"`
//
//	CreationTemplateId *string `type:"string"`
//
//	// Specifies the region where the bucket was created.
//	LocationConstraint *string `type:"string" enum:"BucketLocationConstraint"`
//
//	// The name of the bucket.
//	Name *string `type:"string"`
//
//	noSmithyDocumentSerde
//}
//
//type ListBucketsExtendedInput struct {
//
//	// Sets the IBM Service Instance Id in the request.
//	//
//	// Only Valid for IBM IAM Authentication
//	IBMServiceInstanceId *string `location:"header" locationName:"ibm-service-instance-id" type:"string"`
//
//	// Specifies the bucket to start with when listing all buckets.
//	Marker *string `location:"querystring" locationName:"marker" type:"string"`
//
//	// Sets the maximum number of keys returned in the response. The response might
//	// contain fewer keys but will never contain more.
//	MaxKeys *int64 `location:"querystring" locationName:"max-keys" type:"integer"`
//
//	// Limits the response to buckets that begin with the specified prefix.
//	Prefix *string `location:"querystring" locationName:"prefix" type:"string"`
//
//	noSmithyDocumentSerde
//}
//
//type ListBucketsExtendedOutput struct {
//	Buckets []*BucketExtended `locationNameList:"Bucket" type:"list"`
//
//	// Indicates whether the returned list of buckets is truncated.
//	IsTruncated *bool `type:"boolean"`
//
//	// The bucket at or after which the listing began.
//	Marker *string `type:"string"`
//
//	MaxKeys *int64 `type:"integer"`
//
//	// Container for the owner's display name and ID.
//	Owner *types.Owner `type:"structure"`
//
//	// When a prefix is provided in the request, this field contains the specified
//	// prefix. The result contains only buckets starting with the specified prefix.
//	Prefix *string `type:"string"`
//
//	noSmithyDocumentSerde
//}
