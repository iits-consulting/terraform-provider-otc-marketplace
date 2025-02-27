package resource_product_revision

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*productRevisionResource)(nil)

const productRevisionResourcePath = "/product-revisions"

func NewProductRevisionResource() resource.Resource {
	return &productRevisionResource{}
}

type productRevisionResource struct {
	client *util.MarketplaceAPIClient
}

type ProductRevisionResourceNativeModel struct {
	AdminSuggestion           string                                           `json:"admin_suggestion,omitempty" tfsdk:"admin_suggestion"`
	Eula                      string                                           `json:"eula,omitempty" tfsdk:"eula"`
	Configuration             []ProductRevisionResourceConfigNativeModel       `json:"configuration,omitempty" tfsdk:"product_revision_application_configuration"`
	ContractualDocumentsInfo  []contractualDocsInfoNativeModel                 `json:"contractual_documents_info,omitempty" tfsdk:"contractual_documents_info"`
	Categories                []string                                         `json:"categories,omitempty" tfsdk:"categories"`
	Description               string                                           `json:"description,omitempty" tfsdk:"description"`
	DescriptionShort          string                                           `json:"description_short,omitempty" tfsdk:"description_short"`
	Guidance                  string                                           `json:"guidance,omitempty" tfsdk:"guidance"`
	HelmExternal              string                                           `json:"helm_external,omitempty" tfsdk:"helm_external"`
	Icon                      string                                           `json:"icon,omitempty" tfsdk:"icon"`
	Id                        string                                           `json:"id,omitempty" tfsdk:"id"`
	LicenseFee                string                                           `json:"license_fee,omitempty" tfsdk:"license_fee"`
	LicenseInfo               string                                           `json:"license_info,omitempty" tfsdk:"license_info"`
	PostDeploymentInfo        string                                           `json:"post_deployment_info,omitempty" tfsdk:"post_deployment_info"`
	PreDeploymentInfo         string                                           `json:"pre_deployment_info,omitempty" tfsdk:"pre_deployment_info"`
	PricingInfo               string                                           `json:"pricing_info,omitempty" tfsdk:"pricing_info"`
	ProductId                 string                                           `json:"product_id,omitempty" tfsdk:"product_id"`
	ProposedReleaseDate       string                                           `json:"proposed_release_date,omitempty" tfsdk:"proposed_release_date"`
	ScheduledReleaseDate      string                                           `json:"scheduled_release_date,omitempty" tfsdk:"scheduled_release_date"`
	ScheduledReleaseUntilDate string                                           `json:"scheduled_release_until_date,omitempty" tfsdk:"scheduled_release_until_date"`
	State                     string                                           `json:"state,omitempty" tfsdk:"state"`
	Number                    int64                                            `json:"number,omitempty" tfsdk:"number"`
	UsedSoftware              []productRevisionResourceUsedSoftwareNativeModel `json:"used_software,omitempty" tfsdk:"used_software"`
	Version                   string                                           `json:"version,omitempty" tfsdk:"version"`
	Byol                      productRevisionByolNativeModel                   `json:"byol,omitempty" tfsdk:"byol"`
}

type productRevisionByolNativeModel struct {
	ActivationUrl    string `json:"activation_url,omitempty" tfsdk:"activation_url"`
	FileNameInSecret string `json:"file_name_in_secret,omitempty" tfsdk:"file_name_in_secret"`
	SecretName       string `json:"secret_name,omitempty" tfsdk:"secret_name"`
	WebshopUrl       string `json:"webshop_url,omitempty" tfsdk:"webshop_url"`
}
type contractualDocsInfoNativeModel struct {
	Filename string `json:"file_name,omitempty" tfsdk:"file_name"`
	Url      string `json:"url,omitempty" tfsdk:"url"`
}

type ContractualDocsNativeModel struct {
	Filename  string `json:"file_name,omitempty" tfsdk:"file_name"`
	Content   string `json:"content,omitempty" tfsdk:"content"`
	IsDeleted bool   `json:"is_deleted,omitempty" tfsdk:"is_deleted"`
}

type ProductRevisionResourceConfigNativeModel struct {
	DefaultValue string                                               `json:"default_value,omitempty" tfsdk:"default_value"`
	Confidential bool                                                 `json:"confidential,omitempty" tfsdk:"confidential"`
	Hidden       bool                                                 `json:"hidden,omitempty" tfsdk:"hidden"`
	Hint         string                                               `json:"hint,omitempty" tfsdk:"hint"`
	InputType    string                                               `json:"input_type,omitempty" tfsdk:"input_type"`
	Key          string                                               `json:"key,omitempty" tfsdk:"key"`
	Label        string                                               `json:"label,omitempty" tfsdk:"label"`
	Multiple     bool                                                 `json:"multiple,omitempty" tfsdk:"multiple"`
	Required     bool                                                 `json:"required,omitempty" tfsdk:"required"`
	Tooltip      string                                               `json:"tooltip,omitempty" tfsdk:"tooltip"`
	Validation   []productRevisionResourceConfigValidationNativeModel `json:"validation,omitempty" tfsdk:"validation"`
	Values       []productRevisionResourceConfigValueNativeModel      `json:"values,omitempty" tfsdk:"values"`
}

type ProductRevisionResourceConfigSwitchNativeModel struct {
	DefaultValue bool                                                 `json:"default_value,omitempty" tfsdk:"default_value"`
	Confidential bool                                                 `json:"confidential,omitempty" tfsdk:"confidential"`
	Hidden       bool                                                 `json:"hidden,omitempty" tfsdk:"hidden"`
	Hint         string                                               `json:"hint,omitempty" tfsdk:"hint"`
	InputType    string                                               `json:"input_type,omitempty" tfsdk:"input_type"`
	Key          string                                               `json:"key,omitempty" tfsdk:"key"`
	Label        string                                               `json:"label,omitempty" tfsdk:"label"`
	Multiple     bool                                                 `json:"multiple,omitempty" tfsdk:"multiple"`
	Required     bool                                                 `json:"required,omitempty" tfsdk:"required"`
	Tooltip      string                                               `json:"tooltip,omitempty" tfsdk:"tooltip"`
	Validation   []productRevisionResourceConfigValidationNativeModel `json:"validation,omitempty" tfsdk:"validation"`
	Values       []productRevisionResourceConfigValueNativeModel      `json:"values,omitempty" tfsdk:"values"`
}

type productRevisionResourceConfigValidationNativeModel struct {
	Message string `json:"message,omitempty" tfsdk:"message"`
	Pattern string `json:"pattern,omitempty" tfsdk:"pattern"`
}

type productRevisionResourceConfigValueNativeModel struct {
	Value string `json:"value,omitempty" tfsdk:"value"`
	Label string `json:"label,omitempty" tfsdk:"label"`
}

type productRevisionResourceUsedSoftwareNativeModel struct {
	LicenseName string `json:"license_name,omitempty" tfsdk:"license_name"`
	LicenseUrl  string `json:"license_url,omitempty" tfsdk:"license_url"`
	Name        string `json:"name,omitempty" tfsdk:"name"`
}

func (r *productRevisionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_revision"
}

type productRevisionResourceNativeModModel struct {
	Byol                 productRevisionByolNativeModel                   `json:"byol,omitempty" tfsdk:"byol"`
	Categories           []string                                         `json:"categories,omitempty" tfsdk:"categories"`
	ProductId            string                                           `json:"product_id,omitempty" tfsdk:"product_id"`
	Description          string                                           `json:"description,omitempty" tfsdk:"description"`
	DescriptionShort     string                                           `json:"description_short,omitempty" tfsdk:"description_short"`
	Icon                 string                                           `json:"icon,omitempty" tfsdk:"icon"`
	PreDeploymentInfo    string                                           `json:"pre_deployment_info,omitempty" tfsdk:"pre_deployment_info"`
	PostDeploymentInfo   string                                           `json:"post_deployment_info,omitempty" tfsdk:"post_deployment_info"`
	ProposedReleaseDate  string                                           `json:"proposed_release_date,omitempty" tfsdk:"proposed_release_date"`
	LicenseInfo          string                                           `json:"license_info,omitempty" tfsdk:"license_info"`
	Guidance             string                                           `json:"guidance,omitempty" tfsdk:"guidance"`
	Version              string                                           `json:"version,omitempty" tfsdk:"version"`
	HelmExternal         string                                           `json:"helm_external,omitempty" tfsdk:"helm_external"`
	LicenseFee           string                                           `json:"license_fee,omitempty" tfsdk:"license_fee"`
	PricingInfo          string                                           `json:"pricing_info,omitempty" tfsdk:"pricing_info"`
	UsedSoftware         []productRevisionResourceUsedSoftwareNativeModel `json:"used_software,omitempty" tfsdk:"used_software"`
	Configuration        []ProductRevisionResourceConfigNativeModel       `json:"configuration,omitempty" tfsdk:"product_revision_application_configuration"`
	ContractualDocuments []ContractualDocsNativeModel                     `json:"contractual_documents,omitempty" tfsdk:"contractual_documents"`
}

func (r *productRevisionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProductRevisionResourceSchema(ctx)
}

func (r *productRevisionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = clientPTR
}

// TODO - cleaner to just return diags and check for errors in whatever calls this?
func ProductRevisionMapper(ctx context.Context, newDataNativePTR *ProductRevisionResourceNativeModel) (*ProductRevisionModel, error) {

	var diags diag.Diagnostics

	var tempNewCategories []types.String
	for _, category := range newDataNativePTR.Categories {
		tempNewCategories = append(tempNewCategories, types.StringValue(category))
	}
	categoriesAsList := util.ListValueOrNull[types.String](ctx, types.StringType, tempNewCategories, &diags)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
	}

	var tempNewConfigs []ProductRevisionApplicationConfigurationValue
	confListType := types.ObjectType{
		AttrTypes: ProductRevisionApplicationConfigurationValue{}.AttributeTypes(ctx),
	}
	for _, conf := range newDataNativePTR.Configuration {
		var validationObjs []ValidationValue
		for _, validation := range conf.Validation {
			valObj, valDiags := NewValidationValue(ValidationValue{}.AttributeTypes(ctx), map[string]attr.Value{
				"message": types.StringValue(validation.Message),
				"pattern": types.StringValue(validation.Pattern),
			})
			diags.Append(valDiags...)
			validationObjs = append(validationObjs, valObj)
		}
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}
		validationAsList := util.ListValueOrNull[ValidationValue](ctx, ValidationValue{}.Type(ctx), validationObjs, &diags)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}

		var valueObjs []ValuesValue
		for _, value := range conf.Values {
			valObj, valDiags := NewValuesValue(ValuesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"label": types.StringValue(value.Label),
					"value": types.StringValue(value.Value),
				},
			)
			diags.Append(valDiags...)
			valueObjs = append(valueObjs, valObj)
		}
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}
		valuesAsList := util.ListValueOrNull[ValuesValue](ctx, ValuesValue{}.Type(ctx), valueObjs, &diags)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}

		confObj, confDiags := NewProductRevisionApplicationConfigurationValue(ProductRevisionApplicationConfigurationValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"confidential":  types.BoolValue(conf.Confidential),
			"default_value": types.StringValue(conf.DefaultValue),
			"hidden":        types.BoolValue(conf.Hidden),
			"hint":          types.StringValue(conf.Hint),
			"input_type":    types.StringValue(conf.InputType),
			"key":           types.StringValue(conf.Key),
			"label":         types.StringValue(conf.Label),
			"multiple":      types.BoolValue(conf.Multiple),
			"required":      types.BoolValue(conf.Required),
			"tooltip":       types.StringValue(conf.Tooltip),
			"validation":    validationAsList,
			"values":        valuesAsList,
		})
		diags.Append(confDiags...)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}
		tempNewConfigs = append(tempNewConfigs, confObj)

	}
	configsAsList := util.ListValueOrNull[ProductRevisionApplicationConfigurationValue](ctx, confListType, tempNewConfigs, &diags)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
	}

	var tempNewUsedSoftware []UsedSoftwareValue
	usedSoftListType := types.ObjectType{
		AttrTypes: UsedSoftwareValue{}.AttributeTypes(ctx),
	}
	for _, soft := range newDataNativePTR.UsedSoftware {
		softObj, softDiags := NewUsedSoftwareValue(UsedSoftwareValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"license_name": types.StringValue(soft.LicenseName),
			"license_url":  types.StringValue(soft.LicenseUrl),
			"name":         types.StringValue(soft.Name),
		})
		diags.Append(softDiags...)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}
		tempNewUsedSoftware = append(tempNewUsedSoftware, softObj)
	}
	softAsList := util.ListValueOrNull[UsedSoftwareValue](ctx, usedSoftListType, tempNewUsedSoftware, &diags)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
	}

	var tempContDocsInfo []ContractualDocumentsInfoValue
	contractualDocumentsInfoListType := types.ObjectType{
		AttrTypes: ContractualDocumentsInfoValue{}.AttributeTypes(ctx),
	}
	for _, docs := range newDataNativePTR.ContractualDocumentsInfo {
		docsObj, docsDiags := NewContractualDocumentsInfoValue(ContractualDocumentsInfoValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"file_name": types.StringValue(docs.Filename),
			"url":       types.StringValue(docs.Url),
		})
		diags.Append(docsDiags...)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
		}
		tempContDocsInfo = append(tempContDocsInfo, docsObj)
	}
	docsInfoAsList := util.ListValueOrNull[ContractualDocumentsInfoValue](ctx, contractualDocumentsInfoListType, tempContDocsInfo, &diags)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
	}

	// Backend should remove these and return them as ContractualDocumentInfo
	var emptyDocsArr []ContractualDocumentsValue
	contractualDocumentsListType := types.ObjectType{
		AttrTypes: ContractualDocumentsValue{}.AttributeTypes(ctx),
	}
	emptyContractualDocsList := util.ListValueOrNull[ContractualDocumentsValue](ctx, contractualDocumentsListType, emptyDocsArr, &diags)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
	}

	newData := ProductRevisionModel{
		AdminSuggestion:                         types.StringValue(newDataNativePTR.AdminSuggestion),
		Eula:                                    types.StringValue(newDataNativePTR.Eula),
		Categories:                              categoriesAsList,
		ProductRevisionApplicationConfiguration: configsAsList,
		ContractualDocuments:                    emptyContractualDocsList,
		ContractualDocumentsInfo:                docsInfoAsList,
		Description:                             types.StringValue(newDataNativePTR.Description),
		DescriptionShort:                        types.StringValue(newDataNativePTR.DescriptionShort),
		Guidance:                                types.StringValue(newDataNativePTR.Guidance),
		HelmExternal:                            types.StringValue(newDataNativePTR.HelmExternal),
		Icon:                                    types.StringValue(newDataNativePTR.Icon),
		Id:                                      types.StringValue(newDataNativePTR.Id),
		LicenseFee:                              types.StringValue(newDataNativePTR.LicenseFee),
		LicenseInfo:                             types.StringValue(newDataNativePTR.LicenseInfo),
		Number:                                  types.Int64Value(newDataNativePTR.Number),
		PostDeploymentInfo:                      types.StringValue(newDataNativePTR.PostDeploymentInfo),
		PreDeploymentInfo:                       types.StringValue(newDataNativePTR.PreDeploymentInfo),
		PricingInfo:                             types.StringValue(newDataNativePTR.PricingInfo),
		ProductId:                               types.StringValue(newDataNativePTR.ProductId),
		ProposedReleaseDate:                     types.StringValue(newDataNativePTR.ProposedReleaseDate),
		ScheduledReleaseDate:                    types.StringValue(newDataNativePTR.ScheduledReleaseDate),
		ScheduledReleaseUntilDate:               types.StringValue(newDataNativePTR.ScheduledReleaseUntilDate),
		State:                                   types.StringValue(newDataNativePTR.State),
		UsedSoftware:                            softAsList,
		Version:                                 types.StringValue(newDataNativePTR.Version),
		Byol: ByolValue{
			ActivationUrl:    types.StringValue(newDataNativePTR.Byol.ActivationUrl),
			FileNameInSecret: types.StringValue(newDataNativePTR.Byol.FileNameInSecret),
			SecretName:       types.StringValue(newDataNativePTR.Byol.SecretName),
			WebshopUrl:       types.StringValue(newDataNativePTR.Byol.WebshopUrl),
		},
	}

	return &newData, nil
}

func (r *productRevisionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProductRevisionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	tflog.Debug(ctx, fmt.Sprintf("read req data: %+v", data))

	if data.Id.IsNull() || data.Id.IsUnknown() || data.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"resource id needs to be set", "resource id is either null or unknown")
		return
	}

	url := fmt.Sprintf("%s/%s", productRevisionResourcePath, util.SanitizeString(data.Id.ValueString()))
	newDataNativePTR, err := util.MakePRMarketplaceRequest[ProductRevisionResourceNativeModel](ctx, http.MethodGet, url, nil, r.client) // TODO - Switch to normal one when fixed
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, url, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	dataPTR, err := ProductRevisionMapper(ctx, newDataNativePTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, dataPTR)...)
}

func (r *productRevisionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProductRevisionModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if data.ProductId.IsNull() || data.ProductId.IsUnknown() {
		resp.Diagnostics.AddError(
			"product_id needs to be set", "product_id is either null or unknown")
		return
	}

	body, err := productRevisionModMapper(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Couldn't map plan data into ProductRevisionMod struct", fmt.Sprintf("err: %v", err))
		return
	}

	newProductPTR, err := util.MakePRMarketplaceRequest[ProductRevisionResourceNativeModel](ctx, http.MethodPost, productRevisionResourcePath, bytes.NewReader(body), r.client) // TODO - switch to normal one when fixed
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %s", http.MethodPost, productRevisionResourcePath, body),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	// Set computed vars to null (new fetch required anyway)
	data = ProductRevisionModel{
		AdminSuggestion:                         util.StringSetOrNull(data.AdminSuggestion.ValueString()),
		Byol:                                    byolSetOrNull(data.Byol),
		Categories:                              util.ListSetOrNull(data.Categories, ContractualDocumentsInfoValue{}.Type(ctx)),
		ContractualDocuments:                    util.ListSetOrNull(data.ContractualDocuments, ContractualDocumentsValue{}.Type(ctx)),
		ContractualDocumentsInfo:                util.ListSetOrNull(data.ContractualDocumentsInfo, ContractualDocumentsInfoValue{}.Type(ctx)),
		Description:                             util.StringSetOrNull(data.Description.ValueString()),
		DescriptionShort:                        util.StringSetOrNull(data.DescriptionShort.ValueString()),
		Eula:                                    util.StringSetOrNull(data.Eula.ValueString()),
		Guidance:                                util.StringSetOrNull(data.Guidance.ValueString()),
		HelmExternal:                            util.StringSetOrNull(data.HelmExternal.ValueString()),
		Icon:                                    util.StringSetOrNull(data.Icon.ValueString()),
		Id:                                      util.StringSetOrNull(newProductPTR.Id), // The only thing that's updated from backend call as other fields are expected to change - don't want inconsistent state errors
		LicenseFee:                              util.StringSetOrNull(data.LicenseFee.ValueString()),
		LicenseInfo:                             util.StringSetOrNull(data.LicenseInfo.ValueString()),
		Number:                                  types.Int64Value(data.Number.ValueInt64()), // Uninitialized ints have a value of 0, so this shouldn't ever be null/unknown
		PostDeploymentInfo:                      util.StringSetOrNull(data.PostDeploymentInfo.ValueString()),
		PreDeploymentInfo:                       util.StringSetOrNull(data.PreDeploymentInfo.ValueString()),
		PricingInfo:                             util.StringSetOrNull(data.PricingInfo.ValueString()),
		ProductId:                               util.StringSetOrNull(data.ProductId.ValueString()),
		ProductRevisionApplicationConfiguration: util.ListSetOrNull(data.ProductRevisionApplicationConfiguration, ProductRevisionApplicationConfigurationValue{}.Type(ctx)),
		ProposedReleaseDate:                     util.StringSetOrNull(data.ProposedReleaseDate.ValueString()),
		ScheduledReleaseDate:                    util.StringSetOrNull(data.ScheduledReleaseDate.ValueString()),
		ScheduledReleaseUntilDate:               util.StringSetOrNull(data.ScheduledReleaseUntilDate.ValueString()),
		State:                                   util.StringSetOrNull(data.State.ValueString()),
		UsedSoftware:                            util.ListSetOrNull(data.UsedSoftware, UsedSoftwareValue{}.Type(ctx)),
		Version:                                 util.StringSetOrNull(data.Version.ValueString()),
	}

	if resp.Diagnostics.HasError() {
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	resp.Diagnostics.AddWarning("Expecting state drift", "Read required after Create. Run `terraform apply -refresh-only` now.")
}

func productRevisionModMapper(ctx context.Context, data ProductRevisionModel) ([]byte, error) {
	var tempCategories []string
	if !data.Categories.IsUnknown() { // TODO - only checking these for unknown what if it's null
		diags := data.Categories.ElementsAs(ctx, &tempCategories, false)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("couldn't convert data categories to list of str: elements: %v, tempCategories: %v, err: %+v", data.Categories.Elements(), tempCategories, diags.Errors()))
		}
	}

	var tempUsedSoft []productRevisionResourceUsedSoftwareNativeModel
	if !data.UsedSoftware.IsUnknown() {
		diags := data.UsedSoftware.ElementsAs(ctx, &tempUsedSoft, false)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("couldn't convert data software to list of productRevisionResourceUsedSoftwareNativeModel: elements: %v, tempUsedSoft: %v, err: %+v", data.UsedSoftware.Elements(), tempUsedSoft, diags.Errors()))
		}
	}

	var tempConfig []ProductRevisionResourceConfigNativeModel
	if !data.ProductRevisionApplicationConfiguration.IsUnknown() {
		diags := data.ProductRevisionApplicationConfiguration.ElementsAs(ctx, &tempConfig, false)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("couldn't convert configs to list of productRevisionResourceConfigNativeModel: elements: %v, tempUsedSoft: %v, err: %+v", data.ProductRevisionApplicationConfiguration.Elements(), tempConfig, diags.Errors()))
		}
	}

	var tempContDocs []ContractualDocsNativeModel
	if !data.ContractualDocuments.IsUnknown() {
		diags := data.ContractualDocuments.ElementsAs(ctx, &tempContDocs, false)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("couldn't convert docs to list of contractualDocsNativeModel: elements: %v, tempUsedSoft: %v, err: %+v", data.ContractualDocuments.Elements(), tempContDocs, diags.Errors()))
		}
	}

	body, err := json.Marshal(productRevisionResourceNativeModModel{
		Categories:           tempCategories,
		ProductId:            data.ProductId.ValueString(),
		Description:          data.Description.ValueString(),
		DescriptionShort:     data.DescriptionShort.ValueString(),
		Icon:                 data.Icon.ValueString(),
		PreDeploymentInfo:    data.PreDeploymentInfo.ValueString(),
		PostDeploymentInfo:   data.PostDeploymentInfo.ValueString(),
		ProposedReleaseDate:  data.ProposedReleaseDate.ValueString(),
		LicenseInfo:          data.LicenseInfo.ValueString(),
		Guidance:             data.Guidance.ValueString(),
		Version:              data.Version.ValueString(),
		HelmExternal:         data.HelmExternal.ValueString(),
		LicenseFee:           data.LicenseFee.ValueString(),
		PricingInfo:          data.PricingInfo.ValueString(),
		UsedSoftware:         tempUsedSoft,
		Configuration:        tempConfig,
		ContractualDocuments: tempContDocs,
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("couldn't marshall product into json: %v", err))
	}
	return body, err
}

func (r *productRevisionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProductRevisionModel
	var priorState ProductRevisionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	data.Id = priorState.Id

	if data.Id.IsNull() || data.Id.IsUnknown() {
		resp.Diagnostics.AddError(
			"Id needs to be set", "this resource has no ID")
		return
	}

	body, err := productRevisionModMapper(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Couldn't map plan data into ProductRevisionMod struct", fmt.Sprintf("err: %v", err))
		return
	}

	url := fmt.Sprintf("%s/%s", productRevisionResourcePath, util.SanitizeString(data.Id.ValueString()))
	newProductPTR, err := util.MakePRMarketplaceRequest[ProductRevisionResourceNativeModel](ctx, http.MethodPatch, url, bytes.NewReader(body), r.client) // TODO - switch to the normal one when fixed
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %s", http.MethodPatch, url, body),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	// Remove Unknown
	data = ProductRevisionModel{
		AdminSuggestion:                         util.StringSetOrNull(data.AdminSuggestion.ValueString()),
		Byol:                                    byolSetOrNull(data.Byol),
		Categories:                              util.ListSetOrNull(data.Categories, ContractualDocumentsInfoValue{}.Type(ctx)),
		ContractualDocuments:                    util.ListSetOrNull(data.ContractualDocuments, ContractualDocumentsValue{}.Type(ctx)),
		ContractualDocumentsInfo:                util.ListSetOrNull(data.ContractualDocumentsInfo, ContractualDocumentsInfoValue{}.Type(ctx)),
		Description:                             util.StringSetOrNull(data.Description.ValueString()),
		DescriptionShort:                        util.StringSetOrNull(data.DescriptionShort.ValueString()),
		Eula:                                    util.StringSetOrNull(data.Eula.ValueString()),
		Guidance:                                util.StringSetOrNull(data.Guidance.ValueString()),
		HelmExternal:                            util.StringSetOrNull(data.HelmExternal.ValueString()),
		Icon:                                    util.StringSetOrNull(data.Icon.ValueString()),
		Id:                                      util.StringSetOrNull(newProductPTR.Id), // The only thing that's updated from backend call as other fields are expected to change - don't want inconsistent state errors
		LicenseFee:                              util.StringSetOrNull(data.LicenseFee.ValueString()),
		LicenseInfo:                             util.StringSetOrNull(data.LicenseInfo.ValueString()),
		Number:                                  types.Int64Value(data.Number.ValueInt64()), // Uninitialized ints have a value of 0
		PostDeploymentInfo:                      util.StringSetOrNull(data.PostDeploymentInfo.ValueString()),
		PreDeploymentInfo:                       util.StringSetOrNull(data.PreDeploymentInfo.ValueString()),
		PricingInfo:                             util.StringSetOrNull(data.PricingInfo.ValueString()),
		ProductId:                               util.StringSetOrNull(data.ProductId.ValueString()),
		ProductRevisionApplicationConfiguration: util.ListSetOrNull(data.ProductRevisionApplicationConfiguration, ProductRevisionApplicationConfigurationValue{}.Type(ctx)),
		ProposedReleaseDate:                     util.StringSetOrNull(data.ProposedReleaseDate.ValueString()),
		ScheduledReleaseDate:                    util.StringSetOrNull(data.ScheduledReleaseDate.ValueString()),
		ScheduledReleaseUntilDate:               util.StringSetOrNull(data.ScheduledReleaseUntilDate.ValueString()),
		State:                                   util.StringSetOrNull(data.State.ValueString()),
		UsedSoftware:                            util.ListSetOrNull(data.UsedSoftware, UsedSoftwareValue{}.Type(ctx)),
		Version:                                 util.StringSetOrNull(data.Version.ValueString()),
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Warn(ctx, fmt.Sprintf("Read required after Update. Run `terraform apply -refresh-only` now."))
}

func (r *productRevisionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProductRevisionModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	url := fmt.Sprintf("%s/%s", productRevisionResourcePath, util.SanitizeString(data.Id.ValueString()))
	_, err := util.MakeMarketplaceRequest[struct{}](ctx, http.MethodDelete, url, nil, r.client)
	if err != nil {
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodDelete, url, nil),
			fmt.Sprintf("ignoring error in the hope the parent product will be deleted later as part of the plan."+
				"\nif this resource was to be replaced, please recreate its parent.\nerror: %v", err),
		)
		//return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func byolSetOrNull(b ByolValue) ByolValue {
	if b.IsNull() || b.IsUnknown() {
		return NewByolValueNull()
	}
	return b
}
