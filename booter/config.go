package booter

type ComponentType string

const (
	TypeService ComponentType = "service"
	TypeHandler ComponentType = "handler"
	TypeDAO     ComponentType = "dao"
	TypeDefault ComponentType = ""
)

type ComponentConfig struct {
	Name      string
	Singleton bool
	Type      ComponentType
}

var defaultConfig = &ComponentConfig{
	Name:      "", // Name为空，默认取结构体名
	Singleton: true,
	Type:      TypeDefault,
}

func DefaultConfig(option ...ConfigOption) *ComponentConfig {
	config := &ComponentConfig{
		Name:      defaultConfig.Name,
		Singleton: defaultConfig.Singleton,
		Type:      defaultConfig.Type,
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

func WithType(t ComponentType) ConfigOption {
	return func(c *ComponentConfig) {
		c.Type = t
	}
}
