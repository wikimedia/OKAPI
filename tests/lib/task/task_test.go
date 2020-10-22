package task

import (
	"okapi/helpers/logger"
	"okapi/lib/task"
	"okapi/models"
	"testing"
)

func TestTask(t *testing.T) {
	iterator := 0
	titles := []string{"Okapi", "Ninja", "Jimmy Hendrix", "Cat"}
	finishRan := false
	poolRan := false

	finish := func() (err error) {
		finishRan = true
		return
	}

	pool := func() (queue []task.Payload, err error) {
		if iterator < len(titles) {
			poolRan = true
			title := titles[iterator]
			iterator++
			queue = append(queue, title)
		}

		return
	}

	worker := func(id int, payload task.Payload) (msg string, info map[string]interface{}, err error) {
		title := payload.(string)

		for _, testTitle := range titles {
			if title == testTitle {
				msg = title
				return
			}
		}

		t.Error("payload data is wrong")

		return
	}

	job := func(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
		ctx.State.Set("test", 1)
		return pool, worker, finish, nil
	}

	state := testState{}
	ctx := task.Context{
		State:   &state,
		Project: &models.Project{},
		Log:     logger.Job,
		Params: task.Params{
			Workers: 4,
			DBName:  "test",
			Restart: true,
		},
	}

	err := task.Exec(job, &ctx)

	if err != nil {
		t.Error(err)
	}

	if !poolRan {
		t.Error("'pool' function wasn't executed")
	}

	if !finishRan {
		t.Error("'finish' function wasn't executed")
	}

	if len(state.store) > 0 {
		t.Error("'state` was not cleared")
	}
}

type testState struct {
	store map[string]interface{}
}

func (s *testState) Get(name string) (string, error) {
	if val, ok := s.store[name]; ok {
		return val.(string), nil
	}

	return "", nil
}

func (s *testState) GetInt(name string, initial int) int {
	if val, ok := s.store[name]; ok {
		return val.(int)
	}

	return initial
}

func (s *testState) GetString(name string, initial string) string {
	if val, ok := s.store[name]; ok {
		return val.(string)
	}

	return initial
}

func (s *testState) Set(name string, value interface{}) error {
	s.store[name] = value
	return nil
}

func (s *testState) Exists(name string) bool {
	_, ok := s.store[name]
	return ok
}

func (s *testState) Clear() {
	s.store = map[string]interface{}{}
}
