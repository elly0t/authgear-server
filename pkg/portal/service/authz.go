package service

import (
	"errors"

	"github.com/authgear/authgear-server/pkg/portal/model"
)

type AuthzConfigService interface {
	GetStaticAppIDs() ([]string, error)
}

type AuthzCollaboratorService interface {
	NewCollaborator(appID string, userID string) *model.Collaborator
	CreateCollaborator(c *model.Collaborator) error
	ListCollaboratorsByUser(userID string) ([]*model.Collaborator, error)
}

type AuthzService struct {
	Configs       AuthzConfigService
	Collaborators AuthzCollaboratorService
}

func (s *AuthzService) ListAuthorizedApps(userID string) ([]string, error) {
	appIDs, err := s.Configs.GetStaticAppIDs()
	if errors.Is(err, ErrGetStaticAppIDsNotSupported) {
		var cs []*model.Collaborator
		cs, err = s.Collaborators.ListCollaboratorsByUser(userID)
		if err == nil {
			appIDs = make([]string, len(cs))
			for i, c := range cs {
				appIDs[i] = c.AppID
			}
		}
	}

	if err != nil {
		return nil, err

	}

	return appIDs, nil
}

func (s *AuthzService) AddAuthorizedUser(appID string, userID string) error {
	c := s.Collaborators.NewCollaborator(appID, userID)
	return s.Collaborators.CreateCollaborator(c)
}
