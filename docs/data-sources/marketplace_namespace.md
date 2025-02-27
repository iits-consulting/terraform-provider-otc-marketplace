# Data Source: marketplace_namespace

## Description

No description available.

## Example Usage

```hcl
data "marketplace_namespace" "example" {
  cluster_id = "example string"
  namespaces = {
    cluster_id = "example string"
    name = "example string"
    project_id = "example string"
  }
  project_id = "example string"
}
```

## Argument Reference

- `cluster_id` - ID of the CCE Instance
  (Required)
- `namespaces` - No description available.
  (Computed)
  - `cluster_id` - Id of the cluster to fetch the namespace for
    (Computed)
  - `name` - Name of the namespace (name == namespace.toString())
    (Computed)
  - `project_id` - Id of the project to fetch the namespace for
    (Computed)
- `project_id` - ID of the Project
  (Required)
