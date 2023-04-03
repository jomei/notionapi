package notionapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type AuthenticationService interface {
	CreateToken(ctx context.Context, request *TokenCreateRequest) (*TokenCreateResponse, error)
}

type AuthenticationClient struct {
	apiClient *Client
}

// Create https://developers.notion.com/reference/create-a-token
func (cc *AuthenticationClient) CreateToken(ctx context.Context, request *TokenCreateRequest) (*TokenCreateResponse, error) {
	res, err := cc.apiClient.request(ctx, http.MethodPost, "oauth/token", nil, request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response TokenCreateResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type TokenCreateRequest struct {
	Code            string          `json:"code"`
	GrantType       string          `json:"grant_type"` // Default value of grant_type is always "authorization_code"
	RedirectUri     string          `json:"redirect_uri"`
	ExternalAccount ExternalAccount `json:"external_account,omitempty"`
}

type ExternalAccount struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type TokenCreateResponse struct {
	AccessToken          string `json:"access_token"`
	BotId                string `json:"bot_id"`
	DuplicatedTemplateId string `json:"duplicated_template_id,omitempty"`

	// Owner can be { "workspace": true } OR a User object.
	// Ref: https://developers.notion.com/docs/authorization#step-4-notion-responds-with-an-access_token-and-some-additional-information
	Owner         interface{} `json:"owner,omitempty"`
	WorkspaceIcon string      `json:"workspace_icon"`
	WorkspaceId   string      `json:"workspace_id"`
	WorkspaceName string      `json:"workspace_name"`
}
