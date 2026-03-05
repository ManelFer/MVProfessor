package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Aluno struct {
	ID          int       `json:"id"`
	Nome        string    `json:"nome"`
	Email       string    `json:"email"`
	ProfessorID int       `json:"professor_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

//dto para cadastro
type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateAlunoInput struct {
	Nome  string `json:"nome" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type AlunoResponse struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

//dto para login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

//resposta token
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Modelo de Atividade
type Atividade struct {
	ID          int       `json:"id"`
	Nome        string    `json:"nome"`
	Descricao   string    `json:"descricao"`
	LinkAcesso  string    `json:"link_acesso"`
	ProfessorID int       `json:"professor_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DTO para criar atividade
type CreateAtividadeInput struct {
	Nome       string `json:"nome" binding:"required"`
	Descricao  string `json:"descricao" binding:"required"`
	LinkAcesso string `json:"link_acesso" binding:"required,url"`
	AlunoIDs   []int  `json:"aluno_ids"` // Opcional - pode deixar vazio e adicionar depois
}

// DTO para associar aluno a uma atividade
type AssociarAlunoAtividadeInput struct {
	AlunoID int `json:"aluno_id" binding:"required"`
}

// Resposta da atividade
type AtividadeResponse struct {
	ID         int    `json:"id"`
	Nome       string `json:"nome"`
	Descricao  string `json:"descricao"`
	LinkAcesso string `json:"link_acesso"`
}

// Atividade com alunos associados
type AtividadeComAlunosResponse struct {
	Atividade AtividadeResponse `json:"atividade"`
	Alunos    []AlunoResponse   `json:"alunos"`
}
