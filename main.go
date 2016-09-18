package main

import (
	"fmt"
	"strings"

	"github.com/kataras/iris"
)

var dbInstance *DB

func main() {
	iris.Get("/health", Healthcheck)
	iris.Post("/users/signup", SignupMethod)
	iris.Post("/users/signup/", SignupMethod)
	iris.Post("/users/login", LoginMethod)
	iris.Post("/users/login/", LoginMethod)

	iris.Get("/tasks", NeedsToken, ListTasksMethod)
	iris.Get("/tasks/", NeedsToken, ListTasksMethod)
	iris.Post("/tasks", NeedsToken, CreateTaskMethod)
	iris.Post("/tasks/", NeedsToken, CreateTaskMethod)
	iris.Put("/tasks", NeedsToken, UpdateTaskMethod)
	iris.Put("/tasks/", NeedsToken, UpdateTaskMethod)

	iris.Get("/", iris.StaticHandler("./client/", 0, true, false, []string{"index.html", "css", "js", "vendor", "views", "media"}))
	iris.Static("/css", "./client/css", 1)
	iris.Static("/js", "./client/js", 1)
	iris.Static("/vendor", "./client/vendor", 1)
	iris.Static("/views", "./client/views", 1)
	iris.Static("/media", "./client/media", 1)
	// iris.Static("/", "./client/index.html", 1)

	dbInstance = GetDB(GetMongoAddress())

	iris.Listen(EnvIPAddress + ":" + EnvPort)
}

func GetMongoAddress() string {
	return strings.Replace(EnvDBAddress, "tcp://", "mongodb://", -1)
}

func GetDB(address string) *DB {
	count := 0
	max := 5
	for count < max {
		db, err := NewDB(address, EnvDBName, "users", "token", "tasks")
		count++
		if err != nil {
			fmt.Println("Error connecting to mongo: ", err)
			continue
		}
		return db
	}
	return nil
}

func Healthcheck(context *iris.Context) {
	context.Write("OK")
}

func SignupMethod(context *iris.Context) {
	user := User{}
	err := context.ReadJSON(&user)
	if err != nil {
		context.Write("Error: %v", err)
		return
	}

	userDB, err := dbInstance.Signup(user.Name, user.Username, user.Password)
	if err != nil {
		context.Write("Error: %v", err)
		return
	}
	context.JSON(201, userDB)
}

func LoginMethod(context *iris.Context) {
	user := User{}
	err := context.ReadJSON(&user)
	if err != nil {
		context.Write("Error: %v", err)
		return
	}
	userDB, err := dbInstance.Login(user.Username, user.Password)
	if err != nil {
		context.Write("Error: %v", err)
		return
	}
	token, err := dbInstance.GetToken(userDB)
	if err != nil {
		context.Write("Error: %v", err)
		return
	}
	context.JSON(200, token)
}

func NeedsToken(context *iris.Context) {
	auth := context.RequestHeader("Authorization")
	if len(auth) <= 0 {
		context.SetStatusCode(401)
		context.Write("Not authorized")
		context.Do()
		return
	}
	// result := strings.Split(auth, " ")
	// if len(result) < 2 {
	// 	context.SetStatusCode(401)
	// 	context.Write("Not authorized")
	// 	context.Do()
	// 	return
	// }
	user, err := dbInstance.GetUserByToken(auth)
	if err != nil {
		context.SetStatusCode(401)
		context.Write("Not authorized ", err)
		context.Do()
		return
	}
	context.Set("user", user)
	context.Next()
}

func ListTasksMethod(context *iris.Context) {
	user := GetUser(context)
	if user == nil {
		context.SetStatusCode(404)
		context.Write("Error, user not found")
		context.Do()
		return
	}
	tasks, err := dbInstance.ListTasks(user)
	if err != nil {
		context.SetStatusCode(500)
		context.Write("Error ", err)
		context.Do()
		return
	}
	context.JSON(200, tasks)
}

func CreateTaskMethod(context *iris.Context) {
	task := &Task{}
	err := context.ReadJSON(task)
	if err != nil {
		context.SetStatusCode(400)
		context.Write("Error ", err)
		context.Do()
		return
	}
	user := GetUser(context)
	if user == nil {
		context.SetStatusCode(404)
		context.Write("Error, user not found")
		context.Do()
		return
	}
	task, err = dbInstance.CreateTask(task.Name, user)
	if err != nil {
		context.SetStatusCode(500)
		context.Write("Error ", err)
		context.Do()
		return
	}
	context.JSON(201, task)
}

func UpdateTaskMethod(context *iris.Context) {
	task := &Task{}
	err := context.ReadJSON(task)
	if err != nil {
		context.SetStatusCode(400)
		context.Write("Error ", err)
		context.Do()
		return
	}
	user := GetUser(context)
	if user == nil {
		context.SetStatusCode(404)
		context.Write("Error, user not found")
		context.Do()
		return
	}
	task.UserID = user.ID
	task, err = dbInstance.UpdateTask(task)
	if err != nil {
		context.SetStatusCode(500)
		context.Write("Error ", err)
		context.Do()
		return
	}
	context.JSON(200, task)
}

func GetUser(context *iris.Context) *User {
	userInt := context.Get("user")
	var user *User
	var ok bool
	user, ok = userInt.(*User)
	if ok {
		return user
	} else {
		return nil
	}
}
