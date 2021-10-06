package cloudendure

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// summary:
// tested cleanup of schema, seems to work only with set fields
// tested cookie authentication - still lots of issues (cookiejar problems, etc)
// figured out how to show custom debug outputs and errors in terraform
// working on taking api_key on provider level
// merge datasource blueprint read stuff into one function

// TODO:
// move authentication stuff on provider level
// Move URLs on provider level so they dont duplicate in resources and datasources
// implement a resource for changing blueprint
// move hardcoded tokens on provider level to do them only ones while debugging
// debug authentication cookie problems

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "API Key used for authenticating to Cloudendure",
				DefaultFunc: schema.EnvDefaultFunc("CLOUDENDURE_API_KEY", nil),
			},
			"cloudendure_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloudendure API URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudendure_blueprint": resourceBlueprint(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cloudendure_blueprint": dataSourceBlueprint(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Authenticate in Cloudendure - assign cookie and xsrf token to HTTP client
//func authCloudEndure(httpClient *http.Client, apiKey string, cloudEndureApiURL string) (cookieList []*http.Cookie, authErr error) {
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Define Cloudendure URLs
	//cloudEndureApiURL := "https://console.cloudendure.com/api/latest"
	//cloudEndureProjectId := "projects/d5aed277-b6fb-4c6c-bedf-bb52799c99f2"
	//cloudEndureBluePrintId := "f320947e-1555-4cee-9128-58a6cc4dd99c"

	//cloudEndureApiURL, cloudEndureUrlOk := d.Get("cloudendure_url").(string)
	//if cloudEndureApiURL == "" || cloudEndureUrlOk == false {
	//diags = append(diags, diag.Diagnostic{
	//Severity: diag.Warning,
	//Summary:  fmt.Sprintf("Error retrieving CloudEndure URL from Terraform, current value set to: [%v]", cloudEndureApiURL),
	//Detail:   "Cloudendure URL not provided",
	//})
	//}

	cloudEndureApiURL := d.Get("cloudendure_url").(string)

	// TODO: Taking cookie from authenticated session doesnt work for some reason
	//var cookieList []*http.Cookie
	//for i := range cookieList {
	//request.AddCookie(cookieList[i])
	//}

	// Init HTTP client
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	// TODO: check why api_key appears empty here
	// Check if CloudEndure API Key provided
	apiKey, ok := d.Get("api_key").(string)
	if apiKey == "" || ok == false {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("Error retrieving API Key: %v", apiKey),
			Detail:   "Unable to authenticate to Cloudendure API",
		})
	}

	userApiToken := fmt.Sprintf(`{"userApiToken":"%v"}`, apiKey)

	fmt.Println("api_key: ", apiKey)
	fmt.Println("userApiToken: ", userApiToken)

	var requestBody = []byte(userApiToken)

	fmt.Println("Authenticate into CloudEndure using API Key")
	//var requestBody = []byte(`{"userApiToken":""}`)

	loginURL := fmt.Sprintf("%s/login", cloudEndureApiURL)
	request, defineRequestErr := http.NewRequest("POST", loginURL, bytes.NewBuffer(requestBody))
	if defineRequestErr != nil {
		return nil, diag.FromErr(defineRequestErr)
	}

	request.Header.Set("Content-Type", "application/json")

	fmt.Println("request URL:", request.URL)
	fmt.Println("request Method:", request.Method)
	fmt.Println("request Headers:", request.Header)
	fmt.Println("")

	response, sendRequestErr := client.Do(request)
	if sendRequestErr != nil {
		return nil, diag.FromErr(sendRequestErr)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)

	responseBody, readResponseBodyErr := ioutil.ReadAll(response.Body)
	if readResponseBodyErr != nil {
		return nil, diag.FromErr(readResponseBodyErr)
	}
	fmt.Println("response Body:", string(responseBody))
	fmt.Println("")

	return client, diags
}
