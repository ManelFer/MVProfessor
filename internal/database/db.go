package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // driver PostgreSQL
)

var DB *sql.DB

// Connect estabelece a conexão com o PostgreSQL usando as variáveis de ambiente
func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=America/Sao_Paulo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com PostgreSQL: %v", err)
	}

	// Configurações importantes para pool de conexões (evita problemas em produção)
	DB.SetMaxOpenConns(25)                 // máximo de conexões abertas simultâneas
	DB.SetMaxIdleConns(10)                 // conexões ociosas mantidas
	DB.SetConnMaxLifetime(5 * time.Minute) // tempo máximo de vida de uma conexão
	DB.SetConnMaxIdleTime(1 * time.Minute) // tempo máximo ocioso antes de fechar

	// Testa a conexão de verdade
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Falha no Ping ao PostgreSQL: %v", err)
	}

	log.Println("Conexão com PostgreSQL estabelecida com sucesso")
}

// Close fecha a conexão com o banco (deve ser chamado com defer no main)
func Close() error {
	if DB == nil {
		return nil
	}
	err := DB.Close()
	if err != nil {
		log.Printf("Erro ao fechar conexão com banco: %v", err)
		return err
	}
	log.Println("Conexão com PostgreSQL fechada")
	return nil
}
