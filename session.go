package gocqlmock

import (
	"fmt"
	"testing"
)

// Session represents the mock of the expectations of the
// query, iterators and scan actions
type Session struct {
	expectations expectationsIface
	active       expectIface
	t            *testing.T
}

// Close is an empty wrap
func (m Session) Close() {}

// Query
func (m Session) Query(query string, args ...interface{}) (q *Query) {
	q = &Query{
		expectations: m.expectations,
		args:         args,
		query:        query,
	}

	e := m.expectations.next()
	query = stripQuery(query)
	if e == nil {
		m.printErr("all expectations were already fulfilled, call to query '%s' with args %+v was not expected", query, args)
		return
	}

	eq, ok := e.(*expectQuery)
	if !ok {
		m.printErr("call to query '%s' with args %+v, was not expected, next expectation is %T as %+v", query, args, e, e)
	}

	eq.triggered = true

	if !eq.sqlRegex.MatchString(query) {
		m.printErr("query '%s', does not match regex '%s'", query, eq.sqlRegex.String())
	}

	if err := argsMatch(args, eq.args); err != nil {
		m.printErr(err.Error())
	}

	return
}

// printErr helps to format and write error messages through panic or to the
// testing.T if was defined
func (m *Session) printErr(str string, args ...interface{}) {
	if len(args) == 0 {
		return
	}

	err := fmt.Errorf(str, args...)

	if m.t != nil {
		m.t.Error(err)

		return
	}

	panic(err)
}
