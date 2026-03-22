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
)

func AddLegalHold() {

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
	legalHoldID := "legal-hold-1"

	// Step 1: Create the bucket with Object Lock enabled (required for legal holds)
	fmt.Printf("Creating bucket with Object Lock enabled: %s\n", bucket)
	_, err = cosClient.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket:                     aws.String(bucket),
		IBMServiceInstanceId:       aws.String(serviceInstanceID),
		ObjectLockEnabledForBucket: aws.Bool(true),
	})
	if err != nil {
		fmt.Println("CreateBucket error:", err)
		return
	}
	fmt.Println("Bucket created successfully with Object Lock enabled")

	// Step 2: Upload an object to the bucket
	fmt.Printf("Uploading object: %s\n", key)
	_, err = cosClient.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader("Sample content for legal hold test"),
	})
	if err != nil {
		fmt.Println("PutObject error:", err)
		return
	}
	fmt.Println("Object uploaded successfully")

	// Step 3: Add legal hold to the object
	fmt.Printf("Adding legal hold to object: s3://%s/%s\n", bucket, key)
	out, err := cosClient.AddLegalHold(context.TODO(), &s3.AddLegalHoldInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		RetentionLegalHoldId: aws.String(legalHoldID),
	})
	if err != nil {
		fmt.Println("AddLegalHold error:", err)
		return
	}

	fmt.Println("Legal hold added successfully")
	fmt.Println(out)
}

func main() {
	AddLegalHold()
}
