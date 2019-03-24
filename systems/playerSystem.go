package systems

import (
	"../utils"
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// ジャンプの時間
	jumpDuration int
	// カメラの進んだ距離
	distance int
	// 落ちているかどうか
	ifFalling bool
	// ダメージ
	damage int
}

type PlayerSystem struct {
	world *ecs.World
	playerEntity *Player
	texture *common.Texture
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

func (*PlayerSystem) Remove(ecs.BasicEntity) {}

func (ps *PlayerSystem) Update(dt float32) {
	// 当たり判定
	if (ps.playerEntity.jumpDuration == 0 && utils.Contains(FallPoint,int(ps.playerEntity.SpaceComponent.Position.X)) ){
		ps.playerEntity.ifFalling = true
		ps.playerEntity.SpaceComponent.Position.Y += 5
	}
	if(!ps.playerEntity.ifFalling){
			// 右移動
	if engo.Input.Button("MoveRight").Down()  {	
		// 画面の真ん中より左に位置していれば、カメラを移動せずプレーヤーを移動する
		if (int(ps.playerEntity.SpaceComponent.Position.X) < ps.playerEntity.distance + int(engo.WindowWidth()) / 2){
			ps.playerEntity.SpaceComponent.Position.X += 5
		} else {
			// 画面の右端に達していなければプレーヤーを移動する
			if (int(ps.playerEntity.SpaceComponent.Position.X) < ps.playerEntity.distance + int(engo.WindowWidth()) - 10){
				ps.playerEntity.SpaceComponent.Position.X += 5
			}
			// カメラを移動する
			engo.Mailbox.Dispatch(common.CameraMessage{
				Axis:        common.XAxis,
				Value:       5,
				Incremental: true,
			})
			ps.playerEntity.distance += 5
		}
	}
	// プレーヤーを左に移動
	if engo.Input.Button("MoveLeft").Down()  {
		if int(ps.playerEntity.SpaceComponent.Position.X) > ps.playerEntity.distance + 10{
			ps.playerEntity.SpaceComponent.Position.X -= 5
		}
	}
	// プレーヤーをジャンプ
	if engo.Input.Button("Jump").JustPressed() {
		if ps.playerEntity.jumpDuration == 0 {
			ps.playerEntity.jumpDuration = 1
		}
	}
	if ps.playerEntity.jumpDuration != 0 {
		ps.playerEntity.jumpDuration += 1
		if ps.playerEntity.jumpDuration < 14 {
			ps.playerEntity.SpaceComponent.Position.Y -= 5
		} else if ps.playerEntity.jumpDuration < 26 {
			ps.playerEntity.SpaceComponent.Position.Y += 5
		} else {
			ps.playerEntity.jumpDuration = 0
		}
	}
	}
}

