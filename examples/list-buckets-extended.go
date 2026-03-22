//go:build example
// +build example

package main

import (
	"context"
	"fmt"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/config"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

func ListBucketsExtended() {

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

	prefix := "demo-bucket"

	fmt.Println("Listing buckets with extended information with prefix:", prefix)

	out, err := cosClient.ListBucketsExtended(context.TODO(), &s3.ListBucketsExtendedInput{
		Prefix:               aws.String(prefix),
		IBMServiceInstanceId: aws.String(serviceInstanceID),
		MaxKeys:              aws.Int32(10),
	})
	if err != nil {
		fmt.Println("ListBucketsExtended error:", err)
		return
	}

	fmt.Println("Operation completed successfully")
	fmt.Println(out)
}

func main() {
	ListBucketsExtended()
}
