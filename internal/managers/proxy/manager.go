package proxy

import (
	"context"
	"encoding/json"
	httpCli "github.com/akhmettolegen/proxy/internal/clients"
	"github.com/akhmettolegen/proxy/internal/database/drivers"
	"github.com/akhmettolegen/proxy/internal/managers"
	"github.com/akhmettolegen/proxy/internal/models"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Manager struct {
	ctx        context.Context
	httpClient httpCli.HttpClient
	taskRepo   drivers.TaskRepository
}

func NewManager(ctx context.Context, cli httpCli.HttpClient, taskRepo drivers.TaskRepository) managers.ProxyManager {
	return &Manager{
		ctx:        ctx,
		httpClient: cli,
		taskRepo:   taskRepo,
	}
}

func (m *Manager) ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error) {

	taskId := uuid.NewString()

	go m.ProcessRequest(req, taskId)

	return &models.ProxyResponse{
		Id: taskId,
	}, nil
}

func (m *Manager) ProcessRequest(req *models.ProxyRequest, taskId string) {

	task := &models.Task{
		Id:     taskId,
		Status: models.TaskStatusInProcess,
		Length: 0,
	}

	err := m.taskRepo.Create(m.ctx, task)
	if err != nil {
		log.Println("[ERROR] Create task error:", err.Error())
		return
	}

	reqByte, err := json.Marshal(&req.Body)
	if err != nil {
		log.Println("[ERROR] Json marshal error:", err.Error())
		return
	}

	resp, err := m.httpClient.Request(req.Method, req.Url, req.Headers, reqByte)
	if err != nil {
		log.Println("[ERROR] Http client request error:", err.Error())
		return
	}

	defer resp.Body.Close()

	task.HttpStatusCode = resp.StatusCode
	task.Status = models.TaskStatusDone
	task.Length = int(resp.ContentLength)
	task.Headers = resp.Header

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errTxt := "ProxyRequest error: code=" + strconv.Itoa(resp.StatusCode) + " message=" + resp.Status
		rawBody, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			errResponse := new(models.ErrorResponse)
			if err = json.Unmarshal(rawBody, errResponse); err == nil && errResponse.Error != nil {
				errTxt = "ProxyRequest error:" + errResponse.Error.Message
			}
		}
		log.Println("[ERROR]", errTxt)

		task.Status = models.TaskStatusError
	}

	err = m.taskRepo.Update(m.ctx, task)
	if err != nil {
		log.Println("[ERROR] Update task error:", err.Error())
	}
}
