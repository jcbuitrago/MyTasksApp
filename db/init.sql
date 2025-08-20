-- Limpia si estás iterando local
-- DROP TYPE IF EXISTS task_status;
-- DROP TABLE IF EXISTS tasks, categories, users CASCADE;

-- 1) Enum para estado de tareas (coincide con la guía)
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'task_status') THEN
        CREATE TYPE task_status AS ENUM ('Sin Empezar', 'Empezada', 'Finalizada');
    END IF;
END $$;

-- 2) Usuarios
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    picture_url     TEXT
);

-- 3) Categorías (globales; si prefieres por usuario, añade user_id FK)
CREATE TABLE IF NOT EXISTS categories (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(50) UNIQUE NOT NULL,
    description     TEXT
);

-- 4) Tareas
CREATE TABLE IF NOT EXISTS tasks (
    id                  SERIAL PRIMARY KEY,
    description         TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tentative_due_date  DATE,
    status              task_status NOT NULL DEFAULT 'Sin Empezar',
    category_id         INT REFERENCES categories(id) ON DELETE SET NULL,
    user_id             INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- Datos de ejemplo (opcional)
INSERT INTO categories (name, description) VALUES
('Hogar','Tareas de casa'),('Trabajo','Tareas de la oficina')
ON CONFLICT DO NOTHING;

-- Constraint del estado (si ya existe, ignora el error)
DO $$ BEGIN
  ALTER TABLE tasks ADD CONSTRAINT chk_tasks_status
    CHECK (status IN ('Sin Empezar','Empezada','Finalizada'));
EXCEPTION WHEN duplicate_object THEN NULL; END $$;

-- Índices útiles
CREATE INDEX IF NOT EXISTS idx_tasks_user     ON tasks(user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_category ON tasks(category_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status   ON tasks(status);