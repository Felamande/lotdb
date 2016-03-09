package query

import "fmt"

type DatabaseError struct {
	Uri string
	Err error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("when open %s: %s", e.Uri, e.Err.Error())
}
