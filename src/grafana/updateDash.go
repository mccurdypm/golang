package main

import (
	zts "../ztsclient"
	"./grafana"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ssm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

func getSession(region string, creds *credentials.Credentials) *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds},
	)
	return sess
}

func dashboard(status, region string, apiKey *ssm.GetParameterOutput, params DashParams, data ...grafana.GetDetails) grafana.GetDetails {
	grafanaUrl := "https://dashboard:3000"
	var details grafana.GetDetails
	httpClient := &http.Client{}
	var method string
	var path string
	var auth int
	var a string

	if status == "get" {
		path = fmt.Sprintf("%s/api/dashboards/uid/%s", grafanaUrl, params.UID)
		method = "GET"
	} else if status == "update" {
		path = fmt.Sprintf("%s/api/dashboards/db", grafanaUrl)
		b, _ := json.Marshal(data)
		a = strings.Trim(string(b), "[]") // strip [] from json obj
		method = "POST"
		auth = 1
	}
	dataBuffer := bytes.NewBuffer([]byte(a))
	urlPath := fmt.Sprintf(path)
	request, err := http.NewRequest(method, urlPath, dataBuffer)
	errorCheck(err)

	authHeader := fmt.Sprintf("Bearer %s", *apiKey.Parameter.Value)
	if auth == 1 {
		request.Header.Set("Authorization", authHeader)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	errorCheck(err)

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &details)
	return details
}

func getKey(creds *credentials.Credentials) *ssm.GetParameterOutput {
	// new sess cos paramStore only in us-west-2
	svc := ssm.New(getSession("us-west-2", creds))
	param := "grafana_api_key"
	resp, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           &param,
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
		StackName:         aws.String(stackName)}

	svc := cloudformation.New(s)
	resp, err := svc.DescribeStackResources(stackInput)
	errorCheck(err)

	return resp.StackResources[0].PhysicalResourceId
}

func getDashboardParams(s *session.Session) DashParams {
	param := "/prod/dashboard"
	svc := ssm.New(s)
	resp, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: &param})
	errorCheck(err)

	var params DashParams
	body := *resp.Parameter.Value
	json.Unmarshal([]byte(body), &params)

	return params
}

func main() {
	var region string
	var ztsCfg zts.ZtsConfig
	if len(os.Args) > 1 {
		region = os.Args[1]
	} else {
		fmt.Println("Usage: ./updateGrafana <region>")
		os.Exit(1)
	}

	ztsCfg = zts.ZtsConfig{
		AthensDomain: "",
		AwsRole:      "",
		ZtsUrl:       "zts:4443/zts/v1",
		KeyFile:      "priv.pem",
		CertFile:     "cert.pem",
	}

	creds, err := zts.GetTemporaryCredentials(ztsCfg)
	errorCheck(err)

	sess := getSession(region, creds)
	params := getDashboardParams(sess)
	apiKey := getKey(creds)
	newAsg := cfStack(sess)
	result := dashboard("get", region, apiKey, params)
	panels := result.Dashboard.Panels

	for i, _ := range panels {
		if panels[i].Targets != nil {
			if fmt.Sprintf("%d", panels[i].ID) == params.AsgPanel {
				panels[i].Targets[0].Dimensions.AutoScalingGroupName = *newAsg
			}
		}
	}
	result.Dashboard.Panels = panels
	//update := dashboard("update", region, apiKey, params, result)
	fmt.Println(result)

	/*for i,_ := range panels.([]interface{}) {
	      currentPanel := panels.([]interface{})[i].(map[string]interface{})["id"]
	      if fmt.Sprintf("%.0f",currentPanel) == params.AsgPanel {
	          panels.([]interface{})[i].(map[string]interface{})["targets"].([]interface{})[0].(map[string]interface{})["dimensions"].(map[string]interface{})["AutoScalingGroupName"] = *newAsg
	      }
	  }
	  result["dashboard"].(map[string]interface{})["panels"] = panels
	  update := dashboard("update", region, params, result)
	  fmt.Println(fmt.Sprintf("Updating %s: %s", region, update["status"]))
	*/
}
