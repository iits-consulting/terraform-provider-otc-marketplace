package datasource_sales_history

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const salesHistoryServicePath = "/sales-history"

var _ datasource.DataSource = (*salesHistoryDataSource)(nil)

func NewSalesHistoryDataSource() datasource.DataSource {
	return &salesHistoryDataSource{}
}

type salesHistoryDataSource struct {
	client *util.MarketplaceAPIClient
}

type salesHistoryDataSourceNativeModel struct {
	ProductRevisionId     string `json:"product_revision_id,omitempty"`
	ProductId             string `json:"product_id,omitempty"`
	ProductName           string `json:"product_name,omitempty"`
	CustomerCompanyName   string `json:"customer_company_name,omitempty"`
	CustomerCompanyUrl    string `json:"customer_company_url,omitempty"`
	CustomerContactNumber string `json:"customer_contact_number,omitempty"`
	CustomerContactEmail  string `json:"customer_contact_email,omitempty"`
	DeployedAt            string `json:"deployed_at,omitempty"`
}

func (d *salesHistoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sales_history"
}

func (d *salesHistoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SalesHistoryDataSourceSchema(ctx)
}

func (d *salesHistoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *salesHistoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SalesHistoryModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]salesHistoryDataSourceNativeModel](ctx, http.MethodGet, salesHistoryServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, salesHistoryServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeSales := range *newDataNativeArrPTR {
		sale, diags := NewSalesHistoryValue(SalesHistoryValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"product_revision_id":     types.StringValue(nativeSales.ProductRevisionId),
			"product_id":              types.StringValue(nativeSales.ProductId),
			"product_name":            types.StringValue(nativeSales.ProductName),
			"customer_company_name":   types.StringValue(nativeSales.CustomerCompanyName),
			"customer_company_url":    types.StringValue(nativeSales.CustomerCompanyUrl),
			"customer_contact_number": types.StringValue(nativeSales.CustomerContactNumber),
			"customer_contact_email":  types.StringValue(nativeSales.CustomerContactEmail),
			"deployed_at":             types.StringValue(nativeSales.DeployedAt),
		})
		if diags.HasError() {
			return
		}
		newData = append(newData, sale)
	}

	salesSet, diags := types.SetValue(SalesHistoryValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = SalesHistoryModel{SalesHistory: salesSet}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
