package datasource_projects

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-otc-marketplace/internal/util"
)

const projectServicePath = "/projects"

var _ datasource.DataSource = (*projectDataSource)(nil)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *util.MarketplaceAPIClient
}

type projectDataSourceNativeModel struct {
	Id   string `json:"id,omitempty" tfsdk:"id"`
	Name string `json:"name,omitempty" tfsdk:"name"`
}

func (d *projectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ProjectsDataSourceSchema(ctx)
}

func (d *projectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProjectsModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	newDataNativeArrPTR, err := util.MakeMarketplaceRequest[[]projectDataSourceNativeModel](ctx, http.MethodGet, projectServicePath, nil, d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Couldn't send %s to %s with a body of %v", http.MethodGet, projectServicePath, nil),
			fmt.Sprintf("error: %v", err),
		)
		return
	}

	var newData []attr.Value
	for _, nativeProj := range *newDataNativeArrPTR {
		projObj, diags := NewProjectsValue(ProjectsValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"id":   types.StringValue(nativeProj.Id),
			"name": types.StringValue(nativeProj.Name),
		})
		if diags.HasError() {
			return
		}
		newData = append(newData, projObj)
	}

	projectSet, diags := types.SetValue(ProjectsValue{}.Type(ctx), newData)
	if diags.HasError() {
		return
	}

	data = ProjectsModel{
		Projects: projectSet}

	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
