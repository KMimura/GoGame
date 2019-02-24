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
	common.SetBackground(color.White)
}

func (*myScene) Setup(u engo.Updater){
	engo.Input.RegisterButton("AddPlayer", engo.KeyX)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	kbs := common.NewKeyboardScroller(
		400, engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
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