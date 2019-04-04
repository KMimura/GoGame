package systems

import (
	"../utils"
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"math/rand"
	"time"
)

var Spritesheet *common.Spritesheet

// 落とし穴のあるX座標
var FallPoint []int
// 落とし穴の始まるX座標
var FallStartPoint []int

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
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

func (*TileSystem) Remove(ecs.BasicEntity) {}

func (ts *TileSystem) Update(dt float32) {
}

func (ts *TileSystem) New(w *ecs.World){
	rand.Seed(time.Now().UnixNano())

	ts.world = w
	// 落とし穴作成中の状態を保持（0 => 作成していない、1以上 => 作成中）
	tileMakingState := 0
	// 雲の作成中の状態を保持 (0の場合:作成していない、奇数の場合:{(x+1)/2}番目の雲の前半を作成中、偶数の場合:{x/2}番目の雲の後半を作成中)
	cloudMakingState := 0
	// 雲の高さを保持
	cloudHeight := 0
	// タイルの作成
	tmp := rand.Intn(2)
	var loadTxt string
	if tmp == 0 {
		loadTxt = "tilemap/tilesheet_grass.png"
	} else {
		loadTxt = "tilemap/tilesheet_snow.png"
	}
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	for j := 0; j < 2800; j++ {
		// 地表の作成
		if (j > 10){
			if (tileMakingState > 1 && tileMakingState < 4){
				for t:= 0; t < 8; t++ {
					FallPoint = append(FallPoint,j * 16 - t)
				}
			} else if (tileMakingState == 0){
				// すでに作成中でない場合、たまに落とし穴を作る
				randomNum := rand.Intn(10)
				if (randomNum == 0) {
					FallStartPoint = append(FallStartPoint,j * 16)
					tileMakingState = 1
				}
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
		// タイルEntityの作成
		tile := &Tile{BasicEntity: ecs.NewBasic()}
		// 位置情報の設定
		tile.SpaceComponent.Position = engo.Point{
			X: float32(j * 16),
			Y: float32(237),
		}
		// 見た目の設定
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
	for j := 0; j < 2800; j++ {
		// 雲の作成
		if (cloudMakingState == 0){
			randomNum := rand.Intn(6)
			if (randomNum < 7 && randomNum % 2 == 1) {
				cloudMakingState = randomNum
			}
			cloudHeight = rand.Intn(70) + 10
		}
		if (cloudMakingState != 0){
			// 雲Entityの作成
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
		if (!utils.Contains(FallPoint,j * 16)){
			// 落とし穴の上には作らない
			var grassTile int
			randomNum := rand.Intn(18)
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
				if (utils.Contains(FallStartPoint,j * 16)){
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