package systems

import (
	"../utils"
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
	"math/rand"
)

var Spritesheet *common.Spritesheet

// 落とし穴のあるX座標
var fallPoint []int

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
	// カメラの進んだ距離
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
		// 画面の真ん中より左に位置していれば、カメラを移動せずプレーヤーを移動する
		if (int(ps.playerEntity.SpaceComponent.Position.X) < ps.distance + int(engo.WindowWidth()) / 2){
			ps.playerEntity.SpaceComponent.Position.X += 5
		} else {
			// 画面の右端に達していなければプレーヤーを移動する
			if (int(ps.playerEntity.SpaceComponent.Position.X) < ps.distance + int(engo.WindowWidth()) - 10){
				ps.playerEntity.SpaceComponent.Position.X += 5
			}
			// カメラを移動する
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.XAxis,
				Value:       5,
				Incremental: true,
			})
			ps.distance += 5
		}
	}
	// プレーヤーを左に移動
	if engo.Input.Button("MoveLeft").Down()  {
		if int(ps.playerEntity.SpaceComponent.Position.X) > ps.distance + 10{
			ps.playerEntity.SpaceComponent.Position.X -= 5
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
	common.CameraBounds = engo.AABB{
		Min: engo.Point{X: 0, Y: 0},
		Max: engo.Point{X: 40000, Y: 300},
	}
}

func (ts *TileSystem) New(w *ecs.World){
	ts.world = w
	// 落とし穴作成中の状態を保持（0 => 作成していない、1以上 => 作成中）
	tileMakingState := 0
	// 雲の作成中の状態を保持 (0の場合:作成していない、奇数の場合:{(x+1)/2}番目の雲の前半を作成中、偶数の場合:{x/2}番目の雲の後半を作成中)
	cloudMakingState := 0
	// 雲の高さを保持
	cloudHeight := 0
	// タイルの作成
	Spritesheet = common.NewSpritesheetWithBorderFromFile("tilemap/tilesheet_grass.png", 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	for j := 0; j < 2800; j++ {
		// 地表の作成
		// すでに作成中でない場合、たまに落とし穴を作る
		if (tileMakingState == 0){
			randomNum := rand.Intn(20)
			if (randomNum == 0) {
				fallPoint = append(fallPoint,j)
				tileMakingState = 1
			}
		}
		// 描画するタイルを保持
		var selectedTile int
		// 描画するタイルを選択
		switch tileMakingState {
			case 0: selectedTile = 1
			case 1: selectedTile = 2
			case 2: tileMakingState += 1; continue
			case 3: tileMakingState += 1; continue
			case 4: selectedTile = 0
		}
		tile := &Tile{BasicEntity: ecs.NewBasic()}
		tile.SpaceComponent.Position = engo.Point{
			X: float32(j * 16),
			Y: float32(237),
		}
		tile.RenderComponent.Drawable = Spritesheet.Cell(selectedTile)
		tile.RenderComponent.SetZIndex(0)
		Tiles = append(Tiles, tile)

		if (tileMakingState > 0){
			if (tileMakingState == 4){
				tileMakingState = 0
				continue
			}
			tileMakingState += 1
		}
		// 雲の作成
		if (cloudMakingState == 0){
			randomNum := rand.Intn(12)
			if (randomNum < 7 && randomNum % 2 == 1) {
				cloudMakingState = randomNum
			}
			cloudHeight = rand.Intn(70) + 10
		}
		if (cloudMakingState != 0){
			cloudTile := cloudMakingState + 9
			cloud := &Tile{BasicEntity: ecs.NewBasic()}
			cloud.SpaceComponent.Position = engo.Point{
				X: float32(j * 16),
				Y: float32(cloudHeight),
			}
			cloud.RenderComponent.Drawable = Spritesheet.Cell(cloudTile)
			cloud.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, cloud)
			// 前半を作成中であれば、次は後半を作成する
			if (cloudMakingState % 2 == 1){
				cloudMakingState += 1
			} else {
				cloudMakingState = 0
			}
		}
		//草の作成
		if (!utils.Contains(fallPoint,j)){
			// 落とし穴の上には作らない
			var grassTile int
			randomNum := rand.Intn(42)
			if (randomNum  < 6) {
				grassTile = 26 + randomNum
				grass := &Tile{BasicEntity: ecs.NewBasic()}
				grass.SpaceComponent.Position = engo.Point{
					X: float32(j * 16),
					Y: float32(221),
				}
				grass.RenderComponent.Drawable = Spritesheet.Cell(grassTile)
				grass.RenderComponent.SetZIndex(1)
				Tiles = append(Tiles, grass)
	
			}
		}
	}
	// 地面の描画
	for i := 0; i < 3; i++ {
		tileMakingState = 0
		for j := 0; j < 2800; j++ {
			if (tileMakingState == 0){
				// 落とし穴を作る場合
				if (utils.Contains(fallPoint,j)){
					tileMakingState = 1
				}
			}
			// 描画するタイルを保持
			var selectedTile int
			// 描画するタイルを選択
			switch tileMakingState {
				case 0: selectedTile = 17
				case 1: selectedTile = 18
				case 2: tileMakingState += 1; continue
				case 3: tileMakingState += 1; continue
				case 4: selectedTile = 16
			}
			tile := &Tile{BasicEntity: ecs.NewBasic()}
            tile.SpaceComponent.Position = engo.Point{
				X: float32(j * 16),
				Y: float32(285 - i * 16),
			}
			tile.RenderComponent.Drawable = Spritesheet.Cell(selectedTile)
			tile.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, tile)

			if (tileMakingState > 0){
				if (tileMakingState == 4){
					tileMakingState = 0
					continue
				}
				tileMakingState += 1
			}
		}
	}
	tileMakingState = 0
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
