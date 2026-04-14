-- Ensure extension exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================
-- SEED DATA (user + project + tasks)
-- =========================

WITH new_user AS (
  INSERT INTO users (name, email, password)
  VALUES ('Test User', 'test@example.com', 'password123')
  RETURNING id
),
new_project AS (
  INSERT INTO projects (name, description, owner_id)
  SELECT 'Demo Project', 'This is a seeded project', id
  FROM new_user
  RETURNING id, owner_id
)
INSERT INTO tasks (title, description, status, priority, project_id, assignee_id, due_date)
SELECT
  'First Task',
  'Seeded task for testing',
  'todo',
  'medium',
  new_project.id,
  new_project.owner_id,
  CURRENT_DATE + INTERVAL '7 days'
FROM new_project;

-- Extra tasks
INSERT INTO tasks (title, description, status, priority, project_id)
SELECT
  'Second Task',
  'Another seeded task',
  'in_progress',
  'high',
  p.id
FROM projects p
WHERE p.name = 'Demo Project';