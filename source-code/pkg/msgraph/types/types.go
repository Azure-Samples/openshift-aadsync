package types

// AuthError holds MS Graph auth errors
type AuthError struct {
	Type        string `json:"error"`
	Description string `json:"error_description"`
}

// AuthToken holds the MS Graph access token details
type AuthToken struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

// Group holds the Azure AD Group details
type Group struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Users       []User `json:"value"`
}

// User holds the Azure AD User details
type User struct {
	OdataType         string `json:"@odata.type"`
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	GivenName         string `json:"givenName"`
	Surname           string `json:"surname"`
	UserPrincipalName string `json:"userPrincipalName"`
}
