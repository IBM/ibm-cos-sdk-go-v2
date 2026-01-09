//go:build example
// +build example

package main

/*
import (
	//"bytes"
	"context"
	"fmt"

	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
	"github.com/IBM/ibm-cos-sdk-go-v2/config"
	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
)

func ListBuckets() {

	// create client
	var cosClient *s3.Client
	region := "us-south"
	apiKey := ""
	serviceInstanceID := "crn:v1:bluemix:public:cloud-object-storage:global:a/7a8551ec20234c3f9f246f541907b7:1155e2d-1066-4f0-8390-fda12cb23e9b::"
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
	fmt.Println("Listing all buckets")
	out, err := cosClient.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("ListBuckets error:", err)
		return
	}

	fmt.Println("Buckets:")
	for _, b := range out.Buckets {
		fmt.Println(" -", aws.ToString(b.Name))
	}

	fmt.Println(out)
}
*/
