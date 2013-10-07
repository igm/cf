package cf

import (
	"fmt"
)

type Error struct {
	CfCode         int `json:"code"`
	HttpStatusCode int
	Description    string `json:"description"`
	ErrorMsg       string `json:"error"`
}

func (c *Error) Error() string {
	return fmt.Sprintf("Error code: [%d]\n%s\nDescription: %s", c.CfCode, c.ErrorMsg, c.Description)
}
