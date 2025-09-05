package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "strconv"
	"strings"

	"GoApis/api"

	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	_ "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"

	_ "golang.org/x/crypto/bcrypt"
)

// @title My API
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http

type User struct {
	ID       int    `json:"ID" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string `json:"password"` // don’t expose in JSON
	Address  string `json:"address"`
	Age      int    `json:"age"`
	Email    string `json:"email" gorm:"unique"`
	Role     string `json:"role"` //User or Admin
}

type MyServer struct {
	db *gorm.DB
}

func (s *MyServer) DBConnect() {
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
	err1 := s.db.AutoMigrate(&User{})
	if err1 != nil {
		return
	}
}

func (s *MyServer) insertDB(user User) {
	//s.DBConnect()
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword) // store hashed password

	res := s.db.Create(&user)
	if res.Error != nil {
		fmt.Println("error:", res.Error)
		return
	}
	fmt.Println("User Added Successfully", user)
}
func (s *MyServer) updateDB(user User) {
	//s.DBConnect()

	s.db.Model(&User{}).Where("ID = ?", user.ID).Updates(&user)
}
func (s *MyServer) deleteDB(user *User) {
	//s.DBConnect()

	s.db.Delete(user)
}
func (s *MyServer) printDB() {
	//s.DBConnect()
	var users []User
	result := s.db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("Users Size in DB = ", len(users))
	for _, user := range users {
		fmt.Println("ID: ", user.ID, "Name: ", user.Name)
	}
}
func (s *MyServer) getAllUsersDB() []User {
	//s.DBConnect()
	var users []User
	result := s.db.Find(&users)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Println("Users Size in DB = ", len(users))
	for _, user := range users {
		fmt.Println("ID: ", user.ID, "Name: ", user.Name)
	}
	return users
}

func (s *MyServer) GetUserByID(ID int) (*User, error) {
	s.DBConnect()

	var user User
	result := s.db.First(&user, ID)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

type UserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

// PostUsersSignUp godoc
// @Summary SignUp New User
// @Accept json
// @Tags users
// @Produce json
// @Param user body UserInput true "Create User"
// @Success 200 {object} User
// @Failure 400 {string} string "Invalid input"
// @Router /users/signUp [post]
func (s *MyServer) PostUsersSignUp(w http.ResponseWriter, r *http.Request) {

	var input UserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var user User
	user.Name = input.Name
	user.Email = input.Email
	user.Password = input.Password
	user.Age = input.Age
	user.Address = input.Address

	allUsers := s.getAllUsersDB()

	user.ID = len(allUsers) + 1
	if !strings.Contains(user.Email, "admin") {
		user.Role = "User"
	} else {
		user.Role = "Admin"
	}

	s.insertDB(user)

	newUser, _ := s.GetUserByID(user.ID)

	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(newUser)
	if err != nil {
		return
	}
}

// GetUsersId godoc
// @Summary Get specific User By ID
// @Description Get a single user by its ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (s *MyServer) GetUsersId(w http.ResponseWriter, r *http.Request, ID int) {
	user, ok := s.GetUserByID(ID)
	if ok != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// PutUsersUpdateId godoc
// @Summary Update New User
// @Accept json
// @Produce json
// @Param User body User true "Update User"
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "User not found"
// @Router /users/Update/{id} [put]
func (s *MyServer) PutUsersUpdateId(w http.ResponseWriter, r *http.Request, ID int) {
	userDB, ok := s.GetUserByID(ID)

	if ok != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	user.ID = userDB.ID
	s.updateDB(user)
	//w.WriteHeader(http.StatusCreated)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// DeleteUsersDeleteId godoc
// @Summary Delete New User
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "User Deleted"
// @Failure 404 {string} string "User not found"
// @Router /users/Delete/{id} [delete]
func (s *MyServer) DeleteUsersDeleteId(w http.ResponseWriter, r *http.Request, ID int) {
	user, ok := s.GetUserByID(ID)
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

// GetUsersGetAllUsers godoc
// @Summary Get all Users
// @Produce json
// @Tags users
// @Success 200 {array} User
// @Router /users/GetAllUsers [get]
func (s *MyServer) GetUsersGetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []User = s.getAllUsersDB()

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

	api.HandlerFromMux(server, r)

	http.ListenAndServe(":8080", r)
}
