package initializers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var SVC *s3.S3

const BUCKET_NAME = "gmonitoringdocs"

func ConnectToS3() {

	region := "kz-ast"
	accessKeyID := "GQZ3CU1LKP6AX6IMC3E2"
	secretAccessKey := "Dime61guUW7WrAYgPCGjWJQijegpn935dwbX1AXw"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			"", // token, not needed for static credentials
		),
	})

	if err != nil {
		// Handle error
		panic(err)
	}
	SVC = s3.New(sess)

	_, err = SVC.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(BUCKET_NAME),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeBucketAlreadyExists {
			// Handle the case where the bucket already exists
			fmt.Println("Bucket already exists. Choose a different name.")
		} else {
			// Handle other errors
			panic(err)
		}
	}
}
