package booter

import "reflect"

// isNil 通过反射判断是否是nil
func isNil(o any) bool {
	return o == nil || reflect.ValueOf(o).IsNil()
}

// isPtr 通过反射判断是否是指针
func isPtr(o any) bool {
	return reflect.ValueOf(o).Kind() == reflect.Ptr
}

// getPtrValue 通过反射取指针的值
func getPtrValue(o any) any {
	return reflect.ValueOf(o).Elem().Interface()
}

// newPointerInstance，通过反射创建指针类型的实例
// o是非指针类型，struct类型
func newPointerInstance[T any](o T) T {
	return reflect.New(reflect.TypeOf(o).Elem()).Interface().(T)
}

// getPtrStructName
// o is ptr
func getPtrStructName(o any) string {
	return reflect.TypeOf(o).Elem().Name()
}
