package systems

import (
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
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
	// ジャンプの時間
	jumpDuration int
	// プレーヤーの進んだ距離
	distance int
	playerEntity *Player
	texture *common.Texture
}

type TileSystem struct {
	world *ecs.World
	// x軸座標
	positionX int
	// y軸座標
	positionY int
	tileEntity []*Tile
	texture *common.Texture
}

func (*PlayerSystem) Remove(ecs.BasicEntity) {}

func (*TileSystem) Remove(ecs.BasicEntity) {}

func (ts *TileSystem) Update(dt float32) {
	// 背景を移動させる
	// to do 
}

func (ps *PlayerSystem) Update(dt float32) {
	// 右移動
	if engo.Input.Button("MoveRight").Down()  {	
		// 画面の真ん中より左に位置していれば、移動する
		if (int(ps.playerEntity.SpaceComponent.Position.X) < (ps.distance + int(engo.WindowWidth())) / 2){
			ps.playerEntity.SpaceComponent.Position.X += 5
			// ps.playerEntity.SpaceComponent = common.SpaceComponent{
			// 	Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
			// 	Width:    30,
			// 	Height:   30,
			// }
		} else {
			ps.playerEntity.SpaceComponent.Position.X += 5
			// ps.playerEntity.SpaceComponent = common.SpaceComponent{
			// 	Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
			// 	Width:    30,
			// 	Height:   30,
			// }
			// カメラを移動する
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.XAxis,
				Value:       5,
				Incremental: true,
			})
			// // プレーヤーは画面の真ん中に
			// ps.positionX = int(engo.WindowWidth() / 2)
			// ps.playerEntity.SpaceComponent = common.SpaceComponent{
			// 	Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
			// 	Width:    30,
			// 	Height:   30,
			// }
			ps.distance += 5
		}
	}
	// プレーヤーを左に移動
	if engo.Input.Button("MoveLeft").Down()  {
		if ps.playerEntity.SpaceComponent.Position.X > 10{
			ps.playerEntity.SpaceComponent.Position.X -= 5
			// ps.playerEntity.SpaceComponent = common.SpaceComponent{
			// 	Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
			// 	Width:    30,
			// 	Height:   30,
			// }
		}
	}
	// プレーヤーをジャンプ
	if engo.Input.Button("Jump").JustPressed() {
		if ps.jumpDuration == 0 {
			ps.jumpDuration = 1
		}
	}
	if ps.jumpDuration != 0 {
		ps.jumpDuration += 1

		// ps.playerEntity.SpaceComponent = common.SpaceComponent{
		// 	Position: engo.Point{X:float32(ps.positionX),Y:float32(ps.positionY)},
		// 	Width:    30,
		// 	Height:   30,
		// }
		if ps.jumpDuration < 14 {
			ps.playerEntity.SpaceComponent.Position.Y -= 5
		} else if ps.jumpDuration < 26 {
			ps.playerEntity.SpaceComponent.Position.Y += 5
		} else {
			ps.jumpDuration = 0
		}
	}
}

func (ps *PlayerSystem) New(w *ecs.World){
	ps.world = w
	// プレーヤーの作成
	player := Player{BasicEntity: ecs.NewBasic()}

	// 初期の配置
	positionX := int(engo.WindowWidth() / 2)
	positionY := int(engo.WindowHeight() - 88)
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X:float32(positionX),Y:float32(positionY)},
		Width: 30,
		Height: 30,
	}
	// 画像の読み込み
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
	// タイルの作成
	Spritesheet = common.NewSpritesheetWithBorderFromFile("tilemap/tilesheet_grass.png", 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	// 地面の描画
	for i := 0; i < 3; i++ {
		for j := 0; j < 280; j++ {
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
	// 地表の作成
	for j := 0; j < 280; j++ {
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
				ts.tileEntity = append(ts.tileEntity, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}

}
