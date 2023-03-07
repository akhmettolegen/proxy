package managers

import "github.com/akhmettolegen/proxy/internal/models"

type ProxyManager interface {
	ProxyRequest(req *models.ProxyRequest) (*models.ProxyResponse, error)
	ProcessRequest(req *models.ProxyRequest, taskId string)
}
