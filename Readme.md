# Go User API

A RESTful API built with **Go**, **Chi Router**, **GORM**, and **PostgreSQL**.  
This project demonstrates user management with **authentication**, **role-based authorization**, and **Swagger API documentation**.

---

## 🚀 Features

- **User Management**
    - Sign up new users (hashed password storage using `bcrypt`)
    - User login with **JWT authentication**
    - Update and delete user profiles
    - Fetch single user or all users

- **Authentication & Authorization**
    - **JWT-based authentication**
    - Role-based access control (**Admin** / **User**)
        - Admin: Can get all users, delete users, fetch user by ID
        - User: Can update their own profile

- **Database Integration**
    - PostgreSQL database with **GORM ORM**
    - Auto-migration for user schema

- **API Documentation**
    - Integrated with **Swagger (Swaggo)**
    - Auto-generated `swagger.json` and `swagger.yaml`

- **Environment Config**
    - Secure config using `.env` (DB credentials, JWT secret, etc.)

---

## 📂 Project Structure

```plaintext
GoApis/
 ├── api/                # Auto-generated server bindings
 ├── docs/               # Swagger documentation files
 ├── DB_Gorm.go          # Database setup and configuration
 ├── main.go             # Application entry point
 ├── apiSpec.yaml        # OpenAPI specification
 ├── Instructions.txt    # Notes and instructions
 ├── .env                # Environment variables
 ├── .gitignore          # Git ignore rules
 └── go.mod              # Go module dependencies