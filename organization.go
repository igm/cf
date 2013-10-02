package cf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	orgMetadata struct {
		Guid string `json:"guid"`
	}
	orgEntity struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}
)

type Organization struct {
	orgMetadata `json:"metadata"`
	orgEntity   `json:"entity"`
}

// GetOrtanizations returns a slice of all organizations
func (target *Target) OrganizationsGet() (orgs []Organization, err error) {
	url := fmt.Sprintf("%s/v2/organizations", target.TargetUrl)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := sendRequest(req, target)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res struct {
		Orgs []Organization `json:"resources"`
	}
	err = decoder.Decode(&res)
	orgs = res.Orgs
	return
}
