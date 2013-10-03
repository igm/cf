package cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type routeMetadata struct {
	Guid string `json:"guid"`
}

type routeEntity struct {
	Host   string `json:"host"`
	Domain Domain `json:"domain"`
	Space  Space  `json:"space"`
	Apps   []App  `json:"apps"`
}

type Route struct {
	routeMetadata `json:"metadata"`
	routeEntity   `json:"entity"`
}

func (r Route) String() string { return fmt.Sprintf("%s.%s", r.Host, r.Domain.Name) }

func (target *Target) RoutesGet() (routes []Route, err error) {
	url := fmt.Sprintf("%s/v2/routes?inline-relations-depth=1", target.TargetUrl)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var response struct {
		Routes []Route `json:"resources"`
	}
	err = decoder.Decode(&response)
	routes = response.Routes
	return
}

func (target *Target) RouteCreate(host, domainGUID, spaceGUID string) (err error) {
	url := fmt.Sprintf("%s/v2/routes", target.TargetUrl)

	body, err := json.Marshal(struct {
		Host       string `json:"host"`
		DomainGUID string `json:"domain_guid"`
		SpaceGUID  string `json:"space_guid"`
	}{host, domainGUID, spaceGUID})
	if err != nil {
		return
	}
	req, _ := http.NewRequest("POST", url, closeReader{bytes.NewReader(body)})

	_, err = target.sendRequest(req)
	return
}

func (target *Target) RouteDelete(routeGUID string) (err error) {
	url := fmt.Sprintf("%s/v2/routes/%s", target.TargetUrl, routeGUID)
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err = target.sendRequest(req)
	return
}
