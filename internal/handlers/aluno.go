package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/ManelFer/MVProfessor/internal/database"
	"github.com/ManelFer/MVProfessor/internal/models"
	"github.com/ManelFer/MVProfessor/internal/utils"
)

func CreateAluno(c *gin.Context) {
	// Pegamos o user_id do JWT
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	professorID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuário inválido"})
		return
	}

	var input models.CreateAlunoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ALUNO] Iniciando criação de aluno | Nome: %s | Email: %s | Professor ID: %d", input.Nome, input.Email, professorID)

	var alunoID int
	err := database.DB.QueryRow(
		`INSERT INTO alunos (nome, email, professor_id, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())
		 RETURNING id`,
		input.Nome, input.Email, professorID,
	).Scan(&alunoID)

	if err != nil {
		log.Printf("[ALUNO] ❌ Erro ao inserir aluno no banco de dados | Email: %s | Erro: %v", input.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email já existe para este professor ou erro no banco"})
		return
	}

	log.Printf("[ALUNO] ✅ Aluno criado com sucesso no banco | Aluno ID: %d | Email: %s", alunoID, input.Email)

	// Gera senha aleatória para o aluno
	senhaGerada, err := utils.GenerateSecurePassword(12)
	if err != nil {
		log.Printf("[SENHA] ❌ Erro ao gerar senha para aluno ID: %d | Erro: %v", alunoID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar senha"})
		return
	}
	log.Printf("[SENHA] ✅ Senha gerada com sucesso para aluno ID: %d", alunoID)

	// Hash da senha para guardar (caso queira salvar no futuro)
	_, _ = bcrypt.GenerateFromPassword([]byte(senhaGerada), bcrypt.DefaultCost)

	// Monta corpo do email com credenciais
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:8080"
	}

	emailBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<style>
		body { font-family: Arial, sans-serif; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #F3EFEA; color: #0F172A; padding: 20px; border-radius: 5px 5px 0 0; }
		.content { background-color: #f9f9f9; padding: 20px; }
		.credentials { background-color: #F3EFEA; padding: 15px; border-left: 4px solid #c7a67d; margin: 20px 0; }
		.button { display: inline-block; background-color: #1E40AF; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin-top: 20px; }
		.footer { background-color: #F3EFEA; color: #0F172A; padding: 15px; text-align: center; border-radius: 0 0 5px 5px; font-size: 12px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h2>Bem-vindo ao Sistema BEGIN!</h2>
		</div>
		<div class="content">
			<p>Olá <strong>%s</strong>,</p>
			<p>Sua conta foi criada com sucesso por um professor. Segue abaixo suas credenciais de acesso:</p>
			
			<div class="credentials">
				<h3>Suas Credenciais</h3>
				<p><strong>Email:</strong> %s</p>
				<p><strong>Senha Temporária:</strong> <code style="background: #fff; padding: 2px 5px;">%s</code></p>
			</div>

			<p><strong>⚠️ Importante:</strong> Esta é uma senha unica. Recomendamos compartilhar.</p>

			<a href="%s/login" class="button">Acessar Sistema</a>

			<p style="margin-top: 30px; color: #555;">Se tiver dúvidas, entre em contato com seu professor.</p>
		</div>
		<div class="footer">
			<p>&copy; 2026 Begin - Sistema de Gestão de Alunos</p>
		</div>
	</div>
</body>
</html>
	`, input.Nome, input.Email, senhaGerada, appURL)

	// Envia email
	log.Printf("[EMAIL] Iniciando envio de email para aluno ID: %d, Email: %s, Nome: %s", alunoID, input.Email, input.Nome)
	err = utils.SendEmail(input.Email, "🎓 Bem-vindo ao MVProfessor - Suas Credenciais de Acesso", emailBody)
	if err != nil {
		log.Printf("[EMAIL] ❌ FALHA ao enviar email para %s | Aluno ID: %d | Erro: %v", input.Email, alunoID, err)
		// Não retorna erro, pois o aluno foi criado no banco
	} else {
		log.Printf("[EMAIL] ✅ Sucesso ao enviar email para %s | Aluno ID: %d | Nome: %s", input.Email, alunoID, input.Nome)
	}

	log.Printf("[ALUNO] ✅✅ Aluno criado COMPLETO | ID: %d | Email: %s | Status: Pronto para usar", alunoID, input.Email)

	c.JSON(http.StatusCreated, models.AlunoResponse{
		ID:    alunoID,
		Nome:  input.Nome,
		Email: input.Email,
	})
}

// ListAlunosProfessor retorna todos os alunos do professor autenticado
func ListAlunosProfessor(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	professorID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuário inválido"})
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, nome, email FROM alunos WHERE professor_id = $1 ORDER BY nome ASC`,
		professorID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar alunos"})
		return
	}
	defer rows.Close()

	alunos := []models.AlunoResponse{}
	for rows.Next() {
		var aluno models.AlunoResponse
		err := rows.Scan(&aluno.ID, &aluno.Nome, &aluno.Email)
		if err != nil {
			log.Printf("Erro ao scanear aluno: %v", err)
			continue
		}
		alunos = append(alunos, aluno)
	}

	c.JSON(http.StatusOK, gin.H{"alunos": alunos})
}

// DeleteAluno remove um aluno, apenas se pertencer ao professor autenticado
func DeleteAluno(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	professorID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuário inválido"})
		return
	}

	alunoID := c.Param("id")

	// Verifica se o aluno pertence ao professor
	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM alunos WHERE id = $1 AND professor_id = $2`,
		alunoID, professorID,
	).Scan(&count)

	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Aluno não encontrado ou não pertence a você"})
		return
	}

	_, err = database.DB.Exec(
		`DELETE FROM alunos WHERE id = $1 AND professor_id = $2`,
		alunoID, professorID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar aluno"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Aluno deletado com sucesso"})
}
