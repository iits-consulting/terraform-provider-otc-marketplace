package datasource_profile

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const profileServicePath = "/profiles/profile"

var _ datasource.DataSource = (*profileDataSource)(nil)

func NewProfileDataSource() datasource.DataSource {
	return &profileDataSource{}
}

type profileDataSource struct {
	client *util.MarketplaceAPIClient
}

type profileDataSourceNativeModel struct {
	CustomerSupportNumber     string `json:"customer_support_number,omitempty" tfsdk:"customer_support_number"`
	Description               string `json:"description,omitempty" tfsdk:"description"`
	Email                     string `json:"email,omitempty" tfsdk:"email"`
	Id                        string `json:"id,omitempty" tfsdk:"id"`
	Name                      string `json:"name,omitempty" tfsdk:"name"`
	Status                    string `json:"status,omitempty" tfsdk:"status"`
	SupportEmail              string `json:"support_email,omitempty" tfsdk:"support_email"`
	SupportUrl                string `json:"support_url,omitempty" tfsdk:"support_url"`
	TempCustomerSupportNumber string `json:"temp_customer_support_number,omitempty" tfsdk:"temp_customer_support_number"`
	TempDescription           string `json:"temp_description,omitempty" tfsdk:"temp_description"`
	TempEmail                 string `json:"temp_email,omitempty" tfsdk:"temp_email"`
	TempName                  string `json:"temp_name,omitempty" tfsdk:"temp_name"`
	TempSupportEmail          string `json:"temp_support_email,omitempty" tfsdk:"temp_support_email"`
	TempSupportUrl            string `json:"temp_support_url,omitempty" tfsdk:"temp_support_url"`
}

func (d *profileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (d *profileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProfileDataSourceSchema(ctx)
}

func (d *profileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *profileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProfileModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newDataNativePTR, err := util.MakeMarketplaceRequest[profileDataSourceNativeModel](ctx, http.MethodGet, profileServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, profileServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	data = ProfileModel{
		CustomerSupportNumber:     types.StringValue(newDataNativePTR.CustomerSupportNumber),
		Description:               types.StringValue(newDataNativePTR.Description),
		Email:                     types.StringValue(newDataNativePTR.Email),
		Id:                        types.StringValue(newDataNativePTR.Id),
		Name:                      types.StringValue(newDataNativePTR.Name),
		Status:                    types.StringValue(newDataNativePTR.Status),
		SupportEmail:              types.StringValue(newDataNativePTR.SupportEmail),
		SupportUrl:                types.StringValue(newDataNativePTR.SupportUrl),
		TempCustomerSupportNumber: types.StringValue(newDataNativePTR.TempCustomerSupportNumber),
		TempDescription:           types.StringValue(newDataNativePTR.TempDescription),
		TempEmail:                 types.StringValue(newDataNativePTR.TempEmail),
		TempName:                  types.StringValue(newDataNativePTR.TempName),
		TempSupportEmail:          types.StringValue(newDataNativePTR.TempSupportEmail),
		TempSupportUrl:            types.StringValue(newDataNativePTR.TempSupportUrl),
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
