package models

import "fmt"

type RadioParamError struct {
	Field string
	Msg   string
}

func (e RadioParamError) Error() string {
	return fmt.Sprintf("in field %s: %s", e.Field, e.Msg)
}

type DatabaseError struct {
	Uri string
	Err error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("when open %s: %s", e.Uri, e.Err.Error())
}
