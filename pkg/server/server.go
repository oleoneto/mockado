package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	flogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/oleoneto/mockado/pkg/envloader"
	"github.com/oleoneto/mockado/pkg/helpers"
	"github.com/oleoneto/mockado/pkg/httpclient"
)

type ServerConfig func(*Server)

type Server struct {
	app       *fiber.App
	varloader *envloader.EnvLoader
}

var once sync.Once

var instance *Server = &Server{}

func New(options ...ServerConfig) *Server {
	once.Do(func() {
		app := fiber.New()

		app.Use(flogger.New(
			flogger.Config{
				ForceColors: true,
				Format:      flogger.JSONFormat,
				TimeZone:    "America/Sao_Paulo",
			},
		))

		app.Use(recover.New(
			recover.Config{EnableStackTrace: false},
		))

		instance.app = app

		for _, fn := range options {
			fn(instance)
		}
	})

	return instance
}

func WithRoutes(routes ...httpclient.Route) ServerConfig {
	return func(s *Server) {
		for _, r := range routes {
			s.app.Add(r.Methods, r.Path, s.setupRoute(r))
		}
	}
}

func WithEnvLoader(el *envloader.EnvLoader) ServerConfig {
	return func(s *Server) { s.varloader = el }
}

func (s *Server) setupRoute(r httpclient.Route) func(c fiber.Ctx) error {
	fn := func(c fiber.Ctx) error {
		c.Response().Header.Add("X-Request-Time", time.Now().String())

		if r.Accepts != "" {
			c.Response().Header.Add("Accepts", r.Accepts)
		}

		for key, value := range r.Response.Headers {
			c.Response().Header.Add(key, s.varloader.LoadVar(value))
		}

		latency := helpers.UnwrapOrDefault(r.Response.Latency, 0)
		c.Response().Header.Add("X-Response-Latency", latency.String())

		var variables map[string]string
		if r.Response.EnvAsBody {
			variables = s.varloader.Environment()
		}

		time.Sleep(latency)

		c.Response().Header.Add("X-Response-Time", time.Now().String())

		c.Status(r.Response.StatusCode)
		if r.Response.EnvAsBody {
			return c.JSON(variables)
		}

		if r.Response.Body != nil {
			if b, ok := r.Response.Body.([]byte); ok {
				return c.Send(b)
			}
		}

		return c.Send([]byte(``))
	}

	return fn
}

func (s *Server) App() *fiber.App { return s.app }

func (s *Server) Listen(port int, showRoutes bool) error {
	address := "127.0.0.1"

	s.app.Hooks().OnListen(func(d fiber.ListenData) error {
		fmt.Printf("[Server] Listening on %s:%s [pid: %d]\n", address, d.Port, d.PID)
		return nil
	})

	s.app.Hooks().OnPreShutdown(func() error {
		fmt.Println("[Server] Shutting down")
		return nil
	})

	return s.app.Listen(
		fmt.Sprintf("%s:%d", address, port),
		fiber.ListenConfig{
			EnablePrintRoutes:     showRoutes,
			DisableStartupMessage: true,
		},
	)
}
