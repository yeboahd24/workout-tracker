package repository

import (
	"context"
	"database/sql"

	"github.com/yeboahd24/workout-tracker/model"
)



type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) Create(ctx context.Context, exercise *model.Exercise) error {
	query := `
		INSERT INTO exercises (name, description, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		exercise.Name, exercise.Description, exercise.Category, exercise.CreatedAt, exercise.UpdatedAt,
	).Scan(&exercise.ID)

	return err
}

func (r *ExerciseRepository) GetByID(ctx context.Context, id int) (*model.Exercise, error) {
	query := `
		SELECT id, name, description, category, created_at, updated_at
		FROM exercises
		WHERE id = $1`

	var exercise model.Exercise
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&exercise.ID, &exercise.Name, &exercise.Description, &exercise.Category,
		&exercise.CreatedAt, &exercise.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (r *ExerciseRepository) GetAll(ctx context.Context) ([]*model.Exercise, error) {
	query := `
		SELECT id, name, description, category, created_at, updated_at
		FROM exercises
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []*model.Exercise
	for rows.Next() {
		var exercise model.Exercise
		err := rows.Scan(
			&exercise.ID, &exercise.Name, &exercise.Description, &exercise.Category,
			&exercise.CreatedAt, &exercise.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, &exercise)
	}

	return exercises, nil
}
