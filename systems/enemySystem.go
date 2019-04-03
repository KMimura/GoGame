package systems

import (
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
	"engo.io/engo/common"
	"math/rand"
)

type Enemy struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	// ジャンプの状態(0 => 着地中, 1 => １ジャンプ中, 2 => 降下中)
	jumpState int
	// ジャンプの残り時間
	jumpDuration int
	// 移動の速度(0 ~ 2, 数値が高いほど早い)
	velocity int
	// 画面から消えているか
	ifDissappearing bool
}

type EnemySystem struct {
	world *ecs.World
	enemyEntity []*Enemy
	texture *common.Texture
}

func (es *EnemySystem) Remove(entity ecs.BasicEntity) {
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

func (es *EnemySystem) Update(dt float32) {
	// カメラとプレーヤーの位置を取得
	var cameraPosition float32
	var playerPositionX float32
	for _, system := range es.world.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			cameraPosition = sys.X()
		case *PlayerSystem:
			playerPositionX = sys.playerEntity.SpaceComponent.Position.X
		}
	}
	for _, o := range es.enemyEntity{
		// 画面に描画されていないオブジェクトは移動処理をしない
		if (o.SpaceComponent.Position.X > cameraPosition - 240 && o.SpaceComponent.Position.X < cameraPosition + 200 && !o.ifDissappearing){
			// プレーヤーとの当たり判定
			if (o.SpaceComponent.Position.X == playerPositionX) {
				for _, system := range es.world.Systems() {
					switch sys := system.(type) {
					case *PlayerSystem:
						sys.playerEntity.damage += 1
					}
				}
			}
			o.SpaceComponent.Position.X -= float32(o.velocity + 1)
			// ジャンプをしていない場合
			if (o.jumpState == 0){
				o.jumpState = rand.Intn(2) + 1
				jumpTemp := rand.Intn(3)
				switch (jumpTemp) {
					case 0: o.jumpDuration = 15
					case 1: o.jumpDuration = 25
					case 2: o.jumpDuration = 35
				}
			}
			// ジャンプ処理
			if (o.jumpState == 1){
				// ジャンプをし終わっていない場合
				if (o.jumpDuration > 0){
					o.SpaceComponent.Position.Y -= 3
					o.jumpDuration -= 1
				} else {
					// ジャンプをし終わった場合
					o.jumpState = 2
				}
			} else {
				// 降下をし終わっていない場合
				if (o.SpaceComponent.Position.Y < 212){
					o.SpaceComponent.Position.Y += 3
				} else {
					// 降下し終わった場合
					o.jumpState = 0
				}
			}
		}else if (o.ifDissappearing){
			o.SpaceComponent.Position.Y += 3
		}
	}
}

func (es *EnemySystem) New(w *ecs.World){
	es.world = w
	Enemies := make([]*Enemy, 0)
	// ランダムで配置
	for i := 0; i < 44800; i++ {
		randomNum := rand.Intn(400)
		if (randomNum == 0){
			// 敵の作成
			enemy := Enemy{BasicEntity: ecs.NewBasic()}
			enemy.SpaceComponent = common.SpaceComponent{
				Position: engo.Point{X:float32(i),Y:float32(212)},
				Width: 30,
				Height: 30,
			}
			// 画像の読み込み
			texture, err := common.LoadedSprite("pics/ghost.png")
			if err != nil {
				fmt.Println("Unable to load texture: " + err.Error())
			}
			enemy.RenderComponent = common.RenderComponent{
				Drawable: texture,
				Scale: engo.Point{X:1.1, Y:1.1},
			}
			enemy.RenderComponent.SetZIndex(1)
			es.texture = texture
			for _, system := range es.world.Systems() {
				switch sys := system.(type) {
				case *common.RenderSystem:
					sys.Add(&enemy.BasicEntity, &enemy.RenderComponent, &enemy.SpaceComponent)
				}
			}
			enemy.velocity = rand.Intn(3)
			Enemies = append(Enemies,&enemy)
		}
		es.enemyEntity = Enemies
	}
}