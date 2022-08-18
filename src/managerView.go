package src

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/replugged-org/installer/middle"

	"github.com/lexisother/frenyard/framework"
)

func (app *UpApplication) ShowManagerView(installed bool, back framework.ButtonBehavior) {
	if !installed {
		showInstallScreen(app)
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
			log += "\nChecking for app folder..."
			progress(log)
			resources, _ := os.Stat(path.Join(app.Config.DiscordPath, "app"))
			if resources != nil {
				log += "\nRenaming app folder..."
				progress(log)
				os.Rename(path.Join(app.Config.DiscordPath, "app"), path.Join(app.Config.DiscordPath, "plug"))
			}
			os.Mkdir(path.Join(app.Config.DiscordPath, "app"), 0755)
			index, _ := os.Create(path.Join(app.Config.DiscordPath, "app/index.js"))
			packageJson, _ := os.Create(path.Join(app.Config.DiscordPath, "app/package.json"))
			log += "\nWriting package.json..."
			progress(log)
			packageJson.WriteString(`{"name": "plug","main":"index.js"}`)
			log += "\nWriting index.js..."
			progress(log)
			index.WriteString(fmt.Sprintf("require('%s')", strings.ReplaceAll(filepath.FromSlash(middle.GetDataPath()+"/repository/src/patcher.js"), "\\", "\\\\")))

			log += "\nChecking for Replugged folder..."
			progress(log)
			replugged, _ := os.Stat(path.Join(middle.GetDataPath(), "repository"))
			if replugged == nil {
				log += "\nReplugged not found. Cloning Replugged..."
				progress(log)
				_, err := git.PlainClone(path.Join(middle.GetDataPath(), "repository"), false, &git.CloneOptions{
					URL:   "https://github.com/replugged-org/replugged",
					Depth: 1,
				})

				// We are running as root, so fix the file permissions
				if middle.IsLinux() {
					userData := middle.GetUserData()
					middle.ChownR(path.Join(middle.GetDataPath(), "repository"), userData.Ruid, userData.Rgid)
				}

				if err != nil {
					errorLog += "\n  cloning Replugged: " + err.Error()
				}
			}

			replugged, _ = os.Stat(path.Join(middle.GetDataPath(), "repository"))
			if replugged != nil {
				if depFolder, _ := os.Stat(path.Join(middle.GetDataPath(), "repository", "node_modules")); depFolder == nil {
					log += "\nInstalling dependencies... (this can take a while)"
					progress(log)

					cmd := exec.Command("npm", "i")
					cmd.Dir = path.Join(middle.GetDataPath(), "repository")
					_, err := cmd.Output()
					if err != nil {
						errorLog += "\n  installing dependencies: " + err.Error()
						return
					}
				}
			}

			if errorLog != "Errors:" {
				log += "\n-- Errors occurred during installation. --\n" + errorLog
			} else {
				log += "\n-- Complete; Restart your Discord client! --"
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
