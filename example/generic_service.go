package example

import "context"

// GenericService is a service that works with a generic type
//
//go:generate go run github.com/stijmetkii/validation-gen/cmd/gen
type GenericService[T any] struct {
	Repository Repository[T]
	Logger     Logger
	Options    *GenericOptions
}

// Repository is a generic repository interface
type Repository[T any] interface {
	Get(ctx context.Context, id string) (T, error)
	Save(ctx context.Context, item T) error
	List(ctx context.Context) ([]T, error)
}

// GenericOptions contains options for generic services
type GenericOptions struct {
	Timeout    int
	MaxRetries int
	CacheTTL   int
}

// CacheService uses both generics and interfaces
//
//go:generate go run github.com/stijmetkii/validation-gen/cmd/gen
type CacheService[K comparable, V any] struct {
	Store      KeyValueStore[K, V]
	Serializer Serializer[V]
	TTL        int
	MaxSize    *int
}

// KeyValueStore is a generic key-value store interface
type KeyValueStore[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V) error
	Delete(key K) error
}

// Serializer can convert between a type and bytes
type Serializer[T any] interface {
	Marshal(item T) ([]byte, error)
	Unmarshal(data []byte) (T, error)
}

// EventProcessor processes events with constraints on the generic type
//
//go:generate go run github.com/stijmetkii/validation-gen/cmd/gen
type EventProcessor[E Event] struct {
	Handler    EventHandler[E]
	Queue      *EventQueue[E]
	MaxWorkers int
	Config     *ProcessorConfig
}

// Event is an interface for all event types
type Event interface {
	ID() string
	Type() string
	Payload() []byte
}

// EventHandler processes events
type EventHandler[E Event] interface {
	Handle(ctx context.Context, event E) error
}

// EventQueue manages event queues
type EventQueue[E Event] struct {
	// Queue implementation details
}

// ProcessorConfig contains configuration for event processors
type ProcessorConfig struct {
	BatchSize   int
	RetryPolicy string
}
