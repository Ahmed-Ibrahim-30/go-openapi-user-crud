package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "strconv"

	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"GoApis/api"

	"github.com/go-chi/chi/v5"
)

type MyServer struct {
	db *gorm.DB
}

func (s *MyServer) DBConnect() {
	//env File
	//Create in yaml File
	//auth in JWT

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	name := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", name, pass, host, port, dbName)
	fmt.Println("DSN = " + dsn)
	s.db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Turn Off Gorm Self Naming
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	err1 := s.db.AutoMigrate(&api.User{})
	if err1 != nil {
		return
	}
}

func (s *MyServer) insertDB(user api.User) {
	//s.DBConnect()
	res := s.db.Create(&user)
	if res.Error != nil {
		fmt.Println("error:", res.Error)
		return
	}
	fmt.Println("User Added Successfully", user)
}
func (s *MyServer) updateDB(user api.User) {
	//s.DBConnect()

	s.db.Model(&api.User{}).Where("id = ?", user.Id).Updates(&user)
}
func (s *MyServer) deleteDB(user *api.User) {
	//s.DBConnect()

	s.db.Delete(user)
}
func (s *MyServer) printDB() {
	//s.DBConnect()
	var users []api.User
	result := s.db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("Users Size in DB = ", len(users))
	for _, user := range users {
		fmt.Println("ID: ", user.Id, "Name: ", user.Name)
	}
}
func (s *MyServer) getAllUsersDB() []api.User {
	//s.DBConnect()
	var users []api.User
	result := s.db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("Users Size in DB = ", len(users))
	for _, user := range users {
		fmt.Println("ID: ", user.Id, "Name: ", user.Name)
	}
	return users
}

func (s *MyServer) GetUserByID(id int) (*api.User, error) {
	s.DBConnect()

	var user api.User
	result := s.db.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// POST /CreateUser
func (s *MyServer) PostUsersCreateUser(w http.ResponseWriter, r *http.Request) {
	var user api.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	s.insertDB(user)

	newUser, _ := s.GetUserByID(user.Id)

	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(newUser)
	if err != nil {
		return
	}
}

// GET /users/{id}
func (s *MyServer) GetUsersId(w http.ResponseWriter, r *http.Request, id int) {
	user, ok := s.GetUserByID(id)
	if ok != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// PUT /users/{id}
func (s *MyServer) PutUsersUpdateId(w http.ResponseWriter, r *http.Request, id int) {
	userDB, ok := s.GetUserByID(id)

	if ok != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var user api.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	user.Id = userDB.Id
	s.updateDB(user)
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// Delete /users/{id}
func (s *MyServer) DeleteUsersDeleteId(w http.ResponseWriter, r *http.Request, id int) {
	user, ok := s.GetUserByID(id)
	if ok != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	s.deleteDB(user)
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// Get / users/GetAllUsers
func (s *MyServer) GetUsersGetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []api.User = s.getAllUsersDB()

	for _, user := range users {
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			return
		}
	}

	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)

}

func main() {
	//DataBase
	r := chi.NewRouter()
	server := &MyServer{}
	server.DBConnect()

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
