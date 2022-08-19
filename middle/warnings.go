package middle

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type WarningID int

const (
	// NullActionWarningID cannot be automatically fixed.
	NullActionWarningID WarningID = iota
	// InstallOrUpdatePackageWarningID warnings can be solved by installing/updating the package Parameter.
	InstallOrUpdatePackageWarningID
	// URLAndCloseWarningID warnings can be solved manually by the user given navigation to a URL. The application closes.
	URLAndCloseWarningID
)

type Warning struct {
	Text       string
	Action     WarningID
	ActionText string
	Parameter  string
}

var remoteVersion string
var hasAlreadyCheckedUpdate = false

func checkUpdate() {
	if !hasAlreadyCheckedUpdate {
		hasAlreadyCheckedUpdate = true
	} else {
		return
	}

	res, _ := http.Get("https://raw.githubusercontent.com/replugged-org/installer/main/middle/version.go")

	data, _ := ioutil.ReadAll(res.Body)
	remoteVersion = strings.Trim(string(data)[33:38], "\r\n")
}

var npm = false
var hasAlreadyCheckedNpm = false

func checkNpm() {
	if !hasAlreadyCheckedNpm {
		hasAlreadyCheckedNpm = true
	} else {
		return
	}

	var cmd *exec.Cmd
	var err error
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("zsh", "-c", "'which npm'")
		err = cmd.Run()
		if err == nil {
			npm = true
		} else {
			fmt.Println("Error finding NPM (#1): " + err.Error())
			fmt.Println("Trying again using brew...")
			_, err := os.Stat("/opt/homebrew/bin/npm")
			if err == nil {
				npm = true
			} else {
				fmt.Println("Error finding NPM (#2): " + err.Error())
			}
		}
	} else {
		cmd = exec.Command("npm", "help")
		err = cmd.Run()
		if err == nil {
			npm = true
		} else {
			fmt.Println("Error finding NPM: " + err.Error())
		}
	}
}

func FindWarnings(config Config) []Warning {
	warnings := []Warning{}

	if !hasAlreadyCheckedUpdate {
		checkUpdate()
	}
	if remoteVersion != version {
		warnings = append(warnings, Warning{
			Text:       "A new version of the installer is available! (v" + remoteVersion + ")",
			Action:     URLAndCloseWarningID,
			ActionText: "UPDATE",
			Parameter:  "https://github.com/replugged-org/installer/releases",
		})
	}

	if !hasAlreadyCheckedNpm {
		checkNpm()
	}
	if !npm {
		warnings = append(warnings, Warning{
			Text:       "NPM is not installed.",
			Action:     URLAndCloseWarningID,
			ActionText: "INSTALL",
			Parameter:  "https://docs.npmjs.com/downloading-and-installing-node-js-and-npm",
		})
	}

	return warnings
}
