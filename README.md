# OTC Marketplace Provider

🚨 Note: This project is still under development and is not recommended for use in production environments. 🚨

## Example Usage

1. Create a variables.tf
```hcl
variable "otc_domain_name" {
  description = "The OTC domain_name (example OTC-EU-DE-00000000..) for authentication"
  type        = string
}

variable "username" {
  description = "The username of the otc-marketplace account"
  type        = string
}

variable "password" {
  description = "The password of the otc-marketplace account"
  type        = string
}
```

2. Create providers.tf
```hcl
provider "otc-marketplace" {
domain_name = var.otc_domain_name
username    = var.username
password    = var.password
}
```

3. Create datasources.tf and replace _eu-de_my_project_ and _my-cce-clustername_ with correct values
```hcl
locals {
  projects = tomap({
    for project in data.marketplace_project.all_projects.projects : project.name => project
  })

  clusters = tomap({
    for cluster in data.marketplace_cluster.all_clusters.clusters : cluster.name => cluster
  })

  categories = tomap({
    for category in data.marketplace_category.all_categories.categories : category.name => category
  })
  #
  selected_project  = local.projects["eu-de_my_project"]
  selected_cluster  = local.clusters["my-cce-clustername"]
}

data "otc-marketplace_cluster" "all_clusters" {
  project_id = local.selected_project.id
}

data "otc-marketplace_namespace" "all_namespaces" {
  project_id = local.selected_project.id
  cluster_id = local.selected_cluster.id
}

data "otc-marketplace_category" "all_categories" {}

data "otc-marketplace_project" "all_projects" {}

data "otc-marketplace_profile" "me" {}
```

3. create a main.tf file with such a content. In this example we create a otc-prometheus-exporter service
```hcl
resource "otc-marketplace_product" "iits_otc_prometheus_exporter" {
  name = "OTC prometheus-exporter" // TODO - will always need to be new (create and delete means this name can never be used again)
  license_type = "opensource"
  type = "container"
  eol = false
  weight = 0
  lifecycle {
    ignore_changes = [
      "seller"
    ]
  }
}

resource "otc-marketplace_product_revision" "iits_otc_prometheus_exporter_revision" {
  depends_on = [marketplace_product.iits_otc_prometheus_exporter]
  categories            = [local.categories["Monitoring"].id]
  product_id            = marketplace_product.iits_otc_prometheus_exporter.id
  description           = local.description
  description_short     = "This software gathers metrics from the Open Telekom Cloud (OTC) for Prometheus."
  icon                  = local.iits_icon
  pre_deployment_info   = local.description
  post_deployment_info  = "It pushes data to prometheus-stack"
  helm_external         = "${local.helm_chart_link}:${local.helm_chart_version}" # after : the version is provided
  license_fee           = "0.00"
  license_info          = "GNU General Public License v3.0"
  pricing_info          = "Free"
  number = 1
  proposed_release_date                      = formatdate("YYYY-MM-DD'T'HH:mm:ss'Z'", timestamp())
  state                                      = "approved"
  used_software                              = [
    {
      name         = "Prometheus Exporter"
      license_name = "Apache License 2.0"
      license_url  = "https://raw.githubusercontent.com/prometheus/node_exporter/refs/heads/master/LICENSE"
    },
  ]
  version                                    = "1.0.0"
  product_revision_application_configuration = local.product_revision_application_configuration
  lifecycle {
    ignore_changes = [
      "id",
    ]
  }
}

locals {
  iits_icon = "https://iits-consulting.de/wp-content/uploads/2024/03/iits-favicon.png" #FIXME does not work yet
  helm_chart_link = "oci://registry-1.docker.io/iits/otc-prometheus-exporter" # Needs to be oci
  helm_chart_version = "1.2.1"
  product_revision_application_configuration = [
    {
      confidential = true #boolean value, if this is an confidential field
      hint = "User in the OTC with access to the API"
      input_type = "text" #string value, needs to be selection, text or switch
      key = "deployment.env.OS_USERNAME" #string values, defines the path inside the values.yaml
      label = "OS_USERNAME"  #string value, will be displayed to the client
      validation = []
      tooltip = "User in the OTC with access to the API"
      values = [
        {
          label = "Username"
          value = "" #string value
        }
      ]
      hidden = false #boolean value, needs everytime to be false
      required = true #boolean value, if the field is required
      default_value = "" #string value, needs to be one of the value of values list section
      multiple = false #boolean value, Allowing multiple selections, can only be true when input_type is "selection"
    },
    {
      confidential = false #boolean value, if this is an confidential field
      hint = "Specifies whether a service account should be created"
      input_type = "switch" #string value, needs to be selection, text or switch
      validation = []
      key = "serviceAccount.create" #string values, defines the path inside the values.yaml
      label = "Service Account Create"  #string value, will be displayed to the client
      tooltip = "Specifies whether a service account should be created"
      values = [] #needs to be empty when input_type is switch
      hidden = false #boolean value, needs everytime to be false
      required = false #boolean value, if the field is required
      default_value = "false" #string value, needs to be one of the value of values list section
      multiple = false #boolean value, Allowing multiple selections, can only be true when input_type is "selection"
    },
    {
      confidential = false #boolean value, if this is an confidential field
      hint = "Choose available features to enable."
      input_type = "selection" #string value, needs to be selection, text or switch
      validation = []
      key = "features" #string values, defines the path inside the values.yaml
      label = "Select Features" #string value, will be displayed to the client
      tooltip = "Choose available features to enable."
      values = [
        {
          label = "Feature A"
          value = "feature_a" #string value, will be displayed to the client
        },
        {
          label = "Feature B"
          value = "feature_b" #string value, will be displayed to the client
        }
      ]
      hidden = false #boolean value, needs everytime to be false
      required = true #boolean value, if the field is required
      default_value = "feature_a" #string value, needs to be one of the value of values list section
      multiple = false #boolean value, Allowing multiple selections, can only be true when input_type is "selection"
    }
  ]
  pre_deployment_info   = <<EOT
  MY DOCS 

  EOT
}
```

4. Terraform apply and it should be created

## Known limitation / Issues
Take a look at TODO.md

## Development Setup

### Terraform local provider configuration

execute first this:

Create a new file called .terraformrc in your home directory (~), then add the dev_overrides block below

Replace $USER with your linux/MAC username

```hcl
provider_installation {

  dev_overrides {
      "registry.terraform.io/iits-consulting/otc-marketplace" = "/home/$USER/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```


### Apply the terraform test code

First set the credentials inside terraform-provider-otc-marketplace/test/.envrc

Then execute the following

```bash
cd terraform-provider-otc-marketplace
go mod tidy
go install .

cd test
# Terraform init will fail, just apply directly
terraform apply
```

## Docs

A reference of the resources and how to use them can either be found in this repo's 
[GitHub Pages](https://iits-consulting.github.io/terraform-provider-otc-marketplace/) or in the [docs folder](https://github.com/iits-consulting/terraform-provider-otc-marketplace/tree/main/docs)

## Contributing

After updating the code, make sure to remove the old docs and run 
`terraform providers schema -json > schema.json && python3 gen_docs.py`  
to regen the docs (and check to see if `tfplugindocs` now supports the new 
protocol, so we don't need to use this!).
