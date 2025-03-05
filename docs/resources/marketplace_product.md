# Resource: marketplace_product

## Description

No description available.

## Example Usage

```hcl
resource "otc-marketplace_product" "example" {
  active_revision_id = "example string"
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
  state = "example string"
  type = "example string"
  weight = 123
}
```

## Argument Reference

- `active_revision_id` - Default kind of id for most objects defined in this project
  (Computed)
- `created_at` - The date and time when the product was created
  (Optional)
- `eol` - Set product to EOL. The data will be calculated on backend
  (Optional)
- `eol_date` - End-of-life of the product
  (Optional)
- `id` - Default kind of id for most objects defined in this project
  (Optional)
- `license_type` - The type of license, MVP is only unpaid licenses
  (Required)
- `llm_hub` - No description available.
  (Optional)
  - `external_api` - Link to external API for the LLM
    (Required)
- `name` - The product name which is shown in the teaser and used as a title on the product offering page
  (Required)
- `seller` - The entity responsible for selling the product on the Marketplace
  (Optional)
  - `description` - An optional seller description
    (Optional)
  - `id` - Default kind of id for most objects defined in this project
    (Required)
  - `name` - The seller name
    (Required)
  - `state` - State of the Seller. Can be either `active` or `suspended`
    (Optional)
  - `support_email` - The seller's email address
    (Optional)
  - `support_url` - The seller's website
    (Optional)
- `state` - State of the Product's publishing status. Either `published` or `de-published`
  (Computed)
- `type` - The service deployment type in MVP this is container (CCE), post MVP this will expand to other types
  (Required)
- `weight` - The weight of the product, the higher the number the better the recommendation
  (Optional)
