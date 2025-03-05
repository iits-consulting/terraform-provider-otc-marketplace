# Data Source: otc-marketplace_profile

## Description

No description available.

## Example Usage

```hcl
data "otc-marketplace_profile" "example" {
  customer_support_number = "example string"
  description = "example string"
  email = "example string"
  id = "example string"
  name = "example string"
  status = "example string"
  support_email = "example string"
  support_url = "example string"
  temp_customer_support_number = "example string"
  temp_description = "example string"
  temp_email = "example string"
  temp_name = "example string"
  temp_support_email = "example string"
  temp_support_url = "example string"
}
```

## Argument Reference

- `customer_support_number` - Number at which the Seller can be reached for customer support
  (Computed)
- `description` - Description of the Seller
  (Computed)
- `email` - Email where the Seller can be reached
  (Computed)
- `id` - Default kind of id for most objects defined in this project
  (Computed)
- `name` - The name of the Seller
  (Computed)
- `status` - Enum showing the Seller's status for selling products on the Marketplace
  (Computed)
- `support_email` - Email where the Customer can reach the Seller for support with the Product
  (Computed)
- `support_url` - Link to webpage where Seller can give customer support
  (Computed)
- `temp_customer_support_number` - Temporary number at which the Seller can be reached for customer support
  (Computed)
- `temp_description` - Temporary description of the Seller
  (Computed)
- `temp_email` - Temporary email where the Seller can be reached
  (Computed)
- `temp_name` - Temporary name (Unsure) Might be set when the Seller is being verified?
  (Computed)
- `temp_support_email` - Temporary email where the Customer can reach the Seller for support with the Product
  (Computed)
- `temp_support_url` - Temporary link to webpage where Seller can give customer support
  (Computed)
