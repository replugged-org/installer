package src

import (
	"os"
	"path"

	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func (app *UpApplication) ShowPrimaryView() {
	var installStatus string
	if _, installedOrNot := os.Stat(path.Join(app.Config.DiscordPath, "app/plugged.txt")); installedOrNot == nil {
		installStatus = "installed!"
	} else {
		installStatus = "not installed."
	}

	slots := []framework.FlexboxSlot{
		{
			Grow: 1,
		},
		{
			Element: framework.NewUILabelPtr(
				integration.NewTextTypeChunk("Welcome to the Replugged installer!", design.GlobalFont),
				0xFFFFFFFF,
				0,
				frenyard.Alignment2i{},
			),
		},
		{
			Element: framework.NewUILabelPtr(
				integration.NewTextTypeChunk("Replugged is currently: "+installStatus, design.GlobalFont),
				0xFFFFFFFF,
				0,
				frenyard.Alignment2i{},
			),
		},
		{
			Basis:  frenyard.Scale(design.DesignScale, 32),
			Shrink: 1,
		},
		{
			Element: framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
				DirVertical: false,
				Slots: []framework.FlexboxSlot{
					{
						Grow: 1,
					},
					{
						Element: design.ButtonAction(design.ThemeOkActionButton, "Install", func() {
							app.GSRightwards()
							app.ShowManagerView(false, func() {
								app.GSLeftwards()
								app.ShowPrimaryView()
							})
						}),
						Shrink: 1,
					},
					{
						Basis:  frenyard.Scale(design.DesignScale, 32),
						Shrink: 1,
					},
					{
						Element: design.ButtonAction(design.ThemeRemoveActionButton, "Uninstall", func() {}),
						Shrink:  1,
					},
					{
						Grow: 1,
					},
				},
			}),
		},
		{
			Grow: 1,
		},
	}

	app.Teleport(design.LayoutDocument(design.Header{
		Title: "Replugged Installer",
		Back: func() {
			app.CachedPrimaryView = nil
			app.GSLeftwards()
			app.ResetWithDiscordInstance(false, "computer://")
		},
		BackIcon:    design.BackIconID,
		ForwardIcon: design.MenuIconID,
		Forward: func() {
			app.GSRightwards()
			app.ShowOptionsMenu(func() {
				app.GSLeftwards()
				app.ShowPrimaryView()
			})
		},
	}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots:       slots,
	}), true))
}
