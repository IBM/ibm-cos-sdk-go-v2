//go:build example
// +build example

package main

/*
import (
	//"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/config"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

func GetObject() {

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
	bucket := "demo-bucket" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	key := "sdk-test"
	fmt.Printf("Downloading object: s3://%s/%s\n", bucket, key)

	out, err := cosClient.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("GetObject error:", err)
		return
	}
	defer out.Body.Close()

	data, _ := io.ReadAll(out.Body)

	fmt.Println("Object downloaded")
	fmt.Printf("Size: %d bytes\n", len(data))
	fmt.Println(out)
}
*/
