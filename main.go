package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io/ioutil"
	"bytes"
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
	request.Header.Add("X-XSRF-TOKEN", "faM9oN50+uHrpDLUBeBrVg==\\012")
	request.Header.Add("Cookie", "XSRF-TOKEN=\"faM9oN50+uHrpDLUBeBrVg==\\012\"; session=.eJxNkG1rwjAUhf_KuJ-L9H2zIEyGE8oaqailHSOkbazRJpEkLVbxv6-CjH28D-ccHu4NSFXJThjcdayG6AYvJUSQXOtjwXM3dxcGuYjng20jngZfm8ZH18SgZeoXm1OQHBd-fkx9uFtAzkzjM1WYM9EZCpFr2xa0RBtMKsN6ig3jI3ZCz3NfnbfAm4SOE7qhBWOLM62ZFBqi76fDlu8OxXw2G7efZOfF53q5_UfSrG5LsY5Lgew8u2i4_1igKKmxFO2AW9kwAdGetJpa0GmqBHkoQCkV05OBnBiX_TuXSjDRaEPUpJIcLOipetiMyT4Yz4tWe2zkiYq_DxX8E62y6Qpt5qreHoZqsU5227hNRZwVV3v60TzM77-E_XMA.FDtRrQ.uBQ0rx35ci9AQnXv5EEGz8T8rB4")

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
