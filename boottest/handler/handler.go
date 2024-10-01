package handler

import (
	"fmt"
	"github.com/xiaolongliu98/goboot/booter"
	"github.com/xiaolongliu98/goboot/boottest/service"
)

type TestHandler struct {
	Msg string

	TestService1 service.StringGetter  `autowired:"TestService1"`
	TestService2 *service.TestService2 `autowired:""`
}

func (t *TestHandler) Initialize(ctx *booter.BootContext) error {
	getString := t.TestService1.GetString()
	fmt.Println("TestHandler NewComponent", getString)
	return nil
}
