//go:build example
// +build example

package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/config"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3/types"
)

func CompleteMultipartUpload() {

	// create client
	var cosClient *s3.Client
	region := "us-south"
	apiKey := "your-api-key"
	serviceInstanceID := "crn:v1:bluemix:public:cloud-object-storage:global:a/example_account_id:example_instance_id::"
	authEndpoint := "https://iam.cloud.ibm.com/identity/token"
	endpoint := "https://s3.us-south.cloud-object-storage.appdomain.cloud"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(ibmiam.NewStaticCredentials(authEndpoint, apiKey, serviceInstanceID)),
		config.WithRegion(region),
		config.WithEndpoint(endpoint))

	if err != nil {
		panic(err)
	}

	cosClient = s3.NewFromConfig(cfg)

	bucket := "demo-bucket" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	key := "sdk-test"

	// Step 1: Create the bucket
	fmt.Printf("Creating bucket: %s\n", bucket)
	_, err = cosClient.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:               aws.String(bucket),
		IBMServiceInstanceId: aws.String(serviceInstanceID),
	})
	if err != nil {
		fmt.Println("CreateBucket error:", err)
		return
	}
	fmt.Println("Bucket created successfully")

	// Step 2: Initiate multipart upload
	fmt.Printf("Initiating multipart upload for: s3://%s/%s\n", bucket, key)
	createResp, err := cosClient.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("CreateMultipartUpload error:", err)
		return
	}
	uploadID := aws.ToString(createResp.UploadId)
	fmt.Printf("Multipart upload initiated with UploadId: %s\n", uploadID)

	// Step 3: Upload parts (each part must be at least 5MB except the last one)
	var completedParts []types.CompletedPart

	// Create parts with minimum 5MB size
	minPartSize := 5 * 1024 * 1024 // 5MB
	part1Data := strings.Repeat("A", minPartSize)
	part2Data := strings.Repeat("B", minPartSize)

	parts := []string{part1Data, part2Data}

	for i, partData := range parts {
		partNumber := int32(i + 1)
		fmt.Printf("Uploading part %d (size: %d bytes)\n", partNumber, len(partData))

		uploadResp, err := cosClient.UploadPart(context.TODO(), &s3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(key),
			PartNumber: aws.Int32(partNumber),
			UploadId:   aws.String(uploadID),
			Body:       strings.NewReader(partData),
		})
		if err != nil {
			fmt.Printf("UploadPart %d error: %v\n", partNumber, err)
			// Abort the multipart upload on error
			cosClient.AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      aws.String(key),
				UploadId: aws.String(uploadID),
			})
			return
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadResp.ETag,
			PartNumber: aws.Int32(partNumber),
		})
		fmt.Printf("Part %d uploaded with ETag: %s\n", partNumber, aws.ToString(uploadResp.ETag))
	}

	// Step 4: Complete the multipart upload
	fmt.Printf("Completing multipart upload: s3://%s/%s\n", bucket, key)
	out, err := cosClient.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		fmt.Println("CompleteMultipartUpload error:", err)
		return
	}

	fmt.Println("Multipart upload completed successfully")
	fmt.Printf("Location: %s\n", aws.ToString(out.Location))
	fmt.Printf("ETag: %s\n", aws.ToString(out.ETag))
}

func main() {
	CompleteMultipartUpload()
}
