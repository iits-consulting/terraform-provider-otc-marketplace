package datasource_products

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const productServicePath = "/products"

var _ datasource.DataSource = (*productDataSource)(nil)

func NewProductDataSource() datasource.DataSource {
	return &productDataSource{}
}

type productDataSource struct {
	client *util.MarketplaceAPIClient
}

func (d *productDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product"
}

func (d *productDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProductsDataSourceSchema(ctx)
}

func (d *productDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *productDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProductsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]util.ProductDataSourceNativeModel](ctx, http.MethodGet, productServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, productServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeProducts := range *newDataNativeArrPTR {
		sellerObj, diags := NewSellerValue(SellerValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"description":   types.StringValue(nativeProducts.Seller.Description),
			"id":            types.StringValue(nativeProducts.Seller.Id),
			"name":          types.StringValue(nativeProducts.Seller.Name),
			"state":         types.StringValue(nativeProducts.Seller.State),
			"support_email": types.StringValue(nativeProducts.Seller.SupportEmail),
			"support_url":   types.StringValue(nativeProducts.Seller.SupportUrl),
		})
		if diags.HasError() {
			return
		}

		productObj, diags := NewProductsValue(ProductsValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"id":                 types.StringValue(nativeProducts.Id),
			"name":               types.StringValue(nativeProducts.Name),
			"created_at":         types.StringValue(nativeProducts.CreatedAt),
			"eol":                types.BoolValue(nativeProducts.EOL),
			"eol_date":           types.StringValue(nativeProducts.EOLDate),
			"license_type":       types.StringValue(nativeProducts.LicenseType),
			"seller":             sellerObj,
			"state":              types.StringValue(nativeProducts.State),
			"weight":             types.Int64Value(nativeProducts.Weight),
			"type":               types.StringValue(nativeProducts.Type),
			"active_revision_id": types.StringValue(nativeProducts.ActiveRevisionId),
		})
		if diags.HasError() {
			return
		}

		newData = append(newData, productObj)
	}

	productSet, diags := types.SetValue(ProductsValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = ProductsModel{Products: productSet}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
