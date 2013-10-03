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
	Name string `json:"name"`
}

type Space struct {
	spaceMetadata `json:"metadata"`
	spaceEntity   `json:"entity"`
}

func (target *Target) SpacesGet(orgGUID string) (spaces []Space, err error) {
	url := fmt.Sprintf("%s/v2/organizations/%s/spaces", target.TargetUrl, orgGUID)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := target.sendRequest(req)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)

	var res struct {
		Spaces []Space `json:"resources"`
	}
	err = decoder.Decode(&res)
	spaces = res.Spaces
	return
}
