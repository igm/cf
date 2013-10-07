package cf

import (
	"fmt"
)

type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	ErrorMsg    string `json:"error"`
}

func (c *Error) Error() string {
	return fmt.Sprintf("Error code: [%d]\n%s\nDescription: %s", c.Code, c.ErrorMsg, c.Description)
}
