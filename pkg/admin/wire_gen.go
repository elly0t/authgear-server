// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package admin

import (
	"github.com/authgear/authgear-server/pkg/admin/graphql"
	"github.com/authgear/authgear-server/pkg/admin/loader"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/oob"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/password"
	service2 "github.com/authgear/authgear-server/pkg/lib/authn/authenticator/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/authenticator/totp"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/anonymous"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/loginid"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/oauth"
	"github.com/authgear/authgear-server/pkg/lib/authn/identity/service"
	"github.com/authgear/authgear-server/pkg/lib/authn/user"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/feature/verification"
	"github.com/authgear/authgear-server/pkg/lib/infra/db"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"net/http"
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

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	appProvider := p.AppProvider
	config := appProvider.Config
	secretConfig := config.SecretConfig
	databaseCredentials := deps.ProvideDatabaseCredentials(secretConfig)
	appConfig := config.AppConfig
	appID := appConfig.ID
	sqlBuilder := db.ProvideSQLBuilder(databaseCredentials, appID)
	request := p.Request
	context := deps.ProvideRequestContext(request)
	handle := appProvider.Database
	sqlExecutor := db.SQLExecutor{
		Context:  context,
		Database: handle,
	}
	store := &user.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticationConfig := appConfig.Authentication
	identityConfig := appConfig.Identity
	loginidStore := &loginid.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	loginIDConfig := identityConfig.LoginID
	rootProvider := appProvider.RootProvider
	reservedNameChecker := rootProvider.ReservedNameChecker
	typeCheckerFactory := &loginid.TypeCheckerFactory{
		Config:              loginIDConfig,
		ReservedNameChecker: reservedNameChecker,
	}
	checker := &loginid.Checker{
		Config:             loginIDConfig,
		TypeCheckerFactory: typeCheckerFactory,
	}
	normalizerFactory := &loginid.NormalizerFactory{
		Config: loginIDConfig,
	}
	provider := &loginid.Provider{
		Store:             loginidStore,
		Config:            loginIDConfig,
		Checker:           checker,
		NormalizerFactory: normalizerFactory,
	}
	oauthStore := &oauth.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	clock := _wireSystemClockValue
	oauthProvider := &oauth.Provider{
		Store: oauthStore,
		Clock: clock,
	}
	anonymousStore := &anonymous.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	anonymousProvider := &anonymous.Provider{
		Store: anonymousStore,
		Clock: clock,
	}
	serviceService := &service.Service{
		Authentication: authenticationConfig,
		Identity:       identityConfig,
		LoginID:        provider,
		OAuth:          oauthProvider,
		Anonymous:      anonymousProvider,
	}
	factory := appProvider.LoggerFactory
	logger := verification.NewLogger(factory)
	verificationConfig := appConfig.Verification
	passwordStore := &password.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticatorConfig := appConfig.Authenticator
	authenticatorPasswordConfig := authenticatorConfig.Password
	passwordLogger := password.NewLogger(factory)
	historyStore := &password.HistoryStore{
		Clock:       clock,
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	passwordChecker := password.ProvideChecker(authenticatorPasswordConfig, historyStore)
	queue := appProvider.TaskQueue
	passwordProvider := &password.Provider{
		Store:           passwordStore,
		Config:          authenticatorPasswordConfig,
		Clock:           clock,
		Logger:          passwordLogger,
		PasswordHistory: historyStore,
		PasswordChecker: passwordChecker,
		TaskQueue:       queue,
	}
	totpStore := &totp.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	authenticatorTOTPConfig := authenticatorConfig.TOTP
	totpProvider := &totp.Provider{
		Store:  totpStore,
		Config: authenticatorTOTPConfig,
		Clock:  clock,
	}
	authenticatorOOBConfig := authenticatorConfig.OOB
	oobStore := &oob.Store{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	oobProvider := &oob.Provider{
		Config: authenticatorOOBConfig,
		Store:  oobStore,
		Clock:  clock,
	}
	service3 := &service2.Service{
		Password: passwordProvider,
		TOTP:     totpProvider,
		OOBOTP:   oobProvider,
	}
	redisHandle := appProvider.Redis
	storeRedis := &verification.StoreRedis{
		Redis: redisHandle,
		AppID: appID,
		Clock: clock,
	}
	verificationService := &verification.Service{
		Logger:         logger,
		Config:         verificationConfig,
		LoginID:        loginIDConfig,
		Clock:          clock,
		Authenticators: service3,
		Store:          storeRedis,
	}
	queries := &user.Queries{
		Store:        store,
		Identities:   serviceService,
		Verification: verificationService,
	}
	userLoader := &loader.UserLoader{
		Users: queries,
	}
	graphqlContext := &graphql.Context{
		Users: userLoader,
	}
	graphQLHandler := &GraphQLHandler{
		GraphQLContext: graphqlContext,
		Database:       handle,
	}
	return graphQLHandler
}

var (
	_wireSystemClockValue = clock.NewSystemClock()
)
