package main

import (
	"engo.io/engo"
	"engo.io/engo/common"
	"engo.io/ecs"
	"fmt"
	"image/color"
)

type myScene struct {}

type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (*myScene) Type() string { return "myGame" }

func (*myScene) Preload() {
	engo.Files.Load("pics/greenoctocat.png")
	common.SetBackground(color.White)
}

func (*myScene) Setup(u engo.Updater){
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	player := Player{BasicEntity: ecs.NewBasic()}
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{10, 10},
		Width:    10,
		Height:   10,
	}
	texture, err := common.LoadedSprite("pics/greenoctocat.png")
	if err != nil {
		fmt.Println("Unable to load texture: " + err.Error())
	}
	player.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{1, 1},
	}
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
}

func main(){
	opts := engo.RunOptions{
		Title:"myGame",
		Width:400,
		Height:400,
	}
	engo.Run(opts,&myScene{})
}