package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/ssm"
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
    flag.StringVar(&exec, "exec", "", "shell cmd to exec")
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

func instancesInput(key, value string) *ec2.DescribeInstancesInput {
    input := &ec2.DescribeInstancesInput{
        Filters: []*ec2.Filter{
            {
                Name: aws.String(fmt.Sprintf("tag:%s", key)),
                Values: []*string{aws.String(value)},
            },
            {
                Name: aws.String("instance-state-name"),
                Values: []*string{aws.String("running")},
            },
        },
    }
    return input
}

func listInstances(s *session.Session, input *ec2.DescribeInstancesInput) *ec2.DescribeInstancesOutput {
    svc := ec2.New(s)
    resp, err := svc.DescribeInstances(input)
    errorCheck(err)

    return resp
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
    input := instancesInput(k,v)
    instances := listInstances(sess,input)

    if opts["op"] == "list" {
        for _, i := range instances.Reservations {
            for _, data := range i.Instances {
                fmt.Println(*data.PrivateIpAddress)
            }
        }
    } else if opts["op"] == "exec" {
        var instanceIds []string
        params := make(map[string][]*string)

        for _, i := range instances.Reservations {
            for _, data := range i.Instances {
                instanceIds = append(instanceIds, *data.InstanceId)
            }
        }

        svc := ssm.New(sess)
        cmd := aws.String(opts["cmd"])
        params["commands"] = append(params["commands"], cmd)
        fmt.Println(fmt.Sprintf("Executing: %s", *cmd))

        r,_ := svc.SendCommand(&ssm.SendCommandInput{
            DocumentName: aws.String("AWS-RunShellScript"),
            InstanceIds: aws.StringSlice(instanceIds),
            Parameters: params})

        cmdId := r.Command.CommandId
        var instanceList []*string
        passed := 0
        nopass := 0
        for {
            status, _ := svc.ListCommands(&ssm.ListCommandsInput{
                CommandId: cmdId})
            if len(instanceList) == len(instanceIds) {
                break
            }

            if *status.Commands[0].Status == "Success" {
                fmt.Println(fmt.Sprintf("%s: cmd Success!", *status.Commands[0].InstanceIds[passed]))
                instanceList = append(instanceList, status.Commands[0].InstanceIds[passed])
                passed += 1
            } else if *status.Commands[0].Status == "TimedOut" || *status.Commands[0].Status == "Failed" || *status.Commands[0].Status == "Cancelled" {
                instanceList = append(instanceList, status.Commands[0].InstanceIds[nopass])
                nopass += 1
            }
        }
        if passed != len(instanceList) {
            fmt.Println(fmt.Sprintf("Error while executing: %s", *cmd))
            for _, i := range instanceList{
                fmt.Println(*i)
            }
        }
    }
}
