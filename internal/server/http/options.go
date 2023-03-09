package http

import (
	"github.com/akhmettolegen/proxy/internal/managers"
	"github.com/akhmettolegen/proxy/internal/managers/auth"
)

type APIServerOption func(srv *APIServer)

func WithProxyManager(proxyManager managers.ProxyManager) APIServerOption {
	return func(srv *APIServer) {
		srv.proxyManager = proxyManager
	}
}

func WithAuthManager(authManager *auth.AuthManager) APIServerOption {
	return func(srv *APIServer) {
		srv.authManager = authManager
	}
}
