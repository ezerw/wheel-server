package service

import (
	"context"
	"time"

	"github.com/ezerw/wheel/db"
)

// Turns is the service in charge of interact with the turns table in the database.
type Turns struct {
	store db.Store
}

// TurnAPI is the representation returned to the client
type TurnAPI struct {
	ID        int64     `json:"id"`
	PersonID  int64     `json:"person_id"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTurns creates a new TeamsService instance.
func NewTurns(store db.Store) *Turns {
	return &Turns{store: store}
}

// ListTurns gets turns from the DB based on passed params.
func (s *Turns) ListTurns(
	ctx context.Context,
	teamID int64,
	dateFrom time.Time,
	dateTo time.Time,
	limit int64,
	offset int64,
) ([]TurnAPI, error) {
	if !dateFrom.IsZero() && !dateTo.IsZero() {
		args := db.ListTurnsWithBothDatesParams{
			TeamID: teamID,
			Date:   dateFrom,
			Date_2: dateTo,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		dbTurns, err := s.store.ListTurnsWithBothDates(ctx, args)
		if err != nil {
			return nil, err
		}

		var turns []TurnAPI

		for _, turn := range dbTurns {
			turns = append(turns, TurnAPI{
				ID:        turn.ID,
				PersonID:  turn.PersonID,
				Date:      turn.Date,
				CreatedAt: turn.CreatedAt.Time,
			})
		}

		return turns, nil
	}

	if !dateFrom.IsZero() {
		args := db.ListTurnsWithDateFromParams{
			Date:   dateFrom,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		dbTurns, err := s.store.ListTurnsWithDateFrom(ctx, args)
		if err != nil {
			return nil, err
		}

		var turns []TurnAPI
		for _, turn := range dbTurns {
			turns = append(turns, TurnAPI{
				ID:        turn.ID,
				PersonID:  turn.PersonID,
				Date:      turn.Date,
				CreatedAt: turn.CreatedAt.Time,
			})
		}
		return turns, nil
	}

	if !dateTo.IsZero() {
		args := db.ListTurnsWithDateToParams{
			TeamID: teamID,
			Date:   dateTo,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		dbTurns, err := s.store.ListTurnsWithDateTo(ctx, args)
		if err != nil {
			return nil, err
		}

		var turns []TurnAPI
		for _, turn := range dbTurns {
			turns = append(turns, TurnAPI{
				ID:        turn.ID,
				PersonID:  turn.PersonID,
				Date:      turn.Date,
				CreatedAt: turn.CreatedAt.Time,
			})
		}

		return turns, nil
	}

	args := db.ListTurnsParams{
		TeamID: teamID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	dbTurns, err := s.store.ListTurns(ctx, args)
	if err != nil {
		return nil, err
	}

	var turns []TurnAPI
	for _, turn := range dbTurns {
		turns = append(turns, TurnAPI{
			ID:        turn.ID,
			PersonID:  turn.PersonID,
			Date:      turn.Date,
			CreatedAt: turn.CreatedAt.Time,
		})
	}

	return turns, nil
}

// GetTurn gets one turn from the DB using id and teamID as params.
func (s *Turns) GetTurn(ctx context.Context, args db.GetTurnParams) (*TurnAPI, error) {
	turn, err := s.store.GetTurn(ctx, args)
	if err != nil {
		return nil, err
	}

	apiTurn := &TurnAPI{
		ID:        turn.ID,
		PersonID:  turn.PersonID,
		Date:      turn.Date,
		CreatedAt: turn.CreatedAt.Time,
	}

	return apiTurn, nil
}

// GetTurnByDate gets one turn from the DB using date and teamID as params.
func (s *Turns) GetTurnByDate(ctx context.Context, args db.GetTurnByDateAndTeamParams) (*TurnAPI, error) {
	turn, err := s.store.GetTurnByDateAndTeam(ctx, args)
	if err != nil {
		return nil, err
	}

	apiTurn := &TurnAPI{
		ID:        turn.ID,
		PersonID:  turn.PersonID,
		Date:      turn.Date,
		CreatedAt: turn.CreatedAt.Time,
	}

	return apiTurn, nil
}

// AddTurn adds a turn to the DB for the specified team.
func (s *Turns) AddTurn(ctx context.Context, teamID int64, personID int64, date time.Time) (*TurnAPI, error) {
	args := db.CreateTurnParams{
		PersonID: personID,
		Date:     date,
	}

	result, err := s.store.CreateTurn(ctx, args)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	getTurnArgs := db.GetTurnParams{
		ID:     id,
		TeamID: teamID,
	}

	turn, err := s.GetTurn(ctx, getTurnArgs)
	if err != nil {
		return nil, err
	}

	return turn, nil
}

// UpdateTurn updates a turn in the DB.
func (s *Turns) UpdateTurn(ctx context.Context, teamID int64, args db.UpdateTurnParams) (*TurnAPI, error) {
	_, err := s.store.UpdateTurn(ctx, args)
	if err != nil {
		return nil, err
	}

	return s.GetTurn(ctx, db.GetTurnParams{
		ID:     args.ID,
		TeamID: teamID,
	})

}
