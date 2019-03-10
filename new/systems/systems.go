package systems

import (
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
	// "image"
	// "image/color"
)

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type PlayerSystem struct {
	world *ecs.World
	// x軸座標
	positionX int
	// y軸座標
	positionY int
	// ジャンプの時間
	jumpDuration int
	playerEntity *Player
	texture *common.Texture
}

func (*PlayerSystem) Remove(ecs.BasicEntity) {}

func (ps *PlayerSystem) Update(dt float32) {

	if engo.Input.Button("MoveRight").Down()  {
		if float32(ps.positionX) < engo.WindowWidth() - 10{
			ps.positionX += 10
			ps.playerEntity.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
				Width:    30,
				Height:   30,
			}
		}
	}
	if engo.Input.Button("MoveLeft").Down()  {
		if ps.positionX > 10{
			ps.positionX -= 10
			ps.playerEntity.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
				Width:    30,
				Height:   30,
			}
		}
	}
	if engo.Input.Button("Jump").JustPressed() {
		ps.jumpDuration = 1
	}

	if ps.jumpDuration != 0 {
		ps.jumpDuration += 1

		ps.playerEntity.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
			Width:    30,
			Height:   30,
		}
		if ps.jumpDuration < 12 {
			ps.positionY -= 5
		} else if ps.jumpDuration < 22 {
			ps.positionY += 5
		} else {
			ps.jumpDuration = 0
		}
	}
}

func (ps *PlayerSystem) New(w *ecs.World){
	ps.world = w

	player := Player{BasicEntity: ecs.NewBasic()}

	ps.positionX = int(engo.WindowWidth() / 2)
	ps.positionY = int(engo.WindowHeight() - 100)
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
		Width: 30,
		Height: 30,
	}
	texture, err := common.LoadedSprite("pics/greenoctocat.png")
	if err != nil {
		fmt.Println("Unable to load texture: " + err.Error())
	}
	player.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale: engo.Point{X:0.1, Y:0.1},
	}
	ps.playerEntity = &player
	ps.texture = texture
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}

}