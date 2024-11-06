package route

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mummumgoodboy/verify"
	"github.com/onfirebyte/todo-dumb/internal/auth"
	"github.com/onfirebyte/todo-dumb/internal/dto"
	"github.com/onfirebyte/todo-dumb/internal/model"
	"github.com/onfirebyte/todo-dumb/internal/service"
	"gorm.io/gorm"
)

func JsonHeaderMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func CreateTodoRoute(todoService *service.TodoService, verify *verify.JWTVerifier) {
	http.HandleFunc("GET /todos", JsonHeaderMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claim, ok := getClaim(verify, r, w)
		if !ok {
			return
		}
		todos, err := todoService.GetTodosByUserId(claim.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "error while fetching todos",
				Code:  "internal_error",
			})
			return
		}

		resp := make([]dto.Todo, 0, len(todos))
		for _, todo := range todos {
			resp = append(resp, dto.Todo{
				Id:      todo.ID,
				Title:   todo.Title,
				Content: todo.Content,
				Done:    todo.Done,
			})
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))

	http.HandleFunc("POST /todos", JsonHeaderMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claim, ok := getClaim(verify, r, w)
		if !ok {
			return
		}

		var input dto.CreateTodoInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "invalid input",
				Code:  "bad_request",
			})
			return
		}

		err = todoService.CreateTodo(model.Todo{
			Title:   input.Title,
			Content: input.Content,
			OwnerID: claim.UserId,
			Done:    false,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "error while creating todo",
				Code:  "internal_error",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	http.HandleFunc("PUT /todos", JsonHeaderMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claim, ok := getClaim(verify, r, w)
		if !ok {
			return
		}

		var input dto.Todo
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "invalid input",
				Code:  "bad_request",
			})
			return
		}

		isOwner, err := todoService.IsOwner(input.Id, claim.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "error while checking ownership",
				Code:  "internal_error",
			})
			return
		}
		if !isOwner {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "you are not the owner of this todo",
				Code:  "forbidden",
			})
			return
		}

		err = todoService.UpdateTodoById(model.Todo{
			Model: gorm.Model{
				ID: input.Id,
			},
			Title:   input.Title,
			Content: input.Content,
			OwnerID: claim.UserId,
			Done:    input.Done,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "error while creating todo",
				Code:  "internal_error",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	http.HandleFunc("DELETE /todos/{id}", JsonHeaderMiddleware(func(w http.ResponseWriter, r *http.Request) {
		claim, ok := getClaim(verify, r, w)
		if !ok {
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "invalid input, id must be integer",
				Code:  "bad_request",
			})
			return
		}

		isOwner, err := todoService.IsOwner(uint(id), claim.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "error while checking ownership",
				Code:  "internal_error",
			})
			return
		}
		if !isOwner {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "you are not the owner of this todo",
				Code:  "forbidden",
			})
			return
		}

		err = todoService.DeleteTodoById(uint(id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.Error{
				Error: "error while deleting todo",
				Code:  "internal_error",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))
}

func getClaim(v *verify.JWTVerifier, r *http.Request, w http.ResponseWriter) (verify.Claims, bool) {
	token, found := auth.GetTokenHeader(r)
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.Error{
			Error: "please provide a token",
			Code:  "unauthorized",
		})
		return verify.Claims{}, false
	}
	claim, err := v.Verify(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.Error{
			Error: "unauthorized",
			Code:  "unauthorized",
		})
		return verify.Claims{}, false
	}

	return claim, true
}
