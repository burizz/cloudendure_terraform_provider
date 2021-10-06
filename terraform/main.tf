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
  api_key         = ""
  cloudendure_url = "https://console.cloudendure.com/api/latest"
}

resource "cloudendure_blueprint" "test_blueprint" {
  #blueprint_id = "f320947e-1555-4cee-9128-58a6cc4dd99c"
  blueprint_id = local.blueprint_id
}

data "cloudendure_blueprint" "blueprint" {
  #blueprint_id = "f320947e-1555-4cee-9128-58a6cc4dd99c"
  blueprint_id = local.blueprint_id
}

output "blueprint_machineid" {
  value = data.cloudendure_blueprint.blueprint.machine_id
}

output "blueprint_instance_type" {
  value = data.cloudendure_blueprint.blueprint.instance_type
}

output "blueprint_security_group_ids" {
  value = data.cloudendure_blueprint.blueprint.security_group_ids
}

output "blueprint_subnet_ids" {
  value = data.cloudendure_blueprint.blueprint.subnet_ids
}

