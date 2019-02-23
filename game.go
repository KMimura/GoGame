package main

import (
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/ecs"
	"image/color"
	"./systems"
)

type myScene struct {}

func (*myScene) Type() string { return "myGame" }

func (*myScene) Preload() {
	engo.Files.Load("pics/greenoctocat.png")
	common.SetBackground(color.White)
}

func (*myScene) Setup(u engo.Updater){
	engo.Input.RegisterButton("AddPlayer", engo.KeyA)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	world.AddSystem(&systems.PlayerSystem{})
	// player := Player{BasicEntity: ecs.NewBasic()}
	// player.SpaceComponent = common.SpaceComponent{
	// 	Position: engo.Point{10, 10},
	// 	Width:    30,
	// 	Height:   30,
	// }
	// texture, err := common.LoadedSprite("pics/greenoctocat.png")
	// if err != nil {
	// 	fmt.Println("Unable to load texture: " + err.Error())
	// }
	// player.RenderComponent = common.RenderComponent{
	// 	Drawable: texture,
	// 	Scale:    engo.Point{1, 1},
	// }
	// for _, system := range world.Systems() {
	// 	switch sys := system.(type) {
	// 	case *common.RenderSystem:
	// 		sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
	// 	}
	// }
}

func main(){
	opts := engo.RunOptions{
		Title:"myGame",
		Width:600,
		Height:400,
	}
	engo.Run(opts,&myScene{})
}