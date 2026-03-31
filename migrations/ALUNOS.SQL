-- Adicionar coluna role na tabela users se não existir
ALTER TABLE users
ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'professor';

-- Criar tabela alunos vinculada ao professor
CREATE TABLE IF NOT EXISTS alunos (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL,
    professor_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (professor_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(email, professor_id)
);

-- Criar índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_alunos_professor_id ON alunos(professor_id);
CREATE INDEX IF NOT EXISTS idx_alunos_email ON alunos(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
