package systems

import (
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
)

type MouseTracker struct {
    ecs.BasicEntity
    common.MouseComponent
}

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type PlayerSystem struct {
	mouseTracker MouseTracker
	world *ecs.World
}

func (*PlayerSystem) Remove(ecs.BasicEntity) {}

func (pl *PlayerSystem) Update(dt float32) {
	if engo.Input.Button("AddPlayer").JustPressed()  {
		fmt.Println("The gamer pressed X")
		player := Player{BasicEntity: ecs.NewBasic()}
		player.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{pl.mouseTracker.MouseX, pl.mouseTracker.MouseY},
			Width:    30,
			Height:   30,
		}
		texture, err := common.LoadedSprite("pics/greenoctocat.png")
		if err != nil {
			fmt.Println("Unable to load texture: " + err.Error())
		}
		player.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{0.1, 0.1},
		}
		for _, system := range pl.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
			}
		}
	
	}
}

func (player *PlayerSystem) New(w *ecs.World){
	fmt.Println("IT WORKED!!")

	player.world = w
	player.mouseTracker.BasicEntity = ecs.NewBasic()
	player.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&player.mouseTracker.BasicEntity, &player.mouseTracker.MouseComponent, nil, nil)
		}
	}
}