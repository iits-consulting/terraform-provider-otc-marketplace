# Data Source: marketplace_sales_history

## Description

No description available.

## Example Usage

```hcl
data "marketplace_sales_history" "example" {
  sales_history = {
    customer_company_name = "example string"
    customer_company_url = "example string"
    customer_contact_email = "example string"
    customer_contact_number = "example string"
    deployed_at = "example string"
    product_id = "example string"
    product_name = "example string"
    product_revision_id = "example string"
  }
}
```

## Argument Reference

- `sales_history` - No description available.
  (Computed)
  - `customer_company_name` - Name of the customer's company that bought the product
    (Computed)
  - `customer_company_url` - Link to the customer's webpage
    (Computed)
  - `customer_contact_email` - Email at which the customer can be contacted
    (Computed)
  - `customer_contact_number` - Number at which the customer can be reached
    (Computed)
  - `deployed_at` - Time at which the sold product has been deployed by the customer
    (Computed)
  - `product_id` - Default kind of id for most objects defined in this project
    (Computed)
  - `product_name` - Name of the product sold
    (Computed)
  - `product_revision_id` - Default kind of id for most objects defined in this project
    (Computed)
