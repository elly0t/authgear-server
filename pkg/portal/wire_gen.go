// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package portal

import (
	"github.com/authgear/authgear-server/pkg/lib/admin/authz"
	"github.com/authgear/authgear-server/pkg/lib/infra/mail"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/portal/db"
	"github.com/authgear/authgear-server/pkg/portal/deps"
	"github.com/authgear/authgear-server/pkg/portal/endpoint"
	"github.com/authgear/authgear-server/pkg/portal/graphql"
	"github.com/authgear/authgear-server/pkg/portal/loader"
	"github.com/authgear/authgear-server/pkg/portal/service"
	"github.com/authgear/authgear-server/pkg/portal/session"
	"github.com/authgear/authgear-server/pkg/portal/task"
	"github.com/authgear/authgear-server/pkg/portal/task/tasks"
	"github.com/authgear/authgear-server/pkg/portal/transport"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/intl"
	"github.com/authgear/authgear-server/pkg/util/template"
	"net/http"
)

// Injectors from wire.go:

func newPanicEndMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panicEndMiddleware := &middleware.PanicEndMiddleware{}
	return panicEndMiddleware
}

func newPanicLogMiddleware(p *deps.RequestProvider) httproute.Middleware {
	rootProvider := p.RootProvider
	factory := rootProvider.LoggerFactory
	panicLogMiddlewareLogger := middleware.NewPanicLogMiddlewareLogger(factory)
	panicLogMiddleware := &middleware.PanicLogMiddleware{
		Logger: panicLogMiddlewareLogger,
	}
	return panicLogMiddleware
}

func newPanicWriteEmptyResponseMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panicWriteEmptyResponseMiddleware := &middleware.PanicWriteEmptyResponseMiddleware{}
	return panicWriteEmptyResponseMiddleware
}

func newBodyLimitMiddleware(p *deps.RequestProvider) httproute.Middleware {
	bodyLimitMiddleware := &middleware.BodyLimitMiddleware{}
	return bodyLimitMiddleware
}

func newSentryMiddleware(p *deps.RequestProvider) httproute.Middleware {
	rootProvider := p.RootProvider
	hub := rootProvider.SentryHub
	environmentConfig := rootProvider.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:  hub,
		TrustProxy: trustProxy,
	}
	return sentryMiddleware
}

func newSessionInfoMiddleware(p *deps.RequestProvider) httproute.Middleware {
	sessionInfoMiddleware := &session.SessionInfoMiddleware{}
	return sessionInfoMiddleware
}

func newSessionRequiredMiddleware(p *deps.RequestProvider) httproute.Middleware {
	sessionRequiredMiddleware := &session.SessionRequiredMiddleware{}
	return sessionRequiredMiddleware
}

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	environmentConfig := rootProvider.EnvironmentConfig
	devMode := environmentConfig.DevMode
	authgearConfig := rootProvider.AuthgearConfig
	adminAPIConfig := rootProvider.AdminAPIConfig
	controller := rootProvider.ConfigSourceController
	configSource := deps.ProvideConfigSource(controller)
	clock := _wireSystemClockValue
	adder := &authz.Adder{
		Clock: clock,
	}
	adminAPIService := &service.AdminAPIService{
		AuthgearConfig: authgearConfig,
		AdminAPIConfig: adminAPIConfig,
		ConfigSource:   configSource,
		AuthzAdder:     adder,
	}
	userLoader := loader.NewUserLoader(adminAPIService)
	factory := rootProvider.LoggerFactory
	logger := graphql.NewLogger(factory)
	appServiceLogger := service.NewAppServiceLogger(factory)
	appConfig := rootProvider.AppConfig
	request := p.Request
	context := deps.ProvideRequestContext(request)
	configServiceLogger := service.NewConfigServiceLogger(factory)
	configService := &service.ConfigService{
		Context:      context,
		Logger:       configServiceLogger,
		AppConfig:    appConfig,
		Controller:   controller,
		ConfigSource: configSource,
	}
	databaseConfig := rootProvider.DatabaseConfig
	sqlBuilder := db.NewSQLBuilder(databaseConfig)
	dbLogger := db.NewLogger(factory)
	handle := &db.Handle{
		Context: context,
		Config:  databaseConfig,
		Logger:  dbLogger,
	}
	sqlExecutor := &db.SQLExecutor{
		Context:  context,
		Database: handle,
	}
	mailConfig := rootProvider.MailConfig
	inProcessExecutorLogger := task.NewInProcessExecutorLogger(factory)
	mailLogger := mail.NewLogger(factory)
	smtpConfig := rootProvider.SMTPConfig
	smtpServerCredentials := deps.ProvideSMTPServerCredentials(smtpConfig)
	dialer := mail.NewGomailDialer(smtpServerCredentials)
	sender := &mail.Sender{
		Logger:       mailLogger,
		DevMode:      devMode,
		GomailDialer: dialer,
	}
	sendMessagesLogger := tasks.NewSendMessagesLogger(factory)
	sendMessagesTask := &tasks.SendMessagesTask{
		EmailSender: sender,
		Logger:      sendMessagesLogger,
	}
	inProcessExecutor := task.NewExecutor(inProcessExecutorLogger, sendMessagesTask)
	inProcessQueue := &task.InProcessQueue{
		Executor: inProcessExecutor,
	}
	trustProxy := environmentConfig.TrustProxy
	requestOriginProvider := &endpoint.RequestOriginProvider{
		Request:    request,
		TrustProxy: trustProxy,
	}
	endpointsProvider := &endpoint.EndpointsProvider{
		OriginProvider: requestOriginProvider,
	}
	manager := rootProvider.Resources
	defaultTemplateLanguage := _wireDefaultTemplateLanguageValue
	resolver := &template.Resolver{
		Resources:          manager,
		DefaultLanguageTag: defaultTemplateLanguage,
	}
	engine := &template.Engine{
		Resolver: resolver,
	}
	collaboratorService := &service.CollaboratorService{
		Context:        context,
		Clock:          clock,
		SQLBuilder:     sqlBuilder,
		SQLExecutor:    sqlExecutor,
		MailConfig:     mailConfig,
		TaskQueue:      inProcessQueue,
		Endpoints:      endpointsProvider,
		TemplateEngine: engine,
	}
	authzService := &service.AuthzService{
		Context:       context,
		Configs:       configService,
		Collaborators: collaboratorService,
	}
	domainService := &service.DomainService{
		Context:      context,
		Clock:        clock,
		DomainConfig: configService,
		SQLBuilder:   sqlBuilder,
		SQLExecutor:  sqlExecutor,
	}
	appService := &service.AppService{
		Logger:      appServiceLogger,
		AppConfig:   appConfig,
		AppConfigs:  configService,
		AppAuthz:    authzService,
		AppAdminAPI: adminAPIService,
		AppDomains:  domainService,
	}
	appLoader := &loader.AppLoader{
		Apps:  appService,
		Authz: authzService,
	}
	domainLoader := &loader.DomainLoader{
		Domains: domainService,
		Authz:   authzService,
	}
	collaboratorLoader := &loader.CollaboratorLoader{
		Collaborators: collaboratorService,
		Authz:         authzService,
	}
	graphqlContext := &graphql.Context{
		Users:         userLoader,
		GQLLogger:     logger,
		Apps:          appLoader,
		Domains:       domainLoader,
		Collaborators: collaboratorLoader,
	}
	graphQLHandler := &transport.GraphQLHandler{
		DevMode:        devMode,
		GraphQLContext: graphqlContext,
		Database:       handle,
	}
	return graphQLHandler
}

var (
	_wireSystemClockValue             = clock.NewSystemClock()
	_wireDefaultTemplateLanguageValue = template.DefaultTemplateLanguage(intl.DefaultLanguage)
)

func newRuntimeConfigHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	authgearConfig := rootProvider.AuthgearConfig
	runtimeConfigHandler := &transport.RuntimeConfigHandler{
		AuthgearConfig: authgearConfig,
	}
	return runtimeConfigHandler
}

func newAdminAPIHandler(p *deps.RequestProvider) http.Handler {
	request := p.Request
	context := deps.ProvideRequestContext(request)
	rootProvider := p.RootProvider
	factory := rootProvider.LoggerFactory
	configServiceLogger := service.NewConfigServiceLogger(factory)
	appConfig := rootProvider.AppConfig
	controller := rootProvider.ConfigSourceController
	configSource := deps.ProvideConfigSource(controller)
	configService := &service.ConfigService{
		Context:      context,
		Logger:       configServiceLogger,
		AppConfig:    appConfig,
		Controller:   controller,
		ConfigSource: configSource,
	}
	clockClock := _wireSystemClockValue
	databaseConfig := rootProvider.DatabaseConfig
	sqlBuilder := db.NewSQLBuilder(databaseConfig)
	logger := db.NewLogger(factory)
	handle := &db.Handle{
		Context: context,
		Config:  databaseConfig,
		Logger:  logger,
	}
	sqlExecutor := &db.SQLExecutor{
		Context:  context,
		Database: handle,
	}
	mailConfig := rootProvider.MailConfig
	inProcessExecutorLogger := task.NewInProcessExecutorLogger(factory)
	mailLogger := mail.NewLogger(factory)
	environmentConfig := rootProvider.EnvironmentConfig
	devMode := environmentConfig.DevMode
	smtpConfig := rootProvider.SMTPConfig
	smtpServerCredentials := deps.ProvideSMTPServerCredentials(smtpConfig)
	dialer := mail.NewGomailDialer(smtpServerCredentials)
	sender := &mail.Sender{
		Logger:       mailLogger,
		DevMode:      devMode,
		GomailDialer: dialer,
	}
	sendMessagesLogger := tasks.NewSendMessagesLogger(factory)
	sendMessagesTask := &tasks.SendMessagesTask{
		EmailSender: sender,
		Logger:      sendMessagesLogger,
	}
	inProcessExecutor := task.NewExecutor(inProcessExecutorLogger, sendMessagesTask)
	inProcessQueue := &task.InProcessQueue{
		Executor: inProcessExecutor,
	}
	trustProxy := environmentConfig.TrustProxy
	requestOriginProvider := &endpoint.RequestOriginProvider{
		Request:    request,
		TrustProxy: trustProxy,
	}
	endpointsProvider := &endpoint.EndpointsProvider{
		OriginProvider: requestOriginProvider,
	}
	manager := rootProvider.Resources
	defaultTemplateLanguage := _wireDefaultTemplateLanguageValue
	resolver := &template.Resolver{
		Resources:          manager,
		DefaultLanguageTag: defaultTemplateLanguage,
	}
	engine := &template.Engine{
		Resolver: resolver,
	}
	collaboratorService := &service.CollaboratorService{
		Context:        context,
		Clock:          clockClock,
		SQLBuilder:     sqlBuilder,
		SQLExecutor:    sqlExecutor,
		MailConfig:     mailConfig,
		TaskQueue:      inProcessQueue,
		Endpoints:      endpointsProvider,
		TemplateEngine: engine,
	}
	authzService := &service.AuthzService{
		Context:       context,
		Configs:       configService,
		Collaborators: collaboratorService,
	}
	authgearConfig := rootProvider.AuthgearConfig
	adminAPIConfig := rootProvider.AdminAPIConfig
	adder := &authz.Adder{
		Clock: clockClock,
	}
	adminAPIService := &service.AdminAPIService{
		AuthgearConfig: authgearConfig,
		AdminAPIConfig: adminAPIConfig,
		ConfigSource:   configSource,
		AuthzAdder:     adder,
	}
	adminAPILogger := transport.NewAdminAPILogger(factory)
	adminAPIHandler := &transport.AdminAPIHandler{
		Authz:    authzService,
		AdminAPI: adminAPIService,
		Logger:   adminAPILogger,
	}
	return adminAPIHandler
}
