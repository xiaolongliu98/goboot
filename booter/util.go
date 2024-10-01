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

func foreachSubComponent(component Component, f func(fieldComponent Component, filedIndex int, autowiredName string)) {
	// 递归初始化instance依赖的Component（tag标记autowired:""）
	// 遍历instance的字段，如果字段是Component类型，且tag标记autowired:""，则初始化
	instanceValue := reflect.ValueOf(component).Elem()
	instanceType := reflect.TypeOf(component).Elem()
	for i := 0; i < instanceValue.NumField(); i++ {
		var (
			autowiredName  string
			ok             bool
			fieldComponent Component
		)
		// 获取属性的tag
		if autowiredName, ok = instanceType.Field(i).Tag.Lookup("autowired"); !ok {
			continue
		}
		filed := instanceValue.Field(i)
		// 判断属性类型是否是指针或接口
		if filed.Kind() != reflect.Ptr && filed.Kind() != reflect.Interface {
			continue
		}
		// 判断属性类型是否是Component
		if fieldComponent, ok = filed.Interface().(Component); !ok {
			if autowiredName == "" {
				continue
			}
			// 设置了autoWiredName，说明可能是接口类型，需要通过接口名获取Component
		}
		f(fieldComponent, i, autowiredName)
	}
}
