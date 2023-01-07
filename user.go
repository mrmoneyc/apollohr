package apollohr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type UserInfo struct {
	Data struct {
		IsVerify                 bool          `json:"isVerify"`
		UserModule               []string      `json:"userModule"`
		UserName                 string        `json:"userName"`
		UserRole                 []interface{} `json:"userRole"`
		IsSupervisor             bool          `json:"IsSupervisor"`
		IsSecretary              bool          `json:"IsSecretary"`
		PersonalPicture          string        `json:"PersonalPicture"`
		EmployeeNumber           string        `json:"EmployeeNumber"`
		GroupID                  string        `json:"GroupId"`
		Gender                   string        `json:"Gender"`
		Language                 string        `json:"Language"`
		CompanyCode              string        `json:"CompanyCode"`
		EmployeeID               string        `json:"EmployeeId"`
		CompanyID                string        `json:"CompanyId"`
		IDType                   int           `json:"IDType"`
		IDTypeNameByUserLanguage string        `json:"IDTypeNameByUserLanguage"`
		EmployeeName             string        `json:"EmployeeName"`
		CompanyName              string        `json:"CompanyName"`
		NickName                 string        `json:"NickName"`
		EnglishName              string        `json:"EnglishName"`
		CompanyTimeZone          string        `json:"CompanyTimeZone"`
	} `json:"Data"`
}

func (c *Client) GetUserInfo() (UserInfo, error) {
	userInfo := UserInfo{}
	urlStr := fmt.Sprintf("%s/api/userInfo", c.LinkupBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return userInfo, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return userInfo, err
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Add("User-Agent", c.UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return userInfo, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return userInfo, err
		}
		return userInfo, errors.New(errorResponse.Error.Title)
	}

	err = json.NewDecoder(resp.Body).Decode(&userInfo)

	return userInfo, err
}
