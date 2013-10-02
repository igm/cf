package cf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Target is the centra structure of this library. All the API remote operatins are performed on Target
type Target struct {
	TargetUrl             string
	AuthorizationEndpoint string
	AccessToken           string
	RefreshToken          string
}

func Login(targetUrl, username, pass string) (*Target, error) {
	infoUrl := fmt.Sprintf("%s/v2/info", targetUrl)
	req, err := http.NewRequest("GET", infoUrl, nil)
	if err != nil {
		return nil, err
	}
	traceReq(req)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	traceResp(resp)

	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(resp.Status)
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var info Info

	err = decoder.Decode(&info)
	if err != nil {
		return nil, err
	}

	target := &Target{
		TargetUrl:             targetUrl,
		AuthorizationEndpoint: info.AuthorizationEndpoint,
	}

	return target, getToken(target, url.Values{
		"grant_type": {"password"},
		"scope":      {},
		"username":   {username},
		"password":   {pass},
	})
}

func refreshToken(target *Target) error {
	return getToken(target, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {target.RefreshToken},
		"scope":         {},
	})
}

func getToken(target *Target, values url.Values) error {
	body := strings.NewReader(values.Encode())
	url := fmt.Sprintf("%s/oauth/token", target.AuthorizationEndpoint)
	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json")
	req.Header.Set("User-Agent", "cf90")
	req.Header.Set("authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:")))

	traceReq(req)
	resp, err := HttpClient.Do(req)
	if err != nil {
		return err
	}
	traceResp(resp)

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	var oauthResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		Jti          string `json:"jti"`
	}
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	decoder.Decode(&oauthResponse)
	if oauthResponse.AccessToken == "" {
		return errors.New("no oauth response received.")
	}
	target.AccessToken = oauthResponse.AccessToken
	target.RefreshToken = oauthResponse.RefreshToken
	return err
}

type Info struct {
	Name                  string `json:"name"`
	Build                 string `json:"build"`
	Support               string `json:"support"`
	Version               int    `json:"version"`
	Description           string `json:"description"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	ApiVersion            string `json:"api_version"`
}

func (target *Target) Info() (info Info, err error) {
	infoUrl := fmt.Sprintf("%s/v2/info", target.TargetUrl)
	req, err := http.NewRequest("GET", infoUrl, nil)
	if err != nil {
		return
	}
	traceReq(req)

	resp, err := HttpClient.Do(req)
	if err != nil {
		return
	}
	traceResp(resp)

	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(resp.Status)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&info)
	return
}
