package main

import (
	"context"
)

// ResponseJSON ...
type ResponseJSON struct {
	ActiveTasks   []Task `json:"active"`
	FinishedTasks []Task `json:"finished"`
}

// Task ...
type Task struct {
	ID       int    `json:"id"`
	Desc     string `json:"desc"`
	Status   string `json:"status"`
	Finished bool   `json:"finished"`
}

func createTask(ctx context.Context, usr, desc string) error {
	// TODO: First discover the num of elements on that col to set last id
	lastID := 0

	t := &Task{
		ID:   lastID,
		Desc: desc,
	}
	return mongoInstance.Insert(ctx, t, usr)
}

func updateTask(ctx context.Context, usr string, id int, desc string) error {
	return nil
}

func deleteTask(ctx context.Context, usr string, id int) error {
	return nil
}

func listTasks(ctx context.Context, usr string, status ...string) (*ResponseJSON, error) {
	// TODO: filter from status
	docs, err := mongoInstance.ListAll(ctx, usr)
	if err != nil {
		return nil, err
	}

	act := make([]Task, 0)
	fin := make([]Task, 0)

	for _, t := range docs {
		if t.Finished {
			fin = append(fin, t)
		} else {
			act = append(act, t)
		}
	}

	return &ResponseJSON{
		ActiveTasks:   act,
		FinishedTasks: fin,
	}, nil
}

func completeTask(ctx context.Context, usr string, id int) error {
	// TODO: just set finished bool to true
	return nil
}

func restoreTask(ctx context.Context, usr string, id int) error {
	// TODO: just set finished bool to false
	return nil
}
