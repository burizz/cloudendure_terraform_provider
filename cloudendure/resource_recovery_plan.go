package cloudendure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"terraform-provider-cloudendure/models"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecoveryPlan() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Recovery plan name.",
			},
			"steps": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
						Elem: &schema.Schema{Type: schema.TypeString},
					},
				},
				Optional:    true,
				Description: "A set of recovery plan steps.",
			},
		},
		CreateContext: resourceRecoveryPlanCreate,
		ReadContext:   resourceRecoveryPlanRead,
		UpdateContext: resourceRecoveryPlanUpdate,
		DeleteContext: resourceRecoveryPlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Pet represents a single pet resource.",
	}
}

func marshalRecoveryPlanInput(d *schema.ResourceData) *models.RecoveryPlanInput {
	var securityGroupIDs []string
	for _, securityGroupValue := range d.Get("security_group_ids").([]interface{}) {
		securityGroupIDs = append(securityGroupIDs, securityGroupValue.(string))
	}

	// TODO: fix this
	return &models.RecoveryPlanInput{
		// requried
		Name:  d.Get("machine_name").(string),
		//Steps: ,
	}
}

func resourceRecoveryPlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := meta.(*http.Client)

	cloudEndureApiURL := "https://console.cloudendure.com/api/latest"
	cloudEndureProjectId := "projects/d5aed277-b6fb-4c6c-bedf-bb52799c99f2"

	requestURL := fmt.Sprintf("%s/%s/recoveryPlans", cloudEndureApiURL, cloudEndureProjectId)

	requestBody, jsonMarshalErr := json.Marshal(marshalRecoveryPlanInput(d))
	if jsonMarshalErr != nil {
		return diag.FromErr(jsonMarshalErr)
	}

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("Debug MSG"),
		Detail:   fmt.Sprintf("Request Body: %v", string(requestBody)),
	})

	req, defineRequestErr := http.NewRequest("POST", requestURL, bytes.NewBuffer(requestBody))
	if defineRequestErr != nil {
		return diag.FromErr(defineRequestErr)
	}

	req.Header.Add("X-XSRF-TOKEN", "3kZPD1PXwmipfaJn1Wuljg==")
	req.Header.Add("Cookie", "Cookie_1=value; XSRF-TOKEN=\"3kZPD1PXwmipfaJn1Wuljg==\\012\"; session=.eJxNkG1rgzAURv_KuJ-lWN_WCoWNUQvdKjhqnY4RoqY21iRdEp1t8b_PQmGD--U53PtwuFfARSFarlHb0hL8Kzzk4MPmUtYZS63UWurQCll6Ns2QRe7btnLCy0aHq8jJtkd3Uy-dtI4cGAzAJ6rQiUjEKG81Ad8yTQMarDTChaYdQZqyEU892_bm4zxOPG82m04NGK8YVYoKrsD_vDvEbHfInheLsftOdvb6VK7ifyRKyibn7-uch2aa9AqGLwMkwSUSvDmjRlSUg7_HjSIGtIpIjm8KkAtJ1eSMj5SJ7okJySmvlMZyUggGBnRE3mzGzc4dY6_kHmlxJPzvQ5bG8TLq4-Bg50nzk7HgNa-DjzLpv7OLOX-pbubDL6CUc4U.FEHBgQ.41V_rVmKg__CNId2VTkLyZAUgBs")

	response, sendRequestErr := client.Do(req)
	if sendRequestErr != nil {
		return diag.FromErr(sendRequestErr)
	}
	defer response.Body.Close()

	responseBody, readResponseBodyErr := ioutil.ReadAll(response.Body)
	if readResponseBodyErr != nil {
		fmt.Printf("Unable to read HTTP response body: %s", readResponseBodyErr)
		return diag.FromErr(readResponseBodyErr)
	}

	if response.StatusCode != 200 && response.StatusCode != 201 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Cloudendure HTTP Response: %v", response.Status),
			Detail:   fmt.Sprintf("Error in create blueprint API request: %v", string(responseBody)),
		})
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}

func resourceRecoveryPlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceRecoveryPlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceRecoveryPlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
