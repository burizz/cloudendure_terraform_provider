package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// Init Cookie Jar for HTTP client
	//cookieJar, createCookieJarErr := cookiejar.New(nil)
	//if createCookieJarErr != nil {
	//log.Fatalf("Error Initializing cookiejar: %v", createCookieJarErr)
	//}

	// Init HTTP client
	client := http.Client{
		Timeout: 20 * time.Second,
		//Jar:     cookieJar,
	}

	// Define Cloudendure URLs
	cloudEndureApiURL := "https://console.cloudendure.com/api/latest"
	cloudEndureProjectId := "projects/d5aed277-b6fb-4c6c-bedf-bb52799c99f2"
	cloudEndureBluePrintId := "f320947e-1555-4cee-9128-58a6cc4dd99c"

	// Authenticate to Cloudendure
	cookieList, authErr := authCloudEndure(client, cloudEndureApiURL)
	if authErr != nil {
		log.Fatalf("Authentication Error: %s\n", authErr)
	}

	// Get Blueprint definition from API
	getBluePrintErr := getBluePrint(client, cookieList, cloudEndureApiURL, cloudEndureProjectId, cloudEndureBluePrintId)
	if getBluePrintErr != nil {
		log.Fatalf("Get BluePrint Error: %s\n", getBluePrintErr)
	}

}

// Authenticate in Cloudendure - assign cookie and xsrf token to HTTP client
func authCloudEndure(httpClient http.Client, cloudEndureApiURL string) (cookieList []*http.Cookie, authErr error) {
	// TODO: fix passing API Token from var
	var requestBody = []byte(`{"userApiToken":"B212-1445-FBE4-525A-658D-0885-86FD-4510-8192-EDA1-CA50-7738-AAAB-6D5B-A502-1F07"}`)
	fmt.Println("Authenticate into CloudEndure using API Key")

	loginURL := fmt.Sprintf("%s/login", cloudEndureApiURL)
	request, defineRequestErr := http.NewRequest("POST", loginURL, bytes.NewBuffer(requestBody))
	if defineRequestErr != nil {
		return nil, fmt.Errorf("Unable to define HTTP request: %s", defineRequestErr)
	}

	request.Header.Set("Content-Type", "application/json")

	fmt.Println("request URL:", request.URL)
	fmt.Println("request Method:", request.Method)
	fmt.Println("request Headers:", request.Header)
	fmt.Println("")

	response, sendRequestErr := httpClient.Do(request)
	if sendRequestErr != nil {
		return nil, fmt.Errorf("Unable to send HTTP request: %s", sendRequestErr)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	// TODO: DEBUG :
	for _, c := range response.Cookies() {
		fmt.Println("")
		fmt.Println("HERE: ", c)
		fmt.Println("")
	}

	responseBody, readResponseBodyErr := ioutil.ReadAll(response.Body)
	if readResponseBodyErr != nil {
		return nil, fmt.Errorf("Unable to read HTTP response body: %s", readResponseBodyErr)
	}
	fmt.Println("response Body:", string(responseBody))
	fmt.Println("")

	return response.Cookies(), nil
}

// Get Cloudendure Blueprint config by ID
func getBluePrint(httpClient http.Client, cookieList []*http.Cookie, cloudEndureApiURL string, cloudEndureProjectId string, cloudEndureBluePrintId string) error {
	requestURL := fmt.Sprintf("%s/%s/blueprints/%s", cloudEndureApiURL, cloudEndureProjectId, cloudEndureBluePrintId)

	fmt.Printf("Get Cloudendure Blueprint with ID: [%s]\n", cloudEndureBluePrintId)

	request, defineRequestErr := http.NewRequest("GET", requestURL, nil)
	if defineRequestErr != nil {
		return fmt.Errorf("Unable to define HTTP request: %s", defineRequestErr)
	}

	// TODO: Taking cookie from authenticated session doesnt work for some reason
	//for i := range cookieList {
	//fmt.Println("Current Cookie: ", cookieList[i])
	//request.AddCookie(cookieList[i])
	//}

	request.Header.Add("X-XSRF-TOKEN", "7H7e9E5H1kVL9A0QyUpHeg==")
	request.Header.Add("Cookie", "Cookie_1=value; XSRF-TOKEN=\"7H7e9E5H1kVL9A0QyUpHeg==\\012\"; session=.eJxNkNtqwkAURX-lnOcgMTcxILSUVJBmSorRJqUMk2SMo3ORmUnqBf_dCFL6eBZ7bxbnAqSuVSct7jrWQHyBpwpiSM_NrhSFV3iJRR4Sxcl1kcjC92UboHNq0TwLyuU-THdJUOyyAK4OkAMz-EA1Fkx2lkLsua4DnBiLSW1ZT7FlYsDjyPcDP5pG41E4mQRR6MDQEswYpqSB-PvhkIvVtnyZzYbtB1n5i0Mzz_-RbN3wSn4uKoncYn00cP1xQFPSYCX5CXPVMgnxhnBDHegM1ZLcFaBSmpnRieyZUP2zUFoy2RpL9KhWAhzoqb7bDMk-HM6j0Rts1Z7Kvw8ht_XLJX9DeXsk4zL9yJPf_Iuv6mTLy7M7fW3v5tcbjppzDQ.FD3WoQ.qT9ctnkzprYNnWZg4aQmeNIEI9c")

	fmt.Println("request URL:", request.URL)
	fmt.Println("request Method:", request.Method)
	fmt.Println("request Headers:", request.Header)
	fmt.Println("")

	// TODO: DEBUG
	for _, c := range request.Cookies() {
		fmt.Println("")
		fmt.Println("HERE: ", c)
		fmt.Println("")
	}

	response, sendRequestErr := httpClient.Do(request)
	if sendRequestErr != nil {
		return fmt.Errorf("Unable to send HTTP request: %s", sendRequestErr)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	// Used for Debugging HTTP body contents
	//responseBody, readResponseBodyErr := ioutil.ReadAll(response.Body)
	//if readResponseBodyErr != nil {
	//fmt.Printf("Unable to read HTTP response body: %s", readResponseBodyErr)
	//return readResponseBodyErr
	//}
	//fmt.Println("response Body:", string(responseBody))

	// Parse JSON into map
	bluePrintConfigs := make(map[string]interface{}, 0)
	jsonDecodeErr := json.NewDecoder(response.Body).Decode(&bluePrintConfigs)
	if jsonDecodeErr != nil {
		return fmt.Errorf("Unable to decode JSON from response: %s", sendRequestErr)
	}

	// Convert JSON to follow Terraform's expected structure - put JSON map inside an array
	var updatedItems []interface{}
	updatedItems = append(updatedItems, bluePrintConfigs)

	// Set data source schema values
	//terraformResourceData.Set("machine_id", bluePrintConfigs["machineId"])
	//terraformResourceData.Set("instance_type", bluePrintConfigs["instanceType"])
	//terraformResourceData.Set("security_group_ids", bluePrintConfigs["securityGroupIDs"])
	//terraformResourceData.Set("subnet_ids", bluePrintConfigs["subnetIDs"])

	// Add check if resource doesn't exist to set ID to blank
	//if resourceDoesntExist {
	//terraformResourceData.SetID("")
	//return
	//}

	// SetId sets the ID of the resource. If the value is blank, then the resource is destroyed - always run
	//terraformResourceData.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}
