//go:build example
// +build example

package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/config"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

func PutObjectRetention() {

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

	fmt.Println("Setting object retention for bucket:", bucket)

	out, err := cosClient.PutObjectRetention(context.TODO(), &s3.PutObjectRetentionInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Println("PutObjectRetention error:", err)
		return
	}

	fmt.Println("Operation completed successfully")
	fmt.Println(out)
}

func main() {
	PutObjectRetention()
}
