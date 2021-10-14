terraform {
  required_providers {
    cloudendure = {
      version = "0.1"
      source  = "hashicorp.com/edu/cloudendure"
    }
  }
}

locals {
  blueprint_id = "f320947e-1555-4cee-9128-58a6cc4dd99c"
}

provider "cloudendure" {
  api_key         = "B212-1445-FBE4-525A-658D-0885-86FD-4510-8192-EDA1-CA50-7738-AAAB-6D5B-A502-1F07"
  cloudendure_url = "https://console.cloudendure.com/api/latest"
}

resource "cloudendure_recovery_plan" "" {
  name  = "test_terraform_recovery_plan"
  steps = ""
}

resource "cloudendure_blueprint" "blueprint_resource" {
  machine_name       = "terraform_test_machine"
  subnet_ids         = ["subnet-00741c4d"]
  security_group_ids = ["sg-0244a14e569eaba68", "sg-3247085f", "sg-c54906a8", "sg-d64807bb"]

  byol_on_dedicated_instance = true
  run_after_launch           = true
  force_uefi                 = true
  use_shared_ram             = true
  instance_type              = "COPY_ORIGIN"
  static_ip_action           = "DONT_CREATE" // EXISTING
  public_ip_action           = "ALLOCATE"
  cpus                       = 0
  cores_per_cpu              = 0
  mb_ram                     = 0
  security_group_action      = "FROM_POLICY"
  tenancy                    = "SHARED"
  private_ip_action          = "CREATE_NEW"
  #disk = [{
  #"type" : "COPY_ORIGIN",
  #"iops" : 0,
  #"throughput" : 0,
  #"name" : "disk1"
  #}]
}

data "cloudendure_blueprint" "blueprint_data" {
  blueprint_id = local.blueprint_id
}

output "blueprint_machineid" {
  value = data.cloudendure_blueprint.blueprint_data.machine_id
}

output "blueprint_instance_type" {
  value = data.cloudendure_blueprint.blueprint_data.instance_type
}

output "blueprint_security_group_ids" {
  value = data.cloudendure_blueprint.blueprint_data.security_group_ids
}

output "blueprint_subnet_ids" {
  value = data.cloudendure_blueprint.blueprint_data.subnet_ids
}

