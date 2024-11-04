package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yeboahd24/workout-tracker/model"
)

type WorkoutRepository struct {
	db *sql.DB
}

func NewWorkoutRepository(db *sql.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (r *WorkoutRepository) Create(ctx context.Context, workout *model.Workout) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert workout
	query := `
		INSERT INTO workouts (user_id, name, description, scheduled_for, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		workout.UserID, workout.Name, workout.Description, workout.ScheduledFor,
		workout.CreatedAt, workout.UpdatedAt,
	).Scan(&workout.ID)
	if err != nil {
		return err
	}

	// Insert workout exercises
	for _, exercise := range workout.Exercises {
		query := `
			INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight, notes)
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.ExecContext(ctx, query,
			workout.ID, exercise.ExerciseID, exercise.Sets, exercise.Reps, exercise.Weight, exercise.Notes,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *WorkoutRepository) GetByID(ctx context.Context, id int) (*model.Workout, error) {
	query := `
		SELECT w.id, w.user_id, w.name, w.description, w.scheduled_for, w.created_at, w.updated_at,
			   we.id, we.exercise_id, we.sets, we.reps, we.weight, we.notes
		FROM workouts w
		LEFT JOIN workout_exercises we ON w.id = we.workout_id
		WHERE w.id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workout *model.Workout
	for rows.Next() {
		if workout == nil {
			workout = &model.Workout{}
		}

		var we model.WorkoutExercise
		err := rows.Scan(
			&workout.ID, &workout.UserID, &workout.Name, &workout.Description, &workout.ScheduledFor,
			&workout.CreatedAt, &workout.UpdatedAt,
			&we.ID, &we.ExerciseID, &we.Sets, &we.Reps, &we.Weight, &we.Notes,
		)
		if err != nil {
			return nil, err
		}

		if we.ID != 0 {
			workout.Exercises = append(workout.Exercises, we)
		}
	}

	return workout, nil
}

func (r *WorkoutRepository) GetByUserID(ctx context.Context, userID int) ([]*model.Workout, error) {
	query := `
		SELECT id, user_id, name, description, scheduled_for, created_at, updated_at
		FROM workouts
		WHERE user_id = $1
		ORDER BY scheduled_for DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []*model.Workout
	for rows.Next() {
		var w model.Workout
		err := rows.Scan(
			&w.ID, &w.UserID, &w.Name, &w.Description, &w.ScheduledFor,
			&w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, &w)
	}

	return workouts, nil
}

func (r *WorkoutRepository) Update(ctx context.Context, workout *model.Workout) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update workout
	query := `
		UPDATE workouts
		SET name = $1, description = $2, scheduled_for = $3, updated_at = $4
		WHERE id = $5`

	_, err = tx.ExecContext(ctx, query,
		workout.Name, workout.Description, workout.ScheduledFor, workout.UpdatedAt, workout.ID,
	)
	if err != nil {
		return err
	}

	// Delete existing workout exercises
	_, err = tx.ExecContext(ctx, "DELETE FROM workout_exercises WHERE workout_id = $1", workout.ID)
	if err != nil {
		return err
	}

	// Insert updated workout exercises
	for _, exercise := range workout.Exercises {
		query := `
			INSERT INTO workout_exercises (workout_id, exercise_id, sets, reps, weight, notes)
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.ExecContext(ctx, query,
			workout.ID, exercise.ExerciseID, exercise.Sets, exercise.Reps, exercise.Weight, exercise.Notes,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *WorkoutRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM workouts WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *WorkoutRepository) GenerateReport(ctx context.Context, userID int, startDate, endDate string) (map[string]interface{}, error) {
	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	// Fetch workouts within the date range
	query := `
        SELECT w.id, w.name, w.scheduled_for, we.exercise_id, we.sets, we.reps, we.weight
        FROM workouts w
        JOIN workout_exercises we ON w.id = we.workout_id
        WHERE w.user_id = $1 AND w.scheduled_for BETWEEN $2 AND $3
        ORDER BY w.scheduled_for
    `
	rows, err := r.db.QueryContext(ctx, query, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	workouts := make(map[int]map[string]interface{})
	var totalWorkouts, totalExercises int
	for rows.Next() {
		var workoutID int
		var workoutName string
		var scheduledFor time.Time
		var exerciseID, sets, reps int
		var weight float64

		err := rows.Scan(&workoutID, &workoutName, &scheduledFor, &exerciseID, &sets, &reps, &weight)
		if err != nil {
			return nil, err
		}

		if _, exists := workouts[workoutID]; !exists {
			workouts[workoutID] = map[string]interface{}{
				"name":          workoutName,
				"scheduled_for": scheduledFor,
				"exercises":     []map[string]interface{}{},
			}
			totalWorkouts++
		}

		workouts[workoutID]["exercises"] = append(workouts[workoutID]["exercises"].([]map[string]interface{}), map[string]interface{}{
			"exercise_id": exerciseID,
			"sets":        sets,
			"reps":        reps,
			"weight":      weight,
		})
		totalExercises++
	}

	// Prepare the final report
	report := map[string]interface{}{
		"start_date":      startDate,
		"end_date":        endDate,
		"total_workouts":  totalWorkouts,
		"total_exercises": totalExercises,
		"workouts":        workouts,
	}

	return report, nil
}
