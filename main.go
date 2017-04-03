package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

var dbInstance *DB

func main() {

	for _, key := range os.Environ() {
		fmt.Println(key, os.Getenv(key))

	}

	app := iris.New()
	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	app.Adapt(crs)

	app.Get("/", iris.StaticHandler("./client", true, true))
	app.StaticWeb("/vendor", "./client/vendor")
	app.StaticWeb("/media", "./client/media")
	app.StaticWeb("/js", "./client/js")
	app.StaticWeb("/views", "./client/views")
	app.StaticWeb("/css", "./client/css")

	app.Get("/health", Healthcheck)
	app.Post("/users/signup", SignupMethod)
	app.Post("/users/signup/", SignupMethod)
	app.Post("/users/login", LoginMethod)
	app.Post("/users/login/", LoginMethod)

	app.Get("/tasks", NeedsToken, ListTasksMethod)
	app.Get("/tasks/", NeedsToken, ListTasksMethod)
	app.Post("/tasks", NeedsToken, CreateTaskMethod)
	app.Post("/tasks/", NeedsToken, CreateTaskMethod)
	app.Put("/tasks", NeedsToken, UpdateTaskMethod)
	app.Put("/tasks/", NeedsToken, UpdateTaskMethod)
	mongoAddr := GetMongoAddress()
	fmt.Println("mongo", mongoAddr)
	dbInstance = GetDB(mongoAddr)

	app.Listen("0.0.0.0:80")
}

func GetMongoAddress() string {
	if strings.Contains(EnvDBAddress, "tcp://") {
		return strings.Replace(EnvDBAddress, "tcp://", "mongodb://", -1)
	} else {
		return fmt.Sprintf("mongodb://%s:27017", EnvDBAddress)
	}

}

func GetDB(address string) *DB {
	timeout := time.Now().Add(time.Minute * 5)

	for time.Now().Before(timeout) {
		db, err := NewDB(address, EnvDBName, "users", "token", "tasks")
		if err != nil {
			fmt.Println("Error connecting to mongo: ", err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("Connection to mongo succeeded")
		return db
	}
	fmt.Println("Connection to mongo timed out")
	return nil
}

func Healthcheck(context *iris.Context) {
	context.WriteString("OK")
}

func SignupMethod(context *iris.Context) {
	user := User{}
	err := context.ReadJSON(&user)
	if err != nil {
		context.WriteString(fmt.Sprintf("Error: %v", err))
		return
	}

	userDB, err := dbInstance.Signup(user.Name, user.Username, user.Password)
	if err != nil {
		context.WriteString(fmt.Sprintf("Error: %v", err))
		return
	}
	context.JSON(201, userDB)
}

func LoginMethod(context *iris.Context) {
	user := User{}
	err := context.ReadJSON(&user)
	if err != nil {
		context.WriteString(fmt.Sprintf("Error: %v", err))
		return
	}
	userDB, err := dbInstance.Login(user.Username, user.Password)
	if err != nil {
		context.WriteString(fmt.Sprintf("Error: %v", err))
		return
	}
	token, err := dbInstance.GetToken(userDB)
	if err != nil {
		context.WriteString(fmt.Sprintf("Error: %v", err))
		return
	}
	context.JSON(200, token)
}

func NeedsToken(context *iris.Context) {
	auth := context.RequestHeader("Authorization")
	if len(auth) <= 0 {
		context.SetStatusCode(401)
		context.WriteString("Not authorized")
		return
	}
	user, err := dbInstance.GetUserByToken(auth)
	if err != nil || user == nil {
		context.SetStatusCode(401)
		context.WriteString("Not authorized")
		return
	}
	context.Set("user", user)
	context.Next()
}

func ListTasksMethod(context *iris.Context) {
	user := GetUser(context)
	tasks, err := dbInstance.ListTasks(user)
	if err != nil {
		context.SetStatusCode(500)
		context.WriteString(fmt.Sprintf("Error: %v", err))
		context.Next()
		return
	}
	context.JSON(200, tasks)
	context.Next()
}

func CreateTaskMethod(context *iris.Context) {
	task := &Task{}
	err := context.ReadJSON(task)
	if err != nil {
		context.SetStatusCode(400)
		context.WriteString(fmt.Sprintf("Error: %v", err))
		context.Next()
		return
	}
	user := GetUser(context)
	task, err = dbInstance.CreateTask(task.Name, user)
	if err != nil {
		context.SetStatusCode(500)
		context.WriteString(fmt.Sprintf("Error: %v", err))
		context.Next()
		return
	}
	context.JSON(201, task)
	context.Next()
}

func UpdateTaskMethod(context *iris.Context) {
	task := &Task{}
	err := context.ReadJSON(task)
	if err != nil {
		context.SetStatusCode(400)
		context.WriteString(fmt.Sprintf("Error: %v", err))
		context.Next()
		return
	}
	user := GetUser(context)
	task.UserID = user.ID
	fmt.Println("got this task", task)
	task, err = dbInstance.UpdateTask(task)
	if err != nil {
		context.SetStatusCode(500)
		context.WriteString(fmt.Sprintf("Error: %v", err))
		context.Next()
		return
	}
	context.JSON(200, task)
	context.Next()
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
