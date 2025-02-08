package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	_ "github.com/lib/pq"
)

type UserDTO struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

var db *sql.DB

func main() {
	connStr := "postgres://user:password@postgres:5432/mydb?sslmode=disable"
	
	var err error
	
	db, err = sql.Open("postgres", connStr)
	
	if err != nil {
		log.Fatalf("Failed to connect to Postgres : %v", err)
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		age INT
	);`)

	if err != nil {
		log.Fatalf("Failed to Create Table : %v", err)
	}

	CreateUser("a8m", 30)

	app := fiber.New(fiber.Config{
		AppName: "Test App v1.0.2",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
	}))

	app.Get("/user/:name", func (c fiber.Ctx) error {
		u, err := QueryUserByName(c.Params("name"))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		return c.JSON(u)
	})

	app.Get("/user", func (c fiber.Ctx) error {
		u, err := QueryAllUser()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		return c.JSON(u)
	})

	app.Post("/user", func (c fiber.Ctx) error {
		user := new(UserDTO)
		body := c.Body()

		if err := json.Unmarshal(body, &user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON input",
			})
		}

		u, _ := CreateUser(user.Name, user.Age)

		return c.JSON(u)
	})

	log.Fatal(app.Listen(":4000"))
}

func CreateUser(name string, age int) (map[string]interface{}, error) {
	var id int
	err := db.QueryRow(`INSERT INTO users(name, age) VALUES ($1, $2) RETURNING id`, name, age).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("Failed to Create User : %v", err)
	}

	user := map[string]interface{} {
		"id": id,
		"name": name,
		"age": age,
	}

	return user, nil
}

func QueryUserByName(name string) (map[string]interface{}, error) {
	var id int
	var nameTmp string
	var age int

	row := db.QueryRow(`SELECT id, name, age FROM users WHERE name=$1`, name)
	err := row.Scan(&id, &nameTmp, &age)

	if err != nil {
		return nil, fmt.Errorf("Failed Querying user by name : %v", err)
	}

	user := make(map[string]interface{})
	user["id"] = id
	user["name"] = nameTmp
	user["age"] = age

	return user, nil
}

func QueryUserByID(id int) (map[string]interface{}, error) {
	var idTmp int
	var name string
	var age int

	row := db.QueryRow(`SELECT id, name, age FROM users WHERE id=$1`, id)
	err := row.Scan(&idTmp, &name, &age)

	if err != nil {
		return nil, fmt.Errorf("Failed Querying user by name : %v", err)
	}

	user := make(map[string]interface{})
	user["id"] = id
	user["name"] = name
	user["age"] = age

	return user, nil
}

func QueryAllUser() ([]map[string]interface{}, error) {
	rows, err := db.Query(`SELECT id, name, age FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed fetching all users: %v", err)
	}
	defer rows.Close()

	var users []map[string]interface{}

	var id int
	var name string
	var age int

	for rows.Next() {
		if err := rows.Scan(&id, &name, &age); err != nil {
			return nil, fmt.Errorf("failed scanning user: %v", err)
		}

		user := make(map[string]interface{})

		user["id"] = id
		user["name"] = name
		user["age"] = age

		users = append(users, user)
	}

	return users, nil
}