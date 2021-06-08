package handler

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/middleware"
	"github.com/ezerw/wheel/service"
	"github.com/ezerw/wheel/util"
)

// Server serves HTTP requests for our wheel api.
type Server struct {
	config        util.Config
	router        *gin.Engine
	peopleService *service.People
	teamsService  *service.Teams
	turnsService  *service.Turns
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config:        config,
		peopleService: service.NewPeople(store),
		teamsService:  service.NewTeams(store),
		turnsService:  service.NewTurns(store),
	}

	server.setupRouter()
	return server, nil
}

// Start runs the HTTP server on a specific address.
func (s *Server) Start(address string, port string) error {
	return s.router.Run(address + ":" + port)
}

// setupRouter defines a router and add routes to it.
func (s *Server) setupRouter() {
	ginMode := gin.ReleaseMode
	if s.config.AppDebug {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	r := gin.Default()

	// TODO: authenticate requests
	api := r.
		Group("/api").
		Use(middleware.Cors())

	// teams
	api.GET("/teams", s.HandleListTeams)
	api.GET("/teams/:team-id", s.HandleShowTeam)
	api.POST("/teams", s.HandleAddTeam)
	api.PUT("/teams/:team-id", s.HandleUpdateTeam)
	api.DELETE("/teams/:team-id", s.HandleDeleteTeam)

	// team people
	api.GET("/teams/:team-id/people", s.HandleListPeople)
	api.GET("/teams/:team-id/people/:person-id", s.HandleShowPerson)
	api.POST("/teams/:team-id/people", s.HandleAddPerson)
	api.PUT("/teams/:team-id/people/:person-id", s.HandleUpdatePerson)
	api.DELETE("/teams/:team-id/people/:person-id", s.HandleDeletePerson)

	// team turns
	api.GET("/teams/:team-id/turns", s.HandleListTurns)
	api.POST("/teams/:team-id/turns", s.HandleUpsertTurn)

	s.router = r
}

// teamExists checks if the specified teamID exists in the DB.
func (s *Server) teamExists(ctx context.Context, teamID int64) (bool, error) {
	_, err := s.teamsService.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, nil
	}
	return true, nil
}
