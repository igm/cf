package cf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type spaceMetadata struct {
	Guid string `json:"guid"`
}

type spaceEntity struct {
	Name         string       `json:"name"`
	Organization Organization `json:"organization"`
}

type Space struct {
	spaceMetadata `json:"metadata"`
	spaceEntity   `json:"entity"`
}

func (s Space) String() string { return fmt.Sprintf("%s@%s", s.Name, s.Organization.Name) }

func (target *Target) SpacesGet() (spaces []Space, err error) {
	url := fmt.Sprintf("%s/v2/spaces?inline-relations-depth=1", target.TargetUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	resp, err := target.sendRequest(req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(resp.Body)

	var res struct {
		Spaces []Space `json:"resources"`
	}
	err = decoder.Decode(&res)
	spaces = res.Spaces
	return
}
