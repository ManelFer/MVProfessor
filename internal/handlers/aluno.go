package handlers

import (
	"net/http"
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/ManelFer/MVProfessor/internal/database"
	"github.com/ManelFer/MVProfessor/internal/models"
	"github.com/ManelFer/MVProfessor/internal/utils"
)

func CreateAluno(c *gin.Context) {
	// Pegamos o user_id do JWT (vamos adicionar middleware depois)
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Não autorizado"})
		return
	}

	// Aqui você pode checar se é professor (depois no middleware ou query)
	// Por enquanto assumimos que quem chega aqui é professor

	var input models.CreateAlunoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gera senha aleatória
	senhaGerada, err := utils.GenerateSecurePassword(12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar senha"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(senhaGerada), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
		return
	}

	var alunoID int
	err = database.DB.QueryRow(
		`INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
		 VALUES ($1, $2, $3, 'aluno', NOW(), NOW())
		 RETURNING id`,
		input.Nome, input.Email, string(hashed),
	).Scan(&alunoID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email já existe ou erro no banco"})
		return
	}

	// Monta corpo do email
	appURL := os.Getenv("APP_URL")
	body := fmt.Sprintf(`
		<h2>Bem-vindo ao sistema, %s!</h2>
		<p>Seu acesso foi criado por um professor.</p>
		<p><strong>Email:</strong> %s</p>
		<p><strong>Senha temporária:</strong> %s</p>
		<p>Faça login em: <a href="%s/login">%s/login</a></p>
		<p>Recomendamos alterar a senha imediatamente após o primeiro acesso.</p>
		<br>
		Atenciosamente,<br>Sistema da Escola
	`, input.Nome, input.Email, senhaGerada, appURL, appURL)

	err = utils.SendEmail(input.Email, "Acesso ao Sistema - Credenciais", body)
	if err != nil {
		// Aqui você pode logar e continuar (ou rollback se quiser)
		log.Printf("Aluno criado, mas falhou ao enviar email: %v", err)
	}

	c.JSON(http.StatusCreated, models.AlunoResponse{
		ID:    alunoID,
		Nome:  input.Nome,
		Email: input.Email,
	})
}