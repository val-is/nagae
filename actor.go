package nagae

type Actor struct {
	actorId     ActorId
	parentScene *Scene

	componentMask ComponentList
	components    map[ComponentId]Component
}

func NewActor(actorId ActorId) *Actor {
	return &Actor{
		actorId:    actorId,
		components: make(map[ComponentId]Component),
	}
}

func (a Actor) Id() ActorId         { return a.actorId }
func (a Actor) ParentScene() *Scene { return a.parentScene }

func (a Actor) GetComponentBySystemType(componentType ComponentSystem) (Component, bool) {
	if !a.componentMask.CheckComponent(componentType) {
		return nil, false
	}
	for _, component := range a.components {
		if component.SystemType() == componentType {
			return component, true
		}
	}
	return nil, false
}

func (a Actor) GetComponentByType(componentType ComponentType) (Component, bool) {
	for _, component := range a.components {
		if component.ComponentType() == componentType {
			return component, true
		}
	}
	return nil, false
}

func (a Actor) GetComponentById(componentId ComponentId) (Component, bool) {
	if component, present := a.components[componentId]; !present {
		return nil, false
	} else {
		return component, true
	}
}

func (a *Actor) AddComponent(component Component) error {
	if a.componentMask.CheckComponent(component.SystemType()) {
		return ErrComponentPresent
	} else if _, present := a.GetComponentByType(component.ComponentType()); present {
		return ErrComponentPresent
	} else if _, present := a.GetComponentById(component.Id()); present {
		return ErrComponentPresent
	}
	component.SetParent(a)
	a.components[component.Id()] = component
	a.componentMask = a.componentMask.AddComponent(component.SystemType())
	return nil
}

func (a *Actor) RemoveComponentByType(componentType ComponentType) error {
	if component, present := a.GetComponentByType(componentType); !present {
		return ErrComponentNotPresent
	} else {
		delete(a.components, component.Id())
		return nil
	}
}

func (a *Actor) RemoveComponentById(componentId ComponentId) error {
	if component, present := a.GetComponentById(componentId); !present {
		return ErrComponentNotPresent
	} else {
		delete(a.components, component.Id())
		return nil
	}
}

func (a *Actor) Init() error {
	for _, component := range a.components {
		if err := component.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Actor) Update(dt float64) error {
	for _, component := range a.components {
		if err := component.Update(dt); err != nil {
			return err
		}
	}
	return nil
}
