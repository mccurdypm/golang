package main

import (
    "os"
    "fmt"
    "bytes"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ssm"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudformation"

)

type DashParams struct {
    UID      string `json:"uid"`
    AsgPanel string `json:"asgPanel"`
}

func errorCheck(err error) {
    if err != nil {
        fmt.Println("Error!")
        fmt.Println(err)
    }
}

func getSession(region string) *session.Session {
    sess, _ := session.NewSession(&aws.Config{
        Region: aws.String(region)},
    )
    return sess
}

func dashboard(status, region string, apiKey *ssm.GetParameterOutput,  params DashParams, data ...map[string]interface{}) map[string]interface{} {
    grafanaUrl := "https://dashboard:3000"
    var result map[string]interface{}
    httpClient := &http.Client{}
    var method string
    var path string
    var a string

    if status == "get" {
        path = fmt.Sprintf("%s/api/dashboards/uid/%s", grafanaUrl, params.UID)
        method = "GET"
    } else if status == "update"{
        path = fmt.Sprintf("%s/api/dashboards/db", grafanaUrl)
        b,_ := json.Marshal(data)
        a = strings.Trim(string(b), "[]") // strip [] from json obj
        method = "POST"
    }
    dataBuffer := bytes.NewBuffer([]byte(a))
    urlPath := fmt.Sprintf(path)
    request, err := http.NewRequest(method, urlPath, dataBuffer)
    errorCheck(err)

    authHeader := fmt.Sprintf("Bearer %s", *apiKey.Parameter.Value)
    request.Header.Set("Authorization", authHeader)
    request.Header.Set("Accept", "application/json")
    request.Header.Set("Content-Type", "application/json")
    response, err := httpClient.Do(request)
    errorCheck(err)

    defer response.Body.Close()
    body,_ := ioutil.ReadAll(response.Body)
    json.Unmarshal(body, &result)
    return result
}

func getKey() *ssm.GetParameterOutput {
    // new sess cos paramStore only in us-west-2
    svc := ssm.New(getSession("us-west-2"))
    param := "grafana_api_key"
    resp,err := svc.GetParameter(&ssm.GetParameterInput{
        Name: &param,
        WithDecryption: aws.Bool(true)},
    )
    errorCheck(err)
    return resp
}

func cfStack(s *session.Session) *string {
    stackName := "prod-worker"
    logicalId := "WorkerASG"

    stackInput := &cloudformation.DescribeStackResourcesInput{
        LogicalResourceId: aws.String(logicalId),
        StackName: aws.String(stackName)}

    svc := cloudformation.New(s)
    resp,err := svc.DescribeStackResources(stackInput)
    errorCheck(err)

    return resp.StackResources[0].PhysicalResourceId
}

func getDashboardParams(s *session.Session) DashParams {
    param := "/prod/dashboard"
    svc := ssm.New(s)
    resp,err := svc.GetParameter(&ssm.GetParameterInput{
        Name: &param})
    errorCheck(err)

    var params DashParams
    body := *resp.Parameter.Value
    json.Unmarshal([]byte(body), &params)

    return params
}

func main(){
    var region string
    var profile string
    if len(os.Args) > 2 {
        region = os.Args[1]
        profile = os.Args[2]
        os.Setenv("AWS_PROFILE", profile)
    } else {
        fmt.Println("Usage: ./updateGrafana <region> <AWS_PROFILE>")
        os.Exit(1)
    }

    sess := getSession(region)
    params := getDashboardParams(sess)
    newAsg := cfStack(sess)
    apiKey := getKey()
    result := dashboard("get", region, apiKey, params)
    panels := result["dashboard"].(map[string]interface{})["panels"]

    for i,_ := range panels.([]interface{}) {
        currentPanel := panels.([]interface{})[i].(map[string]interface{})["id"]
        if fmt.Sprintf("%.0f",currentPanel) == params.AsgPanel {
            panels.([]interface{})[i].(map[string]interface{})["targets"].([]interface{})[0].(map[string]interface{})["dimensions"].(map[string]interface{})["AutoScalingGroupName"] = *newAsg
        }
    }
    result["dashboard"].(map[string]interface{})["panels"] = panels
    update := dashboard("update", region, apiKey, params, result)
    fmt.Println(fmt.Sprintf("Updating %s: %s", region, update["status"]))
}
