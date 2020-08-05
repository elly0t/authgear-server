package server

import (
	authtask "github.com/authgear/authgear-server/pkg/auth/task"
	"github.com/authgear/authgear-server/pkg/deps"
	"github.com/authgear/authgear-server/pkg/task"
)

func setupTasks(registry task.Registry, p *deps.RootProvider) {
	authtask.ConfigurePwHousekeeperTask(registry, p.Task(newPwHousekeeperTask))
	authtask.ConfigureSendMessagesTask(registry, p.Task(newSendMessagesTask))
}