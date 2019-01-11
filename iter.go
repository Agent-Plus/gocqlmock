package gocqlmock

import (
	"github.com/gocql/gocql"
)

// Iter represents an expectation with defined rows data to iterate
type Iter struct {
	err         error
	query       string
	expectation expectIface
}

// Close returns error for the expectation or scan
func (i Iter) Close() error {
	if err := i.expectation.getError(); err != nil {
		return err
	}

	return i.err
}

// Scan iterates through the rows
func (i *Iter) Scan(dest ...interface{}) bool {
	if i.err != nil {
		return false
	}

	if i.err = i.expectation.scan(dest...); i.err != nil {
		if i.err == gocql.ErrNotFound {
			i.err = nil
		}

		return false
	}

	return true
}
