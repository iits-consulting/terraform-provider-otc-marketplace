package datasource_whoami

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const whoamiServicePath = "/whoami"

var _ datasource.DataSource = (*whoamiDataSource)(nil)

func NewWhoamiDataSource() datasource.DataSource {
	return &whoamiDataSource{}
}

type whoamiDataSource struct {
	client *util.MarketplaceAPIClient
}

// Needed for JSON decoding. Reflection is too expensive and complicated as it will need to go both ways.
type whoamiDataSourceNativeModel struct {
	DomainName    string `json:"domain_name,omitempty" tfsdk:"domain_name"`
	LastProjectId string `json:"last_project_id,omitempty" tfsdk:"last_project_id"`
	Username      string `json:"username,omitempty" tfsdk:"username"`
	LLMHub        bool   `json:"llm_hub,omitempty" tfsdk:"llm_hub"`
}

func (d *whoamiDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_whoami"
}

func (d *whoamiDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = WhoamiDataSourceSchema(ctx)
}

func (d *whoamiDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *whoamiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WhoamiModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newDataNativePTR, err := util.MakeMarketplaceRequest[whoamiDataSourceNativeModel](ctx, http.MethodGet, whoamiServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, whoamiServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	data.DomainName = types.StringValue(newDataNativePTR.DomainName)
	data.Username = types.StringValue(newDataNativePTR.Username)
	data.LlmHub = types.BoolValue(newDataNativePTR.LLMHub)
	data.LastProjectId = types.StringValue(newDataNativePTR.LastProjectId)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
