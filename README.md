
# Mock gocql

gocqlmock was designed to be useful with gocql library to mock common commands to cassandra database

## Before

To mock gocql easily Your project should use gocql library through the interface.
For example how to use base functionality: Query, Scan, Iter, Exec.

``` go

import (
	"github.com/gocql/gocql"
)

// SessionInterface should describe commonly used functions of the
// gocql.Session
type SessionInterface interface {
	Close()
	Query(string, ...interface{}) QueryInterface
    // Put more You need
}

// QueryInterface should describe commonly used functions of the
// gocql.Query
type QueryInterface interface {
	Exec() error
	Iter() IterInterface
	Scan(...interface{}) error
    // Put more You need
}

// IterInterface should describe commonly used functions of the
// gocql.Iter
type IterInterface interface {
	Close() error
	Scan(...interface{}) bool
    // Put more You need
}

// Sessions is a wrapper for a docql.Session for mockability
type Session struct {
	session *gocql.Session
}

// Query is a wrapper for a gocql.Query for mockability.
type Query struct {
	query *gocql.Query
}

// Iter is a wrapper for an gocql.Iter for mockability.
type Iter struct {
	iter *gocql.Iter
}

// NewSession instantiates a new Session
func NewSession(session *gocql.Session) SessionInterface {
	return &Session{session: session}
}

// NewQuery instantiates a new Query
func NewQuery(query *gocql.Query) QueryInterface {
	return &Query{query}
}

// NewIter instantiates a new Iter
func NewIter(iter *gocql.Iter) IterInterface {
	return &Iter{iter}
}

// Close wraps the session's close method
func (s *Session) Close() {
	s.session.Close()
}

// Exec wraps the query's Exec method
func (q *Query) Exec() error {
	return q.query.Exec()
}

// Iter wraps the query's Iter method
func (q *Query) Iter() IterInterface {
	return NewIter(q.query.Iter())
}

// Scan wraps the query's Scan method
func (q *Query) Scan(dest ...interface{}) error {
	return q.query.Scan(dest...)
}

// Close is a wrapper for the iter's Close method
func (i *Iter) Close() error {
	return i.iter.Close()
}

// Scan is a wrapper for the iter's Scan method
func (i *Iter) Scan(dest ...interface{}) bool {
	return i.iter.Scan(dest...)
}
```

Initialize gocql.Session

``` go
func NewConnection(hosts []string, keyspace string) (SessionInterface, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Consistency = gocql.One

	if keyspace != "" {
		cluster.Keyspace = keyspace
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return NewSession(session), err
}
```

## Install

```
go get github.com/Agent-Plus/gocqlmock
```

## Use

Once session interface was used in the project let's replace gocql with mock in the test.
Create fake structs which will implement gocqlmock Session, Query, Iter.

``` go
import (
	"github.com/Agent-Plus/gocqlmock"
	"testing"
)

type MockSession struct {
	*gocqlmock.Session
}

type MockQuery struct {
	*gocqlmock.Query
}

type MockIter struct {
	*gocqlmock.Iter
}

func (s MockSession) Query(query string, args ...interface{}) dbapi.QueryInterface {
	return &MockQuery{s.Session.Query(query, args...)}
}

func (q MockQuery) Iter() dbapi.IterInterface {
	return &MockIter{q.Query.Iter()}
}

func (q MockQuery) Exec() error {
	return q.Query.Exec()
}
```

Test some handler

``` go
func TestGetUsers(t *testing.T) {
	testSession := &MockSession{gocqlmock.New(t)}
	testSession.ExpectQuery("SELECT.+FROM users").
		ExpectIter().
		WithResult(
			gocqlmock.NewRows().
				AddRow(int64(1)),
		)

	router := Server(testSession)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
        t.Error("expecgted response code 200")
    }
}
```

### Expectations

|                | Description                                              |
| :------------- | :------------------------------------------------------- |
| ExpectQuery    | Expected Query(...) call, will check query or arguments passed to the function Query(...)
| ExpectScan     | Expected Query(...).Scan(...) call, this expectation will deliver fake row to the arguments passed through Scan(...) |
| ExpectIter     | Expected Query(...).Iter().Scan(...) call, this expectation will deliver fake rows to the arguments passed through Scan(...) |
| ExpectExec     | Expected Query(...).Exec() call |

### Row data

To create fake rows use WithResult, which follows with ExpectScan and ExpectQuery

``` go
	testSession.ExpectQuery("SELECT").
		ExpectIter().
		WithResult(
			gocqlmock.NewRows().
				AddRow(int64(1), "Foo", nil, 0, "8100000000", "foo@localhost").
				AddRow(int64(5), "Foo5", &AntStruct{Text: "F", BgColor: "#fff"}, 1, "8100000000", "foo5@localhost"),
		)
```

``` go
	testSession.ExpectQuery("SELECT.+WHERE.+id").
		ExpectScan().
		WithResult(
			gocqlmock.NewRows().
				AddRow(1),
		)
```

## Thanks

Thanks to [DATA-DOG](https://github.com/DATA-DOG) for the [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) taken as idea for this 
