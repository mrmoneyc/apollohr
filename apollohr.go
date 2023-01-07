package apollohr

import (
	"net/http"
	"net/url"
)

const (
	defaultHRMBaseURL    = "https://pt-be.mayohr.com"
	defaultLinkupBaseURL = "https://linkup-be.mayohr.com"
	defaultAuthBaseURL   = "https://asiaauth.mayohr.com"
	defaultUserAgent     = "APOLLO/1.1.43 (com.mayohr.tube; build:1; iOS 16.1.1) Alamofire/1.1.43"
)

type Client struct {
	client        *http.Client
	token         string
	UserAgent     string
	HRMBaseURL    *url.URL
	LinkupBaseURL *url.URL
	AuthBaseURL   *url.URL
}

type ErrorResponse struct {
	Error struct {
		Status       string      `json:"Status"`
		Title        string      `json:"Title"`
		Detail       string      `json:"Detail"`
		MulitiDetail interface{} `json:"MulitiDetail"`
	} `json:"Error"`
}

func NewClient(companyCode string, employeeNo string, password string) *Client {
	hrmBaseURL, _ := url.Parse(defaultHRMBaseURL)
	linkupBaseURL, _ := url.Parse(defaultLinkupBaseURL)
	authBaseURL, _ := url.Parse(defaultAuthBaseURL)

	c := &Client{
		client:        &http.Client{},
		UserAgent:     defaultUserAgent,
		HRMBaseURL:    hrmBaseURL,
		LinkupBaseURL: linkupBaseURL,
		AuthBaseURL:   authBaseURL,
	}

	tokenInfo, err := c.getTokenInfo(companyCode, employeeNo, password)
	if err != nil {
		panic(err)
	}

	accessToken, err := c.getAccessToken(tokenInfo.Code)
	if err != nil {
		panic(err)
	}

	c.token = accessToken.IDToken

	return c
}
