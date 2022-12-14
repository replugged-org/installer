package src

import (
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
)

func (app *UpApplication) ShowCredits(back framework.ButtonBehavior) {
	items := []design.ListItemDetails{
		{
			Text:    "Alyxia Sother",
			Subtext: "Developer of the Replugged Installer.",
		},
	}

	listSlots := []framework.FlexboxSlot{}
	for _, item := range items {
		listSlots = append(listSlots, framework.FlexboxSlot{
			Element: design.ListItem(item),
		})
	}

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Credits",
		Back:  back,
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots:       listSlots,
	}), true))
}
