package gocqlmock

import (
	"github.com/gocql/gocql"
	"testing"
	"time"
)

func TestScanNotFound(t *testing.T) {
	e := &expectScan{}
	var a int
	err := e.scan(&a)

	if err != gocql.ErrNotFound {
		t.Errorf("Expected not found error")
	}
}

func TestScanPrimitives(t *testing.T) {
	e := &expectScan{}
	s := []interface{}{
		int8(1), int16(2), int(3), int64(4),
		float32(2.2), float64(12.2),
		uint8(5), uint16(6), uint(7), uint64(8),
		"foo",
	}

	for _, v := range s {
		var (
			err error

			r = NewRows().AddRow(v)
		)

		e.rows = r.(*rows)

		switch v.(type) {
		case int8:
			var a int8
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case int16:
			var a int16
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case int:
			var a int
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case int64:
			var a int64
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case uint8:
			var a uint8
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case uint16:
			var a uint16
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case uint:
			var a uint
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case uint64:
			var a uint64
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %d", v, a)
			}

		case float32:
			var a float32
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %f", v, a)
			}

		case float64:
			var a float64
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %f", v, a)
			}

		case string:
			var a string
			if err = e.scan(&a); err != nil && err != gocql.ErrNotFound {
				t.Error(err)
			}
			if a != v {
				t.Errorf("expected %d, but got %s", v, a)
			}

		default:
			t.Errorf("unexpected %T with value %v", v, v)
		}
	}
}

func TestScanStruct(t *testing.T) {
	type SomeFoo struct {
		Raw []byte
		Tm  *time.Time
	}

	type SomeBar struct {
		Foo  *SomeFoo
		Name string
	}

	tm := time.Now()
	r := NewRows().
		AddRow("Namw 1", (*SomeBar)(nil)).
		AddRow("Name 2", &SomeFoo{Raw: []byte("Hello"), Tm: &tm})

	e := &expectScan{}
	e.rows = r.(*rows)

	var i int

	for {
		var result SomeBar

		if err := e.scan(&result.Name, &result.Foo); err != nil {
			if err != gocql.ErrNotFound {
				t.Error(err)
			}
			break
		}

		switch i {
		case 0:
			if result.Foo != nil {
				t.Errorf("parameter Foo of the struct SomeBar expected nil, but got %v", result.Foo)
			}
		case 1:
			if result.Foo == nil {
				t.Errorf("parameter Foo of the struct SomeBar expected not nil")
			}
			if result.Foo.Tm != &tm {
				t.Errorf("expected valid Tm parameter in the SomeFoo struct, but got %v", result.Foo.Tm)
			}
		}
		i++
	}

}
