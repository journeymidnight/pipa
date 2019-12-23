package pipa

import (
	"encoding/json"
	"fmt"
	. "github.com/journeymidnight/pipa/error"
	"github.com/journeymidnight/pipa/helper"
	. "github.com/journeymidnight/pipa/library"
	"github.com/journeymidnight/pipa/redis"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
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

			data, err := downloadImage(imgTask.downloadUrl)
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
			ReturnQ <- FinishedTask{200,taskData.UUID,taskData.Url,data}
		}
	}
}

func downloadImage(downloadUrl string) ([]byte, error) {
	helper.Log.Info(fmt.Sprintf("Start to download %s\n", downloadUrl))

	httpClient := &http.Client{Timeout: time.Second * 5}
	resp, err := httpClient.Get(downloadUrl)

	if err != nil {
		helper.Log.Error("Download failed!", err)
		return nil, err
	}
	defer resp.Body.Close()

	//check header
	if resp.StatusCode != 200 {
		helper.Log.Info("Request is not 200")
		return nil, ErrDownloadCode
	}

	mimeType := resp.Header.Get("Content-Type")

	if strings.Contains(mimeType, "image") == false {
		if ok, _ := regexp.MatchString("(jpeg|jpg|png|gif|bmp|webp|tiff)", downloadUrl); ok == false {
			helper.Log.Info(fmt.Sprintf("MIME TYPE is %s not an image\n", mimeType))
			return nil, StatusUnsupportedMediaType //415
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewImageProcessTask(lib Library, taskData Task) (ImageProcessTask, error) {
	downloadUrl, operations, err := ParseUrl(taskData.Url)
	if err != nil {
		return ImageProcessTask{}, err
	}

	return ImageProcessTask{
		lib:         lib,
		ops:         operations,
		downloadUrl: downloadUrl,
	}, nil
}

func listenFinishedTask(resultQ chan FinishedTask) {
	c := redis.Pool.Get()
	defer c.Close()
	for r := range resultQ {
		//put data back to redis
		if r.code == 200 {
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
	result := FinishedTask{code, t.UUID, t.Url, nil}
	sendResult(result)
}

func sendResult(t FinishedTask) {
	ReturnQ <- t
}
