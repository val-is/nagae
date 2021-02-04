package nagae

import "github.com/hajimehoshi/ebiten"

type SceneManager struct {
	scenes       map[SceneId]*Scene
	currentScene SceneId
	sceneStack   []SceneId

	sharedData map[string]interface{}
}

func NewSceneManager(startScene *Scene) *SceneManager {
	manager := &SceneManager{
		scenes:       make(map[SceneId]*Scene),
		currentScene: startScene.sceneId,
		sceneStack:   make([]SceneId, 0),

		sharedData: make(map[string]interface{}),
	}
	err := manager.AddScene(startScene)
	if err != nil {
		panic("something went (impossibly) wrong when creating a scene manager")
	}
	return manager
}

func (s SceneManager) CurrentScene() SceneId   { return s.currentScene }
func (s SceneManager) Scene(id SceneId) *Scene { return s.scenes[s.currentScene] }

func (s *SceneManager) Transition() error {
	s.currentScene = s.sceneStack[0]
	if len(s.sceneStack) == 1 {
		s.sceneStack = make([]SceneId, 0)
	} else {
		s.sceneStack = s.sceneStack[1:]
	}
	return s.scenes[s.currentScene].Init()
}

func (s *SceneManager) Update(dt float64) error {
	return s.scenes[s.currentScene].Update(dt)
}

func (s *SceneManager) Draw(screen *ebiten.Image) error {
	return s.scenes[s.currentScene].Draw(screen)
}

func (s *SceneManager) PushSceneIdToStack(sceneId SceneId) bool {
	if _, present := s.scenes[sceneId]; !present {
		return false
	}
	s.sceneStack = append(s.sceneStack, sceneId)
	return true
}

func (s *SceneManager) AddScene(scene *Scene) error {
	if _, present := s.scenes[scene.Id()]; present {
		return ErrScenePresent
	}
	scene.manager = s
	s.scenes[scene.Id()] = scene
	return nil
}

func (s SceneManager) GetSharedData(key string) (interface{}, bool) {
	data, present := s.sharedData[key]
	if !present {
		return nil, false
	}
	return data, true
}

func (s *SceneManager) PutSharedData(key string, data interface{}) {
	s.sharedData[key] = data
}
