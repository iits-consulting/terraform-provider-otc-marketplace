package util

type ProductDataSourceNativeModel struct {
	Id               string            `json:"id,omitempty"`
	CreatedAt        string            `json:"created_at,omitempty"`
	EOLDate          string            `json:"eol_date,omitempty"`
	EOL              bool              `json:"eol,omitempty"`
	LicenseType      string            `json:"license_type,omitempty"`
	Name             string            `json:"name,omitempty"`
	Seller           SellerNativeModel `json:"seller,omitempty"`
	State            string            `json:"state,omitempty"` // TODO - Second time this shows up, might be a good idea to make an enum?
	Weight           int64             `json:"weight,omitempty"`
	Type             string            `json:"type,omitempty"`
	ActiveRevisionId string            `json:"active_revision_id,omitempty"`
	LlmHub           LlmHubNativeModel `json:"llm_hub,omitempty"`
}

type LlmHubNativeModel struct {
	ExternalApi string `json:"external_api,omitempty"`
}

type SellerNativeModel struct {
	Description  string `json:"description,omitempty"`
	Id           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	State        string `json:"state,omitempty"` // TODO - Second time this shows up, might be a good idea to make an enum?
	SupportEmail string `json:"support_email,omitempty"`
	SupportUrl   string `json:"support_url,omitempty"`
}

type MarketplaceAPIClient struct {
	BaseURL string
	Token   string
}
