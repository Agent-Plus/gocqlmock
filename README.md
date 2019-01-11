
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

## Thanks

Thanks to [DATA-DOG](https://github.com/DATA-DOG) for the [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) taken as idea for this 
