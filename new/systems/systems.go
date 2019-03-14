package systems

import (
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
	"reflect"
)

var Spritesheet *common.Spritesheet

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Tile struct {
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

type TileSystem struct {
	world *ecs.World
	// x軸座標
	positionX int
	// y軸座標
	positionY int
	tileEntity *Tile
	texture *common.Texture
}

func (*PlayerSystem) Remove(ecs.BasicEntity) {}

func (*TileSystem) Remove(ecs.BasicEntity) {}

func (ts *TileSystem) Update(dt float32) {
	for _, system := range ts.world.Systems() {
		fmt.Println(reflect.TypeOf(system))
	}
}

func (ps *PlayerSystem) Update(dt float32) {

	if engo.Input.Button("MoveRight").Down()  {
		if float32(ps.positionX) < engo.WindowWidth() - 10{
			ps.positionX += 5
			ps.playerEntity.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
				Width:    30,
				Height:   30,
			}
		}
	}
	if engo.Input.Button("MoveLeft").Down()  {
		if ps.positionX > 10{
			ps.positionX -= 5
			ps.playerEntity.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
				Width:    30,
				Height:   30,
			}
		}
	}
	if engo.Input.Button("Jump").JustPressed() {
		if ps.jumpDuration == 0 {
			ps.jumpDuration = 1
		}
	}

	if ps.jumpDuration != 0 {
		ps.jumpDuration += 1

		ps.playerEntity.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
			Width:    30,
			Height:   30,
		}
		if ps.jumpDuration < 14 {
			ps.positionY -= 5
		} else if ps.jumpDuration < 26 {
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
	ps.positionY = int(engo.WindowHeight() - 88)
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
	player.RenderComponent.SetZIndex(1)

	ps.playerEntity = &player
	ps.texture = texture
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
}

func (ts *TileSystem) New(w *ecs.World){
	ts.world = w

	Spritesheet = common.NewSpritesheetWithBorderFromFile("tilemap/tilesheet_grass.png", 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	// 地面の描画
	for i := 0; i < 3; i++ {
		for j := 0; j < 28; j++ {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
            tile.SpaceComponent.Position = engo.Point{
				X: float32(j * 16),
				Y: float32(285 - i * 16),
			}
			tile.RenderComponent.Drawable = Spritesheet.Cell(17)
			tile.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, tile)
		}
	}
	for j := 0; j < 28; j++ {
		tile := &Tile{BasicEntity: ecs.NewBasic()}
		tile.SpaceComponent.Position = engo.Point{
			X: float32(j * 16),
			Y: float32(237),
		}
		tile.RenderComponent.Drawable = Spritesheet.Cell(1)
		tile.RenderComponent.SetZIndex(0)
		Tiles = append(Tiles, tile)
	}
	for _, system := range ts.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}

}
