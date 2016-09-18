package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (db *DB) CreateTask(name string, user *User) (*Task, error) {

	collection := db.GetTaskCollection()

	task := Task{
		ID:        bson.NewObjectId(),
		CreatedOn: time.Now(),
		Name:      name,
		UserID:    user.ID,
		Done:      false,
	}
	err := collection.Insert(&task)
	return &task, err
}

func (db *DB) UpdateTask(task *Task) (*Task, error) {

	collection := db.GetTaskCollection()

	err := collection.UpdateId(task.ID, task)
	// _, err := collection.Upsert(bson.M{
	// 	"name": task.Name,
	// }, bson.M{
	// 	"name": task.Name,
	// 	"done": task.Done,
	// })

	return task, err
}

func (db *DB) ListTasks(user *User) ([]Task, error) {
	collection := db.GetTaskCollection()
	tasks := make([]Task, 0, 10)
	err := collection.Find(bson.M{
		"userid": user.ID,
	}).Sort("-createdon").All(&tasks)
	return tasks, err
}

type Task struct {
	ID        bson.ObjectId `json:"_id",bson:"_id"`
	UserID    bson.ObjectId `json:"userid",bson:"userid"`
	Name      string        `json:"name"`
	CreatedOn time.Time     `json:"createdon"`
	Done      bool          `json:"done"`
}
