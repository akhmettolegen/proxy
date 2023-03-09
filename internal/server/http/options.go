package http

import (
	"github.com/akhmettolegen/proxy/internal/managers"
	"github.com/akhmettolegen/proxy/internal/managers/auth"
)

type APIServerOption func(srv *APIServer)

func WithTaskManager(taskManager managers.TaskManager) APIServerOption {
	return func(srv *APIServer) {
		srv.taskManager = taskManager
	}
}

func WithAuthManager(authManager *auth.AuthManager) APIServerOption {
	return func(srv *APIServer) {
		srv.authManager = authManager
	}
}
