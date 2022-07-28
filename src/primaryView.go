package src

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func (app *UpApplication) ShowPrimaryView() {
	slots := []framework.FlexboxSlot{}

	slots = append(slots, framework.FlexboxSlot{
		Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("Welcome to the Replugged installer!", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
		Grow:    1,
	})

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Replugged Installer",
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots:       slots,
	}), true));
}
