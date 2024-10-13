package service

import (
	"context"
	"time"

	"github.com/yeboahd24/workout-tracker/model"
	"github.com/yeboahd24/workout-tracker/repository"
)

type WorkoutService struct {
	workoutRepo *repository.WorkoutRepository
}

func NewWorkoutService(workoutRepo *repository.WorkoutRepository) *WorkoutService {
	return &WorkoutService{workoutRepo: workoutRepo}
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, userID int, name, description string, scheduledFor time.Time, exercises []model.WorkoutExercise) (*model.Workout, error) {
	workout := model.NewWorkout(userID, name, description, scheduledFor)
	workout.Exercises = exercises
	err := s.workoutRepo.Create(ctx, workout)
	if err != nil {
		return nil, err
	}
	return workout, nil
}

func (s *WorkoutService) GetWorkoutByID(ctx context.Context, id int) (*model.Workout, error) {
	return s.workoutRepo.GetByID(ctx, id)
}

func (s *WorkoutService) GetWorkoutsByUserID(ctx context.Context, userID int) ([]*model.Workout, error) {
	return s.workoutRepo.GetByUserID(ctx, userID)
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, workout *model.Workout) error {
	return s.workoutRepo.Update(ctx, workout)
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, id int) error {
	return s.workoutRepo.Delete(ctx, id)
}
