package provider_marketplace

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"net/http"
	"terraform-provider-otc-marketplace/internal/datasource_applications"
	"terraform-provider-otc-marketplace/internal/datasource_categories"
	"terraform-provider-otc-marketplace/internal/datasource_clusters"
	"terraform-provider-otc-marketplace/internal/datasource_namespaces"
	"terraform-provider-otc-marketplace/internal/datasource_product_revisions"
	"terraform-provider-otc-marketplace/internal/datasource_products"
	"terraform-provider-otc-marketplace/internal/datasource_profile"
	"terraform-provider-otc-marketplace/internal/datasource_projects"
	"terraform-provider-otc-marketplace/internal/datasource_sales_history"
	"terraform-provider-otc-marketplace/internal/datasource_whoami"
	"terraform-provider-otc-marketplace/internal/resource_application"
	"terraform-provider-otc-marketplace/internal/resource_product"
	"terraform-provider-otc-marketplace/internal/resource_product_revision"
	"terraform-provider-otc-marketplace/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*marketplaceProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &marketplaceProvider{}
	}
}

type marketplaceProvider struct {
	DomainName string `tfsdk:"domain_name"`
	Username   string `tfsdk:"username"`
	Password   string `tfsdk:"password"`
}

func (p *marketplaceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "The domain name for authentication.",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The username for authentication.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for authentication.",
			},
		},
	}
}

func getAuthedMarketplaceClient(ctx context.Context, config marketplaceProvider) (*util.MarketplaceAPIClient, error) {
	payload := map[string]string{
		"domain_name": config.DomainName,
		"username":    config.Username,
		"password":    config.Password,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	marketplaceClient := util.NewMarketplaceAPIClient()

	// Make the HTTP request
	url := fmt.Sprintf("%s/login", marketplaceClient.BaseURL)
	reqHttp, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	reqHttp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resHttp, err := client.Do(reqHttp)
	if err != nil {
		return nil, err
	}
	defer resHttp.Body.Close()

	// Handle response
	if resHttp.StatusCode != http.StatusOK {
		return nil, err
	}

	var response struct {
		Token string `json:"token"`
	}
	if err = json.NewDecoder(resHttp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Token == "" {
		return nil, errors.New("token is missing from the API response")
	}
	marketplaceClient.Token = response.Token
	return &marketplaceClient, nil
}

func (p *marketplaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config marketplaceProvider
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.DomainName == "" || config.Username == "" || config.Password == "" {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"All of 'domain_name', 'username', and 'password' must be provided.",
		)
		return
	}

	marketplaceClient, err := getAuthedMarketplaceClient(ctx, config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Couldn't authenticate",
			fmt.Sprintf("Couldn't get instance of marketplaceClient: %s", err.Error()),
		)
	}

	resp.DataSourceData = marketplaceClient
	resp.ResourceData = marketplaceClient
}

func (p *marketplaceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "otc-marketplace"
}

func (p *marketplaceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasource_whoami.NewWhoamiDataSource,
		datasource_categories.NewCategoryDataSource,
		datasource_clusters.NewClusterDataSource,
		datasource_namespaces.NewNamespaceDataSource,
		datasource_projects.NewProjectDataSource,
		datasource_sales_history.NewSalesHistoryDataSource,
		datasource_products.NewProductDataSource,
		datasource_product_revisions.NewProductRevisionDataSource,
		datasource_applications.NewApplicationDataSource,
		datasource_profile.NewProfileDataSource,
	}
}

func (p *marketplaceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resource_application.NewApplicationResource,
		resource_product.NewProductResource,
		resource_product_revision.NewProductRevisionResource,
	}
}
