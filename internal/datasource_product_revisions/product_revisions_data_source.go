package datasource_product_revisions

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"net/http"
	"terraform-provider-otc-marketplace/internal/resource_product_revision"
	"terraform-provider-otc-marketplace/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const productRevisionServicePath = "/product-revisions"

var _ datasource.DataSource = (*productRevisionDataSource)(nil)

func NewProductRevisionDataSource() datasource.DataSource {
	return &productRevisionDataSource{}
}

type productRevisionDataSource struct {
	client *util.MarketplaceAPIClient
}

type productRevisionNativeModel struct {
	AdminSuggestion           string                                                               `json:"admin_suggestion,omitempty" tfsdk:"admin_suggestion"`
	Eula                      string                                                               `json:"eula,omitempty" tfsdk:"eula"`
	Byol                      byolNativeModel                                                      `json:"byol,omitempty" tfsdk:"byol"`
	Categories                []string                                                             `json:"categories,omitempty" tfsdk:"categories"`
	Configuration             []resource_product_revision.ProductRevisionResourceConfigNativeModel `json:"configuration,omitempty" tfsdk:"configuration"` // TODO naming
	ContractualDocumentsInfo  []contractualDocsInfoNativeModel                                     `json:"contractual_documents_info,omitempty" tfsdk:"contractual_documents_info"`
	Description               string                                                               `json:"description,omitempty" tfsdk:"description"`
	DescriptionShort          string                                                               `json:"description_short,omitempty" tfsdk:"description_short"`
	Guidance                  string                                                               `json:"guidance,omitempty" tfsdk:"guidance"`
	HelmExternal              string                                                               `json:"helm_external,omitempty" tfsdk:"helm_external"`
	Icon                      string                                                               `json:"icon,omitempty" tfsdk:"icon"`
	Id                        string                                                               `json:"id,omitempty" tfsdk:"id"`
	LicenseFee                string                                                               `json:"license_fee,omitempty" tfsdk:"license_fee"`
	LicenseInfo               string                                                               `json:"license_info,omitempty" tfsdk:"license_info"`
	Number                    int64                                                                `json:"number,omitempty" tfsdk:"number"`
	PostDeploymentInfo        string                                                               `json:"post_deployment_info,omitempty" tfsdk:"post_deployment_info"`
	PreDeploymentInfo         string                                                               `json:"pre_deployment_info,omitempty" tfsdk:"pre_deployment_info"`
	PricingInfo               string                                                               `json:"pricing_info,omitempty" tfsdk:"pricing_info"`
	ProductId                 string                                                               `json:"product_id,omitempty" tfsdk:"product_id"`
	ProposedReleaseDate       string                                                               `json:"proposed_release_date,omitempty" tfsdk:"proposed_release_date"`
	ScheduledReleaseDate      string                                                               `json:"scheduled_release_date,omitempty" tfsdk:"scheduled_release_date"`
	ScheduledReleaseUntilDate string                                                               `json:"scheduled_release_until_date,omitempty" tfsdk:"scheduled_release_until_date"`
	State                     string                                                               `json:"state,omitempty" tfsdk:"state"` // TODO - enum?
	UsedSoftware              []productRevisionUsedNativeSoftwareModel                             `json:"used_software,omitempty" tfsdk:"used_software"`
	Version                   string                                                               `json:"version,omitempty" tfsdk:"version"`
}

type byolNativeModel struct {
	ActivationUrl    string `json:"activation_url,omitempty" tfsdk:"activation_url"`
	FilenameInSecret string `json:"file_name_in_secret,omitempty" tfsdk:"file_name_in_secret"`
	SecretName       string `json:"secret_name,omitempty" tfsdk:"secret_name"`
	WebshopUrl       string `json:"webshop_url,omitempty" tfsdk:"webshop_url"`
}

type productRevisionUsedNativeSoftwareModel struct {
	LicenseName string `json:"license_name,omitempty" tfsdk:"license_name"`
	LicenseUrl  string `json:"license_url,omitempty" tfsdk:"license_url"`
	Name        string `json:"name,omitempty" tfsdk:"name"`
}

type contractualDocsInfoNativeModel struct {
	Filename string `json:"file_name,omitempty" tfsdk:"file_name"`
	Url      string `json:"url,omitempty" tfsdk:"url"`
}

func (d *productRevisionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_revision"
}
func (d *productRevisionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProductRevisionsDataSourceSchema(ctx)
}

func (d *productRevisionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *productRevisionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProductRevisionsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]productRevisionNativeModel](ctx, http.MethodGet, productRevisionServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, productRevisionServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	var tempNewConfigs []attr.Value
	var tempNewCategories []attr.Value
	var tempNewUsedSoftware []attr.Value
	var tempContDocsInfo []attr.Value

	for _, nativePRs := range *newDataNativeArrPTR {
		for _, config := range nativePRs.Configuration {
			var tempValidation []attr.Value
			var tempValues []attr.Value
			var diags diag.Diagnostics

			for _, val := range config.Validation {
				valObj, valDiags := NewValidationValue(ValidationValue{}.AttributeTypes(ctx), map[string]attr.Value{
					"message": types.StringValue(val.Message),
					"pattern": types.StringValue(val.Pattern),
				})
				diags.Append(valDiags...)
				if diags.HasError() {
					return
				}
				tempValidation = append(tempValidation, valObj)
			}

			validationList, diags := types.ListValue(ValidationValue{}.Type(ctx), tempValidation)
			if diags.HasError() {
				return
			}

			for _, val := range config.Values {
				valObj, valDiags := NewValuesValue(ValuesValue{}.AttributeTypes(ctx), map[string]attr.Value{
					"label": types.StringValue(val.Label),
					"value": types.StringValue(val.Value),
				})
				diags.Append(valDiags...)
				if diags.HasError() {
					return
				}
				tempValues = append(tempValues, valObj)
			}

			valuesList, diags := types.ListValue(ValuesValue{}.Type(ctx), tempValues)
			if diags.HasError() {
				return
			}

			confObj, diags := NewProductRevisionApplicationConfigurationValue(ProductRevisionApplicationConfigurationValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"confidential":  types.BoolValue(config.Confidential),
				"default_value": types.StringValue(config.DefaultValue),
				"hidden":        types.BoolValue(config.Hidden),
				"hint":          types.StringValue(config.Hint),
				"input_type":    types.StringValue(config.InputType),
				"key":           types.StringValue(config.Key),
				"label":         types.StringValue(config.Label),
				"multiple":      types.BoolValue(config.Multiple),
				"required":      types.BoolValue(config.Required),
				"tooltip":       types.StringValue(config.Tooltip),
				"validation":    validationList,
				"values":        valuesList,
			})
			if diags.HasError() {
				return
			}

			tempNewConfigs = append(tempNewConfigs, confObj)
		}

		configsList, diags := types.ListValue(ProductRevisionApplicationConfigurationValue{}.Type(ctx), tempNewConfigs)
		if diags.HasError() {
			return
		}

		for _, category := range nativePRs.Categories {
			tempNewCategories = append(tempNewCategories, types.StringValue(category))
		}

		categoryList, diags := types.ListValue(types.StringType, tempNewCategories)

		for _, usedSoftwareSingle := range nativePRs.UsedSoftware {
			softObj, softDiags := NewUsedSoftwareValue(UsedSoftwareValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"license_name": types.StringValue(usedSoftwareSingle.LicenseName),
				"license_url":  types.StringValue(usedSoftwareSingle.LicenseUrl),
				"name":         types.StringValue(usedSoftwareSingle.Name),
			})
			diags.Append(softDiags...)
			if diags.HasError() {
				return
			}
			tempNewUsedSoftware = append(tempNewUsedSoftware, softObj)
		}

		softList, diags := types.ListValue(UsedSoftwareValue{}.Type(ctx), tempNewUsedSoftware)

		byol, diags := NewByolValue(ByolValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"activation_url":      types.StringValue(nativePRs.Byol.ActivationUrl),
			"file_name_in_secret": types.StringValue(nativePRs.Byol.FilenameInSecret),
			"secret_name":         types.StringValue(nativePRs.Byol.SecretName),
			"webshop_url":         types.StringValue(nativePRs.Byol.WebshopUrl),
		})

		for _, docInfo := range nativePRs.ContractualDocumentsInfo {
			docObj, docDiags := NewContractualDocumentsInfoValue(ContractualDocumentsInfoValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"file_name": types.StringValue(docInfo.Filename),
				"url":       types.StringValue(docInfo.Url),
			})
			diags.Append(docDiags...)
			if diags.HasError() {
				return
			}
			tempContDocsInfo = append(tempContDocsInfo, docObj)
		}

		docInfoList, diags := types.ListValue(ContractualDocumentsInfoValue{}.Type(ctx), tempContDocsInfo)
		if diags.HasError() {
			return
		}

		productObj, diags := NewProductRevisionsValue(ProductRevisionsValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"admin_suggestion":           types.StringValue(nativePRs.AdminSuggestion),
			"byol":                       byol,
			"categories":                 categoryList,
			"contractual_documents":      types.ListNull(ContractualDocumentsValue{}.Type(ctx)), // Backend converts these to contractual_documents_info
			"contractual_documents_info": docInfoList,
			"description":                types.StringValue(nativePRs.Description),
			"description_short":          types.StringValue(nativePRs.Description),
			"eula":                       types.StringValue(nativePRs.Eula),
			"guidance":                   types.StringValue(nativePRs.Guidance),
			"helm_external":              types.StringValue(nativePRs.HelmExternal),
			"icon":                       types.StringValue(nativePRs.Icon),
			"id":                         types.StringValue(nativePRs.Id),
			"license_fee":                types.StringValue(nativePRs.LicenseFee),
			"license_info":               types.StringValue(nativePRs.LicenseInfo),
			"number":                     types.Int64Value(nativePRs.Number),
			"post_deployment_info":       types.StringValue(nativePRs.PostDeploymentInfo),
			"pre_deployment_info":        types.StringValue(nativePRs.PreDeploymentInfo),
			"pricing_info":               types.StringValue(nativePRs.PricingInfo),
			"product_id":                 types.StringValue(nativePRs.ProductId),
			"product_revision_application_configuration": configsList,
			"proposed_release_date":                      types.StringValue(nativePRs.ProposedReleaseDate),
			"scheduled_release_date":                     types.StringValue(nativePRs.ScheduledReleaseDate),
			"scheduled_release_until_date":               types.StringValue(nativePRs.ScheduledReleaseUntilDate),
			"state":                                      types.StringValue(nativePRs.State),
			"used_software":                              softList,
			"version":                                    types.StringValue(nativePRs.Version),
		})
		if diags.HasError() {
			return
		}

		newData = append(newData, productObj)
	}

	productRevisionsSet, diags := types.SetValue(ProductRevisionsValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = ProductRevisionsModel{ProductRevisions: productRevisionsSet}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
