package service

import (
	"context"

	"github.com/yeboahd24/workout-tracker/model"
	"github.com/yeboahd24/workout-tracker/repository"
)

type ExerciseService struct {
	exerciseRepo *repository.ExerciseRepository
}

func NewExerciseService(exerciseRepo *repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{exerciseRepo: exerciseRepo}
}

func (s *ExerciseService) CreateExercise(ctx context.Context, name, description, category string) (*model.Exercise, error) {
	exercise := model.NewExercise(name, description, category)
	err := s.exerciseRepo.Create(ctx, exercise)
	if err != nil {
		return nil, err
	}
	return exercise, nil
}

func (s *ExerciseService) GetExerciseByID(ctx context.Context, id int) (*model.Exercise, error) {
	return s.exerciseRepo.GetByID(ctx, id)
}

func (s *ExerciseService) GetAllExercises(ctx context.Context) ([]*model.Exercise, error) {
	return s.exerciseRepo.GetAll(ctx)
}
