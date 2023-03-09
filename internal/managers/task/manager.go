package task

import (
	"context"
	"encoding/json"
	"errors"
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

var (
	ErrInValidTaskId = errors.New("invalid task id")
)

type Manager struct {
	ctx        context.Context
	httpClient httpCli.HttpClient
	taskRepo   drivers.TaskRepository
}

func NewManager(ctx context.Context, cli httpCli.HttpClient, taskRepo drivers.TaskRepository) managers.TaskManager {
	return &Manager{
		ctx:        ctx,
		httpClient: cli,
		taskRepo:   taskRepo,
	}
}

func (m *Manager) TaskCreate(req *models.TaskRequest) (*models.TaskResponse, error) {

	taskId := uuid.NewString()

	go m.Create(req, taskId)

	return &models.TaskResponse{
		Id: taskId,
	}, nil
}

func (m *Manager) Create(req *models.TaskRequest, taskId string) {

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

	err = m.ProcessRequest(req, task)
	if err != nil {
		log.Println("[ERROR] Process request error:", err.Error())
	}

	err = m.taskRepo.Update(m.ctx, task)
	if err != nil {
		log.Println("[ERROR] Update task error:", err.Error())
	}
}

func (m *Manager) ProcessRequest(req *models.TaskRequest, task *models.Task) error {
	reqByte, err := json.Marshal(&req.Body)
	if err != nil {
		log.Println("[ERROR] Json marshal error:", err.Error())
		return err
	}

	resp, err := m.httpClient.Request(req.Method, req.Url, req.Headers, reqByte)
	if err != nil {
		log.Println("[ERROR] Http client request error:", err.Error())
		task.Status = models.TaskStatusError
		return err
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
	return nil
}

func (m *Manager) TaskById(id string) (*models.Task, error) {
	task, err := m.taskRepo.TaskById(m.ctx, id)
	if err != nil {
		log.Println("[ERROR] Task by id error:", err.Error())
		return nil, err
	}

	return &models.Task{
		Id:             task.Id,
		Status:         task.Status,
		HttpStatusCode: task.HttpStatusCode,
		Headers:        task.Headers,
		Length:         task.Length,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
	}, nil
}
