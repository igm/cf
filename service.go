package cf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type serviceMetadata struct {
	Guid string `json:"guid"`
}

type serviceEntity struct {
	Label       string `json:"label"`
	Provider    string `json:"provider"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Active      bool   `json:"active"`
	Bindable    bool   `json:"bindable"`
}

type Service struct {
	serviceMetadata `json:"metadata"`
	serviceEntity   `json:"entity"`
}

func (target *Target) GetServices() (services []Service, err error) {
	url := fmt.Sprintf("%s/v2/services", target.TargetUrl)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var response struct {
		Services []Service `json:"resources"`
	}
	err = decoder.Decode(&response)
	services = response.Services
	return
}
