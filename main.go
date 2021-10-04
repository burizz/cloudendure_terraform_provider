package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	cloudEndureApiURL := "https://console.cloudendure.com/api/latest"
	cloudEndureProjectId := "projects/d5aed277-b6fb-4c6c-bedf-bb52799c99f2"
	cloudEndureBluePrintId := "f320947e-1555-4cee-9128-58a6cc4dd99c"

	cookieList, authErr := authCloudEndure(client, cloudEndureApiURL)
	if authErr != nil {
		fmt.Printf("Authentication Error: %s", authErr)
	}

	bluePrintConfig, getBluePrintErr := getBluePrint(client, cookieList, cloudEndureApiURL, cloudEndureProjectId, cloudEndureBluePrintId)
	if getBluePrintErr != nil {
		fmt.Printf("Get BluePrint Error: %s", getBluePrintErr)
	}

	fmt.Println("")
	fmt.Printf("HERE %s", bluePrintConfig)

	//// Pass map to Terraform ResourceData schema
	//if setResourceDataErr := terraformResourceData.Set("blueprint_config", bluePrintConfig); setResourceDataErr != nil {
	//return diag.FromErr(setResourceDataErr)
	//}

	//// SetId sets the ID of the resource. If the value is blank, then the resource is destroyed.
	//// always run
	//terraformResourceData.SetId(strconv.FormatInt(time.Now().Unix(), 10))
}

// Authenticate in Cloudendure - assign cookie and xsrf token to HTTP client
func authCloudEndure(httpClient *http.Client, cloudEndureApiURL string) (cookieList []*http.Cookie, authErr error) {
	// TODO: fix passing API Token from var
	var requestBody = []byte(`{"userApiToken":"B212-1445-FBE4-525A-658D-0885-86FD-4510-8192-EDA1-CA50-7738-AAAB-6D5B-A502-1F07"}`)
	fmt.Println("Authenticate into CloudEndure using API Key")

	loginURL := fmt.Sprintf("%s/login", cloudEndureApiURL)
	request, defineRequestErr := http.NewRequest("POST", loginURL, bytes.NewBuffer(requestBody))
	if defineRequestErr != nil {
		fmt.Println(defineRequestErr)
		return nil, defineRequestErr
	}

	request.Header.Set("Content-Type", "application/json")

	fmt.Println("request URL:", request.URL)
	fmt.Println("request Method:", request.Method)
	fmt.Println("request Headers:", request.Header)
	fmt.Println("")

	response, sendRequestErr := httpClient.Do(request)
	if sendRequestErr != nil {
		fmt.Println(sendRequestErr)
		return nil, sendRequestErr
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	responseBody, readResponseBodyErr := ioutil.ReadAll(response.Body)
	if readResponseBodyErr != nil {
		fmt.Printf("Unable to read HTTP response body: %s", readResponseBodyErr)
		return nil, readResponseBodyErr
	}
	fmt.Println("response Body:", string(responseBody))
	fmt.Println("")

	return response.Cookies(), nil
}

// Get Cloudendure Blueprint by ID
func getBluePrint(httpClient *http.Client, cookieList []*http.Cookie, cloudEndureApiURL string, cloudEndureProjectId string, cloudEndureBluePrintId string) (bluePrintConfig map[string]interface{}, getBluePrintErr error) {
	requestURL := fmt.Sprintf("%s/%s/blueprints/%s", cloudEndureApiURL, cloudEndureProjectId, cloudEndureBluePrintId)

	fmt.Printf("Get Cloudendure Blueprint with ID: [%s]\n", cloudEndureBluePrintId)

	request, defineRequestErr := http.NewRequest("GET", requestURL, nil)
	if defineRequestErr != nil {
		fmt.Println(defineRequestErr)
	}

	// TODO: Taking cookie from authenticated session doesnt work for some reason
	//for i := range cookieList {
	//request.AddCookie(cookieList[i])
	//}
	request.Header.Add("X-XSRF-TOKEN", "0uc71sHHiyUnwKpBERWSPw==")
	request.Header.Add("Cookie", "Cookie_1=value; XSRF-TOKEN=\"0uc71sHHiyUnwKpBERWSPw==\\012\"; session=.eJxNkFFrwjAUhf_KuM9FYlvdLAgbwzlkBjq1UscIaRtr2iYpSdpppf99FYTt8X6cc_i4V6BpqhppSdPwDIIrPCQQwLrLioOI3dhdWOxiEV8QwiKcfGxzH3dri5ehf9iWk3Wx8OMi9KF3gNbckJppIrhsLIPARciBihpLaGp5y4jlYsDjqed5U4Sm7ujRGz-5MweGluDGcCUNBF93h52IToeX-XzYvpPIW9XZcvePhPusSuTnKpEYxfuzgf7bAc1oRpSsLqRSOZcQHGllmAONYVrSmwIkSnMzutCSC9U-C6Ull7mxVI9SJcCBlumbzZBsJ8N5NvpIrCqZ_PvQe1Tg7q3bLE41i6ImQ_YnLKNNNMZh1qHZa34z738BkWpzZg.FDyqPg.8WtHlqQjTJWDFiCkqbYVfVP7I_E")

	fmt.Println("request URL:", request.URL)
	fmt.Println("request Method:", request.Method)
	fmt.Println("request Headers:", request.Header)
	fmt.Println("")

	response, sendRequestErr := httpClient.Do(request)
	if sendRequestErr != nil {
		fmt.Println(sendRequestErr)
		return nil, sendRequestErr
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	// Used for Debugging HTTP body contents
	//responseBody, readResponseBodyErr := ioutil.ReadAll(response.Body)
	//if readResponseBodyErr != nil {
	//fmt.Printf("Unable to read HTTP response body: %s", readResponseBodyErr)
	//return nil, readResponseBodyErr
	//}
	//fmt.Println("response Body:", string(responseBody))

	// Parse JSON into map
	//items := make(map[string]interface{}, 0)
	//jsonDecodeErr := json.NewDecoder(response.Body).Decode(&bluePrintConfig)
	jsonDecodeErr := json.NewDecoder(response.Body).Decode(&bluePrintConfig)
	if jsonDecodeErr != nil {
		fmt.Println(jsonDecodeErr)
		return nil, jsonDecodeErr
	}

	//// Convert JSON to follow Terraform's expected structure
	//var updatedItems []interface{}
	//updatedItems = append(updatedItems, items)

	return bluePrintConfig, nil
	//return updatedItems, nil
}
