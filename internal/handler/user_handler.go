package handler

import (
	"UserRepository/internal/domain"
	"UserRepository/internal/repository/postgres"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	repo *postgres.UserRepository
}

func NewUserHandler(repo *postgres.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// @Summary Создать пользователя
// @Description Создает нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "Данные пользователя"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по его ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} domain.User
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// @Summary Обновить пользователя
// @Description Обновляет данные существующего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Param user body domain.User true "Данные пользователя"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = chi.URLParam(r, "id")
	if err := h.repo.Update(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Удалить пользователя
// @Description Помечает пользователя как удаленного
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Список пользователей
// @Description Возвращает список пользователей с пагинацией
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Количество записей на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} domain.User
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 10
	}

	conditions := domain.Conditions{
		Limit:  limit,
		Offset: offset,
	}

	users, err := h.repo.List(r.Context(), conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
