package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
)

// HandleListTeams handles GET request to /api/teams
func (s *Server) HandleListTeams(c *gin.Context) {
	teams, err := s.teamsService.ListTeams(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": teams})
}

// HandleShowTeam handles GET request to /api/teams/:team-id
func (s *Server) HandleShowTeam(c *gin.Context) {
	queryTeamID := c.Param("team-id")
	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	team, err := s.teamsService.GetTeam(c.Request.Context(), teamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Team not found.",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add people to the team response
	composite := struct {
		TeamID   int64              `json:"id"`
		TeamName string             `json:"name"`
		People   []db.ListPeopleRow `json:"people"`
	}{}

	composite.TeamName = team.Name
	composite.TeamID = team.ID

	people, _ := s.peopleService.ListPeople(c.Request.Context(), teamID)
	composite.People = people

	c.JSON(http.StatusOK, gin.H{"data": composite})
}

// HandleAddTeam handles POST request to /api/teams
func (s *Server) HandleAddTeam(c *gin.Context) {
	binding := struct {
		Name string `json:"name" binding:"required"`
	}{}
	err := c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := s.teamsService.AddTeam(c.Request.Context(), binding.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": team})
}

// HandleUpdateTeam handles PUT request to /api/teams/:team-id
func (s *Server) HandleUpdateTeam(c *gin.Context) {
	queryTeamID := c.Param("team-id")
	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binding := struct {
		Name string `json:"name,omitempty"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := s.teamExists(c.Request.Context(), teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Team not found."})
		return
	}

	updateTeamArgs := db.UpdateTeamParams{
		Name: binding.Name,
		ID:   teamID,
	}

	err = s.teamsService.UpdateTeam(c.Request.Context(), updateTeamArgs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// HandleDeleteTeam handles DELETE request to /api/teams/:team-id
func (s *Server) HandleDeleteTeam(c *gin.Context) {
	queryTeamID := c.Param("team-id")
	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	exists, err := s.teamExists(c.Request.Context(), teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Team not found."})
		return
	}

	err = s.teamsService.DeleteTeam(c.Request.Context(), teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
