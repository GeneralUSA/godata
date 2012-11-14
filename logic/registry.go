package logic

type logicRegistry struct {
	subjectRegistry registry
}

var Registry logicRegistry = newLogicRegistry()

type registry struct {
	base   loadCreate
	core   map[string]loadCreate
	custom map[string]loadCreate
}

func newLogicRegistry() logicRegistry {
	subjectFunctions := loadCreate{
		SubjectLoader(loadDefaultSubject),
		SubjectCreator(createDefaultSubject),
	}
	return logicRegistry{
		subjectRegistry: newRegistry(subjectFunctions),
	}
}

func newRegistry(base loadCreate) registry {
	return registry{
		base:   base,
		core:   make(map[string]loadCreate),
		custom: make(map[string]loadCreate),
	}
}

type loadCreate struct {
	loader  interface{}
	creator interface{}
}

func (r registry) getRegistryItem(key string) loadCreate {
	if custom, ok := r.custom[key]; ok {
		return custom
	} else if core, ok := r.core[key]; ok {
		return core
	}
	return r.base
}

func (l logicRegistry) registerCoreSubjectType(name string, create SubjectCreator, load SubjectLoader) {
	l.subjectRegistry.core[name] = loadCreate{create, load}
}

func (l logicRegistry) RegisterSubjectType(name string, create SubjectCreator, load SubjectLoader) {
	l.subjectRegistry.custom[name] = loadCreate{create, load}
}

func (l logicRegistry) SubjectCreator(name string) SubjectCreator {
	return l.subjectRegistry.getRegistryItem(name).creator.(SubjectCreator)
}

func (l logicRegistry) SubjectLoader(name string) SubjectLoader {
	return l.subjectRegistry.getRegistryItem(name).loader.(SubjectLoader)
}
