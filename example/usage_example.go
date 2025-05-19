package example

import (
	"context"
	"fmt"
)

// StringEvent implements the Event interface for string payloads
type StringEvent struct {
	id      string
	evtType string
	data    string
}

func (e StringEvent) ID() string      { return e.id }
func (e StringEvent) Type() string    { return e.evtType }
func (e StringEvent) Payload() []byte { return []byte(e.data) }

// JsonHandler implements EventHandler for StringEvent
type JsonHandler struct{}

func (h JsonHandler) Handle(ctx context.Context, event StringEvent) error {
	fmt.Printf("Handling event %s of type %s with data: %s\n", event.ID(), event.Type(), event.data)
	return nil
}

// StringQueue implements EventQueue for string events
type StringQueue struct {
	events []StringEvent
}

func NewStringQueue() *EventQueue[StringEvent] {
	return &EventQueue[StringEvent]{
		// Implementation details
	}
}

// Creating a UsageExample function to demonstrate how to use the generic services
func UsageExample() {
	// Example 1: Using EventProcessor with StringEvent type
	processorConfig := &ProcessorConfig{
		BatchSize:   10,
		RetryPolicy: "exponential",
	}

	// Create the EventProcessor using the generated constructor
	processor, err := NewEventProcessor[StringEvent](EventProcessorParams[StringEvent]{
		Handler:    JsonHandler{},
		Queue:      NewStringQueue(),
		MaxWorkers: 5,
		Config:     processorConfig,
	})
	if err != nil {
		fmt.Printf("Error creating processor: %v\n", err)
		return
	}

	fmt.Printf("Created processor with %d workers\n", processor.MaxWorkers)

	// Example 2: Using GenericService with string type
	stringRepo := &InMemoryRepository[string]{}
	logger := ConsoleLogger{}
	options := &GenericOptions{
		Timeout:    30,
		MaxRetries: 3,
		CacheTTL:   300,
	}

	// Create the GenericService using the generated constructor
	stringService, err := NewGenericService[string](GenericServiceParams[string]{
		Repository: stringRepo,
		Logger:     logger,
		Options:    options,
	})
	if err != nil {
		fmt.Printf("Error creating string service: %v\n", err)
		return
	}

	fmt.Printf("Created string service with repository: %v\n", stringService.Repository)

	// Example 3: Using CacheService with string key and any value type
	maxSize := 1000
	cacheService, err := NewCacheService[string, interface{}](CacheServiceParams[string, interface{}]{
		Store:      &MemoryStore[string, interface{}]{},
		Serializer: &JsonSerializer[interface{}]{},
		TTL:        60,
		MaxSize:    &maxSize,
	})
	if err != nil {
		fmt.Printf("Error creating cache service: %v\n", err)
		return
	}

	fmt.Printf("Created cache service with TTL: %d\n", cacheService.TTL)

	// Example 4: Using AnotherService with ConsoleLogger
	anotherService, err := NewAnotherService(AnotherServiceParams{
		Logger:  ConsoleLogger{},
		Options: Options{Region: "us-west-2", Endpoint: "https://api.example.com"},
		Timeout: 60,
	})
	if err != nil {
		fmt.Printf("Error creating another service: %v\n", err)
		return
	}

	fmt.Printf("Created another service with timeout: %d\n", anotherService.Timeout)
}

// Additional implementations for the example

// InMemoryRepository is a simple in-memory implementation of Repository
type InMemoryRepository[T any] struct {
	items map[string]T
}

func (r *InMemoryRepository[T]) Get(ctx context.Context, id string) (T, error) {
	item, ok := r.items[id]
	if !ok {
		var zero T
		return zero, fmt.Errorf("item not found: %s", id)
	}
	return item, nil
}

func (r *InMemoryRepository[T]) Save(ctx context.Context, item T) error {
	if r.items == nil {
		r.items = make(map[string]T)
	}
	// Implementation
	return nil
}

func (r *InMemoryRepository[T]) List(ctx context.Context) ([]T, error) {
	// Implementation
	var result []T
	return result, nil
}

// ConsoleLogger implements Logger
type ConsoleLogger struct{}

func (l ConsoleLogger) Debug(msg string) { fmt.Println("DEBUG:", msg) }
func (l ConsoleLogger) Info(msg string)  { fmt.Println("INFO:", msg) }
func (l ConsoleLogger) Error(msg string) { fmt.Println("ERROR:", msg) }

// MemoryStore implements KeyValueStore
type MemoryStore[K comparable, V any] struct {
	data map[K]V
}

func (s *MemoryStore[K, V]) Get(key K) (V, bool) {
	value, ok := s.data[key]
	return value, ok
}

func (s *MemoryStore[K, V]) Set(key K, value V) error {
	if s.data == nil {
		s.data = make(map[K]V)
	}
	s.data[key] = value
	return nil
}

func (s *MemoryStore[K, V]) Delete(key K) error {
	delete(s.data, key)
	return nil
}

// JsonSerializer implements Serializer
type JsonSerializer[T any] struct{}

func (s *JsonSerializer[T]) Marshal(item T) ([]byte, error) {
	// Implementation
	return []byte{}, nil
}

func (s *JsonSerializer[T]) Unmarshal(data []byte) (T, error) {
	// Implementation
	var zero T
	return zero, nil
}
