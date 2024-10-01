package booter

func Register(component ...Component) {
	defaultBootContext.Register(component...)
}

func GetInstance[T Component](component T) T {
	return defaultBootContext.GetInstance(component).(T)
}
func GetInstanceByName[T Component](name string) T {
	return defaultBootContext.GetInstanceByName(name).(T)
}

func CleanupAll() {
	defaultBootContext.CleanupAll()
}
