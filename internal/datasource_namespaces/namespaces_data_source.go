package datasource_namespaces

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const namespaceServicePath = "/namespaces"

var _ datasource.DataSource = (*namespaceDataSource)(nil)

func NewNamespaceDataSource() datasource.DataSource {
	return &namespaceDataSource{}
}

type namespaceDataSource struct {
	client *util.MarketplaceAPIClient
}

type namespaceDataSourceNativeModel struct {
	Name      string `json:"name,omitempty" tfsdk:"name"`
	ClusterId string `json:"cluster_id,omitempty" tfsdk:"cluster_id"`
	ProjectId string `json:"project_id,omitempty" tfsdk:"project_id"`
}

func (d *namespaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

func (d *namespaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = NamespacesDataSourceSchema(ctx)
}

func (d *namespaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *namespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NamespacesModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	url := fmt.Sprintf("%s?project_id=%s&cluster_id=%s", namespaceServicePath, data.ProjectId.ValueString(), data.ClusterId.ValueString())
	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]namespaceDataSourceNativeModel](ctx, http.MethodGet, url, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, url, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeNamespaces := range *newDataNativeArrPTR {
		namespaceObj, diags := NewNamespacesValue(NamespacesValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"name":       types.StringValue(nativeNamespaces.Name),
			"project_id": types.StringValue(nativeNamespaces.ProjectId),
			"cluster_id": types.StringValue(nativeNamespaces.ClusterId),
		})
		if diags.HasError() {
			return
		}
		newData = append(newData, namespaceObj)
	}

	namespaceSet, diags := types.SetValue(NamespacesValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = NamespacesModel{
		ClusterId:  data.ClusterId,
		ProjectId:  data.ProjectId,
		Namespaces: namespaceSet,
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
