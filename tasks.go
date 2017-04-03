package main

import (
	"time"

	"fmt"

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

	err := collection.Update(
		bson.M{"_id": task.ID},
		bson.M{
			"$set": bson.M{"name": task.Name, "done": task.Done},
		},
	)
	return task, err
}

func (db *DB) ListTasks(user *User) ([]Task, error) {
	collection := db.GetTaskCollection()
	tasks := make([]Task, 0, 10)
	err := collection.Find(bson.M{
		"userid": user.ID,
	}).Sort("-createdon").All(&tasks)

	if len(tasks) > 0 {
		fmt.Println(tasks[0].ID)
	}
	return tasks, err
}

type Task struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	UserID    bson.ObjectId `bson:"userid" json:"userid"`
	Name      string        `json:"name"`
	CreatedOn time.Time     `json:"createdon"`
	Done      bool          `json:"done"`
}
