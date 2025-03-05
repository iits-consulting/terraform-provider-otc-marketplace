# Data Source: marketplace_whoami

## Description

No description available.

## Example Usage

```hcl
data "otc-marketplace_whoami" "example" {
  domain_name = "example string"
  last_project_id = "example string"
  llm_hub = true
  username = "example string"
}
```

## Argument Reference

- `domain_name` - The OTC domain name as visible in the console
  (Computed)
- `last_project_id` - (Unsure) Last project to have been used with this account
  (Computed)
- `llm_hub` - (Unsure) If access to the LLM Hub has been granted or not
  (Computed)
- `username` - The OTC user name as visible in the console
  (Computed)
