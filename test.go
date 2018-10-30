package main

import (
    zts "go.corp.yahoo.com/athens/zts-go-client"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "crypto/tls"
    "io/ioutil"
    "net/http"
    "fmt"
)

func errorCheck(err error){
    if err != nil {
        fmt.Println("an error occured")
        fmt.Println(err)
    }
}

func ZtsClient(ztsUrl, keyFile, certFile string) (*zts.ZTSClient, error) {
    keypem, err := ioutil.ReadFile(keyFile)
    errorCheck(err)

    certpem, err := ioutil.ReadFile(certFile)
    errorCheck(err)

    config, err := tlsConfiguration(keypem, certpem)
    errorCheck(err)
    tr := &http.Transport{
        TLSClientConfig: config,
    }
    client := zts.NewClient(ztsUrl, tr)
    return &client, nil
}

func tlsConfiguration(keypem, certpem []byte) (*tls.Config, error) {
    config := &tls.Config{}
    if certpem != nil && keypem != nil {
        mycert, err := tls.X509KeyPair(certpem, keypem)
        errorCheck(err)

        config.Certificates = make([]tls.Certificate, 1)
        config.Certificates[0] = mycert
    }
    return config, nil
}

func GetTemporaryCredentials(ztsUrl, keyFile, certFile, domainName, roleName string) (*credentials.Credentials, error) {
    ztsClient, err := ZtsClient(ztsUrl, keyFile, certFile)
    errorCheck(err)

    awsCreds, err := ztsClient.GetAWSTemporaryCredentials(zts.DomainName(domainName), zts.CompoundName(roleName))
    errorCheck(err)

    creds := credentials.NewStaticCredentials(awsCreds.AccessKeyId, awsCreds.SecretAccessKey, awsCreds.SessionToken)
    return creds, nil
}

func main() {
    athensDomain := "video-platform.aws-yxs"
    awsRole := "yxs2-celery-worker"
    ztsUrl := "https://zts:4443/zts/v1"
    keyFile := "/home/y/conf/aws_auth/yxs2_priv.pem"
    certFile := "/home/y/conf/aws_auth/yxs2_svc_cert.pem"

    creds, err := GetTemporaryCredentials(ztsUrl, keyFile, certFile, athensDomain, awsRole)
    errorCheck(err)

    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2"),
        Credentials: creds},
    )

    input := &ec2.DescribeTagsInput{
        Filters: []*ec2.Filter{
            {
                Name: aws.String("key"),
                Values: []*string{
                    aws.String("cluster"),
                },
            },
            {
                Name: aws.String("value"),
                Values: []*string{
                    aws.String("yxs2v2-prod-api"),
                },
            },
        },
    }

    svc := ec2.New(sess)
    resp, err := svc.DescribeTags(input)
    errorCheck(err)

    fmt.Println(resp)

}
