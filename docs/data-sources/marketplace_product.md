# Data Source: marketplace_product

## Description

No description available.

## Example Usage

```hcl
data "otc-marketplace_product" "example" {
  products = {
    active_revision_id = "example string"
    created_at = "example string"
    eol = true
    eol_date = "example string"
    id = "example string"
    license_type = "example string"
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
}
```

## Argument Reference

- `products` - No description available.
  (Computed)
  - `active_revision_id` - Default kind of id for most objects defined in this project
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
  - `state` - State of the Product's publishing status. Either `published` or `de-published`
    (Computed)
  - `type` - The service deployment type in MVP this is container (CCE), post MVP this will expand to other types
    (Computed)
  - `weight` - The weight of the product, the higher the number the better the recommendation
    (Computed)
