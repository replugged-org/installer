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
		{
			Text:    "toonlink",
			Subtext: "Wrote the Discord finder code",
		},
		{
			Text:    "20kdc",
			Subtext: "Fixed up Discord detection from VFS",
		},
	}

	var listSlots []framework.FlexboxSlot
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
