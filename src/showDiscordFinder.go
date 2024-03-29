package src

import (
	"path/filepath"
	"sort"

	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/replugged-org/installer/middle"
)

func (app *UpApplication) ShowDiscordFinder(back framework.ButtonBehavior, vfsPath string) {
	var vfsList []middle.FinderLocation

	app.ShowWaiter("Reading", func(progress func(string)) {
		progress("Scanning to find all of the context in:\n" + vfsPath)
		vfsList = middle.DiscordFinderVFSList(vfsPath)
	}, func() {
		var items []design.ListItemDetails

		for _, v := range vfsList {
			thisLocation := v.Location
			ild := design.ListItemDetails{
				Icon: design.DirectoryIconID,
				Text: filepath.Base(thisLocation),
			}
			ild.Click = func() {
				app.GSRightwards()
				app.ShowDiscordFinder(func() {
					app.GSLeftwards()
					app.ShowDiscordFinder(back, vfsPath)
				}, thisLocation)
			}
			if v.Instance != nil {
				ild.Click = func() {
					app.GSRightwards()
					app.ResetWithDiscordInstance(true, thisLocation)
				}
				ild.Text = "Discord " + v.Instance.Channel
				ild.Subtext = thisLocation
				ild.Icon = design.GameIconID
			} else if v.Drive != "" {
				ild.Text = v.Drive
				ild.Subtext = v.Location
				ild.Icon = design.DriveIconID
			}
			items = append(items, ild)
		}

		sort.Sort(design.SortListItemDetails(items))
		primary := design.LayoutDocument(design.Header{
			Back:  back,
			Title: "Enter Discord's Location",
		}, design.NewUISearchBoxPtr("Discord name...", items), true)
		app.Teleport(primary)
	})
}
