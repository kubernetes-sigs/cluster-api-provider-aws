package quill

import (
	"github.com/wagoodman/go-partybus"

	"github.com/anchore/go-logger"
	"github.com/anchore/go-logger/adapter/redact"
	"github.com/anchore/quill/internal/bus"
	"github.com/anchore/quill/internal/log"
	intRedact "github.com/anchore/quill/internal/redact"
)

// SetLogger sets the logger object used for all logging calls.
func SetLogger(logger logger.Logger) {
	useOrAddRedactor()
	log.Set(logger)
}

// SetBus sets the event bus for all library bus publish events onto (in-library subscriptions are not allowed).
func SetBus(b *partybus.Bus) {
	useOrAddRedactor()
	bus.Set(b)
}

func useOrAddRedactor() {
	// since it is possible to read secrets from the environment during lib calls, we want to ensure that the logger
	// is redacted even if the user did not explicitly set a redaction store.
	store := intRedact.Get()
	if store == nil {
		intRedact.Set(redact.NewStore())
	}
}
