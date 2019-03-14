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
	engo.Files.Load("tilemap/tilesheet_grass.png")
	common.SetBackground(color.RGBA{255, 250, 220, 0})
}

func (*myScene) Setup(u engo.Updater){
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("Jump", engo.KeySpace)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.TileSystem{})
	world.AddSystem(&systems.PlayerSystem{})
}

func main(){
	opts := engo.RunOptions{
		Title:"myGame",
		Width:400,
		Height:300,
		StandardInputs: true,
		NotResizable:true,
	}
	engo.Run(opts,&myScene{})
}