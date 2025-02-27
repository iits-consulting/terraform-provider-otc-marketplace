# Looking at the OpenAPI yaml spec, generate Terraform provider code
tfplugingen-openapi generate --config ./generator_config.yml --output ./provider-code-spec.json ./openapi.yml  && \
tfplugingen-framework generate all --input ./provider-code-spec.json --output ./internal

# Update the docs using the schema reported by the provider (should be replaced, see README.md)
terraform providers schema -json > schema.json && python3 gen_docs.py