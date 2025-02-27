package resource_product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

var _ resource.Resource = (*productResource)(nil)

const productResourcePath = "/products"

func NewProductResource() resource.Resource {
	return &productResource{}
}

type productResource struct {
	client *util.MarketplaceAPIClient
}

func (r *productResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

func (r *productResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProductResourceSchema(ctx)
}

func (r *productResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *productResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProductModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if data.LicenseType.IsNull() || data.LicenseType.IsUnknown() {
		resp.Diagnostics.AddError(
			"license_type needs to be set", "license_type is either null or unknown")
		return
	}

	if data.Name.IsNull() || data.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"name needs to be set", "name is either null or unknown")
		return
	}

	if data.Type.IsNull() || data.Type.IsUnknown() {
		resp.Diagnostics.AddError(
			"type needs to be set", "type is either null or unknown")
		return
	}

	if data.Weight.IsNull() || data.Weight.IsUnknown() {
		resp.Diagnostics.AddError(
			"weight needs to be set", "weight is either null or unknown")
		return
	}

	// TODO - send the whole Product? - Potential inconsistent state issues with stuff like time
	type CreateJSONRequest struct {
		LicenseType string `json:"license_type"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Weight      int64  `json:"weight"`
	}

	body, err := json.Marshal(CreateJSONRequest{
		LicenseType: util.SanitizeString(data.LicenseType.String()),
		Name:        util.SanitizeString(data.Name.String()),
		Type:        util.SanitizeString(data.Type.String()),
		Weight:      data.Weight.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Couldn't marshall product into json", fmt.Sprintf("error: %v", err))
		return
	}

	newProductPTR, err := util.MakeMarketplaceRequest[util.ProductDataSourceNativeModel](ctx, http.MethodPost, productResourcePath, bytes.NewReader(body), r.client)
	if err != nil {
		var potentialReusedName string
		if err.Error() == "unexpected status code: 500" { // TODO - fragile
			potentialReusedName = "This might mean you're trying to create a product with a previously used name. \nProduct names on the OTC must be new and unique."
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %s", http.MethodPost, productResourcePath, body),
			fmt.Sprintf("%s %+v", potentialReusedName, err),
		)
		return
	}

	dataPTR, err := ProductResourceMapper(ctx, newProductPTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	data = *dataPTR

	if data.EolDate.IsNull() || data.EolDate.IsUnknown() || data.EolDate.ValueString() == "" {
		data.Eol = types.BoolValue(false)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *productResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProductModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	url := fmt.Sprintf("%s/%s", productResourcePath, util.SanitizeString(data.Id.ValueString()))
	newProductPTR, err := util.MakeMarketplaceRequest[util.ProductDataSourceNativeModel](ctx, http.MethodGet, url, nil, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, url, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	dataPTR, err := ProductResourceMapper(ctx, newProductPTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	data = *dataPTR

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *productResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProductModel
	var priorState ProductModel

	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	data.Id = priorState.Id

	if data.Id.IsNull() || data.Id.IsUnknown() {
		resp.Diagnostics.AddError(
			"Id needs to be set", "this resource has no ID")
		return
	}

	if data.LicenseType.IsNull() || data.LicenseType.IsUnknown() {
		resp.Diagnostics.AddError(
			"license_type needs to be set", "license_type is either null or unknown")
		return
	}

	if data.Name.IsNull() || data.Name.IsUnknown() {
		resp.Diagnostics.AddError(
			"name needs to be set", "name is either null or unknown")
		return
	}

	if data.Type.IsNull() || data.Type.IsUnknown() {
		resp.Diagnostics.AddError(
			"type needs to be set", "type is either null or unknown")
		return
	}

	if data.Weight.IsNull() || data.Weight.IsUnknown() {
		resp.Diagnostics.AddError(
			"weight needs to be set", "weight is either null or unknown")
		return
	}

	data.LicenseType = util.SanitizeStringValue(data.LicenseType)
	data.Name = util.SanitizeStringValue(data.Name)
	data.Type = util.SanitizeStringValue(data.Type)

	type UpdateJSONRequest struct {
		EOL         bool   `json:"eol"`
		LicenseType string `json:"license_type"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Weight      int64  `json:"weight"`
	}

	body, err := json.Marshal(UpdateJSONRequest{
		EOL:         data.Eol.ValueBool(),
		LicenseType: util.SanitizeString(data.LicenseType.String()),
		Name:        util.SanitizeString(data.Name.String()),
		Type:        util.SanitizeString(data.Type.String()),
		Weight:      data.Weight.ValueInt64(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Couldn't marshall product into json", fmt.Sprintf("error: %v", err))
		return
	}

	url := fmt.Sprintf("%s/%s", productResourcePath, util.SanitizeString(data.Id.ValueString()))
	newProductPTR, err := util.MakeMarketplaceRequest[util.ProductDataSourceNativeModel](ctx, http.MethodPatch, url, bytes.NewReader(body), r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %s", http.MethodPatch, url, body),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	dataPTR, err := ProductResourceMapper(ctx, newProductPTR)
	if err != nil {
		resp.Diagnostics.AddError("Couldn't map response to product resource", fmt.Sprintf("error: %v", err))
		return
	}

	dataPTR.Eol = data.Eol // TODO - better idea to check if EolDate has been set instead?

	data = *dataPTR // Done for "EOL" to prevent state mismatch

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func ProductResourceMapper(ctx context.Context, newProductPTR *util.ProductDataSourceNativeModel) (*ProductModel, error) {
	sellerObj, diags := NewSellerValue(SellerValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"description":   util.StringSetOrNull(newProductPTR.Seller.Description),
		"id":            util.StringSetOrNull(newProductPTR.Seller.Id),
		"name":          util.StringSetOrNull(newProductPTR.Seller.Name),
		"state":         util.StringSetOrNull(newProductPTR.Seller.State),
		"support_email": util.StringSetOrNull(newProductPTR.Seller.SupportEmail),
		"support_url":   util.StringSetOrNull(newProductPTR.Seller.SupportUrl),
	})
	if diags.HasError() {
		return nil, errors.New(fmt.Sprintf("%+v", diags.Errors()))
	}
	data := ProductModel{
		ActiveRevisionId: util.StringSetOrNull(newProductPTR.ActiveRevisionId),
		CreatedAt:        util.StringSetOrNull(newProductPTR.CreatedAt),
		Eol:              types.BoolValue(newProductPTR.EOL),
		EolDate:          util.StringSetOrNull(newProductPTR.EOLDate),
		Id:               util.StringSetOrNull(newProductPTR.Id),
		LicenseType:      util.StringSetOrNull(newProductPTR.LicenseType),
		Name:             util.StringSetOrNull(newProductPTR.Name),
		Seller:           sellerObj,
		State:            util.StringSetOrNull(newProductPTR.State),
		Type:             util.StringSetOrNull(newProductPTR.Type),
		Weight:           types.Int64Value(newProductPTR.Weight),
	}

	return &data, nil
}

func (r *productResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProductModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	url := fmt.Sprintf("%s/%s", productResourcePath, util.SanitizeString(data.Id.ValueString()))
	_, err := util.MakeMarketplaceRequest[struct{}](ctx, http.MethodDelete, url, nil, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodDelete, url, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

}
