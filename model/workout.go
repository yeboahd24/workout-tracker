package model

import "time"

type Workout struct {
	ID           int               `json:"id"`
	UserID       int               `json:"user_id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	ScheduledFor time.Time         `json:"scheduled_for"`
	Exercises    []WorkoutExercise `json:"exercises"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type WorkoutExercise struct {
	ID         int     `json:"id"`
	WorkoutID  int     `json:"workout_id"`
	ExerciseID int     `json:"exercise_id"`
	Sets       int     `json:"sets"`
	Reps       int     `json:"reps"`
	Weight     float64 `json:"weight"`
	Notes      string  `json:"notes"`
}

func NewWorkout(userID int, name, description string, scheduledFor time.Time) *Workout {
	return &Workout{
		UserID:       userID,
		Name:         name,
		Description:  description,
		ScheduledFor: scheduledFor,
		Exercises:    make([]WorkoutExercise, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (w *Workout) AddExercise(exerciseID, sets, reps int, weight float64, notes string) {
	w.Exercises = append(w.Exercises, WorkoutExercise{
		ExerciseID: exerciseID,
		Sets:       sets,
		Reps:       reps,
		Weight:     weight,
		Notes:      notes,
	})
}
