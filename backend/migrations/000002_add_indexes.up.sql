CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

CREATE INDEX IF NOT EXISTS idx_projects_owner_id_created_at ON projects(owner_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tasks_project_id_status ON tasks(project_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_project_id_assignee_id ON tasks(project_id, assignee_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id_status ON tasks(assignee_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date) WHERE due_date IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_status_priority ON tasks(status, priority);
CREATE INDEX IF NOT EXISTS idx_tasks_updated_at ON tasks(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id_project_id ON tasks(assignee_id, project_id);
