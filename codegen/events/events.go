package events

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenEventsPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenEventsPlugin) Init() error {
	// Any initialization logic for the plugin
	g.Files = map[string]func() string{
		"app/events/events.go": func() string {
			return `package events

import (
    "context"
    "log"
    "sync"
    "time"
)

type Event struct {
    Name string
    Data interface{}
}

type EventHandler func(context.Context, Event) error

type Subscription struct {
    CreatedAt int64
    EventName string
    Handler   EventHandler
}

type EventManager struct {
    mu        sync.RWMutex
    listeners map[string][]Subscription
    eventCh   chan Event
    quitCh    chan struct{}
}

func NewEventManager() *EventManager {
    em := &EventManager{
        listeners: make(map[string][]Subscription),
        eventCh:   make(chan Event, 128),
        quitCh:    make(chan struct{}),
    }
    go em.start()
    return em
}

func (em *EventManager) start() {
    ctx := context.Background()
    for {
        select {
        case <-em.quitCh:
            return
        case event := <-em.eventCh:
            if handlers, found := em.listeners[event.Name]; found {
                for _, sub := range handlers {
                    go func(sub Subscription, event Event) {
                        ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
                        defer cancel()
                        start := time.Now()
                        if err := sub.Handler(ctx, event); err != nil {
                            log.Printf("Error handling event %s: %v", event.Name, err)
                        }
                        log.Printf("Handled event %s in %v", event.Name, time.Since(start))
                    }(sub, event)
                }
            }
        }
    }
}

func (em *EventManager) Stop() {
    em.quitCh <- struct{}{}
}

func (em *EventManager) RegisterListener(eventName string, handler EventHandler) Subscription {
    em.mu.Lock()
    defer em.mu.Unlock()

    sub := Subscription{
        CreatedAt: time.Now().UnixNano(),
        EventName: eventName,
        Handler:   handler,
    }

    em.listeners[eventName] = append(em.listeners[eventName], sub)

    return sub
}

func (em *EventManager) UnregisterListener(sub Subscription) {
    em.mu.Lock()
    defer em.mu.Unlock()

    if handlers, found := em.listeners[sub.EventName]; found {
        for i, s := range handlers {
            if s.CreatedAt == sub.CreatedAt {
                em.listeners[sub.EventName] = append(handlers[:i], handlers[i+1:]...)
                break
            }
        }
        if len(em.listeners[sub.EventName]) == 0 {
            delete(em.listeners, sub.EventName)
        }
    }
}

func (em *EventManager) Emit(event Event) {
    em.eventCh <- event
}
`
		},
	}
	return nil
}

func (g *GenEventsPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenEventsPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenEventsPlugin) Name() string {
	return "GenEventsPlugin"
}

func (g *GenEventsPlugin) Version() string {
	return "1.0.0"
}

func (g *GenEventsPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenEventsPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenEventsPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenEventsPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenEventsPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/events"
}

func (g *GenEventsPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenEventsPlugin(data config.ProjectData) *GenEventsPlugin {
	return &GenEventsPlugin{
		Data: data,
	}
}
