package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func NewMarketplaceAPIClient() MarketplaceAPIClient {
	return MarketplaceAPIClient{
		BaseURL: "https://marketplace.otc.t-systems.com/api/v1/seller",
		Token:   "",
	}
}

// TODO - remove when no longer needed
// This only exists because the backend, when returning PRs, sets the `default_value` attribute to a string
// if the `input_type` is text or selection, but a bool if it's set to switch, which - naturally, causes absolute
// havoc upstream.
func MakePRMarketplaceRequest[T any](ctx context.Context, method string, path string, body io.Reader, marketplaceClient *MarketplaceAPIClient) (*T, error) {
	if body != nil {
		var resultMap map[string]interface{}
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, errors.Join(err, errors.New("couldn't read body"))
		}

		if len(bodyBytes) > 0 {
			err = json.Unmarshal(bodyBytes, &resultMap)
			if err != nil {
				return nil, fmt.Errorf("failed to decode response: %w", err)
			}
		}

		configArray, ok := resultMap["configuration"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("configArray is not an array. configArray: %+v", configArray)
		}

		for i, config := range configArray {
			configMap, ok := config.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("configArray[%d] is not a map[string]interface{} - config: %+v", i, config)
			}

			inputType, ok := configMap["input_type"].(string)
			if !ok {
				return nil, fmt.Errorf("input_type is not a string. input_type: %+v", configMap["input_type"])
			}

			if inputType == "switch" {
				defaultValue, ok := configMap["default_value"]
				if !ok {
					continue // Could *technically* be missing and still pass the openapi spec
				}

				defaultValueBool, ok := defaultValue.(bool) // Unlikely
				if !ok {
					tflog.Warn(ctx, "Forcing conversions of default_value to a bool before sending to backend. See TODO.md")
					defaultValueStrNoQuotes := strings.Replace(strings.Replace(fmt.Sprintf("%v", defaultValue), "\"", "", -1), "'", "", -1)
					defaultValueBool, err = strconv.ParseBool(defaultValueStrNoQuotes)
					if err != nil {
						return nil, fmt.Errorf("couldn't convert defaultValue to a bool. defaultValue: %+v, defaultValueStrNoQuotes: %+v, err: %+v", defaultValue, defaultValueStrNoQuotes, err)
					}
				}

				configMap["default_value"] = defaultValueBool
				configArray[i] = configMap
			}
		}

		resultMap["configuration"] = configArray

		bodyStr, err := json.Marshal(resultMap)
		if err != nil {
			return nil, fmt.Errorf("couldn't marshal resultMap to json. resultMap: %+v, err: %+v", resultMap, err)
		}

		body = bytes.NewReader(bodyStr)
	}

	url := fmt.Sprintf("%s%s", marketplaceClient.BaseURL, path)
	reqHttp, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	reqHttp.Header.Set("Authorization", fmt.Sprintf("Bearer %s", marketplaceClient.Token))
	reqHttp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resHttp, err := client.Do(reqHttp)
	if err != nil {
		return nil, err
	}

	// 2xx to 300
	if !(resHttp.StatusCode >= http.StatusOK && resHttp.StatusCode < http.StatusMultipleChoices) {
		return nil, fmt.Errorf("unexpected status code: %d", resHttp.StatusCode)
	}

	var resultMap map[string]interface{}
	bodyBytes, err := io.ReadAll(resHttp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("method: %s, url: %s, body: %s", method, url, string(bodyBytes)))

	// convert switch bool default_value to string!
	if len(bodyBytes) > 0 {
		err = json.Unmarshal(bodyBytes, &resultMap)
		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	configArray, ok := resultMap["configuration"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("configArray is not an array. configArray: %+v", configArray)
	}

	for i, config := range configArray {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("configArray[%d] is not a map[string]interface{} - config: %+v", i, config)
		}

		inputType, ok := configMap["input_type"].(string)
		if !ok {
			return nil, fmt.Errorf("input_type is not a string. input_type: %+v", configMap["input_type"])
		}

		if inputType == "switch" {
			defaultValue, ok := configMap["default_value"]
			if !ok {
				continue // Could *technically* be missing and still pass the openapi spec
			}

			defaultValueStr, ok := defaultValue.(string) // Unlikely
			if !ok {
				tflog.Warn(ctx, "Forcing conversions of default_value to a string before persisting. See TODO.md")
				defaultValueStr = fmt.Sprintf("%v", defaultValue)
			}

			configMap["default_value"] = defaultValueStr
			configArray[i] = configMap
		}
	}

	resultMap["configuration"] = configArray

	bodyBytes, err = json.Marshal(resultMap)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal resultMap to json. resultMap: %+v, err: %+v", resultMap, err)
	}

	var result T
	if len(bodyBytes) > 0 {
		err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
	} else {
		tflog.Debug(ctx, "skipping body.decode() since len(bodyBytes) is not larger than 0")
	}

	err = resHttp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func MakeMarketplaceRequest[T any](ctx context.Context, method string, path string, body io.Reader, marketplaceClient *MarketplaceAPIClient) (*T, error) {
	url := fmt.Sprintf("%s%s", marketplaceClient.BaseURL, path)
	reqHttp, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	reqHttp.Header.Set("Authorization", fmt.Sprintf("Bearer %s", marketplaceClient.Token))
	reqHttp.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resHttp, err := client.Do(reqHttp)
	if err != nil {
		return nil, err
	}

	// 2xx to 300
	if !(resHttp.StatusCode >= http.StatusOK && resHttp.StatusCode < http.StatusMultipleChoices) {
		return nil, fmt.Errorf("unexpected status code: %d", resHttp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resHttp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("method: %s, url: %s, body: %s", method, url, string(bodyBytes)))

	var result T
	if len(bodyBytes) > 0 {
		err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
	} else {
		tflog.Debug(ctx, "skipping body.decode() since len(bodyBytes) is not larger than 0")
	}

	err = resHttp.Body.Close()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func SanitizeStringValue(in types.String) types.String {
	if in.IsUnknown() || in.IsNull() {
		// Return the input as-is if it's unknown or null
		return in
	}

	// Remove all double quotes from the string value
	sanitized := SanitizeString(in.ValueString())

	// Return a new types.StringValue with the sanitized string
	return types.StringValue(sanitized)
}

func SanitizeString(in string) string {
	return strings.ReplaceAll(in, "\"", "")
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func SlicePtr[T any](s []T) *[]T {
	if len(s) == 0 {
		return nil
	}
	return &s
}

func ApplyIfUnknown(newValue, priorValue attr.Value) attr.Value {
	if newValue.IsUnknown() {
		return priorValue
	}
	return newValue
}

func StringSetOrNull(s string) types.String {
	if s != "" {
		return types.StringValue(s)
	}
	return types.StringNull()
}
func ListSetOrNull(l types.List, elementType attr.Type) types.List {
	if l.IsNull() || l.IsUnknown() || len(l.Elements()) == 0 {
		return types.ListNull(elementType)
	}
	return l
}

// Never return empty lists: https://discuss.hashicorp.com/t/unable-to-handle-empty-list-attribute-in-resource/50952/4
func ListValueOrNull[T any](ctx context.Context, elementType attr.Type, elements []T, diags *diag.Diagnostics) types.List {
	if len(elements) == 0 {
		return types.ListNull(elementType)
	}

	result, d := types.ListValueFrom(ctx, elementType, elements)
	diags.Append(d...)
	return result
}
