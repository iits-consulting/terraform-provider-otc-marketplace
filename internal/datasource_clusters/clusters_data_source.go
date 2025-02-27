package datasource_clusters

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const clusterServicePath = "/clusters"

var _ datasource.DataSource = (*clusterDataSource)(nil)

func NewClusterDataSource() datasource.DataSource {
	return &clusterDataSource{}
}

type clusterDataSource struct {
	client *util.MarketplaceAPIClient
}

type clusterDataSourceNativeModel struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (d *clusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *clusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ClustersDataSourceSchema(ctx)
}

func (d *clusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *clusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClustersModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	url := fmt.Sprintf("%s?project_id=%s", clusterServicePath, data.ProjectId.ValueString())
	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]clusterDataSourceNativeModel](ctx, http.MethodGet, url, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, url, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeClusters := range *newDataNativeArrPTR {
		clustObj, diags := NewClustersValue(ClustersValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"id":   types.StringValue(nativeClusters.Id),
			"name": types.StringValue(nativeClusters.Name),
		})
		if diags.HasError() {
			return
		}
		newData = append(newData, clustObj)
	}

	clusterSet, diags := types.SetValue(ClustersValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = ClustersModel{
		Clusters:  clusterSet,
		ProjectId: data.ProjectId,
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
