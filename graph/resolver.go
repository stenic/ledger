package graph

import "github.com/stenic/ledger/internal/pkg/messagebus"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MessageBus messagebus.MessageBus
}
