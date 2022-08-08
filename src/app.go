package src

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/framework"
	"github.com/replugged-org/installer/middle"
)

type UpApplication struct {
	Config            middle.Config
	DiscordInstance   *middle.DiscordInstance
	MainContainer     *framework.UISlideTransitionContainer
	Window            frenyard.Window
	UpQueued          chan func()
	CachedPrimaryView framework.UILayoutElement
	TeleportSettings  framework.SlideTransition
}

const upTeleportLen float64 = 0.25

func (app *UpApplication) GSLeftwards() {
	app.TeleportSettings.Reverse = true
	app.TeleportSettings.Vertical = false
	app.TeleportSettings.Length = upTeleportLen
}
func (app *UpApplication) GSRightwards() {
	app.TeleportSettings.Reverse = false
	app.TeleportSettings.Vertical = false
	app.TeleportSettings.Length = upTeleportLen
}
func (app *UpApplication) GSUpwards() {
	app.TeleportSettings.Reverse = true
	app.TeleportSettings.Vertical = true
	app.TeleportSettings.Length = upTeleportLen
}
func (app *UpApplication) GSDownwards() {
	app.TeleportSettings.Reverse = false
	app.TeleportSettings.Vertical = true
	app.TeleportSettings.Length = upTeleportLen
}
func (app *UpApplication) GSInstant() {
	// direction doesn't matter
	app.TeleportSettings.Length = 0
}
func (app *UpApplication) Teleport(target framework.UILayoutElement) {
	forkTD := app.TeleportSettings
	forkTD.Element = target
	app.MainContainer.TransitionTo(forkTD)
}
