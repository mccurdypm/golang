package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/aws/session"
    "strings"
    "flag"
    "fmt"
    "os"
)

func errorCheck(err error){
    if err != nil {
        fmt.Println("an error occured")
        fmt.Println(err)
    }
}

func usage() {
    fmt.Println(`
Usage: ec2 --profilel$=<AWS_PROFILE> --region=<AWS_REGION> [optional] --tag=<KEY>:<VALUE> [ --list / --exec=<SHELL CMD> ]
For e.g:

    ec2 --profile=saml --region=us=west=2 --tag=cluster:yxs2v2-prod-api --exec="service uwsgi restart"
    `)
    os.Exit(1)
}

func parseOpts() map[string]string {
    var awsRegion, awsProfile, tag, exec string
    var list *bool

    flag.StringVar(&awsProfile, "profile", "", "Your AWS profile for creds")
    flag.StringVar(&awsRegion, "region", "", "Your AWS region")
    flag.StringVar(&tag, "exec", "", "shell cmd to exec")
    flag.StringVar(&tag, "tag", "", "Host tag")
    list = flag.Bool("list", false, "List host by tags")
    flag.Parse()

    if awsProfile == "" || awsRegion == "" {
        usage()
    }

    opts := make(map[string]string)
    opts["profile"] = awsProfile
    opts["region"] = awsRegion
    opts["tag"] = tag

    if *list == true {
        opts["op"] = "list"
    } else if exec != "" {
        opts["op"] = "exec"
        opts["cmd"] = exec
    }

    return opts
}

func tagsInput(key, value string) *ec2.DescribeTagsInput {
    input := &ec2.DescribeTagsInput{
        Filters: []*ec2.Filter{
            {
                Name: aws.String("key"),
                Values: []*string{aws.String(key)},
            },
            {
                Name: aws.String("value"),
                Values: []*string{aws.String(value)},
            },
        },
    }
    return input
}

func getInstanceIds(s *session.Session, input *ec2.DescribeTagsInput) []string {
    svc := ec2.New(s)
    resp, err := svc.DescribeTags(input)
    errorCheck(err)

    var instanceIds []string
    for _, data := range resp.Tags {
        instanceIds = append(instanceIds, *data.ResourceId)
    }
    return instanceIds
}

func getSession(region string) *session.Session {
    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    return sess
}

func main() {
    opts := parseOpts()
    os.Setenv("AWS_PROFILE", opts["profile"])
    t := strings.Split(opts["tag"], ":")
    k, v := t[0], t[1]

    sess := getSession(opts["region"])
    svc := ec2.New(sess)
    input := tagsInput(k,v)

    if opts["op"] == "list" {
        instanceIds := getInstanceIds(svc,input)
        //fmt.Println(instanceIds)
    }
}
