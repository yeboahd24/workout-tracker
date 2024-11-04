package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/yeboahd24/workout-tracker/model"
	"github.com/yeboahd24/workout-tracker/repository"
	"github.com/yeboahd24/workout-tracker/util"
)

type WorkoutHandler struct {
	workoutRepo *repository.WorkoutRepository
}

func NewWorkoutHandler(workoutRepo *repository.WorkoutRepository) *WorkoutHandler {
	return &WorkoutHandler{workoutRepo: workoutRepo}
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		ScheduledFor time.Time `json:"scheduled_for"`
		Exercises    []struct {
			ExerciseID int     `json:"exercise_id"`
			Sets       int     `json:"sets"`
			Reps       int     `json:"reps"`
			Weight     float64 `json:"weight"`
			Notes      string  `json:"notes"`
		} `json:"exercises"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workout := model.NewWorkout(userID, input.Name, input.Description, input.ScheduledFor)
	for _, e := range input.Exercises {
		workout.AddExercise(e.ExerciseID, e.Sets, e.Reps, e.Weight, e.Notes)
	}

	if err := h.workoutRepo.Create(r.Context(), workout); err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

func (h *WorkoutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(workout)
}

func (h *WorkoutHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workouts, err := h.workoutRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch workouts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(workouts)
}

func (h *WorkoutHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		ID           int       `json:"id"`
		Name         string    `json:"name"`
		Description  string    `json:"description"`
		ScheduledFor time.Time `json:"scheduled_for"`
		Exercises    []struct {
			ExerciseID int     `json:"exercise_id"`
			Sets       int     `json:"sets"`
			Reps       int     `json:"reps"`
			Weight     float64 `json:"weight"`
			Notes      string  `json:"notes"`
		} `json:"exercises"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workout, err := h.workoutRepo.GetByID(r.Context(), input.ID)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if workout.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workout.Name = input.Name
	workout.Description = input.Description
	workout.ScheduledFor = input.ScheduledFor
	workout.Exercises = make([]model.WorkoutExercise, len(input.Exercises))
	for i, e := range input.Exercises {
		workout.Exercises[i] = model.WorkoutExercise{
			ExerciseID: e.ExerciseID,
			Sets:       e.Sets,
			Reps:       e.Reps,
			Weight:     e.Weight,
			Notes:      e.Notes,
		}
	}

	if err := h.workoutRepo.Update(r.Context(), workout); err != nil {
		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(workout)
}

func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Workout not found", http.StatusNotFound)
		return
	}

	if workout.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.workoutRepo.Delete(r.Context(), id); err != nil {
		log.Printf("Error deleting workout: %v", err)
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WorkoutHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	userID, err := util.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters for report customization
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Generate the report
	report, err := h.workoutRepo.GenerateReport(r.Context(), userID, startDate, endDate)
	if err != nil {
		log.Printf("Error generating report: %v", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	// Send the report as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
