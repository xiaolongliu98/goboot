package main

import (
	"fmt"
	"github.com/xiaolongliu98/goboot/booter"
	"github.com/xiaolongliu98/goboot/boottest/handler"
	"github.com/xiaolongliu98/goboot/boottest/service"
)

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

}
