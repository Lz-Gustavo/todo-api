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
	ID       int    `json:"id,omitempty"       bson:"id"`
	Desc     string `json:"desc"               bson:"desc"`
	Status   string `json:"status"             bson:"status"`
	Finished bool   `json:"finished,omitempty" bson:"finished"`
}

func insertTask(ctx context.Context, usr string, t Task) error {
	n, err := mongoInstance.Count(ctx, usr)
	if err != nil {
		return err
	}

	// col contains N tasks, indexed 1 to N. Insert next task with N+1 id
	t.ID = n + 1
	return mongoInstance.Insert(ctx, usr, t)
}

func updateTask(ctx context.Context, usr string, id int, desc, status string) error {
	return mongoInstance.Update(ctx, usr, id, desc, status)
}

func deleteTask(ctx context.Context, usr string, id int) error {
	return mongoInstance.Delete(ctx, usr, id)
}

func listTasks(ctx context.Context, usr string, status ...string) (*ResponseJSON, error) {
	var docs []Task
	var err error

	if len(status) == 0 {
		docs, err = mongoInstance.ListAll(ctx, usr)
		if err != nil {
			return nil, err
		}

	} else {
		docs, err = mongoInstance.List(ctx, usr, status)
		if err != nil {
			return nil, err
		}
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
	return mongoInstance.UpdateFinished(ctx, usr, id, true)
}

func restoreTask(ctx context.Context, usr string, id int) error {
	return mongoInstance.UpdateFinished(ctx, usr, id, false)
}
