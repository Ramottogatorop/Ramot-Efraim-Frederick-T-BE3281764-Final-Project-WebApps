package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	taskes := r.Context().Value("id").(string)
	if taskes == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskId := r.URL.Query().Get("task_id")
	if len(taskId) == 0 {
		conv, _ := strconv.Atoi(taskes)
		tasks1, err := t.taskService.GetTasks(r.Context(), conv)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tasks1)
			return
		}
	} else {
		idLogin, _ := strconv.Atoi(taskId)
		tasks, err := t.taskService.GetTaskByID(r.Context(), idLogin)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tasks)
			return
		}
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	taskId := r.Context().Value("id").(string)
	if taskId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	if task.Title == "" || task.Description == "" || strconv.Itoa(task.CategoryID) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	idlogin, _ := strconv.Atoi(taskId)
	task2 := entity.Task{}
	task2.UserID = idlogin
	task2.Title = task.Title
	task2.Description = task.Description
	task2.CategoryID = task.CategoryID
	taskses, err := t.taskService.StoreTask(r.Context(), &task2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	data := map[string]interface{}{
		"user_id": idlogin,
		"task_id": taskses.ID,
		"message": "success create new task",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(data)
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskes := r.Context().Value("id").(string)
	taskID := r.URL.Query().Get("task_id")

	deleteUserId, _ := strconv.Atoi(taskID)
	conv, _ := strconv.Atoi(taskes)
	err := t.taskService.DeleteTask(r.Context(), deleteUserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	data := map[string]interface{}{
		"user_id": deleteUserId,
		"task_id": conv,
		"message": "success delete task",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id").(string)

	if len(userId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idlogin, _ := strconv.Atoi(userId)
	task2 := entity.Task{}
	task2.ID = task.ID
	task2.Title = task.Title
	task2.Description = task.Description
	task2.CategoryID = task.CategoryID
	task2.UserID = idlogin

	// var updateTask = entity.Task{
	// 	Title:       task.Title,
	// 	Description: task.Description,
	// 	CategoryID:  task.CategoryID,
	// }

	data, err := t.taskService.UpdateTask(r.Context(), &task2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idlogin,
		"task_id": data.ID,
		"message": "success update task",
	})
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
