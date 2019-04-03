package systems

import (
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

// Text is an entity containing text printed to the screen
type Text struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

// HUDTextMessage updates the HUD text based on messages sent from other systems
type HUDTextMessage struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.MouseComponent
	Line1, Line2, Line3, Line4 string
}

// HUDTextEntity is an entity for the text system. This keeps track of the position
// size and text associated with that position.
type HUDTextEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.MouseComponent
	Line1, Line2, Line3, Line4 string
}

// HUDTextSystem prints the text to our HUD based on the current state of the game
type HUDTextSystem struct {
	text1, text2, text3, text4, money Text

	entities []HUDTextEntity

	updateMoney bool
	amount      int
}

func (h *HUDTextSystem) New(w *ecs.World) {
	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.Black,
		Size: 40,
	}
	fnt.CreatePreloaded()

	h.text1 = Text{BasicEntity: ecs.NewBasic()}
	h.text1.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "GAME OVER",
	}
	h.text1.SetShader(common.TextHUDShader)
	h.text1.RenderComponent.SetZIndex(1001)
	h.text1.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 100, Y: engo.WindowHeight() - 200},
	}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&h.text1.BasicEntity, &h.text1.RenderComponent, &h.text1.SpaceComponent)
		}
	}
}

// Add adds an entity to the system
func (h *HUDTextSystem) Add(b *ecs.BasicEntity, s *common.SpaceComponent, m *common.MouseComponent, l1, l2, l3, l4 string) {
	h.entities = append(h.entities, HUDTextEntity{b, s, m, l1, l2, l3, l4})
}

func (h *HUDTextSystem) Update(dt float32) {}

// Remove takes an enitty out of the system.
func (h *HUDTextSystem) Remove(basic ecs.BasicEntity) {}