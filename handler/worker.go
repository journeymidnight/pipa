package handler

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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const TaskQueue = "taskQueue"
const PictureMAxSize = 20 << 20

var (
	finishPipa chan bool
	taskQueue  chan string
	wg         sync.WaitGroup
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
	code          int
	returnMessage string
	uuid          string
	url           string
	blob          []byte
}

var redisErrStop = false

func StartWorker() {
	finishPipa = make(chan bool)
	taskQueue = make(chan string)
	for i := 0; i < helper.Config.WorkersNumber; i++ {
		go slave(i)
	}
	go receiveImageTask()
}

func receiveImageTask() {
	for {
		select {
		case <-finishPipa:
			helper.Log.Info("stop receive:")
			return
		default:
			r, err := redis.RedisConn.BRPop(TaskQueue, 15)
			if err != nil || len(r) < 2 {
				if err == redis.CircuitBroken {
					helper.Log.Error("Redis status is abnormal, exit the main program")
					redisErrStop = true
					SignalQueue <- syscall.SIGQUIT
					return
				}
				continue
			}
			taskQueue <- r[1]
		}
	}
}

func slave(slave_num int) {
	lib := NewLibrary()
	defer CloseLibrary()
	for {
		select {
		case <-finishPipa:
			helper.Log.Info("stop slave:", slave_num)
			return
		case task := <-taskQueue:
			wg.Add(1)
			if len(task) == 0 {
				wg.Done()
				continue
			}

			helper.Log.Info("slave", slave_num, "receive task:", task)

			var taskData Task
			err := json.Unmarshal([]byte(task), &taskData)
			if err != nil {
				returnError(ErrInvalidTaskString, Task{})
				wg.Done()
				continue
			}

			imgTask, err := NewImageProcessTask(lib, taskData)
			if err != nil {
				returnError(err, taskData)
				wg.Done()
				continue
			}

			data, err := downloadImage(imgTask.downloadUrl)
			if err != nil {
				returnError(err, taskData)
				wg.Done()
				continue
			}

			for _, op := range imgTask.ops {
				data, err = op.DoProcess(data)
				if err != nil {
					returnError(err, taskData)
					break
				}
			}
			if err == nil {
				ReturnQ := FinishedTask{200, "200,Process picture success!", taskData.UUID, taskData.Url, data}
				listenFinishedTask(ReturnQ)
			}
			wg.Done()
		default:
			continue
		}
	}
}

func downloadImage(downloadUrl string) ([]byte, error) {
	helper.Log.Info(fmt.Sprintf("Start to download %s\n", downloadUrl))

	httpClient := &http.Client{Timeout: time.Second * 30}
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
	contentLength := resp.Header.Get("Content-Length")
	if len, _ := strconv.Atoi(contentLength); len > PictureMAxSize {
		return nil, StatusRequestEntityTooLarge
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewImageProcessTask(lib Library, taskData Task) (ImageProcessTask, error) {
	downloadUrl, operations, err := ParseUrl(taskData.Url, false)
	if err != nil {
		return ImageProcessTask{}, err
	}

	return ImageProcessTask{
		lib:         lib,
		ops:         operations,
		downloadUrl: downloadUrl,
	}, nil
}

func listenFinishedTask(resultQ FinishedTask) {
	if resultQ.code == 200 {
		redis.RedisConn.LPushSucceed(resultQ.url, resultQ.uuid, resultQ.returnMessage, resultQ.blob)
		resultQ.blob = nil
	} else {
		redis.RedisConn.LPushFailed(resultQ.uuid, resultQ.returnMessage)
	}
	helper.Log.Info(fmt.Sprintf("finishing task [%s] for %s code %s\n", resultQ.uuid, resultQ.url, resultQ.returnMessage))
}

func returnError(err error, t Task) {
	var (
		code    int
		message string
	)
	e, ok := err.(PipaError)
	if ok {
		code, message = e.ErrorCode()
	} else {
		code = 400
	}
	helper.Log.Error(err)
	returnMessage := strconv.Itoa(code) + "," + message
	result := FinishedTask{code, returnMessage, t.UUID, t.Url, nil}
	listenFinishedTask(result)
}

func Stop() {
	helper.Log.Info("Stopping Pipa")
	counts := helper.Config.WorkersNumber
	if !redisErrStop {
		counts += 1
	}
	for i := 0; i < counts; i++ {
		finishPipa <- true
	}
	wg.Wait()
	helper.Log.Info("Done")
	close(taskQueue)
	close(finishPipa)
}
