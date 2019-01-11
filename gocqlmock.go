package gocqlmock

import (
	"regexp"
	"strings"
	"testing"
)

var stripRe *regexp.Regexp

func init() {
	stripRe = regexp.MustCompile("\\s+")
}

// New creates empty Session mock to be filled with expection scanarios
//
// func TestSomeCassandraQuery(t *testing.T) {
//     testSession := gocqlmock.New(t)
//
//     // expecting single query
//     testSession.ExpectQuery("SELECT .+ FROM users").
//         ExpectScan().
//         WithResult(
//             gocqlmock.NewRows().
//                 AddRow(234)
//         )
//
//     router := Server(testSession)
//     w := httptest.NewRecorder()
//     req, _ := http.NewRequest("GET", "/user/234", nil)
//     router.ServeHTTP(w, req)
// }
//
// func TestSomeCassandraQueryWithListResult(t *testing.T) {
//     testSession := gocqlmock.New(t)
//
//     // expecting single query
//     testSession.ExpectQuery("SELECT .+ FROM groups").
//         ExpectIter().
//         WithResult(
//             gocqlmock.NewRows().
//                 AddRow(2, "Group 2", 1.2).
//                 AddRow(5, "Group 5", 4.6)
//         )
//
//     router := Server(testSession)
//     w := httptest.NewRecorder()
//     req, _ := http.NewRequest("GET", "/groups", nil)
//     router.ServeHTTP(w, req)
// }
//
func New(t *testing.T) *Session {
	return &Session{
		expectations: new(expectations),
		t:            t,
	}
}

// ExpectQuery Query(...) to be triggered, which will match
// the given query string as a regular expression
func (m *Session) ExpectQuery(queryRegex string) *Session {
	e := &expectQuery{}
	e.sqlRegex = regexp.MustCompile(queryRegex)

	m.expectations.push(e)
	m.active = e
	return m
}

// ExpectScan Query(...).Scan(...) to be triggered, which will assign
// values from the rows mock to the arguments passed through Scan(...)
func (m *Session) ExpectScan() *Session {
	_, ok := m.active.(*expectQuery)
	if !ok {
		m.printErr("scan may be expected only with query based expectations, current is %T", m.active)
		return m
	}

	e := &expectScan{}
	m.expectations.push(e)
	m.active = e
	return m
}

// ExpectIter Query(...).Iter().Scan(...) to be triggered, which will assign
// values from the rows mock to the arguments passed through Scan(...)
func (m *Session) ExpectIter() *Session {
	_, ok := m.active.(*expectQuery)
	if !ok {
		m.printErr("iter may be expected only with query based expectations, current is %T", m.active)
		return m
	}

	e := &expectIter{}
	m.expectations.push(e)
	m.active = e
	return m
}

// ExpectExec Query(...).Exec() to be triggered
func (m *Session) ExpectExec() *Session {
	_, ok := m.active.(*expectQuery)
	if !ok {
		m.printErr("exec may be expected only with query based expectations, current is %T", m.active)
		return m
	}

	e := &expectExec{}
	m.expectations.push(e)
	m.active = e
	return m
}

// WithResult assignes rows mock to the expectation wich was created through
// ExpectScan or ExpectIter
func (m *Session) WithResult(r Rows) *Session {
	es, ok := m.active.(*expectScan)
	if !ok {
		ei, ok := m.active.(*expectIter)

		if !ok {
			m.printErr("rows may be returned only by scan or iter expectations, current is %T", m.active)
		}

		ei.rows = r.(*rows)
	} else {
		es.rows = r.(*rows)
	}

	return m
}

// WithArgs expectation should be called with given arguments.
// Works with Query expectations
//
// testSession.ExpectQuery("SELECT .+ FROM tb WHERE id").
//     WithArgs(1).
//     ExpectExec()
func (m *Session) WithArgs(args ...interface{}) *Session {
	eq, ok := m.active.(*expectQuery)
	if !ok {
		ex, ok := m.active.(*expectExec)
		if !ok {
			m.printErr("arguments may be expected only with query based expectations, current is %T", m.active)
		}
		ex.args = args
	} else {
		eq.args = args
	}
	return m
}

// WithError expectation will return error
func (m *Session) WithError(args ...interface{}) {
	m.active.setError(args)
}

// strip out new lines and trim spaces
func stripQuery(q string) (s string) {
	return strings.TrimSpace(stripRe.ReplaceAllString(q, " "))
}
