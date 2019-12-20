package pipa

import (
	"github.com/journeymidnight/pipa/helper"
	"github.com/journeymidnight/pipa/redis"
	"fmt"
	"encoding/json"
	. "github.com/journeymidnight/pipa/library"
	. "github.com/journeymidnight/pipa/error"
)

const TaskQueue = "taskQueue"

var (
	TaskQ   chan string
	ReturnQ chan FinishedTask
)

type Task struct {
	UUID string `json:"uuid"`
	Url  string `json:"url"`
}

type ImageProcessTask struct {
	downloadUrl string
	ops         []Operation
	lib         Library
}

type TaskData struct {
	uuid         string
	url          string
	taskType     string
	bucketDomain string
	captures     map[string]string
}

type FinishedTask struct {
	code int
	uuid string
	url  string
	blob []byte
	mime string
}

func StartWorker() {
	TaskQ = make(chan string, helper.Config.MaxTaskNumber)
	ReturnQ = make(chan FinishedTask)

	for i := 0; i < helper.Config.WorkersNumber; i++ {
		go slave(i)
	}
	go listenFinishedTask(ReturnQ)

	//TODO: Use signal channel to quit
	for {
		t, err := receiveImageTask()
		if err != nil {
			helper.Log.Error("receive image task error:", err)
			continue
		}
		helper.Log.Info("receive image task:", t)
		TaskQ <- t
	}
}

func receiveImageTask() (string, error) {
	r, err := redis.BRPop(TaskQueue, 0)
	if err != nil {
		return "", err
	}
	return r[1], nil
}

func slave(slave_num int) {
	for {
		select {
		case task := <-TaskQ:
			helper.Log.Info("slave", slave_num, "receive task:", task)
			lib := NewLibrary()

			var taskData Task
			err := json.Unmarshal([]byte(task), &taskData)
			if err != nil {
				returnError(ErrInvalidTaskString, Task{})
				continue
			}

			imgTask, err := NewImageProcessTask(lib, taskData)
			if err != nil {
				returnError(err, taskData)
				continue
			}

			//TODO: Download Origin Image
			data, err := downloadImage(taskData.Url)
			if err != nil {
				returnError(err, taskData)
				continue
			}

			for _, op := range imgTask.ops {
				data, err = op.DoProcess(data)
				op.Close()
				if err != nil {
					returnError(err, taskData)
					break
				}
			}
		}
	}
}

func downloadImage(downloadUrl string) ([]byte, error) {
	return nil, nil
}

func NewImageProcessTask(lib Library, taskData Task) (ImageProcessTask, error) {
	// TODO: parse Url by using url.Parse(), return operations and error
	// url.Parse(taskData.Url)
	// get downloadUrl, Operations.
	return ImageProcessTask{
		lib: lib,
		ops:,
		downloadUrl:,
	}, nil
}

func listenFinishedTask(resultQ chan FinishedTask) {
	c := redis.Pool.Get()
	defer c.Close()
	for r := range resultQ {
		//put data back to redis
		if r.code == 200 {
			//combined := combineData(r.blob, r.mime)
			c.Do("MULTI")
			c.Do("SET", r.url, r.blob)
			c.Do("LPUSH", r.uuid, r.code)
			c.Do("EXEC")
			r.blob = nil
		} else {
			c.Do("LPUSH", r.uuid, r.code)
		}
		helper.Log.Info(fmt.Sprintf("finishing task [%s] for %s code %d\n", r.uuid, r.url, r.code))
	}
}

func returnError(err error, t Task) {
	var code int
	e, ok := err.(PipaError)
	if ok {
		code = e.ErrorCode()
	} else {
		code = 400
	}
	helper.Log.Error(err)
	result := FinishedTask{code, t.UUID, t.Url, nil, ""}
	sendResult(result)
}

func sendResult(t FinishedTask) {
	ReturnQ <- t
}
