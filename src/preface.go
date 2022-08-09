package src

import (
	"fmt"

	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
	"github.com/replugged-org/installer/middle"
)

func (app *UpApplication) ResetWithDiscordInstance(save bool, location string) {
	app.DiscordInstance = nil
	app.Config.DiscordPath = location
	if save {
		middle.WriteConfig(app.Config)
	}
	app.ShowPreface()
}

func (app *UpApplication) ShowPreface() {
	var discordLocations []middle.DiscordInstance

	app.ShowWaiter("Loading...", func(progress func(string)) {
		progress("Checking local installation...")
		di, err := middle.NewDiscordInstance(app.Config.DiscordPath)
		if err == nil {
			app.DiscordInstance = di
		} else {
			fmt.Printf("Failed check: %s\n", err.Error())
		}
		progress("Not configured ; Autodetecting Discord locations...")
		discordLocations = middle.GetChannels()
	}, func() {
		if app.DiscordInstance == nil {
			app.ShowInstanceFinder(discordLocations)
		} else {
			app.CachedPrimaryView = nil
			app.ShowPrimaryView()
		}
	})
}

func (app *UpApplication) ShowInstanceFinder(locations []middle.DiscordInstance) {
	suggestSlots := []framework.FlexboxSlot{}
	for _, location := range locations {
		channel := location.Channel
		path := location.Path
		suggestSlots = append(suggestSlots, framework.FlexboxSlot{
			Element: design.ListItem(design.ListItemDetails{
				Icon:    design.DirectoryIconID,
				Text:    channel,
				Subtext: path,
				Click: func() {
					app.GSRightwards()
					app.ResetWithDiscordInstance(true, path)
				},
			}),
			RespectMinimumSize: true,
		})
	}

	suggestSlots = append(suggestSlots, framework.FlexboxSlot{
		Grow:   1,
		Shrink: 0,
	})

	foundInstallsScroller := design.ScrollboxV(framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		WrapMode:    framework.FlexboxWrapModeNone,
		Slots:       suggestSlots,
	}))

	content := framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
		DirVertical: true,
		Slots: []framework.FlexboxSlot{
			{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("Welcome to the official Replugged installer. Before we begin, we need to know which Discord instance to install to.", design.GlobalFont), design.ThemeText, 0, frenyard.Alignment2i{}),
			},
			{
				Basis: design.SizeMarginAroundEverything,
			},
			{
				Element:            foundInstallsScroller,
				Grow:               1,
				Shrink:             1,
				RespectMinimumSize: true,
			},
			{
				Basis: design.SizeMarginAroundEverything,
			},
			{
				Element: design.ButtonBar([]framework.UILayoutElement{
					design.ButtonAction(design.ThemeOkActionButton, "LOCATE MANUALLY", func() {
						app.GSDownwards()
						app.ShowDiscordFinder(func() {
							app.GSUpwards()
							app.ShowInstanceFinder(locations)
						}, middle.BrowserVFSPathDefault)
					}),
				}),
			},
		},
	})

	primary := design.LayoutDocument(design.Header{
		Title:    "Welcome",
		BackIcon: design.WarningIconID,
		Back: func() {
			app.GSLeftwards()
			app.ShowCredits(func() {
				app.GSRightwards()
				app.ShowInstanceFinder(locations)
			})
		},
	}, content, true)
	app.Teleport(primary)
}
