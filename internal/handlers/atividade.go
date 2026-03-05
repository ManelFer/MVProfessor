package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ManelFer/MVProfessor/internal/database"
	"github.com/ManelFer/MVProfessor/internal/models"
)

// CreateAtividade cria uma nova atividade vinculada ao professor autenticado
// Opcionalmente já associa alunos na mesma requisição
func CreateAtividade(c *gin.Context) {
	// Pega o user_id do JWT
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

	var input models.CreateAtividadeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ATIVIDADE] Iniciando criação | Nome: %s | Professor ID: %d | Alunos: %d", input.Nome, professorID, len(input.AlunoIDs))

	var atividadeID int
	err := database.DB.QueryRow(
		`INSERT INTO atividades (nome, descricao, link_acesso, professor_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())
		 RETURNING id`,
		input.Nome, input.Descricao, input.LinkAcesso, professorID,
	).Scan(&atividadeID)

	if err != nil {
		log.Printf("[ATIVIDADE] ❌ Erro ao inserir atividade | Erro: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao criar atividade"})
		return
	}

	log.Printf("[ATIVIDADE] ✅ Atividade criada com sucesso | ID: %d | Nome: %s", atividadeID, input.Nome)

	// Se tiver alunos para associar
	alunosInseridos := 0
	alunosErro := 0

	if len(input.AlunoIDs) > 0 {
		log.Printf("[ATIVIDADE] Associando %d aluno(s) à atividade ID: %d", len(input.AlunoIDs), atividadeID)

		for _, alunoID := range input.AlunoIDs {
			// Verifica se o aluno pertence ao professor
			var alunoCount int
			err = database.DB.QueryRow(
				`SELECT COUNT(*) FROM alunos WHERE id = $1 AND professor_id = $2`,
				alunoID, professorID,
			).Scan(&alunoCount)

			if err != nil || alunoCount == 0 {
				log.Printf("[ATIVIDADE] ⚠️ Aluno ID %d não encontrado ou não pertence ao professor", alunoID)
				alunosErro++
				continue
			}

			// Insere a associação
			_, err = database.DB.Exec(
				`INSERT INTO atividade_aluno (atividade_id, aluno_id, created_at)
				 VALUES ($1, $2, NOW())
				 ON CONFLICT DO NOTHING`,
				atividadeID, alunoID,
			)

			if err != nil {
				log.Printf("[ATIVIDADE] ⚠️ Erro ao associar aluno ID %d | Erro: %v", alunoID, err)
				alunosErro++
				continue
			}

			alunosInseridos++
			log.Printf("[ATIVIDADE] ✅ Aluno ID %d associado à atividade ID: %d", alunoID, atividadeID)
		}
	}

	resposta := gin.H{
		"id":          atividadeID,
		"nome":        input.Nome,
		"descricao":   input.Descricao,
		"link_acesso": input.LinkAcesso,
	}

	if len(input.AlunoIDs) > 0 {
		resposta["alunos_associados"] = alunosInseridos
		resposta["alunos_erro"] = alunosErro
		log.Printf("[ATIVIDADE] Resumo: %d aluno(s) associado(s) com sucesso, %d erro(s)", alunosInseridos, alunosErro)
	}

	c.JSON(http.StatusCreated, resposta)
}

// ListAtividadesProfessor lista todas as atividades do professor autenticado
func ListAtividadesProfessor(c *gin.Context) {
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

	log.Printf("[ATIVIDADE] Listando atividades do professor ID: %d", professorID)

	rows, err := database.DB.Query(
		`SELECT id, nome, descricao, link_acesso FROM atividades WHERE professor_id = $1 ORDER BY created_at DESC`,
		professorID,
	)

	if err != nil {
		log.Printf("[ATIVIDADE] ❌ Erro ao buscar atividades | Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar atividades"})
		return
	}
	defer rows.Close()

	atividades := []models.AtividadeResponse{}
	for rows.Next() {
		var atividade models.AtividadeResponse
		err := rows.Scan(&atividade.ID, &atividade.Nome, &atividade.Descricao, &atividade.LinkAcesso)
		if err != nil {
			log.Printf("[ATIVIDADE] Erro ao scanear atividade: %v", err)
			continue
		}
		atividades = append(atividades, atividade)
	}

	log.Printf("[ATIVIDADE] ✅ %d atividades encontradas para professor ID: %d", len(atividades), professorID)

	c.JSON(http.StatusOK, gin.H{"atividades": atividades})
}

// AssociarAlunoAtividade associa um aluno a uma atividade
func AssociarAlunoAtividade(c *gin.Context) {
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

	atividadeID := c.Param("atividade_id")
	var input models.AssociarAlunoAtividadeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ATIVIDADE] Associando aluno | Atividade ID: %s | Aluno ID: %d | Professor ID: %d", atividadeID, input.AlunoID, professorID)

	// Verifica se a atividade pertence ao professor e se o aluno pertence ao professor
	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM atividades 
		 WHERE id = $1 AND professor_id = $2`,
		atividadeID, professorID,
	).Scan(&count)

	if err != nil || count == 0 {
		log.Printf("[ATIVIDADE] ❌ Atividade não encontrada ou não pertence ao professor")
		c.JSON(http.StatusForbidden, gin.H{"error": "Atividade não encontrada"})
		return
	}

	// Verifica se o aluno pertence ao professor
	var alunoCount int
	err = database.DB.QueryRow(
		`SELECT COUNT(*) FROM alunos WHERE id = $1 AND professor_id = $2`,
		input.AlunoID, professorID,
	).Scan(&alunoCount)

	if err != nil || alunoCount == 0 {
		log.Printf("[ATIVIDADE] ❌ Aluno não encontrado ou não pertence ao professor")
		c.JSON(http.StatusForbidden, gin.H{"error": "Aluno não encontrado"})
		return
	}

	// Insere a associação
	_, err = database.DB.Exec(
		`INSERT INTO atividade_aluno (atividade_id, aluno_id, created_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT DO NOTHING`,
		atividadeID, input.AlunoID,
	)

	if err != nil {
		log.Printf("[ATIVIDADE] ❌ Erro ao associar aluno | Erro: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao associar aluno à atividade"})
		return
	}

	log.Printf("[ATIVIDADE] ✅ Aluno associado com sucesso | Atividade ID: %s | Aluno ID: %d", atividadeID, input.AlunoID)

	c.JSON(http.StatusOK, gin.H{"message": "Aluno associado à atividade com sucesso"})
}

// ListAtividadesAluno lista atividades de um aluno específico
func ListAtividadesAluno(c *gin.Context) {
	alunoID := c.Param("aluno_id")

	log.Printf("[ATIVIDADE] Listando atividades do aluno ID: %s", alunoID)

	rows, err := database.DB.Query(
		`SELECT a.id, a.nome, a.descricao, a.link_acesso
		 FROM atividades a
		 INNER JOIN atividade_aluno aa ON a.id = aa.atividade_id
		 WHERE aa.aluno_id = $1
		 ORDER BY a.created_at DESC`,
		alunoID,
	)

	if err != nil {
		log.Printf("[ATIVIDADE] ❌ Erro ao buscar atividades do aluno | Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar atividades"})
		return
	}
	defer rows.Close()

	atividades := []models.AtividadeResponse{}
	for rows.Next() {
		var atividade models.AtividadeResponse
		err := rows.Scan(&atividade.ID, &atividade.Nome, &atividade.Descricao, &atividade.LinkAcesso)
		if err != nil {
			log.Printf("[ATIVIDADE] Erro ao scanear atividade: %v", err)
			continue
		}
		atividades = append(atividades, atividade)
	}

	log.Printf("[ATIVIDADE] ✅ %d atividades encontradas para aluno ID: %s", len(atividades), alunoID)

	c.JSON(http.StatusOK, gin.H{"atividades": atividades})
}

// DeleteAtividade remove uma atividade (apenas professor pode)
func DeleteAtividade(c *gin.Context) {
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

	atividadeID := c.Param("atividade_id")

	log.Printf("[ATIVIDADE] Deletando atividade | ID: %s | Professor ID: %d", atividadeID, professorID)

	// Verifica se a atividade pertence ao professor
	var count int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM atividades WHERE id = $1 AND professor_id = $2`,
		atividadeID, professorID,
	).Scan(&count)

	if err != nil || count == 0 {
		log.Printf("[ATIVIDADE] ❌ Atividade não encontrada ou não pertence ao professor")
		c.JSON(http.StatusForbidden, gin.H{"error": "Atividade não encontrada"})
		return
	}

	// Deleta as associações primeiro (cascade)
	_, err = database.DB.Exec(
		`DELETE FROM atividade_aluno WHERE atividade_id = $1`,
		atividadeID,
	)

	if err != nil {
		log.Printf("[ATIVIDADE] ❌ Erro ao deletar associações | Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar atividade"})
		return
	}

	// Deleta a atividade
	_, err = database.DB.Exec(
		`DELETE FROM atividades WHERE id = $1 AND professor_id = $2`,
		atividadeID, professorID,
	)

	if err != nil {
		log.Printf("[ATIVIDADE] ❌ Erro ao deletar atividade | Erro: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar atividade"})
		return
	}

	log.Printf("[ATIVIDADE] ✅ Atividade deletada com sucesso | ID: %s", atividadeID)

	c.JSON(http.StatusOK, gin.H{"message": "Atividade deletada com sucesso"})
}
