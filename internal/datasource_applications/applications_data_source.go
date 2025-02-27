package datasource_applications

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/resource_product_revision"
	"terraform-provider-otc-marketplace/internal/util"
)

const applicationServicePath = "/applications"

var _ datasource.DataSource = (*applicationDataSource)(nil)

func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

type applicationDataSource struct {
	client *util.MarketplaceAPIClient
}

type applicationDataSourceNativeModel struct {
	Id              string                                                       `json:"id,omitempty"`
	ByolLicense     string                                                       `json:"byol_license,omitempty"`
	Description     string                                                       `json:"description,omitempty"`
	ClusterId       string                                                       `json:"cluster_id,omitempty"`
	CreatedAt       string                                                       `json:"created_at,omitempty"`
	Error           string                                                       `json:"error,omitempty"`
	Namespace       string                                                       `json:"namespace,omitempty"`
	Configuration   []ConfigurationNativeModel                                   `json:"configuration,omitempty"`
	Product         util.ProductDataSourceNativeModel                            `json:"product,omitempty"`
	ProductRevision resource_product_revision.ProductRevisionResourceNativeModel `json:"product_revision,omitempty"`
	ProjectId       string                                                       `json:"project_id,omitempty"`
	ReleaseName     string                                                       `json:"release_name,omitempty"`
	State           string                                                       `json:"state,omitempty"` // TODO - enum?
	Username        string                                                       `json:"username,omitempty"`
	Seller          util.SellerNativeModel                                       `json:"seller,omitempty"`
}

type ConfigurationNativeModel struct {
	Key   string `tfsdk:"key" json:"key"`
	Value string `tfsdk:"value" json:"value"`
}

func (d *applicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (d *applicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ApplicationsDataSourceSchema(ctx)
}
func (d *applicationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		// IMPORTANT: This method is called MULTIPLE times. An initial call might not have configured the Provider yet, so we need
		// to handle this gracefully. It will eventually be called with a configured provider.
		return
	}

	clientPTR, ok := req.ProviderData.(*util.MarketplaceAPIClient)
	if !ok || clientPTR == nil {
		resp.Diagnostics.AddError(
			"Provider Configuration Error",
			"The provider was not configured correctly, or the API client is missing.",
		)
		return
	}
	d.client = clientPTR
}

func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ApplicationsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]applicationDataSourceNativeModel](ctx, http.MethodGet, applicationServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, applicationServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeApplications := range *newDataNativeArrPTR {
		var tempValues []attr.Value
		var tempValidations []attr.Value
		var tempNewConfig []attr.Value
		for _, config := range nativeApplications.Configuration {
			confObj, confDiags := NewConfigurationValue(ConfigurationValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"key":   types.StringValue(config.Key),
				"value": types.StringValue(config.Value),
			})
			resp.Diagnostics.Append(confDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			tempNewConfig = append(tempNewConfig, confObj)
		}

		configList, diags := types.ListValue(ConfigurationValue{}.Type(ctx), tempNewConfig)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		var tempPRCategories []attr.Value
		for _, category := range nativeApplications.ProductRevision.Categories {
			tempPRCategories = append(tempPRCategories, types.StringValue(category))
		}

		categoryList, diags := types.ListValue(types.StringType, tempPRCategories)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		var tempPRConfigs []attr.Value
		for _, prConfig := range nativeApplications.ProductRevision.Configuration {
			for _, val := range prConfig.Validation {
				valObj, valDiags := NewValidationValue(ValidationValue{}.AttributeTypes(ctx), map[string]attr.Value{
					"pattern": types.StringValue(val.Pattern),
					"message": types.StringValue(val.Message),
				})
				diags.Append(valDiags...)
				resp.Diagnostics.Append(diags...)
				if diags.HasError() {
					return
				}
				tempValidations = append(tempValidations, valObj)
			}

			validationList, valDiags := types.ListValue(ValidationValue{}.Type(ctx), tempValidations)
			diags.Append(valDiags...)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			for _, val := range prConfig.Values {
				valObj, valiDiags := NewValuesValue(ValuesValue{}.AttributeTypes(ctx), map[string]attr.Value{
					"label": types.StringValue(val.Label),
					"value": types.StringValue(val.Value),
				})
				diags.Append(valiDiags...)
				resp.Diagnostics.Append(diags...)
				if diags.HasError() {
					return
				}
				tempValues = append(tempValues, valObj)
			}

			valuesList, valDiags := types.ListValue(ValuesValue{}.Type(ctx), tempValues)
			diags.Append(valDiags...)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			confObj, confDiags := NewProductRevisionApplicationConfigurationValue(ProductRevisionApplicationConfigurationValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"confidential":  types.BoolValue(prConfig.Confidential),
				"default_value": types.StringValue(prConfig.DefaultValue),
				"hidden":        types.BoolValue(prConfig.Hidden),
				"hint":          types.StringValue(prConfig.Hint),
				"input_type":    types.StringValue(prConfig.InputType),
				"key":           types.StringValue(prConfig.Key),
				"label":         types.StringValue(prConfig.Label),
				"multiple":      types.BoolValue(prConfig.Multiple),
				"required":      types.BoolValue(prConfig.Required),
				"tooltip":       types.StringValue(prConfig.Tooltip),
				"validation":    validationList,
				"values":        valuesList,
			})
			diags.Append(confDiags...)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}

			tempPRConfigs = append(tempPRConfigs, confObj)
		}

		prConfigList, diags := types.ListValue(ProductRevisionApplicationConfigurationValue{}.Type(ctx), tempPRConfigs)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		var tempPRUsedSoftware []attr.Value
		for _, soft := range nativeApplications.ProductRevision.UsedSoftware {
			softObj, softDiags := NewUsedSoftwareValue(UsedSoftwareValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"license_name": types.StringValue(soft.LicenseName),
				"license_url":  types.StringValue(soft.LicenseUrl),
				"name":         types.StringValue(soft.Name),
			})
			diags.Append(softDiags...)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
			tempPRUsedSoftware = append(tempPRUsedSoftware, softObj)
		}

		softList, diags := types.ListValue(UsedSoftwareValue{}.Type(ctx), tempPRUsedSoftware)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		byol, diags := NewByolValue(ByolValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"activation_url":      types.StringValue(nativeApplications.ProductRevision.Byol.ActivationUrl),
			"file_name_in_secret": types.StringValue(nativeApplications.ProductRevision.Byol.FileNameInSecret),
			"secret_name":         types.StringValue(nativeApplications.ProductRevision.Byol.SecretName),
			"webshop_url":         types.StringValue(nativeApplications.ProductRevision.Byol.WebshopUrl),
		})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		byolObj, diags := byol.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		var tempPRContDocsInfo []attr.Value
		for _, docInfo := range nativeApplications.ProductRevision.ContractualDocumentsInfo {
			docObj, docDiags := NewContractualDocumentsInfoValue(ContractualDocumentsInfoValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"file_name": types.StringValue(docInfo.Filename),
				"url":       types.StringValue(docInfo.Url),
			})
			diags.Append(docDiags...)
			resp.Diagnostics.Append(diags...)
			if diags.HasError() {
				return
			}
			tempPRContDocsInfo = append(tempPRContDocsInfo, docObj)
		}

		contDocsInfoList, diags := types.ListValue(ContractualDocumentsInfoValue{}.Type(ctx), tempPRContDocsInfo)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		productRevision, diags := NewProductRevisionValue(ProductRevisionValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"admin_suggestion":           types.StringNull(), // TODO - does this still exist?
			"byol":                       byolObj,
			"categories":                 categoryList,
			"contractual_documents":      types.ListNull(ContractualDocumentsValue{}.Type(ctx)), // These are converted to contractual_documents_info on the backend
			"contractual_documents_info": contDocsInfoList,
			"description":                types.StringValue(nativeApplications.ProductRevision.Description),
			"description_short":          types.StringValue(nativeApplications.ProductRevision.DescriptionShort),
			"eula":                       types.StringNull(), // TODO - replaced by byol?
			"guidance":                   types.StringValue(nativeApplications.ProductRevision.Guidance),
			"helm_external":              types.StringValue(nativeApplications.ProductRevision.HelmExternal),
			"icon":                       types.StringValue(nativeApplications.ProductRevision.Icon),
			"id":                         types.StringValue(nativeApplications.ProductRevision.Id),
			"license_fee":                types.StringValue(nativeApplications.ProductRevision.LicenseFee),
			"license_info":               types.StringValue(nativeApplications.ProductRevision.LicenseInfo),
			"number":                     types.Int64Value(nativeApplications.ProductRevision.Number),
			"post_deployment_info":       types.StringValue(nativeApplications.ProductRevision.PostDeploymentInfo),
			"pre_deployment_info":        types.StringValue(nativeApplications.ProductRevision.PreDeploymentInfo),
			"pricing_info":               types.StringValue(nativeApplications.ProductRevision.PricingInfo),
			"product_id":                 types.StringValue(nativeApplications.ProductRevision.ProductId),
			"product_revision_application_configuration": prConfigList,
			"proposed_release_date":                      types.StringValue(nativeApplications.ProductRevision.ProposedReleaseDate),
			"scheduled_release_date":                     types.StringValue(nativeApplications.ProductRevision.ScheduledReleaseDate),
			"scheduled_release_until_date":               types.StringValue(nativeApplications.ProductRevision.ScheduledReleaseUntilDate),
			"state":                                      types.StringValue(nativeApplications.ProductRevision.State),
			"used_software":                              softList,
			"version":                                    types.StringValue(nativeApplications.ProductRevision.Version),
		})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		productRevisionObj, diags := productRevision.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		productSeller, diags := NewSellerValue(SellerValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"description":   types.StringValue(nativeApplications.Product.Seller.Description),
			"id":            types.StringValue(nativeApplications.Product.Seller.Id),
			"name":          types.StringValue(nativeApplications.Product.Seller.Name),
			"state":         types.StringValue(nativeApplications.Product.Seller.State),
			"support_email": types.StringValue(nativeApplications.Product.Seller.SupportEmail),
			"support_url":   types.StringValue(nativeApplications.Product.Seller.SupportUrl),
		})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		productSellerObj, diags := productSeller.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		llmHub, diags := NewLlmHubValue(LlmHubValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"external_api": types.StringValue(nativeApplications.Product.LlmHub.ExternalApi),
		})

		llmHubObj, diags := llmHub.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		product, diags := NewProductValue(ProductValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"created_at":   types.StringValue(nativeApplications.Product.CreatedAt),
			"eol":          types.BoolValue(nativeApplications.Product.EOL),
			"eol_date":     types.StringValue(nativeApplications.Product.EOLDate),
			"id":           types.StringValue(nativeApplications.Product.Id),
			"license_type": types.StringValue(nativeApplications.Product.LicenseType),
			"name":         types.StringValue(nativeApplications.Product.Name),
			"type":         types.StringValue(nativeApplications.Product.Type),
			"seller":       productSellerObj,
			"weight":       types.Int64Value(nativeApplications.Product.Weight),
			"llm_hub":      llmHubObj,
		})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		productObj, diags := product.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		applicationSeller, diags := NewApplicationSellerValue(ApplicationSellerValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"description":   types.StringValue(nativeApplications.Seller.Description),
			"id":            types.StringValue(nativeApplications.Seller.Id),
			"name":          types.StringValue(nativeApplications.Seller.Name),
			"state":         types.StringValue(nativeApplications.Seller.State),
			"support_email": types.StringValue(nativeApplications.Seller.SupportEmail),
			"support_url":   types.StringValue(nativeApplications.Seller.SupportUrl),
		})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		applicationSellerObj, diags := applicationSeller.ToObjectValue(ctx)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		applicationObj, diags := NewApplicationsValue(ApplicationsValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"byol_license":       types.StringValue(nativeApplications.ByolLicense),
			"cluster_id":         types.StringValue(nativeApplications.ClusterId),
			"configuration":      configList,
			"created_at":         types.StringValue(nativeApplications.CreatedAt),
			"description":        types.StringValue(nativeApplications.Description),
			"id":                 types.StringValue(nativeApplications.Id),
			"namespace":          types.StringValue(nativeApplications.Namespace),
			"product":            productObj,
			"product_revision":   productRevisionObj,
			"project_id":         types.StringValue(nativeApplications.ProjectId),
			"release_name":       types.StringValue(nativeApplications.ReleaseName),
			"application_seller": applicationSellerObj,
			"state":              types.StringValue(nativeApplications.State),
			"username":           types.StringValue(nativeApplications.Username),
		})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}

		newData = append(newData, applicationObj)
	}

	applicationsSet, diags := types.SetValue(ApplicationsValue{}.Type(ctx), newData)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data = ApplicationsModel{Applications: applicationsSet}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
