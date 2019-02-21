package systems

import (
	"engo.io/ecs"
	"fmt"
	"engo.io/engo"
)

type PlayerSystem struct {}

func (*PlayerSystem) Remove(ecs.BasicEntity) {}

func (*PlayerSystem) Update(dt float32) {
	if engo.Input.Button("AddPlayer").JustPressed()  {
		fmt.Println("The gamer pressed A")
	}
}

func (*PlayerSystem) New(*ecs.World){
	fmt.Println("IT WORKED!!")
}