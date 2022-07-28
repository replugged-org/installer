package main

import (
	"github.com/lexisother/frenyard"
	"github.com/lexisother/frenyard/design"
	"github.com/lexisother/frenyard/framework"
	"github.com/replugged-org/installer/src"
)

func main() {
	frenyard.TargetFrameTime = 0.016
	slideContainer := framework.NewUISlideTransitionContainerPtr(nil)
	slideContainer.FyEResize(design.SizeWindowInit)
	wnd, err := framework.CreateBoundWindow("Replugged Installer", true, design.ThemeBackground, slideContainer)
	if err != nil {
		panic(err)
	}
	design.Setup(frenyard.InferScale(wnd))
	wnd.SetSize(design.SizeWindow)
	app := &src.UpApplication{
		MainContainer:    slideContainer,
		Window:           wnd,
		UpQueued:         make(chan func(), 16),
		TeleportSettings: framework.SlideTransition{},
	}
	app.ShowPrimaryView()
	frenyard.GlobalBackend.Run(func(frametime float64) {
		select {
		case fn := <-app.UpQueued:
			fn()
		default:
		}
	})
}
