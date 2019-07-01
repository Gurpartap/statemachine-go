package dynafunc

import (
	"errors"
	"fmt"
	"reflect"
)

// Usage:
//
// args := make(map[reflect.Type]interface{})
// args[reflect.TypeOf(new(error))] = errors.New("some error occurred")
// args[reflect.TypeOf(new(Transition))] = newBasicTransition("", "", NewMachine(), nil)
//
// if err := dynafunc.NewDynamicFunc(callbackFn, args).Call(); err != nil {
// 	panic(err.Error())
// }
type DynamicFunc struct {
	fn  interface{}
	in  map[reflect.Type]interface{}
	Out []reflect.Value
}

func NewDynamicFunc(fn interface{}, in map[reflect.Type]interface{}) *DynamicFunc {
	return &DynamicFunc{
		fn:  fn,
		in:  in,
		Out: nil,
	}
}

func (f *DynamicFunc) Call() error {
	fnType := reflect.TypeOf(f.fn)
	fnValue := reflect.ValueOf(f.fn)

	if fnType.Kind() != reflect.Func {
		return errors.New("argument must be of the type func")
	}

	in := make([]reflect.Value, fnType.NumIn())

	for i := 0; i < fnType.NumIn(); i++ {
		t := fnType.In(i)

		object, ok := f.in[reflect.PtrTo(t)]
		if !ok {
			return fmt.Errorf("unexpected argument with type '%s'", t)
		}

		// fmt.Println(t, "=>", object)

		in[i] = reflect.ValueOf(object)
	}

	f.Out = fnValue.Call(in)

	return nil
}
