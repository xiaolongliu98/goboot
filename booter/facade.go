package booter

import (
	"reflect"
)

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
		ctx.updateTime()
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

	instanceContainer := getInstance(ctx, component, "", -1, nil)
	ctx.handleCircleRecord()
	return instanceContainer.GetInstance(nil)
}

func (ctx *BootContext) GetInstanceByName(name string) Component {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	ctx.initCircleContext()
	defer ctx.releaseCircleContext()

	instanceContainer := getInstanceByName(ctx, name, "", -1, nil)
	ctx.handleCircleRecord()
	return instanceContainer.GetInstance(nil)
}

// getInstance 通过Component获取实例，存在递归初始化行为
func getInstance(ctx *BootContext, component Component, parentName string, parentFieldIdx int, parentInstance Component) *Container {
	// component是指针类型，需要取指针的值
	if !isPtr(component) {
		panic("component must be a pointer")
	}
	name := ctx.getComponentName(component)
	return getInstanceByName(ctx, name, parentName, parentFieldIdx, parentInstance)
}

// getInstanceByName 通过Component名获取实例，存在递归初始化行为
func getInstanceByName(ctx *BootContext, name string, parentName string, parentFieldIdx int, parentInstance Component) *Container {
	// check register
	component, exists := ctx.components[name]
	if !exists {
		panic("component not register, component:" + name)
	}

	config := ctx.getConfig(component)
	instanceContainer, ok := ctx.container[name]
	if !ok {
		instanceContainer = NewContainer(config.Singleton)
		ctx.container[name] = instanceContainer
	}

	if !instanceContainer.ExistInstanceByParentInstance() || !config.Singleton {
		// check circle
		if _, exists := ctx.circleCheck[name]; exists {
			// 非单例模式下，存在循环依赖，这种情况是不允许的
			if !config.Singleton {
				panic("circle dependency in case of non-singleton, component:" + name)
			}
			// 进入循环依赖，先返回zero，记录下来，后续再初始化
			ctx.circleRecord = append(ctx.circleRecord, []any{parentName, parentFieldIdx, name, parentInstance})
			return nil
		}
		ctx.circleCheck[name] = struct{}{}

		instance := newPointerInstance(component)
		instanceContainer.SetInstance(instance, parentInstance)

		// 递归初始化instance依赖的Component（tag标记autowired:""）
		// 遍历instance的字段，如果字段是Component类型，且tag标记autowired:""，则初始化
		collector := make([][]any, 0)
		foreachSubComponent(component, func(fieldComponent Component, filedIndex int, autowiredName string) {
			// 递归初始化instance依赖的Component
			var filedInstanceContainer *Container
			if autowiredName == "" {
				filedInstanceContainer = getInstance(ctx, fieldComponent, name, filedIndex, instance)
			} else {
				filedInstanceContainer = getInstanceByName(ctx, autowiredName, name, filedIndex, instance)
			}
			collector = append(collector, []any{filedIndex, filedInstanceContainer})
		})

		instanceValue := reflect.ValueOf(instance).Elem()
		for _, val := range collector {
			i := val[0].(int)
			// 循环依赖情况下，属性值为nil，直接跳过
			if isNil(val[1]) {
				continue
			}
			childContainer := val[1].(*Container)
			childInstance := childContainer.GetInstance(instance)
			if isNil(childInstance) {
				continue
			}
			instanceValue.Field(i).Set(reflect.ValueOf(childInstance))
		}
		delete(ctx.circleCheck, name)

		// lifecycle function
		if initializer, ok := instance.(ComponentInitializer); ok {
			if err := initializer.Initialize(ctx); err != nil {
				panic(err)
			}
		}
	}

	return instanceContainer
}

func (ctx *BootContext) handleCircleRecord() {
	for _, record := range ctx.circleRecord {
		parentName := record[0].(string)
		parentFieldIdx := record[1].(int)
		name := record[2].(string)
		parentInstance := record[3].(Component)

		if parentName == "" || name == "" || parentFieldIdx == -1 || isNil(parentInstance) {
			continue
		}

		childContainer, ok := ctx.container[name]
		if !ok {
			panic("get instance failed, component:" + name)
		}
		childInstance := childContainer.GetInstance(parentInstance)
		parentValue := reflect.ValueOf(parentInstance).Elem()
		parentValue.Field(parentFieldIdx).Set(reflect.ValueOf(childInstance))
	}
}

func (ctx *BootContext) ExistComponentByName(name string) bool {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	_, exists := ctx.components[name]
	return exists
}

func (ctx *BootContext) ExistAnyInstanceByName(name string) bool {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	container, exists := ctx.container[name]
	if !exists {
		return false
	}
	return container.NotEmpty()
}

func (ctx *BootContext) ExistComponent(component Component) bool {
	name := ctx.getComponentName(component)
	return ctx.ExistComponentByName(name)
}

func (ctx *BootContext) ExistInstance(target Component) bool {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	name := ctx.getComponentName(target)
	container, exists := ctx.container[name]
	if !exists {
		return false
	}
	return container.ExistInstance(target)
}

func (ctx *BootContext) GetTypeComponents(t any) []Component {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	var res []Component
	for _, component := range ctx.components {
		if isEmbedded(component, t) {
			res = append(res, component)
		}
	}
	return res
}

// GetTypeInstances 获取嵌入了指定类型的所有实例
// t: 指定的嵌入类型，如booter.HandlerComponent，也可以是自定义的接口类型
func (ctx *BootContext) GetTypeInstances(t any) []Component {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	var res []Component
	for _, instanceContainer := range ctx.container {
		instanceContainer.ForeachInstance(func(instance Component) {
			if isEmbedded(instance, t) {
				res = append(res, instance)
			}
		})
	}
	return res
}

// GetParentInstance 获取父实例
func (ctx *BootContext) GetParentInstance(instance Component) Component {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()

	name := ctx.getComponentName(instance)
	container, exists := ctx.container[name]
	if !exists {
		return nil
	}
	return container.GetParentInstance(instance)
}

func (ctx *BootContext) CleanupAll() {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	for _, container := range ctx.cleanSort() {
		container.ForeachInstance(func(instance Component) {
			if cleaner, ok := instance.(ComponentCleaner); ok {
				err := cleaner.ComponentCleanup(ctx)
				if err != nil {
					panic(err)
				}
			}
		})
	}

	ctx.container = map[string]*Container{}
	ctx.components = map[string]Component{}
	ctx.m = map[string]any{}
	ctx.updateTime()
}

func (ctx *BootContext) cleanSort() []*Container {
	var res []*Container

	g := map[string]*Node{}
	for name, _ := range ctx.container {
		buildDependencyGraphForComponent(ctx, g, name)
	}

	// 排除graph中没有被实例化的Node
	toBeRemovedList := make([]*Node, 0)
	for name, node := range g {
		container, ok := ctx.container[name]
		if !ok || !container.NotEmpty() {
			toBeRemovedList = append(toBeRemovedList, node)
		}
	}
	deleteGraphNodeFunc := func(node *Node, g map[string]*Node) {
		for _, childNode := range node.Children {
			delete(childNode.Parents, node.Name)
		}
		for _, parentNode := range node.Parents {
			delete(parentNode.Children, node.Name)
		}
		delete(g, node.Name)
	}
	for _, node := range toBeRemovedList {
		deleteGraphNodeFunc(node, g)
	}

	// 进行拓扑排序
	// 收集没有父依赖的节点
	q := make([]*Node, 0)
	for _, node := range g {
		if len(node.Parents) == 0 {
			q = append(q, node)
		}
	}

	for len(q) > 0 {
		size := len(q)
		for i := 0; i < size; i++ {
			node := q[0]
			q = q[1:]
			// collect
			res = append(res, ctx.container[node.Name])
			// get next
			children := node.Children
			deleteGraphNodeFunc(node, g)
			for _, child := range children {
				if len(child.Parents) == 0 {
					q = append(q, child)
				}
			}
		}
	}

	// 剩下没有排序完的，是形成环的节点
	// 默认形成环的节点的cleanup是无强依赖关系的
	for name, _ := range g {
		res = append(res, ctx.container[name])
	}

	return res
}

func buildDependencyGraphForComponent(ctx *BootContext, graph map[string]*Node, name string) *Node {
	component, ok := ctx.components[name]
	if !ok {
		return nil
	}
	if node, ok := graph[name]; ok {
		return node
	}

	curNode := &Node{
		Name:     name,
		Children: map[string]*Node{},
		Parents:  map[string]*Node{},
	}
	graph[name] = curNode

	// 去重
	subComponents := map[string]struct{}{}
	foreachSubComponent(component, func(fieldComponent Component, _ int, autowiredName string) {
		if autowiredName == "" {
			autowiredName = ctx.getComponentName(fieldComponent)
		}
		subComponents[autowiredName] = struct{}{}
	})

	for subComponentName := range subComponents {
		childNode := buildDependencyGraphForComponent(ctx, graph, subComponentName)
		if childNode == nil {
			continue
		}

		curNode.Children[childNode.Name] = childNode
		childNode.Parents[name] = curNode
	}

	return curNode
}

type Node struct {
	Name     string
	Children map[string]*Node // 依赖的子节点，出边
	Parents  map[string]*Node // 被依赖的父节点，入边
}
