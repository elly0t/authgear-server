// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package admin

import (
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/util/httproute"
)

// Injectors from wire.go:

func newSentryMiddleware(p *deps.RootProvider) httproute.Middleware {
	hub := p.SentryHub
	serverConfig := p.ServerConfig
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:    hub,
		ServerConfig: serverConfig,
	}
	return sentryMiddleware
}

func newRootRecoverMiddleware(p *deps.RootProvider) httproute.Middleware {
	factory := p.LoggerFactory
	recoveryLogger := middleware.NewRecoveryLogger(factory)
	recoverMiddleware := &middleware.RecoverMiddleware{
		Logger: recoveryLogger,
	}
	return recoverMiddleware
}

func newRequestRecoverMiddleware(p *deps.RequestProvider) httproute.Middleware {
	appProvider := p.AppProvider
	factory := appProvider.LoggerFactory
	recoveryLogger := middleware.NewRecoveryLogger(factory)
	recoverMiddleware := &middleware.RecoverMiddleware{
		Logger: recoveryLogger,
	}
	return recoverMiddleware
}
