package apollohr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type PunchResponse struct {
	Data struct {
		AttendanceHistoryID string      `json:"AttendanceHistoryId"`
		PunchDate           time.Time   `json:"punchDate"`
		PunchType           int         `json:"punchType"`
		LocationName        string      `json:"LocationName"`
		Note                string      `json:"Note"`
		ClientIP            interface{} `json:"clientIp"`
		CanOverride         bool        `json:"CanOverride"`
		OverrideMessage     interface{} `json:"OverrideMessage"`
	} `json:"Data"`
}

type LocationResponse struct {
	Data []struct {
		PunchesLocationID      string  `json:"PunchesLocationId"`
		LocationCode           string  `json:"LocationCode"`
		LocationName           string  `json:"LocationName"`
		Latitude               float64 `json:"Latitude"`
		Longitude              float64 `json:"Longitude"`
		RadiusofEffectiveRange float64 `json:"RadiusofEffectiveRange"`
	} `json:"Data"`
}

func (c *Client) Punch(attendanceType int64, locationName string) (PunchResponse, error) {
	punchResponse := PunchResponse{}
	urlStr := fmt.Sprintf("%s/api/checkin/punch/locate", c.HRMBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return punchResponse, err
	}

	loc, err := c.GetLocation()
	if err != nil {
		return punchResponse, err
	}

	lat := 0.0
	lon := 0.0
	locId := "00000000-0000-0000-0000-000000000000"

	for _, v := range loc.Data {
		if v.LocationName == locationName {
			lat = v.Latitude
			lon = v.Longitude
			locId = v.PunchesLocationID
		}
	}

	v := url.Values{}
	v.Set("AttendanceType", fmt.Sprintf("%d", attendanceType))
	v.Set("Latitude", fmt.Sprintf("%f", lat))
	v.Set("Longitude", fmt.Sprintf("%f", lon))
	v.Set("PunchesLocationId", locId)
	v.Set("IsOverride", "true")
	reqBody := bytes.NewBufferString(v.Encode())

	req, err := http.NewRequest("POST", u.String(), reqBody)
	if err != nil {
		return punchResponse, err
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Add("User-Agent", c.UserAgent)
	req.Header.Add("actioncode", "Default")
	req.Header.Add("functioncode", "APP-LocationCheckin")

	resp, err := c.client.Do(req)
	if err != nil {
		return punchResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return punchResponse, err
		}
		return punchResponse, errors.New(errorResponse.Error.Title)
	}

	err = json.NewDecoder(resp.Body).Decode(&punchResponse)

	return punchResponse, err
}

func (c *Client) GetLocation() (LocationResponse, error) {
	locationResponse := LocationResponse{}
	urlStr := fmt.Sprintf("%s/api/locations/AppEnableList", c.HRMBaseURL)
	u, err := url.Parse(urlStr)
	if err != nil {
		return locationResponse, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return locationResponse, err
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Add("User-Agent", c.UserAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return locationResponse, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return locationResponse, err
		}
		return locationResponse, errors.New(errorResponse.Error.Title)
	}

	err = json.NewDecoder(resp.Body).Decode(&locationResponse)

	return locationResponse, err
}
