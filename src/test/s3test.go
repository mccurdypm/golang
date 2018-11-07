package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "flag"
    "fmt"
    "log"
    "os"
)

func errorCheck(err error) {
    if err != nil{
        log.Println("error!")
    }
}

func usage() {
    fmt.Println(`
Usage: s3test --profilel$=<AWS_PROFILE> --region=<AWS_REGION> [optional] --bucket=<AWS_BUCKET>
Note: Not specifying bucket will list all buckets in the AWS Account
For e.g:

    List ALL bucket:
    s3test --profile=saml --region=us-west-2

    List bucket objects
    s3test --profile=saml --=region=us=west=2 --bucket=yxs2-us-west-2-settings
    `)
    os.Exit(1)
}

func parseOpts() map[string]string {
    var awsRegion, awsProfile, bucket string

    flag.StringVar(&awsProfile, "profile", "", "Your AWS profile for creds")
    flag.StringVar(&awsRegion, "region", "", "Your AWS region")
    flag.StringVar(&bucket, "bucket", "", "Your AWS S3 bucket")
    flag.Parse()

    if awsProfile == "" || awsRegion == "" {
        usage()
    }

    opts := make(map[string]string)
    opts["profile"] = awsProfile
    opts["region"] = awsRegion
    opts["bucket"] = bucket
    return opts
}

func main() {
    opts := parseOpts()
    fmt.Println("Profile: ", opts["profile"])
    os.Setenv("AWS_PROFILE", opts["profile"])

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(opts["region"])},
    )
    errorCheck(err)

    // Create S3 service client
    svc := s3.New(sess)

    if opts["bucket"] == "" {
        result, err := svc.ListBuckets(&s3.ListBucketsInput{})
        errorCheck(err)

        log.Println("Buckets:")

        for _, bucket := range result.Buckets {
            log.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
        }
    } else if opts["bucket"] != "" {
        s3Bucket := opts["bucket"]
        resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(s3Bucket)})
        errorCheck(err)

        for _, item := range resp.Contents {
            fmt.Println("Name:         ", *item.Key)
            fmt.Println("Last modified:", *item.LastModified)
            fmt.Println("Size:         ", *item.Size)
            fmt.Println("Storage class:", *item.StorageClass)
            fmt.Println("")
        }

        fmt.Println("Found", len(resp.Contents), "items in bucket", s3Bucket)
        fmt.Println("")
    }
}
