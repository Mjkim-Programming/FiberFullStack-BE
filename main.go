//TODO Remove Template Renders and send json results
//TODO Add React Frontend to this
//TODO Write Dockerfile for BE and FE
//TODO Write docker-compose.yml for running BE, FE, and Redis
//TODO Change this Code to Use Postgres instead of SQLite

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mjkim-Programming/FiberWeb/ent"
	"github.com/Mjkim-Programming/FiberWeb/ent/user"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
	_ "github.com/mattn/go-sqlite3"
)

type UserDTO struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	if err != nil {
		log.Fatalf("Failed opening connectionto sqlite: %v", err)
	}

	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("Failed creating schema resources: %v", err)
	}

	ctx := context.Background()

	CreateUser(ctx, client, "a8m", 30)

	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		AppName: "Test App v1.0.2",
		Views: engine,
	})

	app.Get("/", func (c fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Text" : "Hello, World 🖐️!",
		})
	})

	app.Get("/hello/:name?", func (c fiber.Ctx) error {
		if c.Params("name") != "" {
			return c.SendString("Hello, " + c.Params("name") + "!")
		}

		return c.SendString("Hello, World!")
	})

	app.Get("/surprise/:name", func (c fiber.Ctx) error {
		return c.SendString("Surprise, " + c.Params("name") + "!")
	})

	app.Get("/ent/user/:name", func (c fiber.Ctx) error {
		u, err := QueryUserByName(ctx, client, c.Params("name"))

		if err != nil {
			return c.SendString("Cannot find user with name " + c.Params("name") + ".")
		}

		return c.Render("user", fiber.Map{
			"user": u,
		})
	})

	app.Get("/ent/user", func (c fiber.Ctx) error {
		u, _ :=  QueryAllUser(ctx, client)

		return c.Render("userAll", fiber.Map{
			"userList": u,
		})
	})

	app.Post("/ent/user", func (c fiber.Ctx) error {
		user := new(UserDTO)
		body := c.Body()

		if err := json.Unmarshal(body, &user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON input",
			})
		}

		u, _ := CreateUser(ctx, client, user.Name, user.Age)

		return c.JSON(u)
	})

	log.Fatal(app.Listen(":4000"))
}

func CreateUser(ctx context.Context, client *ent.Client, name string, age int) (*ent.User, error) {
	u, err := client.User.
		Create().
		SetAge(age).
		SetName(name).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %v", err)
	}
	log.Println("User was created: ", u)

	return u, nil
}

func QueryUserByName(ctx context.Context, client *ent.Client, name string) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.Name(name)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}

	log.Println("User returned: ", u)
	return u, nil
}

func QueryUserByID(ctx context.Context, client *ent.Client, id int) (*ent.User, error) {
	u, err := client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying user: %v", err)
	}

	log.Println("User returned: ", u)
	return u, nil
}

func QueryAllUser(ctx context.Context, client *ent.Client) ([]*ent.User, error) {
	u, err := client.User.
		Query().All(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("failed")
	}

	return u, nil
}