package apollohr

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	srvLocationGlobal = "HUANU4G3app"
	srvLocationChina  = "PCRU8J4app"
)

type TokenInfo struct {
	AccessToken           string    `json:"access_token"`
	TokenType             string    `json:"token_type"`
	ExpiresIn             int       `json:"expires_in"`
	RefreshToken          string    `json:"refresh_token"`
	UserName              string    `json:"userName"`
	Issued                string    `json:".issued"`
	Expires               string    `json:".expires"`
	JobInfo               string    `json:"JobInfo"`
	SelectCompanyRequired bool      `json:"SelectCompanyRequired"`
	UserStatus            int       `json:"UserStatus"`
	Code                  string    `json:"code"`
	RefreshExpire         time.Time `json:"refresh_expire"`
}

type AccessToken struct {
	IDToken string `json:"id_token"`
	Apphash string `json:"apphash"`
}

func (c *Client) getTokenInfo(companyCode, employeeNo, password string) (TokenInfo, error) {
	now := time.Now().Unix()
	hash := getMagicHash("POST", "/token", now, srvLocationGlobal)
	userName := fmt.Sprintf("%s-%s", companyCode, employeeNo)

	tokenInfo := TokenInfo{}

	urlStr := fmt.Sprintf("%s/token", c.AuthBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return tokenInfo, err
	}

	v := url.Values{}
	v.Set("grant_type", "password")
	v.Set("userName", userName)
	v.Set("password", password)
	reqBody := bytes.NewBufferString(v.Encode())

	req, err := http.NewRequest("POST", u.String(), reqBody)
	if err != nil {
		return tokenInfo, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("User-Agent", c.UserAgent)

	q := req.URL.Query()
	q.Add("time", strconv.FormatInt(now, 10))
	q.Add("hash", hash)
	q.Add("_sd", "HRM")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return tokenInfo, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return tokenInfo, err
		}
		return tokenInfo, errors.New(errorResponse.Error.Detail)
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)

	return tokenInfo, err
}

func (c *Client) getAccessToken(code string) (AccessToken, error) {
	accessToken := AccessToken{}

	urlStr := fmt.Sprintf("%s/api/auth/checkticket", c.LinkupBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return accessToken, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return accessToken, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("response_type", "id_token")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return accessToken, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return accessToken, err
		}
		return accessToken, errors.New(errorResponse.Error.Detail)
	}

	err = json.NewDecoder(resp.Body).Decode(&accessToken)

	return accessToken, err
}

func getMagicHash(method string, path string, epoch int64, srvLoc string) string {
	s := fmt.Sprintf("%s%s%d%s", method, path, epoch, srvLoc)
	sum := sha256.Sum256([]byte(s))

	return fmt.Sprintf("%x", sum)
}
