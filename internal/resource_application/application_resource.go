package resource_application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"terraform-provider-otc-marketplace/internal/datasource_applications"
	"terraform-provider-otc-marketplace/internal/resource_product"
	"terraform-provider-otc-marketplace/internal/resource_product_revision"
	"terraform-provider-otc-marketplace/internal/util"
)

var _ resource.Resource = (*applicationResource)(nil)

func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

const applicationResourcePath = "/applications"

type applicationResource struct {
	client *util.MarketplaceAPIClient
}

func (r *applicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *applicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApplicationResourceSchema(ctx)
}

func (r *applicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// TODO - Can Read (GET (both), can't POST, can't DELETE)

type ApplicationNativeModel struct {
	ByolLicense       string                                                       `json:"byol_license,omitempty" tfsdk:"byol_license"`
	ClusterId         string                                                       `json:"cluster_id,omitempty" tfsdk:"cluster_id"`
	Configuration     []datasource_applications.ConfigurationNativeModel           `json:"configuration,omitempty" tfsdk:"configuration"` // TODO - naming
	CreatedAt         string                                                       `json:"created_at,omitempty" tfsdk:"created_at"`
	Description       string                                                       `json:"description,omitempty" tfsdk:"description"`
	Id                string                                                       `json:"id" tfsdk:"id"`
	Namespace         string                                                       `json:"namespace,omitempty" tfsdk:"namespace"`
	Product           util.ProductDataSourceNativeModel                            `json:"product,omitempty" tfsdk:"product"` // TODO - naming
	ProductRevision   resource_product_revision.ProductRevisionResourceNativeModel `json:"product_revision,omitempty" tfsdk:"product_revision"`
	ProductRevisionId string                                                       `json:"product_revision_id,omitempty" tfsdk:"product_revision_id"`
	ProjectId         string                                                       `json:"project_id,omitempty" tfsdk:"project_id"`
	ReleaseName       string                                                       `json:"release_name,omitempty" tfsdk:"release_name"`
	Seller            util.SellerNativeModel                                       `json:"seller,omitempty" tfsdk:"seller"`
	State             string                                                       `json:"state,omitempty" tfsdk:"state"`
	Username          string                                                       `json:"username,omitempty" tfsdk:"username"`
}

type ApplicationResourceModNativeModel struct {
	// Required
	ProductRevisionId string `json:"product_revision_id,omitempty" tfsdk:"product_revision_id"`
	ProjectId         string `json:"project_id,omitempty" tfsdk:"project_id"`
	ClusterId         string `json:"cluster_id,omitempty" tfsdk:"cluster_id"`
	Namespace         string `json:"namespace,omitempty" tfsdk:"namespace"`
	// Optional
	ReleaseName   string                                             `json:"release_name,omitempty" tfsdk:"release_name"`
	Description   string                                             `json:"description,omitempty" tfsdk:"description"`
	State         string                                             `json:"state,omitempty" tfsdk:"state"`
	Username      string                                             `json:"username,omitempty" tfsdk:"username"`
	ByolLicense   string                                             `json:"byol_license,omitempty" tfsdk:"byol_license"`
	Configuration []datasource_applications.ConfigurationNativeModel `json:"application_configuration,omitempty" tfsdk:"application_configuration"`
}

func applicationResourceModMapper(ctx context.Context, data ApplicationModel) ([]byte, error) {
	var tempConfigs []datasource_applications.ConfigurationNativeModel
	if !data.Configuration.IsUnknown() {
		diags := data.Configuration.ElementsAs(ctx, &tempConfigs, false)
		if diags.HasError() {
			return nil, errors.New(fmt.Sprintf("couldn't convert data categories to list of str: elements: %v, tempCategories: %v, err: %+v", data.Configuration.Elements(), tempConfigs, diags.Errors()))
		}
	}
	body, err := json.Marshal(ApplicationResourceModNativeModel{
		ProductRevisionId: data.ProductRevisionId.ValueString(),
		ProjectId:         data.ProjectId.ValueString(),
		ClusterId:         data.ClusterId.ValueString(),
		Namespace:         data.Namespace.ValueString(),
		ReleaseName:       data.ReleaseName.ValueString(),
		Description:       data.Description.ValueString(),
		State:             data.State.ValueString(),
		Username:          data.Username.ValueString(),
		ByolLicense:       data.ByolLicense.ValueString(),
		Configuration:     tempConfigs,
	})

	return body, err
}

// TODO - this is dumb, a much cleaner way exists
func productProductModeltoApplicationProductModel(ctx context.Context, in resource_product.ProductModel) (*ProductValue, error) {
	llmObj, diags := in.LlmHub.ToObjectValue(ctx)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("%+v", diags.Errors())) // TODO - gross
	}

	sellerObj, diags := in.Seller.ToObjectValue(ctx)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("%+v", diags.Errors())) // TODO - gross
	}

	return &ProductValue{
		CreatedAt:   in.CreatedAt,
		Eol:         in.Eol,
		EolDate:     in.EolDate,
		Id:          in.Id,
		LicenseType: in.LicenseType,
		LlmHub:      llmObj,
		Seller:      sellerObj,
		Name:        in.Name,
		ProductType: in.Type,
		Weight:      in.Weight,
	}, nil
}

// TODO - this is dumb, a much cleaner way exists
func productRevisionModeltoApplicationProductRevision(ctx context.Context, in resource_product_revision.ProductRevisionModel) (*ProductRevisionValue, error) {
	byolObj, diags := in.Byol.ToObjectValue(ctx)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("err diags: %+v", diags.Errors())) // TODO - gross
	}
	return &ProductRevisionValue{
		AdminSuggestion:                         in.AdminSuggestion,
		Byol:                                    byolObj,
		Categories:                              in.Categories,
		ContractualDocuments:                    in.ContractualDocuments,
		ContractualDocumentsInfo:                in.ContractualDocumentsInfo,
		Description:                             in.Description,
		DescriptionShort:                        in.DescriptionShort,
		Eula:                                    in.Eula,
		Guidance:                                in.Guidance,
		HelmExternal:                            in.HelmExternal,
		Icon:                                    in.Icon,
		Id:                                      in.Id,
		LicenseFee:                              in.LicenseFee,
		LicenseInfo:                             in.LicenseInfo,
		Number:                                  in.Number,
		PostDeploymentInfo:                      in.PostDeploymentInfo,
		PreDeploymentInfo:                       in.PreDeploymentInfo,
		PricingInfo:                             in.PricingInfo,
		ProductId:                               in.ProductId,
		ProductRevisionApplicationConfiguration: in.ProductRevisionApplicationConfiguration,
		ProposedReleaseDate:                     in.ProposedReleaseDate,
		ScheduledReleaseDate:                    in.ScheduledReleaseDate,
		ScheduledReleaseUntilDate:               in.ScheduledReleaseUntilDate,
		State:                                   in.State,
		UsedSoftware:                            in.UsedSoftware,
		Version:                                 in.Version,
	}, nil
}

func applicationResourceMapper(ctx context.Context, newDataPTR *ApplicationNativeModel) (*ApplicationModel, error) {
	var diags diag.Diagnostics

	productPMObj, err := resource_product.ProductResourceMapper(ctx, &newDataPTR.Product)
	if err != nil {
		return nil, err
	}

	productObj, err := productProductModeltoApplicationProductModel(ctx, *productPMObj)
	if err != nil {
		return nil, err
	}

	productRevisionPRM, err := resource_product_revision.ProductRevisionMapper(ctx, &newDataPTR.ProductRevision)
	if err != nil {
		return nil, err
	}
	productRevisionObj, err := productRevisionModeltoApplicationProductRevision(ctx, *productRevisionPRM)
	if err != nil {
		return nil, err
	}

	var sellerObj ApplicationSellerValue
	sellerObj = ApplicationSellerValue{
		Description:  util.StringSetOrNull(newDataPTR.Seller.Description),
		Id:           util.StringSetOrNull(newDataPTR.Seller.Id),
		Name:         util.StringSetOrNull(newDataPTR.Seller.Name),
		State:        util.StringSetOrNull(newDataPTR.Seller.State),
		SupportEmail: util.StringSetOrNull(newDataPTR.Seller.SupportEmail),
		SupportUrl:   util.StringSetOrNull(newDataPTR.Seller.SupportUrl),
	}

	var tempConf []ConfigurationValue
	for _, conf := range newDataPTR.Configuration {
		confObj, confDiags := NewConfigurationValue(ConfigurationValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"key":   types.StringValue(conf.Key),
			"value": types.StringValue(conf.Value),
		})
		diags.Append(confDiags...)
		tempConf = append(tempConf, confObj)
	}
	configsAsList := util.ListValueOrNull[ConfigurationValue](ctx, ConfigurationValue{}.Type(ctx), tempConf, &diags)
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("error: %v", diags.Errors()))
	}

	return &ApplicationModel{
		ByolLicense:       util.StringSetOrNull(newDataPTR.ByolLicense),
		ClusterId:         util.StringSetOrNull(newDataPTR.ClusterId),
		Configuration:     configsAsList,
		CreatedAt:         util.StringSetOrNull(newDataPTR.CreatedAt),
		Description:       util.StringSetOrNull(newDataPTR.Description),
		Id:                util.StringSetOrNull(newDataPTR.Id),
		Namespace:         util.StringSetOrNull(newDataPTR.Namespace),
		Product:           *productObj,
		ProductRevision:   *productRevisionObj,
		ProductRevisionId: util.StringSetOrNull(newDataPTR.ProductRevisionId),
		ProjectId:         util.StringSetOrNull(newDataPTR.ProjectId),
		ReleaseName:       util.StringSetOrNull(newDataPTR.ReleaseName),
		ApplicationSeller: sellerObj,
		State:             util.StringSetOrNull(newDataPTR.State),
		Username:          util.StringSetOrNull(newDataPTR.Username),
	}, nil
}

func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	if data.ProductRevisionId.IsNull() || data.ProductRevisionId.IsUnknown() {
		resp.Diagnostics.AddError(
			"product_revision_id needs to be set", "product_revision_id is either null or unknown")
		return
	}

	if data.ProjectId.IsNull() || data.ProjectId.IsUnknown() {
		resp.Diagnostics.AddError(
			"project_id needs to be set", "project_id is either null or unknown")
		return
	}

	if data.ClusterId.IsNull() || data.ClusterId.IsUnknown() {
		resp.Diagnostics.AddError(
			"cluster_id needs to be set", "cluster_id is either null or unknown")
		return
	}

	if data.Namespace.IsNull() || data.Namespace.IsUnknown() {
		resp.Diagnostics.AddError(
			"namespace needs to be set", "namespace is either null or unknown")
		return
	}
	// Example data value setting
	// data.Id = types.StringValue("example-id")

	body, err := applicationResourceModMapper(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Couldn't map plan data into ApplicationMod struct", fmt.Sprintf("err: %v", err))
		return
	}

	newProductPTR, err := util.MakeMarketplaceRequest[ApplicationNativeModel](ctx, http.MethodPost, applicationResourcePath, bytes.NewReader(body), r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %s", http.MethodPost, applicationResourcePath, body),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	dataPTR, err := applicationResourceMapper(ctx, newProductPTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	dataPTR.ProductRevisionId = data.ProductRevisionId

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &dataPTR)...)
	tflog.Warn(ctx, fmt.Sprintf("Read required after Create. Run `terraform apply -refresh-only` now."))
}

func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if data.Id.IsNull() || data.Id.IsUnknown() || data.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"resource id needs to be set", "resource id is either null or unknown")
		return
	}

	url := fmt.Sprintf("%s/%s", applicationResourcePath, util.SanitizeString(data.Id.ValueString()))
	newDataNativePTR, err := util.MakeMarketplaceRequest[ApplicationNativeModel](ctx, http.MethodGet, url, nil, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, url, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	dataPTR, err := applicationResourceMapper(ctx, newDataNativePTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, dataPTR)...)
}

// TODO - the openapi yaml doesn't define any Update (Patch) methods, so this might just not be implemented on the backend
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if data.ProductRevisionId.IsNull() || data.ProductRevisionId.IsUnknown() {
		resp.Diagnostics.AddError(
			"product_revision_id needs to be set", "product_revision_id is either null or unknown")
		return
	}

	if data.ProjectId.IsNull() || data.ProjectId.IsUnknown() {
		resp.Diagnostics.AddError(
			"project_id needs to be set", "project_id is either null or unknown")
		return
	}

	if data.ClusterId.IsNull() || data.ClusterId.IsUnknown() {
		resp.Diagnostics.AddError(
			"cluster_id needs to be set", "cluster_id is either null or unknown")
		return
	}

	if data.Namespace.IsNull() || data.Namespace.IsUnknown() {
		resp.Diagnostics.AddError(
			"namespace needs to be set", "namespace is either null or unknown")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	body, err := applicationResourceModMapper(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Couldn't map plan data into ApplicationMod struct", fmt.Sprintf("err: %v", err))
		return
	}

	newProductPTR, err := util.MakeMarketplaceRequest[ApplicationNativeModel](ctx, http.MethodPatch, applicationResourcePath, bytes.NewReader(body), r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %s", http.MethodPatch, applicationResourcePath, body),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	dataPTR, err := applicationResourceMapper(ctx, newProductPTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, dataPTR)...)
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	url := fmt.Sprintf("%s/%s", applicationResourcePath, util.SanitizeString(data.Id.ValueString()))
	_, err := util.MakeMarketplaceRequest[struct{}](ctx, http.MethodDelete, url, nil, r.client)
	if err != nil {
		// TODO - 500s when trying to delete an Application with the install still visible on https://marketplace.otc.t-systems.com/dashboard -> Test Deployment Workload
		resp.Diagnostics.AddWarning(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodDelete, url, nil),
			fmt.Sprintf("error: %+v\n ignoring error and carrying on...", err),
		)
		//return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}
