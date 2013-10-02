package cf

import (
	"fmt"
)

type Error struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (c *Error) Error() string {
	return fmt.Sprintf("[%d] %s", c.Code, c.Description)
}
