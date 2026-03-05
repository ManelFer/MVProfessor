package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/ManelFer/MVProfessor/internal/database"
	"github.com/ManelFer/MVProfessor/internal/handlers"
	"github.com/ManelFer/MVProfessor/internal/middleware"
)

func main() {

	// Carrega .env PRIMEIRO
	paths := []string{
		".env",
		"../../.env",
		"../../../.env",
	}
	loaded := false
	for _, path := range paths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Arquivo .env carregado de: %s", path)
			loaded = true
			break
		}
	}
	if !loaded {
		log.Println("Aviso: .env não encontrado (usando variáveis de ambiente do sistema)")
	}

	// Depois conecta ao banco
	database.Connect()
	defer database.Close()

	r := gin.Default()

	// Rotas públicas
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	//rota privada
	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		private.POST("/alunos", handlers.CreateAluno)
	}

	// Teste simples protegido (depois vamos criar middleware)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s", port)
	r.Run(":" + port)
}
