//go:build example
// +build example

package main

/*
import (
	//"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/config"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

func ListObjects() {

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
	// create an object using 'PutObject' method
	bucket := "demo-bucket" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	fmt.Println("Listing objects in bucket:", bucket)

	out, err := cosClient.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Println("ListObjects error:", err)
		return
	}

	if len(out.Contents) == 0 {
		fmt.Println("No objects found")
		return
	}

	fmt.Println("Objects:")
	for _, obj := range out.Contents {
		fmt.Printf(" - %s (%d bytes)\n", aws.ToString(obj.Key), obj.Size)
	}
	fmt.Println(out)
}
*/
