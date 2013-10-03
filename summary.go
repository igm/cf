package cf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Summary struct {
	Space
	Apps     []App     `json:"apps"`
	Services []Service `json:"services"`
}

// Summary returns a Space summary
func (target *Target) Summary(spaceGUID string) (summary *Summary, err error) {
	url := fmt.Sprintf("%s/v2/spaces/%s/summary", target.TargetUrl, spaceGUID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := target.sendRequest(req)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&summary)
	return
}
