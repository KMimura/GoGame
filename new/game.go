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
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("Jump", engo.KeySpace)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.PlayerSystem{})
}

func main(){
	opts := engo.RunOptions{
		Title:"myGame",
		Width:800,
		Height:600,
		StandardInputs: true,
	}
	engo.Run(opts,&myScene{})
}