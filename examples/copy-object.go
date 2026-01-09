//go:build example
// +build example

package main

/*
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

func CopyObject() {
	// create client
	var cosClient *s3.Client
	region := "us-south"
	apiKey := "your-api-key"
	serviceInstanceID := "crn:v1:bluemix:public:cloud-object-storage:global:a/example_crn:999e9d-9999-9d9-9999-fbhgacb23e9b::"
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

	// create a bucket using 'CreateBucket' method
	srcBucket := "demo-bucket" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	dstBucket := "demo-bucket-1" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	// create an object using 'PutObject' method
	srcKey := "sdk-test"
	dstKey := "sdk-test"

	fmt.Printf("Copying s3://%s/%s to s3://%s/%s\n", srcBucket, srcKey, dstBucket, dstKey)

	out, err := cosClient.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(dstBucket),
		Key:        aws.String(dstKey),
		CopySource: aws.String(fmt.Sprintf("%s/%s", srcBucket, srcKey)),
	})
	if err != nil {
		fmt.Println("CopyObject error:", err)
		return
	}

	fmt.Println("Object copied successfully")
	fmt.Println(out)
}
*/
