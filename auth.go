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

type LoginInfo struct {
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

func (c *Client) login(companyCode, employeeNo, password string) (LoginInfo, error) {
	now := time.Now().Unix()
	hash := getMagicHash("POST", "/token", now, srvLocationGlobal)
	userName := fmt.Sprintf("%s-%s", companyCode, employeeNo)

	loginInfo := LoginInfo{}

	urlStr := fmt.Sprintf("%s/token", c.AuthBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return loginInfo, err
	}

	v := url.Values{}
	v.Set("grant_type", "password")
	v.Set("userName", userName)
	v.Set("password", password)
	reqBody := bytes.NewBufferString(v.Encode())

	req, err := http.NewRequest("POST", u.String(), reqBody)
	if err != nil {
		return loginInfo, err
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
		return loginInfo, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return loginInfo, err
		}
		return loginInfo, errors.New(errorResponse.Error.Title)
	}

	err = json.NewDecoder(resp.Body).Decode(&loginInfo)

	return loginInfo, err
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
		return accessToken, errors.New(errorResponse.Error.Title)
	}

	err = json.NewDecoder(resp.Body).Decode(&accessToken)

	return accessToken, err
}

func (c *Client) RefreshToken() error {
	accessToken := AccessToken{}

	urlStr := fmt.Sprintf("%s/api/auth/refreshtoken", c.LinkupBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Add("User-Agent", c.UserAgent)

	q := req.URL.Query()
	q.Add("response_type", "id_token")
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return err
		}
		return errors.New(errorResponse.Error.Title)
	}

	err = json.NewDecoder(resp.Body).Decode(&accessToken)
	if err != nil {
		return err
	}

	c.token = accessToken.IDToken

	return nil
}

func (c *Client) GetToken() string {
	return c.token
}

func getMagicHash(method string, path string, epoch int64, srvLoc string) string {
	s := fmt.Sprintf("%s%s%d%s", method, path, epoch, srvLoc)
	sum := sha256.Sum256([]byte(s))

	return fmt.Sprintf("%x", sum)
}
