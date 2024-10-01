package booter

type Component interface {
}

type ComponentInitializer interface {
	Initialize(ctx *BootContext) error
}

type ComponentConfigurable interface {
	ComponentConfig(ctx *BootContext) *ComponentConfig
}

type ComponentCleaner interface {
	ComponentCleanup(ctx *BootContext) error
}
