# Data Source: marketplace_category

## Description

No description available.

## Example Usage

```hcl
data "marketplace_category" "example" {
  categories = {
    description = "example string"
    id = "example string"
    name = "example string"
    position = 123
    state = "example string"
  }
}
```

## Argument Reference

- `categories` - No description available.
  (Computed)
  - `description` - The long description of this category
    (Computed)
  - `id` - Default kind of id for most objects defined in this project
    (Computed)
  - `name` - Name of the category
    (Computed)
  - `position` - (Unsure) The weighting used to order categories when listed on the frontend. This isn't used for ordering by the provider
    (Computed)
  - `state` - Enum determining if the category can be used (`active`) or not (`suspended`)
    (Computed)
