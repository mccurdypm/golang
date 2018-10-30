package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ssm"
    "github.com/aws/aws-sdk-go/aws/session"
    "flag"
    "fmt"
    "os"
)

func errorCheck(err error){
    if err != nil {
        fmt.Print("an error occured")
    }
}

func main() {
    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(os.Args[1])},
    )

    param := "/prod/routing_policy"

    svc := ssm.New(sess)
    resp, _ := svc.GetParameter(&ssm.GetParameterInput{
        Name: &param,
        },
    )

    fmt.Println(resp)

}
