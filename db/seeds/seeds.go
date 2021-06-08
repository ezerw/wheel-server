package seeds

import (
	"context"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/util"
)

// Seeder represents the task or entity and the action
// for seeding it into the DB.
type Seeder struct {
	Name string
	Run  func(context.Context, *db.Store) error
}

// tradingTeam is the list of people in the Trading team
var tradingTeam = []db.CreatePersonParams{
	{
		FirstName: "Emily",
		LastName:  "James",
		Email:     "emily.james@vendhq.com",
	},
	{
		FirstName: "Eze",
		LastName:  "Rodriguez",
		Email:     "eze.rodriguez@vendhq.com",
	},
	{
		FirstName: "Harrison",
		LastName:  "Fata",
		Email:     "harrison.fata@vendhq.com",
	},
	{
		FirstName: "Hugh",
		LastName:  "He",
		Email:     "xiao.he@vendhq.com",
	},
	{
		FirstName: "Ivy",
		LastName:  "Fu",
		Email:     "ivy.fu@vendhq.com",
	},
	{
		FirstName: "Lucille",
		LastName:  "Tachet",
		Email:     "lucille.tachet@vendhq.com",
	},
	{
		FirstName: "Matt",
		LastName:  "O\\'Grady",
		Email:     "matthew.ogrady@vendhq.com",
	},
	{
		FirstName: "Rebekah",
		LastName:  "Berriman",
		Email:     "rebekah.berriman@vendhq.com",
	},
}

// All returns the list of all seeders available
// add yours here.
func All() []Seeder {
	return []Seeder{
		{
			Name: "CreateTradingTeam",
			Run: func(ctx context.Context, store *db.Store) error {
				return CreateTradingTeam(ctx, *store, "Trading")
			},
		},
	}
}

// CreateTradingTeam inserts trading team and people into the database.
func CreateTradingTeam(ctx context.Context, store db.Store, name string) error {
	logger := util.NewLogger()

	res, err := store.CreateTeam(ctx, name)
	if err != nil {
		return err
	}

	logger.Infof("'%s' team created successfully", name)

	teamID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, person := range tradingTeam {
		person.TeamID = teamID

		_, err = store.CreatePerson(ctx, person)
		if err != nil {
			return err
		}
		logger.Infof("%s was added successfully to the %s team", person.FirstName, name)
	}

	return nil
}
