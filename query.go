package gocqlmock

import (
	"fmt"
)

// Query represents query mock
type Query struct {
	expectations expectationsIface
	err          error
	args         []interface{}
	query        string
}

// Bind
func (q *Query) Bind(args ...interface{}) *Query {
	q.args = args
	return q
}

// Exec
func (m Query) Exec() error {
	e := m.expectations.next()
	if e == nil {
		return fmt.Errorf("all expectations were already fulfilled, call to exec of the query '%s' with args %+v was not expected", m.query, m.args)
	}

	ex, ok := e.(*expectExec)
	if !ok {
		return fmt.Errorf("call to exec of the query '%s' with args %+v, was not expected, next expectation is %T as %+v", m.query, m.args, e, e)
	}

	ex.triggered = true

	if ex.err != nil {
		return ex.err
	}

	if err := argsMatch(m.args, ex.args); err != nil {
		return fmt.Errorf("expection %T with args %+v, got args %+v: %s", e, ex.args, m.args, err.Error())
	}

	return nil
}

// Iter
func (m Query) Iter() *Iter {
	e := m.expectations.next()
	if e == nil {

		return &Iter{
			err:         fmt.Errorf("all expectations were already fulfilled, call to iter of the query '%s' with args %+v was not expected", m.query, m.args),
			expectation: &expectIter{},
			query:       m.query,
		}
	}

	eq, ok := e.(*expectIter)
	if !ok {
		return &Iter{
			err:         fmt.Errorf("call to iter of the query '%s' with args %+v, was not expected, next expectation is %T as %+v", m.query, m.args, e, e),
			expectation: &expectIter{},
			query:       m.query,
		}
	}

	eq.triggered = true

	return &Iter{
		expectation: e,
		query:       m.query,
	}
}

// Scan
func (m Query) Scan(dest ...interface{}) error {
	e := m.expectations.next()
	if e == nil {
		return fmt.Errorf("all expectations were already fulfilled, call to query '%s' scan with args %+v was not expected", m.query, dest)
	}

	eq, ok := e.(*expectScan)
	if !ok {
		return fmt.Errorf("call to scan query '%s' with args %#v, was not expected, next expectation is %T", m.query, dest, eq)
	}

	eq.triggered = true

	if eq.err != nil {
		return eq.err
	}

	return e.scan(dest...)
}
