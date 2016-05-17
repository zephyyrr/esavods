package models

import (
	"fmt"
)

type Error struct {
	HttpStatus int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Internal   error       `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s | %v | %s", e.HttpStatus, e.Message, e.Data, e.Internal)
}
