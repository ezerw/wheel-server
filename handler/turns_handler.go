package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
	_ "time/tzdata" // required

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/ezerw/wheel/db"
	"github.com/ezerw/wheel/util"
)

// HandleListTurns handles GET requests to /api/teams/:team-id/turns
// it accepts the following query params:
// - limit [Default to 10]
// - offset [Default to 0]
// - date_from [Format: YYYY-MM-DD]
// - date_to [Format: YYYY-MM-DD]
func (s *Server) HandleListTurns(c *gin.Context) {
	queryLimit := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(queryLimit, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit invalid format."})
		return
	}

	queryOffset := c.DefaultQuery("offset", "0")
	offset, err := strconv.ParseInt(queryOffset, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "offset invalid format."})
		return
	}

	queryTeamID := c.Param("team-id")
	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "team_id invalid format"})
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

	var (
		dateFrom time.Time
		dateTo   time.Time
	)

	loc, err := time.LoadLocation(s.config.AppTimezone)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load location: " + err.Error()})
		return
	}

	queryDateFrom := c.Query("date_from")
	if queryDateFrom != "" {
		dateFrom, err = time.ParseInLocation("2006-01-02", queryDateFrom, loc)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date_from invalid format."})
			return
		}
	}

	queryDateTo := c.Query("date_to")
	if queryDateTo != "" {
		dateTo, err = time.ParseInLocation("2006-01-02", queryDateTo, loc)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date_to invalid format."})
			return
		}
	}

	turns, err := s.turnsService.ListTurns(c.Request.Context(), teamID, dateFrom, dateTo, limit, offset)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": turns})
}

// HandleUpsertTurn handles POST request to /api/teams/:team-id/turns
// if the specified date has already a turn for the team it will update the person assigned
// if it doesn't exist will create the turn.
// DB unique: (team_id, date) - A team can't have multiple people assigned for the same date.
func (s *Server) HandleUpsertTurn(c *gin.Context) {
	queryTeamID := c.Param("team-id")
	teamID, err := strconv.ParseInt(queryTeamID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "team_id invalid format"})
		return
	}
	teamExists, err := s.teamExists(c.Request.Context(), teamID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !teamExists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Team not found."})
		return
	}

	// Only required person as date will be calculated.
	binding := struct {
		PersonID int64 `json:"person_id" binding:"required"`
	}{}
	err = c.BindJSON(&binding)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	getPersonArgs := db.GetPersonParams{
		ID:     binding.PersonID,
		TeamID: teamID,
	}
	person, err := s.peopleService.GetPerson(c.Request.Context(), getPersonArgs)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Person not found in the specified team"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	date, err := util.GetNextWorkingDay(s.config.AppTimezone)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error getting next working day"})
		return
	}

	// Lookup the specified turn by date and team
	getTurnArgs := db.GetTurnByDateAndTeamParams{
		Date:     *date,
		TeamID: person.TeamID,
	}
	turn, err := s.turnsService.GetTurnByDate(c.Request.Context(), getTurnArgs)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			// if it doesn't exist create it and return
			turn, err = s.turnsService.AddTurn(c.Request.Context(), teamID, person.ID, *date)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": turn})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// if it does exist, update it
	updateTurnArgs := db.UpdateTurnParams{
		PersonID: person.ID,
		Date:     turn.Date,
		ID:       turn.ID,
	}

	turn, err = s.turnsService.UpdateTurn(c.Request.Context(), teamID, updateTurnArgs)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": turn})
}
