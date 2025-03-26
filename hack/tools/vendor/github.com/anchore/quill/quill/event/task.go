package event

import "github.com/wagoodman/go-progress"

type ManualStagedProgress struct {
	progress.Stage
	progress.Manual
}

type Title struct {
	Default      string
	WhileRunning string
	OnSuccess    string
	OnFail       string
}

type Task struct {
	Title   Title
	Context string
}
