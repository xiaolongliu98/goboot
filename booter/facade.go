package booter

import "reflect"

// Register 注册Component
// @param component 默认传递的是指针
func (ctx *BootContext) Register(component ...Component) {
	for _, c := range component {
		// check component是指针
		if !isPtr(c) {
			panic("component must be a pointer")
		}
		name := ctx.getComponentName(c)
		ctx.lock.Lock()
		if _, exists := ctx.components[name]; exists {
			ctx.lock.Unlock()
			panic("component already register, component:" + name)
		}
		if isNil(c) {
			c = newPointerInstance(c)
		}
		ctx.components[name] = c
		ctx.lock.Unlock()
	}
}

func (ctx *BootContext) getConfig(component Component) *ComponentConfig {
	if _, ok := component.(ComponentConfigurable); ok {
		return component.(ComponentConfigurable).ComponentConfig(ctx)
	}
	return defaultConfig
}

func (ctx *BootContext) getComponentName(component Component) string {
	name := ctx.getConfig(component).Name
	if name == "" {
		// 通过反射，取component的结构体名，需要注意component是指针
		name = getPtrStructName(component)
	}
	return name
}

func (ctx *BootContext) initCircleContext() {
	ctx.circleCheck = make(map[string]struct{})
	ctx.circleRecord = make([][]any, 0)
}

func (ctx *BootContext) releaseCircleContext() {
	ctx.circleCheck = nil
	ctx.circleRecord = nil
}

func (ctx *BootContext) GetInstance(component Component) Component {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	ctx.initCircleContext()
	defer ctx.releaseCircleContext()

	instance := getInstance(ctx, component, "", -1)
	ctx.handleCircleRecord()
	return instance
}

func (ctx *BootContext) GetInstanceByName(name string) Component {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	ctx.initCircleContext()
	defer ctx.releaseCircleContext()

	instance := getInstanceByName(ctx, name, "", -1)
	ctx.handleCircleRecord()
	return instance
}

// getInstance 通过Component获取实例，存在递归初始化行为
func getInstance(ctx *BootContext, component Component, parentName string, parentFieldIdx int) Component {
	// component是指针类型，需要取指针的值
	if !isPtr(component) {
		panic("component must be a pointer")
	}
	name := ctx.getComponentName(component)
	return getInstanceByName(ctx, name, parentName, parentFieldIdx)
}

// getInstanceByName 通过Component名获取实例，存在递归初始化行为
func getInstanceByName(ctx *BootContext, name string, parentName string, parentFieldIdx int) Component {
	// check register
	component, exists := ctx.components[name]
	if !exists {
		panic("component not register, component:" + name)
	}

	config := ctx.getConfig(component)
	instance := ctx.container[name]

	if instance == nil || !config.Singleton {
		// check circle
		if _, exists := ctx.circleCheck[name]; exists {
			// 非单例模式下，存在循环依赖，这种情况是不允许的
			if !config.Singleton {
				panic("circle dependency in case of non-singleton, component:" + name)
			}
			// 进入循环依赖，先返回zero，记录下来，后续再初始化
			ctx.circleRecord = append(ctx.circleRecord, []any{parentName, parentFieldIdx, name})
			return nil
		}
		ctx.circleCheck[name] = struct{}{}

		// 递归初始化之前，需要先初始化其声明的前置依赖
		for _, require := range config.Requires {
			getInstance(ctx, require, name, -1)
		}
		// 递归初始化instance依赖的Component（tag标记autowired:""）
		// 遍历instance的字段，如果字段是Component类型，且tag标记autowired:""，则初始化
		instanceValue := reflect.ValueOf(component).Elem()
		instanceType := reflect.TypeOf(component).Elem()

		collector := make([][]any, 0, instanceValue.NumField())
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
			if filed.Kind() != reflect.Ptr {
				continue
			}
			// 判断属性类型是否是Component
			if fieldComponent, ok = filed.Interface().(Component); !ok {
				continue
			}

			// 递归初始化instance依赖的Component
			var filedInstance Component
			if autowiredName == "" {
				filedInstance = getInstance(ctx, fieldComponent, name, i)
			} else {
				filedInstance = getInstanceByName(ctx, autowiredName, name, i)
			}

			// 支持循环依赖，就不需要这个判断了
			//if isNil(filedInstance) {
			//	panic("get instance failed, component:" + name)
			//}
			collector = append(collector, []any{i, filedInstance})
		}

		instance = newPointerInstance(component)
		//instance = component.NewComponent(ctx)
		//if isNil(instance) {
		//	return nil
		//}
		//if !isPtr(instance) {
		//	panic("NewComponent must return a pointer, component:" + name)
		//}
		instanceValue = reflect.ValueOf(instance).Elem()
		for _, val := range collector {
			i := val[0].(int)
			// 循环依赖情况下，属性值为nil，直接跳过
			if val[1] == nil {
				continue
			}
			filedInstance := val[1].(Component)
			instanceValue.Field(i).Set(reflect.ValueOf(filedInstance))
		}
		ctx.container[name] = instance
		delete(ctx.circleCheck, name)

		// lifecycle function
		if initializer, ok := instance.(ComponentInitializer); ok {
			if err := initializer.Initialize(ctx); err != nil {
				panic(err)
			}
		}
	}
	return instance
}

func (ctx *BootContext) handleCircleRecord() {
	for _, record := range ctx.circleRecord {
		parentName := record[0].(string)
		parentFieldIdx := record[1].(int)
		name := record[2].(string)
		if parentName == "" || name == "" || parentFieldIdx == -1 {
			continue
		}

		required, ok := ctx.container[name]
		if !ok {
			panic("get instance failed, component:" + name)
		}
		parent, ok := ctx.container[parentName]
		if !ok {
			panic("get instance failed, component:" + parentName)
		}
		parentValue := reflect.ValueOf(parent).Elem()
		parentValue.Field(parentFieldIdx).Set(reflect.ValueOf(required))
	}
}

func (ctx *BootContext) ExistComponentByName(name string) bool {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	_, exists := ctx.components[name]
	return exists
}

func (ctx *BootContext) ExistInstanceByName(name string) bool {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	_, exists := ctx.container[name]
	return exists
}

func (ctx *BootContext) ExistComponent(component Component) bool {
	name := ctx.getComponentName(component)
	return ctx.ExistComponentByName(name)
}

func (ctx *BootContext) ExistInstance(component Component) bool {
	name := ctx.getComponentName(component)
	return ctx.ExistInstanceByName(name)
}
