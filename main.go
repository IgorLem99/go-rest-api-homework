package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

func handleTasks(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	// Попробуем закодировать задачи в JSON
	err := json.NewEncoder(res).Encode(tasks)
	if err != nil {
		// Если произошла ошибка, вернем статус 500
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)

		return
	}
	res.WriteHeader(http.StatusOK)

}

func createTaskHandler(res http.ResponseWriter, r *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var newTask Task

	// Декодируем JSON из тела запроса в структуру Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(res, "Bad Request", http.StatusBadRequest)

		return
	}

	// Генерируем ID на основе длины карты задач
	newTask.ID = strconv.Itoa(len(tasks) + 1)

	// Сохраняем новую задачу в мапе
	tasks[newTask.ID] = newTask

	// Возвращаем статус 201 Created
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(newTask)
}

func getTaskByIDHandler(res http.ResponseWriter, r *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Получаем параметр id из URL
	taskID := chi.URLParam(r, "id")

	// Ищем задачу в мапе по ID
	task, exists := tasks[taskID]
	if !exists {
		// Если задача не найдена, возвращаем статус 400
		http.Error(res, "Task not found", http.StatusBadRequest)

		return
	}

	// Возвращаем найденную задачу
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(task)
}

func deleteTaskByIDHandler(res http.ResponseWriter, r *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	// Получаем параметр id из URL
	taskID := chi.URLParam(r, "id")

	// Проверяем, существует ли задача с таким ID
	if _, exists := tasks[taskID]; !exists {
		// Если задача не найдена, возвращаем статус 400
		http.Error(res, "Task not found", http.StatusBadRequest)

		return
	}

	// Удаляем задачу из мапы
	delete(tasks, taskID)

	// Возвращаем успешный статус 200 OK
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(map[string]string{
		"message": "Task deleted successfully",
	})
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...

	r.Get("/tasks", handleTasks)
	r.Post("/tasks", createTaskHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Delete("/tasks/{id}", deleteTaskByIDHandler)

	if err := http.ListenAndServe(":8081", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
