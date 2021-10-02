package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io/ioutil"
	"bytes"
	"log"
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
		log.Fatal(authErr)
	}

	getBluePrints(client, cookieList, cloudEndureApiURL, cloudEndureProjectId, cloudEndureBluePrintId)
}

// Authenticate in Cloudendure - assign cookie and xsrf token to HTTP client
func authCloudEndure(httpClient *http.Client, cloudEndureApiURL string) (cookieList []*http.Cookie, authErr error) {
	// TODO: fix passing API Token from var
	var requestBody = []byte(`{"userApiToken":""}`)
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
    responseBody, _ := ioutil.ReadAll(response.Body)
    fmt.Println("response Body:", string(responseBody))

	fmt.Println("")

	return response.Cookies(), nil
}

// Get Cloudendure Blueprint by ID
func getBluePrints(httpClient *http.Client, cookieList []*http.Cookie, cloudEndureApiURL string, cloudEndureProjectId string, cloudEndureBluePrintId string) error {
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
	request.Header.Add("X-XSRF-TOKEN", "v3+advWcP7hISr+0l3+xUA==\\012")
	request.Header.Add("Cookie", "XSRF-TOKEN=\"v3+advWcP7hISr+0l3+xUA==\\012\"; session=.eJxNkG1rgzAURv_KuJ-lpFEHCoWNUcq6GbC0FjtGiJraWJOUJLq-0P8-hTL27T6Hex8O9wasLHWnHO06UUF8g6cCYkiuVbOTOc7x3BFMZH5BiMg0_FzXAbkmjizSYLc-hkkzD_ImDeDuATsJS0_cUClU5zjEGCEPWmYdZaUTPadOyAFPn31_GqEoCibhMOHAg-FKCmuFVhbir4fDRmaH3etsNnQ_SOYvT9Vi84-k26ot1GpZKILy7dnC_dsDw1lFtWovtNW1UBDvWWu5B53lRrFRAQpthJ1c2FFI3b9IbZRQtXXMTEotwYOem9Fm2OzDIZ6t2VOnj1z9fahqEpNvVzjDJCX4sNz47yZZ_Fw__EOWrlH0Vo_m9199K3Jx.FDoV0g.BTB2IiqdfoBCwsjP5EgAjSwYPac")

	fmt.Println("request URL:", request.URL)
	fmt.Println("request Method:", request.Method)
    fmt.Println("request Headers:", request.Header)
	fmt.Println("")

	response, sendRequestErr := httpClient.Do(request)
	if sendRequestErr != nil {
		fmt.Println(sendRequestErr)
		return sendRequestErr
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
    fmt.Println("response Headers:", response.Header)
    responseBody, _ := ioutil.ReadAll(response.Body)
    fmt.Println("response Body:", string(responseBody))

	// Parse JSON into map
	items := make([]map[string]interface{}, 0)
	jsonDecodeErr := json.NewDecoder(response.Body).Decode(&items)
	if jsonDecodeErr != nil {
		fmt.Println(jsonDecodeErr)
		return jsonDecodeErr
	}
	return nil
}
