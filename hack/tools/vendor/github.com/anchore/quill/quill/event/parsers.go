package event

import (
	"fmt"

	"github.com/wagoodman/go-partybus"
	"github.com/wagoodman/go-progress"

	"github.com/anchore/bubbly"
)

type ErrBadPayload struct {
	Type  partybus.EventType
	Field string
	Value interface{}
}

func (e *ErrBadPayload) Error() string {
	return fmt.Sprintf("event='%s' has bad event payload field='%v': '%+v'", string(e.Type), e.Field, e.Value)
}

func newPayloadErr(t partybus.EventType, field string, value interface{}) error {
	return &ErrBadPayload{
		Type:  t,
		Field: field,
		Value: value,
	}
}

func checkEventType(actual, expected partybus.EventType) error {
	if actual != expected {
		return newPayloadErr(expected, "Type", actual)
	}
	return nil
}

func ParseCLIInputPromptType(e partybus.Event) (bubbly.PromptWriter, error) {
	if err := checkEventType(e.Type, CLIInputPromptType); err != nil {
		return nil, err
	}

	p, ok := e.Value.(bubbly.PromptWriter)
	if !ok {
		return nil, newPayloadErr(e.Type, "Value", e.Value)
	}

	return p, nil
}

func ParseCLIReportType(e partybus.Event) (string, string, error) {
	if err := checkEventType(e.Type, CLIReportType); err != nil {
		return "", "", err
	}

	context, ok := e.Source.(string)
	if !ok {
		// this is optional
		context = ""
	}

	report, ok := e.Value.(string)
	if !ok {
		return "", "", newPayloadErr(e.Type, "Value", e.Value)
	}

	return context, report, nil
}

func ParseTaskType(e partybus.Event) (*Task, progress.StagedProgressable, error) {
	if err := checkEventType(e.Type, TaskType); err != nil {
		return nil, nil, err
	}

	cmd, ok := e.Source.(Task)
	if !ok {
		return nil, nil, newPayloadErr(e.Type, "Source", e.Source)
	}

	p, ok := e.Value.(progress.StagedProgressable)
	if !ok {
		return nil, nil, newPayloadErr(e.Type, "Value", e.Value)
	}

	return &cmd, p, nil
}

func ParseCLINotificationType(e partybus.Event) (string, string, error) {
	if err := checkEventType(e.Type, CLINotificationType); err != nil {
		return "", "", err
	}

	context, ok := e.Source.(string)
	if !ok {
		// this is optional
		context = ""
	}

	notification, ok := e.Value.(string)
	if !ok {
		return "", "", newPayloadErr(e.Type, "Value", e.Value)
	}

	return context, notification, nil
}
