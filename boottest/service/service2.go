package service

import (
	"fmt"
	"github.com/xiaolongliu98/goboot/booter"
)

type TestService2 struct {
	Msg          string
	TestService1 *TestService1 `autowired:""`
}

func (s *TestService2) Initialize(ctx *booter.BootContext) error {
	fmt.Println("TestService2 NewComponent")
	return nil
}
