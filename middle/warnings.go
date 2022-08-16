package middle

import "os/exec"

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
	Text      string
	Action    WarningID
	Parameter string
}

var npm = false
var hasAlreadyCheckedNpm = false

func checkNpm() {
	if !hasAlreadyCheckedNpm {
		hasAlreadyCheckedNpm = true
	} else {
		return
	}

	if _, err := exec.LookPath("npm"); err == nil {
		npm = true
	}
}

func FindWarnings(config Config) []Warning {
	warnings := []Warning{}

	if !hasAlreadyCheckedNpm {
		checkNpm()
	}
	if !npm {
		warnings = append(warnings, Warning{
			Text:      "NPM is not installed.",
			Action:    URLAndCloseWarningID,
			Parameter: "https://docs.npmjs.com/downloading-and-installing-node-js-and-npm",
		})
	}

	return warnings
}
