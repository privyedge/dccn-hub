package util

import (
	"testing"
)

func TestTaskList(t *testing.T) {
	userId := 1
	tasks := TaskList(userId)

	if len(tasks) > 0 {
		for i := range tasks {
			task := tasks[i]
			t.Logf("task id : %d \n", task.ID)
		}
	} else {
		t.Errorf("There is no task!")
	}
}

func TestAddUser(t *testing.T) {
	user := User{Name: "John", Token: "ed1605e17374bde6c68864d072c9f5c9", Money: 1000}
	AddUser(user)

	t.Logf("Successfully add taskmgr")
}
