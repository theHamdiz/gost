package middleware

import (
	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
)

type GenMiddlewarePlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenMiddlewarePlugin) Init() error {
	// Initialize Files
	g.Files = map[string]func() string{
		"app/middleware/recoverer.go": func() string {
			return `package middleware

import (
    "log"
    "net/http"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Recovered from panic: %v", err)
            notifyClients("System shutdown unexpectedly")
            notifyByEmail("System Shutdown", fmt.Sprintf("The system was shut down unexpectedly: %v", err))
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        }
    }()
    next.ServeHTTP(w, r)
})
}
`
		},
		"app/middleware/rateLimiter.go": func() string {
			return `package middleware

import (
    "net/http"
    "time"

    "golang.org/x/time/rate"
)

func RateLimiter(limit rate.Limit, burst int) func(http.Handler) http.Handler {
    limiter := rate.NewLimiter(limit, burst)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
`
		},
		"app/middleware/requestId.go": func() string {
			return `package middleware

import (
    "context"
    "net/http"
    "github.com/google/uuid"
)

type key int

const requestIDKey key = 0

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := uuid.New().String()
        ctx := context.WithValue(r.Context(), requestIDKey, id)
        w.Header().Set("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func GetRequestID(r *http.Request) string {
    if id, ok := r.Context().Value(requestIDKey).(string); ok {
        return id
    }
    return ""
}
`
		},
		"app/middleware/auth.go": func() string {
			return `package middleware

import (
    "net/http"
)

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add your authentication logic here
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        // Validate the token
        // ...
        next.ServeHTTP(w, r)
    })
}
`
		},
		"app/middleware/cors.go": func() string {
			return `package middleware

import (
    "net/http"
)

func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
`
		},
		"app/middleware/logger.go": func() string {
			return `package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, time.Since(start))
    })
}
`
		},
		"app/middleware/notifier.go": func() string {
			return `package middleware

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/go-mail/mail"
)

var (
    clients       = make(map[chan string]bool)
    notifyChannel = make(chan string)
)

// Notifier middleware that sends SSE notifications to connected clients
func Notifier(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        flusher, ok := w.(http.Flusher)
        if !ok {
            http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
            return
        }

        ctx, cancel := context.WithCancel(r.Context())
        defer cancel()

        messageChan := make(chan string)
        clients[messageChan] = true

        defer func() {
            delete(clients, messageChan)
            close(messageChan)
        }()

        w.Header().Set("Content-Type", "text/event-stream")
        w.Header().Set("Cache-Control", "no-cache")
        w.Header().Set("Connection", "keep-alive")

        notifyChannel <- "System started"

        go func() {
            <-ctx.Done()
            cancel()
        }()

        for {
            select {
            case <-ctx.Done():
                return
            case msg := <-messageChan:
                fmt.Fprintf(w, "data: %s\n\n", msg)
                flusher.Flush()
            }
        }
    })
}

// Function to notify all connected SSE clients
func notifyClients(message string) {
    for client := range clients {
        client <- message
    }
}

// Function to notify users via email
func notifyByEmail(subject, body string) {
    m := mail.NewMessage()
    m.SetHeader("From", "your-email@example.com")
    m.SetHeader("To", "user@example.com") // add your recipient's email here
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    d := mail.NewDialer("smtp.example.com", 587, "your-username", "your-password")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := d.DialAndSendContext(ctx, m); err != nil {
        log.Printf("Failed to send email: %v", err)
    }
}
`
		},
	}

	return nil
}

func (g *GenMiddlewarePlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenMiddlewarePlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenMiddlewarePlugin) Name() string {
	return "GenMiddlewarePlugin"
}

func (g *GenMiddlewarePlugin) Version() string {
	return "1.0.0"
}

func (g *GenMiddlewarePlugin) Dependencies() []string {
	return []string{}
}

func (g *GenMiddlewarePlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenMiddlewarePlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenMiddlewarePlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenMiddlewarePlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/middleware"
}

func (g *GenMiddlewarePlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenMiddlewarePlugin(data config.ProjectData) *GenMiddlewarePlugin {
	return &GenMiddlewarePlugin{
		Data: data,
	}
}
