package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
)

// HandleListPeople handles GET request to /api/teams/:team-id/people
func (s *Server) HandleListPeople(c *gin.Context) {
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

	people, err := s.peopleService.ListPeople(c.Request.Context(), teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": people})
}

// HandleShowPerson handles GET request to /api/teams/:team-id/people/:person-id
func (s *Server) HandleShowPerson(c *gin.Context) {
	queryTeamID := c.Param("team-id")
	queryPersonID := c.Param("person-id")

	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	personID, err := strconv.ParseInt(queryPersonID, 10, 64)
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

	args := db.GetPersonParams{
		ID:     personID,
		TeamID: teamID,
	}
	person, err := s.peopleService.GetPerson(c.Request.Context(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": person})
}

// HandleAddPerson handles POST request to /api/teams/:team-id/people
func (s *Server) HandleAddPerson(c *gin.Context) {
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

	binding := struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := db.CreatePersonParams{
		FirstName: binding.FirstName,
		LastName:  binding.LastName,
		Email:     binding.Email,
		TeamID:    teamID,
	}

	person, err := s.peopleService.AddPerson(c.Request.Context(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": person})
}

// HandleUpdatePerson handles PUT request to /api/teams/:team-id/people/:person-id
func (s *Server) HandleUpdatePerson(c *gin.Context) {
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

	queryPersonID := c.Param("person-id")
	personID, err := strconv.ParseInt(queryPersonID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	binding := struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		TeamID    int64  `json:"team_id"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the person exists in the specified team
	getPersonArgs := db.GetPersonParams{
		ID:     personID,
		TeamID: teamID,
	}
	person, err := s.peopleService.GetPerson(c.Request.Context(), getPersonArgs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "Person not found in the specified team.",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("Person: %+v\n", person)
	args := db.UpdatePersonParams{
		ID:        person.ID,
		FirstName: binding.FirstName,
		LastName:  binding.LastName,
		Email:     binding.Email,
		TeamID:    binding.TeamID,
	}

	if args.FirstName == "" {
		args.FirstName = person.FirstName
	}

	if args.LastName == "" {
		args.LastName = person.LastName
	}

	if args.Email == "" {
		args.Email = person.Email
	}

	if args.TeamID == 0 {
		args.TeamID = teamID
	}
	fmt.Printf("Args: %+v", args)
	err = s.peopleService.UpdatePerson(c.Request.Context(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// HandleDeletePerson handles DELETE request to /api/teams/:team-id/people/:person-id
func (s *Server) HandleDeletePerson(c *gin.Context) {
	queryTeamID := c.Param("team-id")
	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	queryPersonID := c.Param("person-id")
	personID, err := strconv.ParseInt(queryPersonID, 10, 64)
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

	args := db.DeletePersonParams{
		ID:     personID,
		TeamID: teamID,
	}

	err = s.peopleService.DeletePerson(c.Request.Context(), args)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
