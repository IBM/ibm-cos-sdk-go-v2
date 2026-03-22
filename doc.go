// Package sdk
// ibm-cos-sdk-go-v2 is the the v2 of the IBM COS SDK for the Go programming language.
//
// # Getting started
//
// The best way to get started working with the SDK is to use `go get` to add the
// SDK and desired service clients to your Go dependencies explicitly.
//
//	go get github.com/IBM/ibm-cos-sdk-go-v2
//	go get github.com/IBM/ibm-cos-sdk-go-v2/config
//
// # Hello COS
//
// # This example shows how you can use the v2 SDK to make an API request
//
// package main
//
// import (
//
//	"context"
//	"fmt"
//	"log"
//	"strconv"
//	"time"
//
//	"github.com/IBM/ibm-cos-sdk-go-v2/aws"
//	"github.com/IBM/ibm-cos-sdk-go-v2/config"
//	"github.com/IBM/ibm-cos-sdk-go-v2/credentials/ibmiam"
//	"github.com/IBM/ibm-cos-sdk-go-v2/service/s3"
//
// )
//
// // -------------------------------------------------------------------
// // IBM COS Example - Create Bucket using API KEY with Go v2 SDK
// // -------------------------------------------------------------------
//
// const (
//
//	region            = "us-south"
//	apiKey            = "API_KEY"
//	serviceInstanceID = "RESOURCE_INSTANCE_ID"
//	authEndpoint      = "https://iam.cloud.ibm.com/identity/token"
//	endpoint          = "https://s3.us-south.cloud-object-storage.appdomain.cloud"
//
// )
//
// func main() {
//
//		// Load configuration using IAM Auth (API key)
//		cfg, err := config.LoadDefaultConfig(context.TODO(),
//			config.WithCredentialsProvider(ibmiam.NewStaticCredentials(authEndpoint, apiKey, serviceInstanceID)),
//			config.WithEndpoint(endpoint),
//			config.WithRegion(region),
//		)
//		if err != nil {
//			log.Fatalf("Failed to load configuration: %v", err)
//		}
//
//		client := s3.NewFromConfig(cfg)
//
//		// Create unique bucket name
//		ms := time.Now().UnixMilli()
//		bucketName := "v2-go-bucket-" + strconv.FormatInt(ms, 10)
//
//		// Create bucket input
//		createBucketInput := &s3.CreateBucketInput{
//			Bucket: aws.String(bucketName),
//		}
//
//		// Execute bucket creation
//		_, err = client.CreateBucket(context.Background(), createBucketInput)
//		if err != nil {
//			log.Fatalf("Failed to create bucket: %v", err)
//		}
//		fmt.Println("Bucket Created Successfully!")
//		fmt.Printf("Bucket Name: %s\n", bucketName)
//	}
package sdk
