package gocqlmock

import (
	"fmt"
	"github.com/gocql/gocql"
	"reflect"
	"regexp"
)

// expectIface represents compatibility interface for the
// expectations
type expectIface interface {
	fulfilled() bool
	setError([]interface{})
	getError() error
	scan(...interface{}) error
}

// expectationsIface is an interface to manage
// the collection of expectations
type expectationsIface interface {
	next() expectIface
	push(expectIface)
}

// expect represents common struct to setisfy
// the expectation interface
type expect struct {
	triggered bool
	err       error
	args      interface{}
	rows      *rows
}

// Rows interface manages the collection of the data rows
type Rows interface {
	AddRow(...interface{}) Rows
}

// rows represents struct to keep mock data for the
// Query.Scan and Iter.Scan actions
type rows struct {
	data   []gocql.RowData
	rowpos int
}

// expectations is the collection of the expectation structs
// implements expectationsIface
type expectations []expectIface

// get next unfulfilled expectation
func (s *expectations) next() (e expectIface) {
	for _, e = range *s {
		if !e.fulfilled() {
			return
		}
	}
	return nil
}

// add new expectation to the collection
func (s *expectations) push(e expectIface) {
	*s = append(*s, e)
}

// returns if expectation was fulfilled
func (e expect) fulfilled() bool {
	return e.triggered
}

// scan assignes values from the mock rows collection to the arguments
// passed to original Scan function of the gocal library
func (e *expect) scan(dest ...interface{}) (err error) {
	var r *rows

	if r = e.rows; r == nil || len(r.data) <= r.rowpos {
		return gocql.ErrNotFound
	}

	data := r.data[r.rowpos]
	r.rowpos++

	for i, v := range data.Values {
		err = assignValue(dest[i], v)
	}

	return
}

// setError inserts error to the expection to be returned
// futher
func (e *expect) setError(args []interface{}) {
	if len(args) < 1 {
		return
	}

	switch args[0].(type) {
	case string:
		str := args[0].(string)
		args = args[1:]
		e.err = fmt.Errorf(str, args...)

	case error:
		e.err = args[0].(error)

	}
}

// getError returns error
func (e expect) getError() error {
	return e.err
}

// NewRows returns new mock storage for the rows data
func NewRows() Rows {
	return &rows{}
}

// AddRow pushes values toe the mock storage
func (r *rows) AddRow(values ...interface{}) Rows {
	r.data = append(r.data, gocql.RowData{
		Values: values,
	})
	return r
}

type expectQuery struct {
	expect
	sqlRegex *regexp.Regexp
}

type expectScan struct {
	expect
}

type expectIter struct {
	expect
}

type expectExec struct {
	expect
}

func assignValue(dst, src interface{}) error {
	si := reflect.ValueOf(src)
	di := reflect.ValueOf(dst)

	if k := di.Kind(); k != reflect.Ptr {
		return fmt.Errorf("expected destination argument as pointer, but got %s", k)
	}

	di = reflect.Indirect(di)

	switch si.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.String:
		di.Set(si)
	case reflect.Struct, reflect.Slice, reflect.Array:
		if si.IsValid() && si.Type().AssignableTo(di.Type()) {
			switch src.(type) {
			case []byte:
				di.SetBytes(si.Bytes())
			default:
				di.Set(si)
			}
		} else {
			return fmt.Errorf("can't set destination argument type %s with row data value type %s", di.Kind(), si.Kind())
		}

	case reflect.Ptr:
		if si.IsNil() {
			di.Set(reflect.Zero(di.Type()))
			return nil
		}

		di.Set(reflect.New(di.Type().Elem()))
		return assignValue(di.Interface(), reflect.Indirect(si).Interface())

	default:
		return fmt.Errorf("can't set destination argument type %s with row data value type %s", di.Kind(), si.Kind())
	}

	return nil
}

func argsMatch(qargs, eargs interface{}) error {
	a := reflect.ValueOf(qargs)
	e := reflect.ValueOf(eargs)

	if !e.IsValid() || e.IsNil() {
		return nil
	}

	if al, el := a.Len(), e.Len(); al != el {
		return fmt.Errorf("expected %d query arguments, but got %d", el, al)
	}

	errStr := "argument at %d expected %v, but got %v"
	for i := 0; i < e.Len(); i++ {
		vi := e.Index(i)
		ei := a.Index(i)
		switch vi.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if vi.Int() != ei.Int() {
				return fmt.Errorf(errStr, i, ei.Int(), vi.Int())
			}
		case reflect.Float32, reflect.Float64:
			if vi.Float() != ei.Float() {
				return fmt.Errorf(errStr, i, ei.Float(), vi.Float())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if vi.Uint() != ei.Uint() {
				return fmt.Errorf(errStr, i, ei.Uint(), vi.Uint())
			}
		case reflect.String:
			if vi.String() != ei.String() {
				return fmt.Errorf(errStr, i, ei.String(), vi.String())
			}
		default:
			// compare types like time.Time based on type only
			if vi.Kind() != ei.Kind() {
				return fmt.Errorf(errStr, i, ei.Kind(), vi.Kind())
			}
		}
	}

	return nil
}
