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

func PutObject() {

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
	content := "Hello, IBM Cloud Object Storage!"

	fmt.Printf("Uploading object to s3://%s/%s\n", bucket, key)

	out, err := cosClient.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
	})
	if err != nil {
		fmt.Println("PutObject error:", err)
		return
	}

	fmt.Println("Object uploaded successfully")
	fmt.Printf("ETag: %s\n", aws.ToString(out.ETag))
}

func main() {
	PutObject()
}
