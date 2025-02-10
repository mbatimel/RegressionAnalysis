package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mbatimel/RegressionAnalysis/internal/config"
)

type HealthServer struct {
	server *fiber.App
}

func NewHealthServer() *HealthServer {
	health := HealthServer{
		server: fiber.New(),
	}
	health.registerHandlers()

	return &health
}

func (h *HealthServer) registerHandlers() {
	h.server.Get("/liveness", probesHandler)
	h.server.Get("/readiness", probesHandler)
	h.server.Get("/build-info", buildInfoHandler)
}

func probesHandler(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func buildInfoHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(BuildInfo{
		ServiceName: config.ServiceName(),
		GitSHA:      config.GitSHA(),
		Version:     config.Version(),
		BuildStamp:  config.BuildStamp(),
		BuildNumber: config.BuildNumber(),
		NodeName:    config.NodeName(),
	})
}

func (h *HealthServer) Start(bindURL string) error {
	return h.server.Listen(bindURL)
}

func (h *HealthServer) Stop() error {
	return h.server.Shutdown()
}

type BuildInfo struct {
	ServiceName string `json:"serviceName"`
	GitSHA      string `json:"gitSha"`
	Version     string `json:"version"`
	BuildStamp  string `json:"buildStamp"`
	BuildNumber string `json:"buildNumber"`
	NodeName    string `json:"nodeName"`
}
