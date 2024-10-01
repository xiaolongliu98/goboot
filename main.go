package main

import (
	"fmt"
	"github.com/xiaolongliu98/goboot/booter"
	"github.com/xiaolongliu98/goboot/boottest/handler"
	"github.com/xiaolongliu98/goboot/boottest/service"
	"reflect"
)

type TestHandler2 struct {
	booter.HandlerComponent
}

func main() {
	booter.Register((*handler.TestHandler)(nil))
	booter.Register((*service.TestService1)(nil))
	booter.Register((*service.TestService2)(nil))

	h := booter.GetInstance((*handler.TestHandler)(nil))
	svc1 := booter.GetInstanceByName[*service.TestService1]("TestService1")
	svc2 := booter.GetInstanceByName[*service.TestService2]("TestService2")

	fmt.Printf("svc1: %p\n", svc1)
	fmt.Printf("svc1: %p\n", svc2)

	fmt.Printf("h.svc1: %p\n", h.TestService1)
	fmt.Printf("h.svc2: %p\n", h.TestService2)

	fmt.Println(svc1)
	fmt.Println(svc2)

	booter.CleanupAll()

	handler2 := &TestHandler2{}
	_, ok := interface{}(handler2).(booter.HandlerComponent)
	fmt.Println(ok)
	ok = isEmbedded(handler2, &booter.HandlerComponent{})
	fmt.Println(ok)
	ok = isEmbedded(handler2, &booter.ServiceComponent{})
	fmt.Println(ok)
}

func isEmbedded(structType interface{}, embeddedType interface{}) bool {
	st := reflect.TypeOf(structType)
	et := reflect.TypeOf(embeddedType)
	if et.Kind() == reflect.Ptr {
		et = et.Elem()
	}
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}

	// 遍历 struct 的所有字段
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)

		// 判断是否是嵌入字段，并且类型匹配
		if field.Anonymous && field.Type == et {
			return true
		}
	}
	return false
}
