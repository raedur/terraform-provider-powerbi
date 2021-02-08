package powerbiapi

import (
	"fmt"
	"net/url"
)

// GetGatewayDatasourcesInGroupResponse represents the response from get gateway datasources
type GetGatewayDatasourcesInGroupResponse struct {
	Value []GetGatewayDatasourcesInGroupResponseItem
}

// GetDatasourcesInGroupResponseItem represents a single datasource
type GetGatewayDatasourcesInGroupResponseItem struct {
	ID                string
	GatewayID         string
	DatasourceType    string
	ConnectionDetails string
}

// UpdateBasicAuthCredentialsForDatasourceRequest represents the request to update the basic auth credentials on a datasources
type UpdateBasicAuthCredentialsForDatasourceRequest struct {
	CredentialDetails UpdateBasicAuthCredentialsForDatasourceRequestCredentialDetails `json:"credentialDetails"`
}

type UpdateBasicAuthCredentialsForDatasourceRequestCredentialDetails struct {
	CredentialType              string `json:"credentialType"`
	Credentials                 string `json:"credentials"`
	EncryptedConnection         string `json:"encryptedConnection"`
	EncryptionAlgorithm         string `json:"encryptionAlgorithm"`
	PrivacyLevel                string `json:"privacyLevel"`
	UseEndUserOAuth2Credentials string `json:"useEndUserOAuth2Credentials"`
}

// GetDatasourcesInGroup gets datasources in a dataset that exists within a group.
func (client *Client) GetGatewayDatasourcesInGroup(groupID string, datasetID string) (*GetGatewayDatasourcesInGroupResponse, error) {

	var respObj GetGatewayDatasourcesInGroupResponse
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/groups/%s/datasets/%s/Default.GetBoundGatewayDatasources", url.PathEscape(groupID), url.PathEscape(datasetID))
	err := client.doJSON("GET", url, nil, &respObj)

	return &respObj, err
}

func (client *Client) UpdateBasicAuthCredentialsForDatasource(datasourceID string, gatewayID string, request UpdateBasicAuthCredentialsForDatasourceRequest) error {
	url := fmt.Sprintf("https://api.powerbi.com/v1.0/myorg/gateways/%s/datasources/%s", url.PathEscape(gatewayID), url.PathEscape(datasourceID))
	err := client.doJSON("PATCH", url, &request, nil)

	return err
}
