package booter

type ComponentType string

type ComponentConfig struct {
	Name      string
	Singleton bool
}

var defaultConfig = &ComponentConfig{
	Name:      "", // Name为空，默认取结构体名
	Singleton: true,
}

func DefaultConfig(option ...ConfigOption) *ComponentConfig {
	config := &ComponentConfig{
		Name:      defaultConfig.Name,
		Singleton: defaultConfig.Singleton,
	}
	for _, opt := range option {
		opt(config)
	}
	return config
}

type ConfigOption func(*ComponentConfig)

func WithName(name string) ConfigOption {
	return func(c *ComponentConfig) {
		c.Name = name
	}
}

func WithSingleton(singleton bool) ConfigOption {
	return func(c *ComponentConfig) {
		c.Singleton = singleton
	}
}
