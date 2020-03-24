// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package webapp

import (
	"github.com/skygeario/skygear-server/pkg/auth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/audit"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authn"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/hook"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/mfa"
	pq3 "github.com/skygeario/skygear-server/pkg/auth/dependency/mfa/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/passwordhistory/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/principal"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/principal/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/principal/password"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/webapp"
	"github.com/skygeario/skygear-server/pkg/core/async"
	pq2 "github.com/skygeario/skygear-server/pkg/core/auth/authinfo/pq"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/mail"
	"github.com/skygeario/skygear-server/pkg/core/sms"
	"github.com/skygeario/skygear-server/pkg/core/time"
	"net/http"
)

// Injectors from wire.go:

func newRootHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	validateProvider := webapp.ProvideValidateProvider(tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	renderProvider := auth.ProvideWebAppRenderProvider(m, tenantConfiguration, engine)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	provider := time.NewProvider()
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	store := pq.ProvidePasswordHistoryStore(provider, sqlBuilder, sqlExecutor)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	passwordProvider := password.ProvidePasswordProvider(sqlBuilder, sqlExecutor, provider, store, factory, tenantConfiguration, reservedNameChecker)
	oauthProvider := oauth.ProvideOAuthProvider(sqlBuilder, sqlExecutor)
	v := auth.ProvidePrincipalProviders(oauthProvider, passwordProvider)
	identityProvider := principal.ProvideIdentityProvider(sqlBuilder, sqlExecutor, v)
	authenticateProcess := authn.ProvideAuthenticateProcess(factory, provider, passwordProvider, oauthProvider, identityProvider)
	passwordChecker := audit.ProvidePasswordChecker(tenantConfiguration, store)
	loginIDChecker := loginid.ProvideLoginIDChecker(tenantConfiguration, reservedNameChecker)
	authinfoStore := pq2.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(provider, sqlBuilder, sqlExecutor)
	urlprefixProvider := urlprefix.NewProvider(r)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, urlprefixProvider, txContext, provider, authinfoStore, userprofileStore, passwordProvider, factory)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	signupProcess := authn.ProvideSignupProcess(passwordChecker, loginIDChecker, identityProvider, passwordProvider, oauthProvider, provider, authinfoStore, userprofileStore, hookProvider, tenantConfiguration, urlprefixProvider, queue)
	oAuthCoordinator := &authn.OAuthCoordinator{
		Authn:  authenticateProcess,
		Signup: signupProcess,
	}
	mfaStore := pq3.ProvideStore(tenantConfiguration, sqlBuilder, sqlExecutor, provider)
	client := sms.ProvideSMSClient(tenantConfiguration)
	sender := mail.ProvideMailSender(tenantConfiguration)
	mfaSender := mfa.ProvideMFASender(tenantConfiguration, client, sender, engine)
	mfaProvider := mfa.ProvideMFAProvider(mfaStore, tenantConfiguration, provider, mfaSender)
	sessionStore := redis.ProvideStore(context, tenantConfiguration, provider, factory)
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, tenantConfiguration)
	authnSessionProvider := authn.ProvideSessionProvider(mfaProvider, sessionProvider, tenantConfiguration, provider, authinfoStore, userprofileStore, identityProvider, hookProvider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	mfaInsecureCookieConfig := auth.ProvideMFAInsecureCookieConfig(m)
	bearerTokenCookieConfiguration := mfa.ProvideBearerTokenCookieConfiguration(r, mfaInsecureCookieConfig, tenantConfiguration)
	authnProvider := &authn.Provider{
		OAuth:                   oAuthCoordinator,
		Authn:                   authenticateProcess,
		Signup:                  signupProcess,
		AuthnSession:            authnSessionProvider,
		Session:                 sessionProvider,
		SessionCookieConfig:     cookieConfiguration,
		BearerTokenCookieConfig: bearerTokenCookieConfiguration,
	}
	authenticateProvider := webapp.ProvideAuthenticateProvider(validateProvider, renderProvider, authnProvider)
	handler := provideRootHandler(authenticateProvider)
	return handler
}

func newSettingsHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	renderProvider := auth.ProvideWebAppRenderProvider(m, tenantConfiguration, engine)
	handler := provideSettingsHandler(renderProvider)
	return handler
}

// wire.go:

func provideRootHandler(authenticateProvider webapp.AuthenticateProvider) http.Handler {
	return &RootHandler{
		AuthenticateProvider: authenticateProvider,
	}
}

func provideSettingsHandler(renderProvider webapp.RenderProvider) http.Handler {
	return &SettingsHandler{RenderProvider: renderProvider}
}
