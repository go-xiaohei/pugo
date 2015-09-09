package service

import (
	"errors"
	"fmt"
	"gopkg.in/inconshreveable/log15.v2"
	"reflect"
	"runtime"
	"strings"
)

var (
	ErrResultSetNeedPointer    = errors.New("ServiceResult need set a pointer")
	ErrResultSetUnknownPointer = func(rv reflect.Value) error {
		return fmt.Errorf("ServiceResult can't assign to %s", rv.Type().String())
	}
	ErrServiceFuncNeedType = func(fn, v interface{}) error {
		return fmt.Errorf("%s need param %s", funcName(fn), reflect.TypeOf(v).String())
	}
)

type Func func(v interface{}) (*Result, error)

// get function name
func funcName(fn interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	nameData := strings.Split(name, "/")
	if len(nameData) > 2 {
		nameData = nameData[len(nameData)-1:]
	}
	if runtime.GOOS == "windows" {
		return strings.TrimSuffix(strings.Join(nameData, "."), "-fm")
	}
	return strings.TrimSuffix(strings.Join(nameData, "."), "·fm")
}

type Result struct {
	funcName string
	data     map[string]reflect.Value
}

func newResult(fn Func) *Result {
	return &Result{
		funcName: funcName(fn),
		data:     make(map[string]reflect.Value),
	}
}

func (r *Result) Set(values ...interface{}) {
	for _, v := range values {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr {
			panic(ErrResultSetNeedPointer)
		}
		r.data[rv.Type().String()] = rv
	}
}

func (r *Result) SetTo(values ...interface{}) {
	for _, v := range values {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr {
			panic(ErrResultSetNeedPointer)
		}
		resRv, ok := r.data[rv.Type().String()]
		if !ok {
			panic(ErrResultSetUnknownPointer(rv))
		}
		rv.Elem().Set(resRv.Elem())
	}
}

func (r *Result) toDataTypes() []string {
	s := make([]string, 0, len(r.data))
	for name, _ := range r.data {
		s = append(s, name)
	}
	return s
}

func Call(fn Func, v interface{}, values ...interface{}) error {
	res, err := CallResult(fn, v)
	if err != nil {
		return err
	}
	if res != nil {
		res.SetTo(values...)
	}
	return nil
}

func CallResult(fn Func, v interface{}) (*Result, error) {
	var vt string = "nil"
	if v != nil {
		vt = reflect.TypeOf(v).String()
	}
	res, err := fn(v)
	if err != nil {
		log15.Error("Service.Call."+funcName(fn), "in", vt, "error", err)
		return res, err
	}
	var result interface{} = nil
	if res != nil {
		result = res.toDataTypes()
	}
	log15.Debug("Service.Call."+funcName(fn), "in", vt, "out", result)
	return res, err
}
