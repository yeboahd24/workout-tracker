package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yeboahd24/workout-tracker/model"
	"github.com/yeboahd24/workout-tracker/repository"
)

type ExerciseHandler struct {
	exerciseRepo *repository.ExerciseRepository
}

func NewExerciseHandler(exerciseRepo *repository.ExerciseRepository) *ExerciseHandler {
	return &ExerciseHandler{exerciseRepo: exerciseRepo}
}

func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exercise := model.NewExercise(input.Name, input.Description, input.Category)

	if err := h.exerciseRepo.Create(r.Context(), exercise); err != nil {
		http.Error(w, "Failed to create exercise", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercise)
}

func (h *ExerciseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	exercise, err := h.exerciseRepo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Exercise not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(exercise)
}

func (h *ExerciseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	exercises, err := h.exerciseRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(exercises)
}
