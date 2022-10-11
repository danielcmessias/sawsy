package context

import (
	"github.com/danielcmessias/sawsy/config"
	"github.com/danielcmessias/sawsy/utils"
)

type ProgramContext struct {
	ScreenHeight      int
	ScreenWidth       int
	MainContentWidth  int
	MainContentHeight int
	AwsAccountId      string
	AwsService        string
	Config            *config.Config

	Keys utils.KeyMap
	// View              config.ViewTyp

	LockKeyboardCapture bool
}
