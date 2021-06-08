package service

import (
	"context"

	"github.com/ezerw/wheel/db"
)

// Teams is the service in charge of interact with the teams table in the database.
type Teams struct {
	store db.Store
}

// NewTeams creates a new TeamsService instance.
func NewTeams(store db.Store) *Teams {
	return &Teams{store: store}
}

// ListTeams gets a list of teams from the DB.
func (s *Teams) ListTeams(ctx context.Context) ([]db.ListTeamsRow, error) {
	return s.store.ListTeams(ctx)
}

// GetTeam gets a team from the DB.
func (s *Teams) GetTeam(ctx context.Context, teamID int64) (*db.GetTeamRow, error) {
	team, err := s.store.GetTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

// AddTeam adds a team to the DB.
func (s *Teams) AddTeam(ctx context.Context, teamName string) (*db.GetTeamRow, error) {
	result, err := s.store.CreateTeam(ctx, teamName)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &db.GetTeamRow{
		ID:   id,
		Name: teamName,
	}, nil
}

// UpdateTeam updates a team name in the DB.
func (s *Teams) UpdateTeam(ctx context.Context, args db.UpdateTeamParams) error {
	_, err := s.store.UpdateTeam(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTeam deletes a team from the DB
func (s *Teams) DeleteTeam(ctx context.Context, teamID int64) error {
	return s.store.DeleteTeam(ctx, teamID)
}
