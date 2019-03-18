package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	msgraph "pkg/msgraph/types"

	logrus "github.com/sirupsen/logrus"
)

// Client contains internal MS Graph Client details
type Client struct {
	Log       *logrus.Entry
	Config    *Config
	Client    *http.Client
	AuthToken *msgraph.AuthToken
}

// NewClient creates a new MS Graph Client with default configuration
func NewClient(log *logrus.Entry) *Client {

	config, err := NewConfigFromEnvironmentVariables()
	if err != nil {
		log.Fatal(err)
	}

	return NewClientForConfig(config, log)
}

// NewClientForConfig creates a new MS Graph Client with the specified configuration
func NewClientForConfig(config *Config, log *logrus.Entry) *Client {

	client := &Client{
		Log:    log,
		Config: config,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	log.Info("Created msgraph client")

	authToken, err := client.GetAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	client.AuthToken = authToken

	return client
}

// GetAccessToken leverages the Auth API to return an AuthToken that provides authenticated access to the MS Graph API
func (c *Client) GetAccessToken() (*msgraph.AuthToken, error) {

	// Obtain Access Token from Auth API
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", c.Config.TenantID)
	c.Log.Infof("Fetching access token: %s", tokenURL)
	tokenResponse, err := c.Client.PostForm(tokenURL, url.Values{
		"client_id":     {c.Config.ClientID},
		"scope":         {"https://graph.microsoft.com/.default"},
		"client_secret": {c.Config.ClientSecret},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		return nil, err
	}
	defer tokenResponse.Body.Close()

	if tokenResponse.StatusCode != 200 {
		var error msgraph.AuthError
		json.NewDecoder(tokenResponse.Body).Decode(&error)
		err := fmt.Errorf("StatusCode: %d, Error: %s, Description:%s", tokenResponse.StatusCode, error.Type, error.Description)
		return nil, err
	}

	var authToken *msgraph.AuthToken
	if err := json.NewDecoder(tokenResponse.Body).Decode(&authToken); err != nil {
		return nil, err
	}

	c.Log.Infof("Successfully obtained access token")
	c.Log.Debugf("TokenType: %s, ExpiresIn: %d, AccessToken: %s", authToken.TokenType, authToken.ExpiresIn, authToken.AccessToken)

	return authToken, nil
}

// GetGroup leverages the MS Graph API to return the AzureAD group specified by the group id
func (c *Client) GetGroup(groupID string) (*msgraph.Group, error) {

	// Obtain Group from MS Graph API
	groupURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s", groupID)
	c.Log.Infof("Fetching group: %s", groupURL)
	groupRequest, _ := http.NewRequest("GET", groupURL, nil)
	groupRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken.AccessToken))
	groupResponse, err := c.Client.Do(groupRequest)
	if err != nil {
		return nil, err
	}
	defer groupResponse.Body.Close()

	if groupResponse.StatusCode != 200 {
		var error msgraph.AuthError
		json.NewDecoder(groupResponse.Body).Decode(&error)
		err := fmt.Errorf("StatusCode: %d, Error: %s, Description:%s", groupResponse.StatusCode, error.Type, error.Description)
		return nil, err
	}

	group := &msgraph.Group{}
	if err := json.NewDecoder(groupResponse.Body).Decode(&group); err != nil {
		return nil, err
	}
	c.Log.Info("Successfully obtained group")

	// Obtain Members of Group
	groupMemberListURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/members", groupID)
	c.Log.Infof("Fetching group members: %s", groupMemberListURL)
	groupMemberListRequest, _ := http.NewRequest("GET", groupMemberListURL, nil)
	groupMemberListRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken.AccessToken))
	groupMemberListResponse, err := c.Client.Do(groupMemberListRequest)
	if err != nil {
		return nil, err
	}
	defer groupMemberListResponse.Body.Close()

	if groupMemberListResponse.StatusCode != 200 {
		var error msgraph.AuthError
		json.NewDecoder(groupMemberListResponse.Body).Decode(&error)
		err := fmt.Errorf("StatusCode: %d, Error: %s, Description:%s", groupMemberListResponse.StatusCode, error.Type, error.Description)
		return nil, err
	}

	if err := json.NewDecoder(groupMemberListResponse.Body).Decode(&group); err != nil {
		return nil, err
	}

	c.Log.Info("Successfully obtained group members")
	c.Log.Debugf("ID: %s, DisplayName: %s, Description: %s, UserCount: %d", group.ID, group.DisplayName, group.Description, len(group.Users))
	for _, user := range group.Users {
		c.Log.Debugf("Type: %s, ID: %s, DisplayName: %s, GivenName: %s, Surname: %s, UserPrincipalName: %s", user.OdataType, user.ID, user.DisplayName, user.GivenName, user.Surname, user.UserPrincipalName)
	}

	return group, nil
}
