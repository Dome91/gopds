package util

import (
	"context"
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
)

//go:generate mockgen -destination=../mock/util/bus.go -source=bus.go

const (
	GenerateCoverTopic   = "catalogEntry.generateCover"
	ExtractMetadataTopic = "catalogEntry.extractMetadata"
)

type Bus interface {
	RegisterHandler(key string, handler *bus.Handler)
	Emit(ctx context.Context, topic string, data interface{}) (*bus.Event, error)
}

func NewBus() *bus.Bus {
	node := uint64(1)
	initialTime := uint64(1577865600000)
	m, err := monoton.New(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		panic(err)
	}

	// init an id generator
	var idGenerator bus.Next = (*m).Next

	// create a new bus instance
	b, err := bus.NewBus(idGenerator)
	if err != nil {
		panic(err)
	}

	b.RegisterTopics(GenerateCoverTopic)

	return b
}

type GenerateCoverEvent struct {
	ID string
}

type ExtractMetadataEvent struct {
	ID string
}
