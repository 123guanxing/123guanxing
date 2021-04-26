package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
}

func NewStudent(name string) (Student, error) {
	return Student{Name: name}, nil
}

func main() {
	ctor := as(NewStudent, new(Student))

	data := reflect.ValueOf(ctor).Call([]reflect.Value{reflect.ValueOf("fff")})
	fmt.Println(reflect.TypeOf(data[0].Interface()))
	//d := NewStudent
	//fmt.Println(reflect.TypeOf(d))
	////fmt.Println(d)
	////fmt.Println(d.Name)
	//fmt.Println(reflect.TypeOf(ctor))
}

func as(in interface{}, as interface{}) interface{} {
	outType := reflect.TypeOf(as)
	if outType.Kind() != reflect.Ptr {
		panic("outType is not a pointer")
	}

	if reflect.TypeOf(in).Kind() != reflect.Func {
		ctype := reflect.FuncOf(nil, []reflect.Type{outType.Elem()}, false)

		return reflect.MakeFunc(ctype, func(args []reflect.Value) (results []reflect.Value) {
			out := reflect.New(outType.Elem())
			out.Elem().Set(reflect.ValueOf(in))

			return []reflect.Value{out.Elem()}
		}).Interface()
	}

	inType := reflect.TypeOf(in)

	ins := make([]reflect.Type, inType.NumIn())
	outs := make([]reflect.Type, inType.NumOut())

	for i := range ins {
		ins[i] = inType.In(i)
	}
	outs[0] = outType.Elem()

	fmt.Println(outs)
	for i := range outs[1:] {
		outs[i+1] = inType.Out(i + 1)
	}

	ctype := reflect.FuncOf(ins, outs, false)

	return reflect.MakeFunc(ctype, func(args []reflect.Value) (results []reflect.Value) {
		outs := reflect.ValueOf(in).Call(args)

		out := reflect.New(outType.Elem())
		if outs[0].Type().AssignableTo(outType.Elem()) {
			// Out: Iface = In: *Struct; Out: Iface = In: OtherIface
			out.Elem().Set(outs[0])
		} else {
			// Out: Iface = &(In: Struct)
			t := reflect.New(outs[0].Type())
			t.Elem().Set(outs[0])
			out.Elem().Set(t)
		}
		outs[0] = out.Elem()

		return outs
	}).Interface()
}
