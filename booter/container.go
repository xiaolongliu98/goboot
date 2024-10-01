package booter

type Container struct {
	Singleton bool
	Instances map[Component]Component // parent -> child
}

func NewContainer(singleton bool) *Container {
	c := &Container{
		Singleton: singleton,
		Instances: make(map[Component]Component, 1),
	}
	return c
}

// GetInstance parent为空时或者为单例模式时，默认的parent为nil
func (c *Container) GetInstance(parent ...Component) Component {
	if c.Singleton {
		return c.Instances[nil]
	}
	if len(parent) == 0 {
		return c.Instances[nil]
	}
	return c.Instances[parent[0]]
}

// SetInstance parent为空时或者为单例模式时，默认的parent为nil
func (c *Container) SetInstance(instance Component, parent ...Component) {
	if isNil(instance) {
		return
	}
	if len(parent) == 0 || c.Singleton {
		c.Instances[nil] = instance
	} else {
		c.Instances[parent[0]] = instance
	}
}

// ExistInstanceByParentInstance parent为空时或者为单例模式时，默认的parent为nil
func (c *Container) ExistInstanceByParentInstance(parent ...Component) bool {
	var ok bool
	if len(parent) == 0 || c.Singleton {
		_, ok = c.Instances[nil]
	} else {
		_, ok = c.Instances[parent[0]]
	}
	return ok
}

func (c *Container) NotEmpty() bool {
	return len(c.Instances) > 0
}

// RemoveInstance parent为空时或者为单例模式时，默认的parent为nil
func (c *Container) RemoveInstance(parent ...Component) {
	if len(parent) == 0 || c.Singleton {
		delete(c.Instances, nil)
	} else {
		delete(c.Instances, parent[0])
	}
}

func (c *Container) ExistInstance(target Component) bool {
	for _, instance := range c.Instances {
		if instance == target {
			return true
		}
	}
	return false
}

// GetParentInstance 获取instance的parent
func (c *Container) GetParentInstance(instance Component) Component {
	for parent, child := range c.Instances {
		if child == instance {
			return parent
		}
	}
	return nil
}

func (c *Container) ForeachInstance(fn func(instance Component)) {
	for _, instance := range c.Instances {
		fn(instance)
	}
}
