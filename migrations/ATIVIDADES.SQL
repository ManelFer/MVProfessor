-- Criar tabela de atividades
CREATE TABLE IF NOT EXISTS atividades (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT NOT NULL,
    link_acesso VARCHAR(500) NOT NULL,
    professor_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (professor_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Criar tabela de associação entre atividades e alunos
CREATE TABLE IF NOT EXISTS atividade_aluno (
    id SERIAL PRIMARY KEY,
    atividade_id INTEGER NOT NULL,
    aluno_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (atividade_id) REFERENCES atividades(id) ON DELETE CASCADE,
    FOREIGN KEY (aluno_id) REFERENCES alunos(id) ON DELETE CASCADE,
    UNIQUE(atividade_id, aluno_id)
);

-- Criar índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_atividades_professor_id ON atividades(professor_id);
CREATE INDEX IF NOT EXISTS idx_atividade_aluno_atividade_id ON atividade_aluno(atividade_id);
CREATE INDEX IF NOT EXISTS idx_atividade_aluno_aluno_id ON atividade_aluno(aluno_id);
