package booter

//
//type Elem struct {
//	CurrentName string
//	ParentName  string
//	CurrentIdx  int
//}
//
//type Node struct {
//	Name string
//	Out  map[string]*NextNode // 依赖的组件
//	In   map[string]*NextNode // 被依赖的组件
//}
//
//type NextNode struct {
//	Node  *Node
//	Idx   int
//	Score int
//}
//
//func BuildGraph(name string) map[string]*Node {
//	// idx-0: current component name
//	// idx-1: parent component name
//	q := []Elem{
//		{name, "", -1},
//	}
//
//	g := map[string]*Node{}
//
//	for len(q) > 0 {
//		size := len(q)
//		for i := 0; i < size; i++ {
//			currentName := q[0].CurrentName
//			parentName := q[0].ParentName
//			currentIdx := q[0].CurrentIdx
//			q = q[1:]
//
//			currentComponent, exists := ctx.components[currentName]
//			// check register
//			if !exists {
//				panic("component not register, component:" + currentName)
//			}
//
//			created := true
//			currentNode := g[currentName]
//			if currentNode == nil {
//				currentNode = &Node{Name: currentName, Out: map[string]*NextNode{}, In: map[string]*NextNode{}}
//				g[currentName] = currentNode
//				created = false
//			}
//
//			if parentName != "" {
//				parentNode := g[parentName]
//				parentNode.Out[currentName].Node = currentNode
//				parentNode.Out[currentName].Idx = currentIdx
//
//				currentNode.In[parentName].Node = parentNode
//				currentNode.In[parentName].Idx = currentIdx
//
//				if currentIdx == -2 {
//					// require 权重更高
//					parentNode.Out[currentName].Score += 16384
//					currentNode.In[parentName].Score += 16384
//				} else if currentIdx >= 0 {
//					parentNode.Out[currentName].Score++
//					currentNode.In[parentName].Score++
//				}
//			}
//			if created {
//				continue
//			}
//
//			if configurable := currentComponent.(ComponentConfigurable); configurable != nil {
//				for _, require := range configurable.ComponentConfig(ctx).Requires {
//					requireName := getComponentName(require)
//					q = append(q, Elem{requireName, currentName, -2})
//				}
//			}
//			foreachFieldComponent(currentComponent, func(idx int, name string) {
//				q = append(q, Elem{name, currentName, idx})
//			})
//		}
//	}
//	return g
//}
//
//func GetInstanceByName2[T Component](name string) T {
//	// TODO check
//
//	g := BuildGraph(name)
//
//	q := make([]string, 0)
//	// 先一遍拓扑排序，将无环的节点都先创建
//	for keyName, node := range g {
//		if len(node.Out) == 0 {
//			// 无依赖
//			q = append(q, keyName)
//		}
//	}
//
//	for len(q) > 0 {
//		size := len(q)
//		for i := 0; i < size; i++ {
//			currentName := q[0]
//			q = q[1:]
//
//			currentNode := g[currentName]
//			for nextName, nextNode := range currentNode.In {
//				delete(nextNode.Node.Out, currentName)
//			}
//		}
//	}
//
//}
//
//func foreachFieldComponent(parent Component, f func(idx int, name string)) {
//	currentValue := reflect.ValueOf(parent).Elem()
//	currentType := reflect.TypeOf(parent).Elem()
//	for i := 0; i < currentValue.NumField(); i++ {
//		field := currentValue.Field(i)
//		if field.Kind() != reflect.Ptr {
//			continue
//		}
//		// 获取属性的tag
//		var (
//			autowiredName  string
//			ok             bool
//			fieldComponent Component
//		)
//		if autowiredName, ok = currentType.Field(i).Tag.Lookup("autowired"); !ok {
//			continue
//		}
//		filed := currentValue.Field(i)
//		if filed.Kind() != reflect.Ptr {
//			continue
//		}
//		// 判断属性类型是否是Component
//		if fieldComponent, ok = filed.Interface().(Component); !ok {
//			continue
//		}
//		if autowiredName == "" {
//			autowiredName = getComponentName(fieldComponent)
//		}
//		f(i, autowiredName)
//	}
//}
