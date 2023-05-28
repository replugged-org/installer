package middle

import (
	"io"
	"net/http"
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

	data, _ := io.ReadAll(res.Body)
	remoteVersion = strings.Trim(string(data)[33:38], "\r\n")
}

func FindWarnings() []Warning {
	var warnings []Warning

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

	return warnings
}
