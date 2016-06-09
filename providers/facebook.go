package providers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type LinkedInProvider struct {
	*ProviderData
}

func NewFacebookInProvider(p *ProviderData) *FacebookInProvider {
	p.ProviderName = "Facebook"
	if p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{Scheme: "https",
			Host: "www.facebook.com",
			Path: "/v2.5/dialog/oauth"}
	}
	if p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{Scheme: "https",
			Host: "www.facebook.com",
			Path: "/oauth/access_token"}
	}
	if p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{Scheme: "https",
			Host: "www.facebook.com",
			Path: "/me"}
	}
	if p.ValidateURL.String() == "" {
		p.ValidateURL = p.ProfileURL
	}
	if p.Scope == "" {
		p.Scope = "public_profile email"
	}
	return &FacebookInProvider{ProviderData: p}
}

func getFacebookHeader(access_token string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	header.Set("x-li-format", "json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return header
}

func (p *FacebookProvider) GetEmailAddress(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String()+"?format=json", nil)
	if err != nil {
		return "", err
	}
	req.Header = getLinkedInHeader(s.AccessToken)

	json, err := api.Request(req)
	if err != nil {
		return "", err
	}

	email, err := json.String()
	if err != nil {
		return "", err
	}
	return email, nil
}

func (p *FacebookProvider) ValidateSessionState(s *SessionState) bool {
	return validateToken(p, s.AccessToken, getFacebookHeader(s.AccessToken))
}
