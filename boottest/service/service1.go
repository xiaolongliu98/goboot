package service

import (
	"fmt"
	"github.com/xiaolongliu98/goboot/booter"
)

type StringGetter interface {
	GetString() string
}

type TestService1 struct {
	Msg          string
	TestService2 *TestService2 `autowired:""`
}

func (s *TestService1) Initialize(ctx *booter.BootContext) error {
	fmt.Println("TestService1 NewComponent")
	return nil
}

func (s *TestService1) GetString() string {
	return "test1"
}
