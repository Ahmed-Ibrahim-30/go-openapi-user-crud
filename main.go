package main

import (
	"encoding/json"
	"net/http"
	_ "strconv"

	"github.com/go-chi/cors"

	"GoApis/api"

	"github.com/go-chi/chi/v5"
)

type MyServer struct {
	users map[int]api.User
}

// GET /users/{id}
func (s *MyServer) GetUsersId(w http.ResponseWriter, r *http.Request, id int) {
	user, ok := s.users[id]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// POST /CreateUser
func (s *MyServer) PostUsersCreateUser(w http.ResponseWriter, r *http.Request) {
	var user api.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.users[user.Id] = user
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// PUT /users/{id}
func (s *MyServer) PutUsersUpdateId(w http.ResponseWriter, r *http.Request, id int) {
	userUpdated, ok := s.users[id]

	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var user api.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.users[userUpdated.Id] = user
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// Delete /users/{id}
func (s *MyServer) DeleteUsersDeleteId(w http.ResponseWriter, r *http.Request, id int) {
	var user api.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	delete(s.users, user.Id)
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// Get / users/GetAllUsers
func (s *MyServer) GetUsersGetAllUsers(w http.ResponseWriter, r *http.Request) {
	for _, user := range s.users {
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			return
		}
	}

	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)

}

func main() {
	r := chi.NewRouter()
	server := &MyServer{users: make(map[int]api.User)}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by browsers
	}))

	// Register generated handlers
	api.HandlerFromMux(server, r)

	http.ListenAndServe(":8080", r)
}
