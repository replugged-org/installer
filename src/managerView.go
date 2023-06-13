package src

import (
	"fmt"
	"github.com/replugged-org/installer/middle"
	"path/filepath"
	"strings"

	"os"
	"path"
	"time"

	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/lexisother/frenyard/integration"
)

func (app *UpApplication) ShowManagerView(installed bool, back framework.ButtonBehavior) {
	if !installed {
		showInstallScreen(app)
	} else {
		showUninstallScreen(app, back)
	}
}

func showInstallScreen(app *UpApplication) {
	if _, err := os.Stat(path.Join(app.Config.DiscordPath, "app/plugged.txt")); err == nil {
		app.MessageBox("Already installed!", "Replugged is already installed. Please restart your client.", func() {
			app.CachedPrimaryView = nil
			app.ShowPrimaryView()
		})
	} else {
		log := "--Log started at " + time.Now().Format(time.RFC1123) + " --"
		errorLog := "Errors:"
		app.ShowWaiter("Installing...", func(progress func(string)) {
			// Why these loose error declarations? If you put a goto above a variable declaration, the compiler will error.
			// As I've been told by 2767mr:
			//   "You cannot jump over variable declarations.
			//   Goto statements were primarily intended for generated code because they are useful but hard to understand.
			//   If you really want it, you will need to move your variable declarations above the goto and the label."
			var downloadErr error
			var renameErr error
			var resources os.FileInfo
			var index *os.File
			var packageJson *os.File

			asarPath := path.Join(middle.GetDataPath(), "replugged.asar")

			log += "\nDownload replugged.asar..."
			progress(log)
			downloadErr = middle.DownloadFile(asarPath, "https://replugged.dev/api/v1/store/dev.replugged.Replugged.asar")
			if downloadErr != nil {
				errorLog += "\n  Downloading replugged.asar: " + downloadErr.Error()
				goto FinishEarly
			}

			log += "\nRenaming app.asar..."
			progress(log)
			renameErr = os.Rename(path.Join(app.Config.DiscordPath, "app.asar"), path.Join(app.Config.DiscordPath, "app.orig.asar"))
			if renameErr != nil {
				errorLog += "\n  Renaming app.asar to app.orig.asar: " + renameErr.Error()
				goto FinishEarly
			}

			log += "\nChecking for app folder..."
			progress(log)
			resources, _ = os.Stat(path.Join(app.Config.DiscordPath, "app"))
			if resources != nil {
				errorLog += "\n  App folder already exists! Is another mod currently installed?"
				goto FinishEarly
			}

			os.Mkdir(path.Join(app.Config.DiscordPath, "app"), 0755)
			index, _ = os.Create(path.Join(app.Config.DiscordPath, "app/index.js"))
			packageJson, _ = os.Create(path.Join(app.Config.DiscordPath, "app/package.json"))
			log += "\nWriting package.json..."
			progress(log)
			packageJson.WriteString(`{"name": "discord","main":"index.js"}`)
			log += "\nWriting index.js..."
			progress(log)
			index.WriteString(fmt.Sprintf("require('%s')", strings.ReplaceAll(filepath.FromSlash(asarPath), "\\", "\\\\")))

		FinishEarly:
			if errorLog != "Errors:" {
				log += "\n\n-- Errors occurred during installation. --\n" + errorLog
			} else {
				log += "\n\n-- Complete; Restart your Discord client! --"
			}
			progress(log)
		}, func() {
			pluggedFile, _ := os.Create(path.Join(app.Config.DiscordPath, "app/plugged.txt"))
			pluggedFile.WriteString("this file was added to indicate that replugged is installed here.")
			app.GSInstant()
			app.MessageBox("Install Complete", log, func() {
				app.CachedPrimaryView = nil
				app.GSLeftwards()
				app.ShowPrimaryView()
			})
		})
	}
}

func showUninstallScreen(app *UpApplication, back framework.ButtonBehavior) {
	if _, err := os.Stat(path.Join(app.Config.DiscordPath, "app/plugged.txt")); err != nil {
		app.MessageBox("Not installed!", "Replugged is not installed. Please install it before trying to remove it.", func() {
			app.CachedPrimaryView = nil
			app.ShowPrimaryView()
		})
	} else {
		listSlots := []framework.FlexboxSlot{
			{
				Grow: 1,
			},
			{
				Element: framework.NewUILabelPtr(integration.NewTextTypeChunk("Are you sure you want to uninstall Replugged?", design.GlobalFont), 0xFFFFFFFF, 0, frenyard.Alignment2i{}),
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
							Element: design.ButtonAction(design.ThemeRemoveActionButton, "Uninstall", func() {
								if _, err := os.Stat(path.Join(app.Config.DiscordPath, "app/plugged.txt")); err != nil {
								} else {
									log := "-- Log started at " + time.Now().Format(time.RFC1123) + " --"
									errorLog := "Errors:"
									app.ShowWaiter("Uninstalling...", func(progress func(string)) {
										// Same reason as before...
										var removeErr error
										var renameErr error

										log += "\nDeleting the app directory..."
										progress(log)
										removeErr = os.RemoveAll(path.Join(app.Config.DiscordPath, "app"))
										if removeErr != nil {
											errorLog += "\n  failed deleting app directory: " + removeErr.Error()
											goto FinishEarly
										}

										log += "\nRename app.orig.asar to app.asar..."
										progress(log)
										renameErr = os.Rename(path.Join(app.Config.DiscordPath, "app.orig.asar"), path.Join(app.Config.DiscordPath, "app.asar"))
										if renameErr != nil {
											errorLog += "\n  Renaming app.orig.asar to app.asar: " + renameErr.Error()
											goto FinishEarly
										}

									FinishEarly:
										if errorLog != "Errors:" {
											log += "\n-- Errors occurred during installation. --\n" + errorLog
										} else {
											log += "\n-- Complete; Restart your Discord client! --"
										}
										progress(log)
									}, func() {
										app.GSInstant()
										app.MessageBox("Uninstall Complete", log, func() {
											app.CachedPrimaryView = nil
											app.GSRightwards()
											app.ShowPrimaryView()
										})
									})
								}
							}),
						},
						{
							Basis: frenyard.Scale(design.DesignScale, 32),
						},
						{
							Element: design.ButtonAction(design.ThemeOkActionButton, "Cancel", back),
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
			Title: "Replugged",
			Back:  back,
		}, framework.NewUIFlexboxContainerPtr(framework.FlexboxContainer{
			DirVertical: true,
			Slots:       listSlots,
		}), true))
	}
}
