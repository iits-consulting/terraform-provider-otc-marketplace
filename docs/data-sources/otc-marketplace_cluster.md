# Data Source: otc-marketplace_cluster

## Description

No description available.

## Example Usage

```hcl
data "otc-marketplace_cluster" "example" {
  clusters = {
    id = "example string"
    name = "example string"
  }
  project_id = "example string"
}
```

## Argument Reference

- `clusters` - No description available.
  (Computed)
  - `id` - Unique id for this cluster
    (Computed)
  - `name` - Name of the cluster
    (Computed)
- `project_id` - ID of the Open Telekom Cloud project
  (Required)
