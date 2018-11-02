package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "./grafana"
)

func errorCheck(err error) {
    if err != nil {
        fmt.Println("Error!")
        fmt.Println(err)
    }
}

func main() {
    grafanaUrl := "https://dashboard:3000"
    var dashboard grafana.GetDetails
    httpClient := &http.Client{}

    urlPath := fmt.Sprintf("%s/api/dashboards/uid/%s", grafanaUrl, "9QtTdtFmz")
    request, err := http.NewRequest("GET", urlPath, nil)
    errorCheck(err)

    request.Header.Set("Accept", "application/json")
    request.Header.Set("Content-Type", "application/json")
    response, err := httpClient.Do(request)
    errorCheck(err)

    defer response.Body.Close()
    body,_ := ioutil.ReadAll(response.Body)
    json.Unmarshal(body, &dashboard)

    fmt.Println(dashboard.Dashboard.Panels[0].ID)
}
