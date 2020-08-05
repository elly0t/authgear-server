//+build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/authgear/authgear-server/pkg/auth/dependency/auth"
	"github.com/authgear/authgear-server/pkg/auth/dependency/webapp"
	"github.com/authgear/authgear-server/pkg/core/sentry"
	"github.com/authgear/authgear-server/pkg/deps"
	"github.com/authgear/authgear-server/pkg/httproute"
	"github.com/authgear/authgear-server/pkg/middlewares"
)

var rootMiddlewareDependencySet = wire.NewSet(
	deps.RootDependencySet,
)

var middlewareDependencySet = wire.NewSet(
	deps.RequestDependencySet,
)

func newSentryMiddleware(p *deps.RootProvider) httproute.Middleware {
	panic(wire.Build(
		rootMiddlewareDependencySet,
		sentry.DependencySet,
		wire.Bind(new(httproute.Middleware), new(*sentry.Middleware)),
	))
}

func newRootRecoverMiddleware(p *deps.RootProvider) httproute.Middleware {
	panic(wire.Build(
		rootMiddlewareDependencySet,
		middlewares.DependencySet,
		wire.Bind(new(httproute.Middleware), new(*middlewares.RecoverMiddleware)),
	))
}

func newRequestRecoverMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*middlewares.RecoverMiddleware)),
	))
}

func newCORSMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*middlewares.CORSMiddleware)),
	))
}

func newCSPMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*webapp.CSPMiddleware)),
	))
}

func newCSRFMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*webapp.CSRFMiddleware)),
	))
}

func newAuthEntryPointMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*webapp.AuthEntryPointMiddleware)),
	))
}

func newSessionMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*auth.Middleware)),
	))
}

func newWebAppStateMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panic(wire.Build(
		middlewareDependencySet,
		wire.Bind(new(httproute.Middleware), new(*webapp.StateMiddleware)),
	))
}