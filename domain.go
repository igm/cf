package cf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	domainMetadata struct {
		Guid string `json:"guid"`
	}
	domainEntity struct {
		Name   string  `json:"name"`
		Spaces []Space `json:"spaces"`
	}
)

type Domain struct {
	domainMetadata `json:"metadata"`
	domainEntity   `json:"entity"`
}

func (d Domain) String() string { return fmt.Sprint(d.Name) }

// GetDomains returns a slice of registered domains for given space
func (target *Target) DomainsGet() (domains []Domain, err error) {
	url := fmt.Sprintf("%s/v2/domains?inline-relations-depth=1", target.TargetUrl)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)

	var res struct {
		Domains []Domain `json:"resources"`
	}
	err = decoder.Decode(&res)
	domains = res.Domains
	return
}
