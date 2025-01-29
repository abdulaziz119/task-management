package tasks

import (
	"context"
	"encoding/json"
)

type TaskStats struct {
	PendingTasks    int        `json:"pending_tasks"`
	InProgressTasks int        `json:"in_progress_tasks"`
	CompletedTasks  int        `json:"completed_tasks"`
	TotalTasks      int        `json:"total_tasks"`
	Tasks           []TaskItem `json:"tasks"`
}

type TaskItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserTasks(ctx context.Context, userId int) (*TaskStats, error) {
	query := `
		SELECT 
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_tasks,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_tasks,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
			COUNT(*) as total_tasks,
			json_agg(
				json_build_object(
					'id', t.id,
					'name', t.name,
					'description', t.description,
					'status', t.status,
					'priority', t.priority,
					'due_date', t.due_date at time zone 'UTC',
					'created_at', t.created_at at time zone 'UTC'
				) ORDER BY t.created_at DESC
			) as tasks
		FROM tasks t
		WHERE t.deleted_at IS NULL AND t.assigned_to = ?
		GROUP BY assigned_to`

	var stats TaskStats
	var tasksJson []byte

	err := r.db.QueryRowContext(ctx, query, userId).Scan(
		&stats.PendingTasks,
		&stats.InProgressTasks,
		&stats.CompletedTasks,
		&stats.TotalTasks,
		&tasksJson,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting user tasks: %v", err)
	}

	if err := json.Unmarshal(tasksJson, &stats.Tasks); err != nil {
		return nil, fmt.Errorf("error parsing tasks: %v", err)
	}

	return &stats, nil
}
