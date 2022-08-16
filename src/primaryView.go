package src

import (
	"github.com/replugged-org/installer/middle"
	"os"
	"path"
	"strings"

	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func (app *UpApplication) ShowPrimaryView() {
	warnings := middle.FindWarnings(app.Config)
	npm := true
	for _, v := range warnings {
		if strings.Contains(v.Text, "NPM") {
			npm = false
		}
	}

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
						Element: design.ButtonAction(If(npm, design.ThemeOkActionButton, design.ThemeImpossibleActionButton), "Install", func() {
							if npm {
								app.GSRightwards()
								app.ShowManagerView(false, func() {
									app.GSLeftwards()
									app.ShowPrimaryView()
								})
							}
						}),
						Shrink: 1,
					},
					{
						Basis:  frenyard.Scale(design.DesignScale, 32),
						Shrink: 1,
					},
					{
						Element: design.ButtonAction(If(npm, design.ThemeRemoveActionButton, design.ThemeImpossibleActionButton), "Uninstall", func() {
							if npm {
							}
						}),
						Shrink: 1,
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

	for _, v := range warnings {
		fixAction := framework.ButtonBehavior(nil)
		if v.Action == middle.URLAndCloseWarningID {
			url := v.Parameter
			fixAction = func() {
				middle.OpenURL(url)
				os.Exit(0)
			}
		}
		slots = append([]framework.FlexboxSlot{
			{
				Element: design.InformationPanel(design.InformationPanelDetails{
					Text:       v.Text,
					ActionText: v.ActionText,
					Action:     fixAction,
				}),
			},
		}, slots...)
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
