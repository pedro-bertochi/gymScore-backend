package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gynScore-backend/internal/config"
	"gynScore-backend/internal/controllers"
	"gynScore-backend/internal/middlewares"
	"gynScore-backend/internal/repositories"
	"gynScore-backend/internal/routes"
	"gynScore-backend/internal/services"
)

func main() {
	// ─── Carregamento de configurações ───────────────────────────────────────────
	cfg := config.Load()

	// ─── Conexão com o banco de dados ────────────────────────────────────────────
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Não foi possível conectar ao banco de dados: %v", err)
	}

	// Auto-migração removida conforme solicitação do usuário.
	// O banco de dados já possui as tabelas e procedures necessárias.
	log.Println("[DB] Conexão com o banco de dados concluída (Auto-migração desativada)")

	// ─── Injeção de dependências ──────────────────────────────────────────────────

	// Repositories
	usuarioRepo := repositories.NovoUsuarioRepository(db)
	desafioRepo := repositories.NovoDesafioRepository(db)
	amizadeRepo := repositories.NovoAmizadeRepository(db)

	// Services
	usuarioSvc := services.NovoUsuarioService(usuarioRepo)
	desafioSvc := services.NovoDesafioService(desafioRepo, usuarioRepo)
	amizadeSvc := services.NovoAmizadeService(amizadeRepo, usuarioRepo)

	// Controllers
	usuarioCtrl := controllers.NovoUsuarioController(usuarioSvc, cfg)
	desafioCtrl := controllers.NovoDesafioController(desafioSvc)
	amizadeCtrl := controllers.NovoAmizadeController(amizadeSvc)

	// ─── Configuração do servidor Fiber ──────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName:      "GymScore API v1.0.0",
		ErrorHandler: errorHandler,
	})

	// Middlewares globais
	app.Use(middlewares.RecoverMiddleware())
	app.Use(middlewares.LoggerMiddleware())
	app.Use(middlewares.CORSMiddleware())

	// Registro das rotas
	routes.Setup(app, cfg, usuarioCtrl, desafioCtrl, amizadeCtrl)

	// ─── Inicialização do servidor ────────────────────────────────────────────────
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("[SERVER] GymScore API iniciando na porta %s (ambiente: %s)", cfg.AppPort, cfg.AppEnv)

	// Graceful shutdown: aguarda sinal de interrupção antes de encerrar
	go func() {
		if err := app.Listen(addr); err != nil {
			log.Fatalf("[FATAL] Erro ao iniciar servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("[SERVER] Encerrando servidor graciosamente...")
	if err := app.Shutdown(); err != nil {
		log.Printf("[SERVER] Erro ao encerrar servidor: %v", err)
	}
	log.Println("[SERVER] Servidor encerrado com sucesso")
}

// errorHandler é o handler global de erros do Fiber
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Erro interno do servidor"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
