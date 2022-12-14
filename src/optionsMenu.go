package src

import (
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
)

func (app *UpApplication) ShowOptionsMenu(back framework.ButtonBehavior) {
	backHere := func() {
		app.GSLeftwards()
		app.ShowOptionsMenu(back)
	}

	listSlots := []framework.FlexboxSlot{
		{
			Element: design.ListItem(design.ListItemDetails{
				Text:    "Credits",
				Subtext: "See who is behind the Replugged installer and related projects",
				Click: func() {
					app.GSRightwards()
					app.ShowCredits(backHere)
				},
			}),
		},
	}

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Options",
		Back:  back,
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots:       listSlots,
	}), true))
}
