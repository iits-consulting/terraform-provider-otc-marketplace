package datasource_categories

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const categoryServicePath = "/categories"

var _ datasource.DataSource = (*categoryDataSource)(nil)

func NewCategoryDataSource() datasource.DataSource {
	return &categoryDataSource{}
}

type categoryDataSource struct {
	client *util.MarketplaceAPIClient
}

type categoryDataSourceNativeModel struct {
	Id          string `json:"id,omitempty" tfsdk:"id"`
	Description string `json:"description,omitempty" tfsdk:"description"`
	Name        string `json:"name,omitempty" tfsdk:"name"`
	State       string `json:"state,omitempty" tfsdk:"state"`
	Position    int64  `json:"position,omitempty" tfsdk:"position"`
}

func (d *categoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_category"
}

func (d *categoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = CategoriesDataSourceSchema(ctx)
}

func (d *categoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *categoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CategoriesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]categoryDataSourceNativeModel](ctx, http.MethodGet, categoryServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, categoryServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeCategories := range *newDataNativeArrPTR {
		catObj, diags := NewCategoriesValue(CategoriesValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"id":          types.StringValue(nativeCategories.Id),
			"description": types.StringValue(nativeCategories.Description),
			"name":        types.StringValue(nativeCategories.Name),
			"position":    types.Int64Value(nativeCategories.Position),
			"state":       types.StringValue(nativeCategories.State),
		})
		if diags.HasError() {
			return
		}
		newData = append(newData, catObj)
	}

	categorySet, diags := types.SetValue(CategoriesValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = CategoriesModel{Categories: categorySet}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
