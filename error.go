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
	errorTmpl := `
CF code          : [%d]
Http Status Code : [%d]
Description      : %s 
Error msg        : %s
`
	return fmt.Sprintf(errorTmpl, c.CfCode, c.HttpStatusCode, c.Description, c.ErrorMsg)
}
