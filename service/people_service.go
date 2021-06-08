package service

import (
	"context"
	"github.com/ezerw/wheel/db"
)

// People is the service in charge of interact with the people table in the database.
type People struct {
	store db.Store
}

// NewPeople creates a new PeopleService instance.
func NewPeople(store db.Store) *People {
	return &People{store: store}
}

// ListPeople gets people of a team from the DB.
func (s *People) ListPeople(ctx context.Context, teamID int64) ([]db.ListPeopleRow, error) {
	return s.store.ListPeople(ctx, teamID)
}

// GetPerson gets one person of the team from the DB.
func (s *People) GetPerson(ctx context.Context, args db.GetPersonParams) (*db.GetPersonRow, error) {
	person, err := s.store.GetPerson(ctx, args)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

// AddPerson add one person to the team in the DB.
func (s *People) AddPerson(ctx context.Context, args db.CreatePersonParams) (*db.GetPersonRow, error) {
	result, err := s.store.CreatePerson(ctx, args)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetPerson(ctx, db.GetPersonParams{
		ID:     id,
		TeamID: args.TeamID,
	})
}

// UpdatePerson updates a person from the team in the DB.
func (s *People) UpdatePerson(ctx context.Context, args db.UpdatePersonParams) error {
	_, err := s.store.UpdatePerson(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

// DeletePerson deletes a person from the team from the DB.
func (s *People) DeletePerson(ctx context.Context, args db.DeletePersonParams) error {
	return s.store.DeletePerson(ctx, args)
}
