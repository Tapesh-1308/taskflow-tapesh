DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_users_created_at;

DROP INDEX IF EXISTS idx_projects_owner_id_created_at;
DROP INDEX IF EXISTS idx_projects_created_at;

DROP INDEX IF EXISTS idx_tasks_project_id_status;
DROP INDEX IF EXISTS idx_tasks_project_id_assignee_id;
DROP INDEX IF EXISTS idx_tasks_assignee_id_status;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_status_priority;
DROP INDEX IF EXISTS idx_tasks_updated_at;
DROP INDEX IF EXISTS idx_tasks_assignee_id_project_id;