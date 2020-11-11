package main

import "context"

// Task ...
type Task struct {
	id       int
	desc     string
	status   string
	finished bool
}

func createTask(ctx context.Context, usr, desc string) error {
	// TODO: First discover the num of elements on that col to set last id
	lastID := 0

	t := &Task{
		id:   lastID,
		desc: desc,
	}
	return mongoInstance.Insert(ctx, t, usr)
}

func updateTask(ctx context.Context, usr string, id int, desc string) {}

func deleteTask(ctx context.Context, usr string) {}

func listTasks(ctx context.Context, usr string) {
	// TODO: separate finished and unfinished ones
}

func completeTask(ctx context.Context, usr string, id int) {
	// TODO: just set finished bool to true
}

func restoreTask(ctx context.Context, usr string, id int) {
	// TODO: just set finished bool to false
}
