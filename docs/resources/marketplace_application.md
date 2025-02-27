# Resource: marketplace_application

## Description

No description available.

## Example Usage

```hcl
resource "marketplace_application" "example" {
  application_seller = {
    description = "example string"
    id = "example string"
    name = "example string"
    state = "example string"
    support_email = "example string"
    support_url = "example string"
  }
  byol_license = "example string"
  cluster_id = "example string"
  configuration = {
    key = "example string"
    value = "example string"
  }
  created_at = "example string"
  description = "example string"
  id = "example string"
  namespace = "example string"
  product = {
    created_at = "example string"
    eol = true
    eol_date = "example string"
    id = "example string"
    license_type = "example string"
    llm_hub = {
      external_api = "example string"
    }
    name = "example string"
    seller = {
      description = "example string"
      id = "example string"
      name = "example string"
      state = "example string"
      support_email = "example string"
      support_url = "example string"
    }
    type = "example string"
    weight = 123
  }
  product_revision = {
    admin_suggestion = "example string"
    byol = {
      activation_url = "example string"
      file_name_in_secret = "example string"
      secret_name = "example string"
      webshop_url = "example string"
    }
    categories = "value"
    contractual_documents = {
      content = "example string"
      file_name = "example string"
      is_deleted = true
    }
    contractual_documents_info = {
      file_name = "example string"
      url = "example string"
    }
    description = "example string"
    description_short = "example string"
    eula = "example string"
    guidance = "example string"
    helm_external = "example string"
    icon = "example string"
    id = "example string"
    license_fee = "example string"
    license_info = "example string"
    number = 123
    post_deployment_info = "example string"
    pre_deployment_info = "example string"
    pricing_info = "example string"
    product_id = "example string"
    product_revision_application_configuration = {
      confidential = true
      default_value = "example string"
      hidden = true
      hint = "example string"
      input_type = "example string"
      key = "example string"
      label = "example string"
      multiple = true
      required = true
      tooltip = "example string"
      validation = {
        message = "example string"
        pattern = "example string"
      }
      values = {
        label = "example string"
        value = "example string"
      }
    }
    proposed_release_date = "example string"
    scheduled_release_date = "example string"
    scheduled_release_until_date = "example string"
    state = "example string"
    used_software = {
      license_name = "example string"
      license_url = "example string"
      name = "example string"
    }
    version = "example string"
  }
  product_revision_id = "example string"
  project_id = "example string"
  release_name = "example string"
  state = "example string"
  username = "example string"
}
```

## Argument Reference

- `application_seller` - The entity responsible for selling the product on the Marketplace
  (Computed)
  - `description` - An optional seller description
    (Computed)
  - `id` - Default kind of id for most objects defined in this project
    (Computed)
  - `name` - The seller name
    (Computed)
  - `state` - State of the Seller. Can be either `active` or `suspended`
    (Computed)
  - `support_email` - The seller's email address
    (Computed)
  - `support_url` - The seller's website
    (Computed)
- `byol_license` - Base64 encoded string containing the license text of the byol license
  (Optional)
- `cluster_id` - The cluster ID within which the namespace for deployment can be found
  (Required)
- `configuration` - List of Configuration for the deployed application
  (Optional)
  - `key` - The key or name of the property to be set
    (Required)
  - `value` - The value of the property to be set
    (Required)
- `created_at` - Time and date of application deployment
  (Optional)
- `description` - user defined application description
  (Optional)
- `id` - Default kind of id for most objects defined in this project
  (Optional)
- `namespace` - The CCE cluster namespace within which the application should be deployed
  (Required)
- `product` - The Product (offering) to be sold on the marketplace
  (Computed)
  - `created_at` - The date and time when the product was created
    (Computed)
  - `eol` - Set product to EOL. The data will be calculated on backend
    (Computed)
  - `eol_date` - End-of-life of the product
    (Computed)
  - `id` - Default kind of id for most objects defined in this project
    (Computed)
  - `license_type` - The type of license, MVP is only unpaid licenses
    (Computed)
  - `llm_hub` - No description available.
    (Computed)
    - `external_api` - Link to external API for the LLM
      (Computed)
  - `name` - The product name which is shown in the teaser and used as a title on the product offering page
    (Computed)
  - `seller` - The entity responsible for selling the product on the Marketplace
    (Computed)
    - `description` - An optional seller description
      (Computed)
    - `id` - Default kind of id for most objects defined in this project
      (Computed)
    - `name` - The seller name
      (Computed)
    - `state` - State of the Seller. Can be either `active` or `suspended`
      (Computed)
    - `support_email` - The seller's email address
      (Computed)
    - `support_url` - The seller's website
      (Computed)
  - `type` - The service deployment type in MVP this is container (CCE), post MVP this will expand to other types
    (Computed)
  - `weight` - The weight of the product, the higher the number the better the recommendation
    (Computed)
- `product_revision` - Combo object made of ProductRevision and the Seller in the ProductRevision
  (Computed)
  - `admin_suggestion` - Admin suggestion is for product revision which got rejected
    (Computed)
  - `byol` - No description available.
    (Computed)
    - `activation_url` - (Unsure) Link that, when visited, registers that the customer has accepted the license
      (Computed)
    - `file_name_in_secret` - filename in secret in which the license data will be stored.
      (Computed)
    - `secret_name` - Name of the secret where byol license will be stored.
      (Computed)
    - `webshop_url` - (Unsure) Link to the webshop where the Product is available
      (Computed)
  - `categories` - Ids correlating to the Categories this Product should be in
    (Computed)
  - `contractual_documents` - Legal documents to be agreed to when using this product. This field is only used during Create (POST)
    (Computed)
    - `content` - base64 encoded file with mimetype
      (Computed)
    - `file_name` - Name of the file
      (Computed)
    - `is_deleted` - Should the file be marked as deleted
      (Computed)
  - `contractual_documents_info` - Legal documents governing the use of this product
    (Computed)
    - `file_name` - Name of the file
      (Computed)
    - `url` - Url to the file
      (Computed)
  - `description` - The Markdown description of the product functionality
    (Computed)
  - `description_short` - The short description of the product which appears in the teasers
    (Computed)
  - `eula` - (Deprecated) The Markdown description of the product EULA
    (Computed)
  - `guidance` - A description of the install process
    (Computed)
  - `helm_external` - The Helm chart URL for the product provided by the seller
    (Computed)
  - `icon` - Base64 encoded image in 16:9 format
    (Computed)
  - `id` - Default kind of id for most objects defined in this project
    (Computed)
  - `license_fee` - The license fee including any details, this may be a either a simple one off license fee in Euro, or a complex annual license fee in Euro and a variable additional cost
    (Computed)
  - `license_info` - Extra info about the license the Customer needs to agree to when buying this Product
    (Computed)
  - `number` - The incremental number of the revision
    (Computed)
  - `post_deployment_info` - The Markdown text for the post deployment screen explaining how to access the product or the next steps
    (Computed)
  - `pre_deployment_info` - The Markdown text for the pre deployment screen explaining what to expect during deployment
    (Computed)
  - `pricing_info` - The pricing information as a guideline for how much the application will cost the user in the OTC
    (Computed)
  - `product_id` - Default kind of id for most objects defined in this project
    (Computed)
  - `product_revision_application_configuration` - Application configuration records
    (Computed)
    - `confidential` - If this config entry should be confidential
      (Computed)
    - `default_value` - Default value of the attribute the customer will need to set
      (Computed)
    - `hidden` - Toggles if the configuration should be hidden or not
      (Computed)
    - `hint` - Extra info describing what the configuration will be used for
      (Computed)
    - `input_type` - The type of input which is used for this element, text, switch or selection (selection may be eiter radio, select or checkbox)
      (Computed)
    - `key` - Name of the attribute the customer will need to set
      (Computed)
    - `label` - Description of the key/value the customer will need to set
      (Computed)
    - `multiple` - When the configuration is a selection type this property defines if multiple options can be selected or if the user can only choose one.
      (Computed)
    - `required` - Toggles if the configuration needs to be set or not
      (Computed)
    - `tooltip` - Extra info to be shown in a tooltip while the customer enters the configuration
      (Computed)
    - `validation` - Describes rules to be used to check the configurations and ensure accuracy (or conformity at the very least)
      (Computed)
      - `message` - No description available.
        (Computed)
      - `pattern` - No description available.
        (Computed)
    - `values` - An array of value objects
      (Computed)
      - `label` - The label of the value being selected
        (Computed)
      - `value` - The value of the selected option which will be used during creation of the application if selected
        (Computed)
  - `proposed_release_date` - When the Seller would like to release this Revision of the Product. Once agreed to, a `scheduled_release_date` and/or `scheduled_release_until_date` will be set
    (Computed)
  - `scheduled_release_date` - When the product is scheduled to be released (usually set after being proposed with the proposed release date)
    (Computed)
  - `scheduled_release_until_date` - Time before the product is scheduled to be released (not after this date, but after scheduled_release_date)
    (Computed)
  - `state` - Enum showing the state this revision is in. Revisions, when persisted, start as `draft`, but can be sent for review by setting this to `ready_for_review` after which this will be set to either `approved` or `rejected`
    (Computed)
  - `used_software` - Entries describing the software used in this Product and the licenses that govern their use
    (Computed)
    - `license_name` - The name of the license used to govern the use of the software
      (Computed)
    - `license_url` - Link to the license text
      (Computed)
    - `name` - The name of the software used
      (Computed)
  - `version` - The version of the release
    (Computed)
- `product_revision_id` - The revision id of the product to enable deployment for review and draft testing
  (Required)
- `project_id` - The project ID within which the CCE cluster for deployment can be found
  (Required)
- `release_name` - The name of the Helm release
  (Optional)
- `state` - Enum showing the Application's deployment state. Starts `pending` on resource creation and is eventually set to `ready` or `error`
  (Optional)
- `username` - (Unsure) Username of the Customer deploying the Application
  (Optional)
