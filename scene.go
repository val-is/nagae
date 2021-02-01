package nagae

import "github.com/hajimehoshi/ebiten"

type Scene struct {
	sceneId SceneId

	actors  map[ActorId]*Actor
	manager *SceneManager

	physicsSystem  PhysicsSystem
	graphicsSystem GraphicsSystem
}

func NewScene(sceneId SceneId) *Scene {
	scene := &Scene{
		sceneId: sceneId,

		actors: make(map[ActorId]*Actor),
	}
	physics := NewPhysicsSystem(scene)
	scene.physicsSystem = physics
	graphics := NewGraphicsSystem(scene)
	scene.graphicsSystem = graphics
	return scene
}

func (s Scene) Id() SceneId            { return s.sceneId }
func (s Scene) Manager() *SceneManager { return s.manager }

func (s *Scene) Init() error {
	if err := s.physicsSystem.Init(); err != nil {
		return err
	}
	if err := s.graphicsSystem.Init(); err != nil {
		return err
	}
	for _, actor := range s.actors {
		if err := actor.Init(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scene) Update(dt float64) error {
	if err := s.physicsSystem.Update(dt); err != nil {
		return err
	}
	if err := s.graphicsSystem.Update(dt); err != nil {
		return err
	}
	for _, actor := range s.actors {
		if err := actor.Update(dt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Scene) Draw(screen *ebiten.Image) error {
	if err := s.graphicsSystem.Draw(screen); err != nil {
		return err
	}
	return nil
}

func (s Scene) GetActor(actorId ActorId) (*Actor, bool) {
	actor, present := s.actors[actorId]
	if !present {
		return nil, false
	}
	return actor, true
}

func (s *Scene) AddActor(actor *Actor) bool {
	if _, present := s.GetActor(actor.Id()); present {
		return false
	}
	actor.parentScene = s
	s.actors[actor.actorId] = actor
	return true
}

func (s *Scene) RemoveActor(actorId ActorId) bool {
	if _, present := s.GetActor(actorId); !present {
		return false
	}
	delete(s.actors, actorId)
	return true
}
