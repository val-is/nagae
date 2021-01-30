package nagae

type Actor struct {
	actorId     ActorId
	parentScene *Scene

	componentMask ComponentList
	components    map[ComponentId]Component
}

func (a Actor) Id() ActorId         { return a.actorId }
func (a Actor) ParentScene() *Scene { return a.parentScene }

func (a Actor) GetComponentByType(componentType ComponentType) (Component, bool) {
	if !a.componentMask.CheckComponent(componentType) {
		return nil, false
	}
	for _, component := range a.components {
		if component.Type() == componentType {
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
	if a.componentMask.CheckComponent(component.Type()) {
		return ErrComponentPresent
	} else if _, present := a.GetComponentById(component.Id()); present {
		return ErrComponentPresent
	}
	component.SetParent(a)
	a.components[component.Id()] = component
	a.componentMask.AddComponent(component.Type())
	return nil
}

func (a *Actor) RemoveComponentByType(componentType ComponentType) error {
	if !a.componentMask.CheckComponent(componentType) {
		return ErrComponentNotPresent
	}
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
