package cloudendure

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlueprint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlueprintCreate,
		//ReadContext:   resourceBlueprintRead,
		//UpdateContext: resourceBlueprintUpdate,
		//DeleteContext: resourceBlueprintDelete,
		//Importer: &schema.ResourceImporter{
		//State: schema.ImportStatePassthrough,
		//},
		Schema: map[string]*schema.Schema{
			"byol_on_dedicated_instance": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "boolean - specifies whether to use byol windows license if dedicated instance tenancy is selected.",
			},
			"dedicated_host_identifier": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Host identifier",
			},
			"disks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"iops": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"throughput": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": &schema.Schema{
							Type:     schema.string,
							Optional: true,
						},
					},
					Description: "array of maps - AWS only. Target machine disk properties",
				},
			},
			"force_uefi": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "boolean - weather to force UEFI",
			},
			"iam_role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - AWS only. Possible values can be fetched from the Region object.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "string - Part of HTTP response - id of created resource",
			},
			"instance_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Possible values can be fetched from the Region object, plus special values COPY_ORIGIN or CUSTOM",
			},
			"machine_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "string - Part of HTTP response - if of created machine object",
			},
			"network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - network interface to be used",
			},
			"placement_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - AWS only. Possible values can be fetched from the Region object.",
			},
			"private_ip_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Valid values: CREATE_NEW; COPY_ORIGIN; CUSTOM_IP; USE_NETWORK_INTERFACE",
			},
			"private_ips": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "array of strings <IP> - List of IP addresses",
			},
			"public_ip_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Valid values: ALLOCATE; DONT_ALLOCATE; AS_SUBNET - Whether to allocate an ephemeral public IP, or not. AS_SUBNET causes CloudEndure to copy this property from the source machine.",
			},
			"recommended_private_ip": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "string - Part of HTTP response - Recommended Priate IP Address",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "string - Part of HTTP response - Region Object ID",
			},
			"run_after_launch": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "boolean - AWS only. Whether to power on the launched target machine after launch. True by default.",
			},
			"security_group_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "array of strings - AWS only. The security groups that will be applied to the target machine. Possible values can be fetched from the Region object.",
			},
			"static_ip": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string <IP> - Possible values can be fetched from the Region object.",
			},
			"static_ip_action": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Valid values: EXISTING; DONT_CREATE; CREATE_NEW; IF_IN_ORIGIN",
			},
			"subnet_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "array of strings - AWS only. Configures a subnets in which the instance network interface will take part. Possible values can be fetched from the Region object.",
			},
			"subnets_host_project": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - GCP only. Host project for corss project network subnet.",
			},
			"tags": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Description: "array of maps - AWS only. Key/Value pair tags that will be applied to the target machine.",
			},
			"tenancy": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Valid values: SHARED; DEDICATED; HOST",
			},
			"scsi_adapter_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "string - Currently relevant for vCenter cloud only. Possible values can be fetched from the Region object.",
			},
			"machine_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "string - Cloudendure machine name",
			},
			"cpus": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer - Number of CPUs per per Target machine; Currently relevant for vCenter cloud only; Max value can be fetched from the maxCpusPerMachine property of the Region object.",
			},
			"mb_ram": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "integer - MB RAM per Target machine; Currently relevant for vCenter cloud only; Max value can be fetched from the maxMbRamPerMachine property of the Region object.",
			},
			"cores_per_cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Description: "integer - Number of CPU cores per CPU in Target machine; Currently relevant for vCenter cloud only."
			},
			"recommended_instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Description: "string - Part of HTTP response - Recommended Instance Type",
			},
			"launch_on_instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "string - Instance id for target machine managed by AMS.",
			},
			"security_group_action": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "FROM_POLICY",
				Description: "string - How to assign a security group to the target machine."

			},
			"compute_location_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "string - TODO",
			},
			"logical_location_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "string - vcenter = vmFolder; relates to $ref LogicalLocation",
			},
			"network_adapter_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "string - Currently relevant for vCenter cloud only. Possible values can be fetched from the Region object.",

			},
			"use_shared_ram": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Description: "boolean - weather to use shared RAM",
			},
		},
	}
}

func resourceBlueprintCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	cloudEndureApiURL := "https://console.cloudendure.com/api/latest"
	cloudEndureProjectId := "projects/d5aed277-b6fb-4c6c-bedf-bb52799c99f2"
	cloudEndureBluePrintId := d.Get("blueprint_id").(string)

	requestURL := fmt.Sprintf("%s/%s/blueprints/%s", cloudEndureApiURL, cloudEndureProjectId, cloudEndureBluePrintId)

	fmt.Printf("Modify Cloudendure Blueprint with ID: [%s]\n", cloudEndureBluePrintId)

	var requestBody = []byte(`{
    "id": "f320947e-1555-4cee-9128-58a6cc4dd99c",
    "machineId": "700628b3-64aa-41c5-a751-e7ed7f4ad8c2",
    "subnetIDs": [
                "subnet-096fff74",
                "subnet-00741c4d"
            ]
	}`)

	req, defineRequestErr := http.NewRequest("PATCH", requestURL, bytes.NewBuffer(requestBody))
	if defineRequestErr != nil {
		return diag.FromErr(defineRequestErr)
	}

	// Init HTTP client
	client := &http.Client{}

	req.Header.Add("X-XSRF-TOKEN", "Tn8cGhv3oCIDVBE17bMV4g==")
	req.Header.Add("Cookie", "Cookie_1=value; XSRF-TOKEN=\"Tn8cGhv3oCIDVBE17bMV4g==\\012\"; session=.eJxNkG1rgzAUhf_KuJ-lRKtdFQoboyuUKTisRccI0aQuapKSROkL_e-zUMY-3odzDg_3CqSu1SAtHgZOIbrCUwURxBfalqLwCm9tEy8RxRmhRKTBR9b4ySW2ySb1y6wL4nbtF23qw80BcuQGH5nGgsvBMog8hBzoibGY1JaPDFsuJuwu5vPAXTyH7mwZhkvXc2BqCW4MV9JA9PVw2In8p3xdrabtB8nn2yPd7P6RdE_7Sn5uK5mgYn8ycPt2QDNCsZL9Gfeq4RKiA-kNc2AwTEtyV4BKaW5mZ9JxocYXobTksjGW6FmtBDgwMn23mZJjMJ0now_Yqo7Jvw_lG98vEFW0TcYU9eu8277HGeVZXqDygsK35m5--wWKiHKp.FD8Odw.5ZpzINgi0Puh7t3OSZErQb0EI0Y")

	response, sendRequestErr := client.Do(req)
	if sendRequestErr != nil {
		return diag.FromErr(sendRequestErr)
	}
	defer response.Body.Close()

	return diags
}
