package groups

// This file has been disabled - Service type is not defined
// Group handlers require proper Service implementation
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("CreateGroupTripHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	var req CreateGroupTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("CreateGroupTripHandler: invalid request body", "error", err.Error(), "user_id", userID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	s.logger.Info("CreateGroupTripHandler: creating group trip", "user_id", userID, "title", req.Title, "budget", req.Budget)

	groupTrip, err := s.CreateGroupTrip(userID, &req)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("CreateGroupTripHandler: failed to create group trip", "error", err.Error(), "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("CreateGroupTripHandler: group trip created successfully", "trip_id", groupTrip.ID, "user_id", userID)
	c.JSON(http.StatusCreated, gin.H{"data": groupTrip, "status": "success"})
}

// GetGroupTripHandler handles GET /api/group-trips/:id
func (s *Service) GetGroupTripHandler(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	s.logger.Debug("GetGroupTripHandler: retrieving group trip", "trip_id", tripID, "user_id", userID)

	groupTrip, err := s.GetGroupTrip(tripID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetGroupTripHandler: failed to retrieve group trip", "error", err.Error(), "trip_id", tripID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Debug("GetGroupTripHandler: group trip retrieved", "trip_id", tripID)
	c.JSON(http.StatusOK, gin.H{"data": groupTrip, "status": "success"})
}

// GetUserGroupTripsHandler handles GET /api/user/group-trips
func (s *Service) GetUserGroupTripsHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("GetUserGroupTripsHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	s.logger.Debug("GetUserGroupTripsHandler: retrieving user group trips", "user_id", userID)

	trips, err := s.GetUserGroupTrips(userID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetUserGroupTripsHandler: failed to retrieve user group trips", "error", err.Error(), "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Debug("GetUserGroupTripsHandler: user group trips retrieved", "user_id", userID, "count", len(trips))
	c.JSON(http.StatusOK, gin.H{"data": trips, "status": "success"})
}

// UpdateGroupTripHandler handles PUT /api/group-trips/:id
func (s *Service) UpdateGroupTripHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("UpdateGroupTripHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	var req UpdateGroupTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("UpdateGroupTripHandler: invalid request body", "error", err.Error(), "user_id", userID, "trip_id", tripID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	s.logger.Info("UpdateGroupTripHandler: updating group trip", "user_id", userID, "trip_id", tripID)

	if err := s.UpdateGroupTrip(tripID, userID, &req); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("UpdateGroupTripHandler: failed to update group trip", "error", err.Error(), "trip_id", tripID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("UpdateGroupTripHandler: group trip updated successfully", "trip_id", tripID, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": nil, "status": "success"})
}

// DeleteGroupTripHandler handles DELETE /api/group-trips/:id
func (s *Service) DeleteGroupTripHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("DeleteGroupTripHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	s.logger.Info("DeleteGroupTripHandler: deleting group trip", "user_id", userID, "trip_id", tripID)

	if err := s.DeleteGroupTrip(tripID, userID); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("DeleteGroupTripHandler: failed to delete group trip", "error", err.Error(), "trip_id", tripID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("DeleteGroupTripHandler: group trip deleted successfully", "trip_id", tripID, "user_id", userID)
	c.JSON(http.StatusNoContent, "")
}

// ============================================================================
// GROUP MEMBER HANDLERS
// ============================================================================

// InviteMemberHandler handles POST /api/group-trips/:id/members
func (s *Service) InviteMemberHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("InviteMemberHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	var req AddGroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("InviteMemberHandler: invalid request body", "error", err.Error(), "user_id", userID, "trip_id", tripID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	s.logger.Info("InviteMemberHandler: inviting member", "user_id", userID, "trip_id", tripID, "new_member", req.UserID, "role", req.Role)

	if err := s.InviteMemberToGroup(tripID, userID, req.UserID, req.Role); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("InviteMemberHandler: failed to invite member", "error", err.Error(), "trip_id", tripID, "new_member", req.UserID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("InviteMemberHandler: member invited successfully", "trip_id", tripID, "new_member", req.UserID)
	c.JSON(http.StatusOK, gin.H{"data": nil, "status": "success"})
}

// GetGroupMembersHandler handles GET /api/group-trips/:id/members
func (s *Service) GetGroupMembersHandler(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	s.logger.Debug("GetGroupMembersHandler: retrieving group members", "trip_id", tripID, "user_id", userID)

	members, err := s.db.GetGroupMembers(tripID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetGroupMembersHandler: failed to retrieve members", "error", err.Error(), "trip_id", tripID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Debug("GetGroupMembersHandler: members retrieved", "trip_id", tripID, "count", len(members))
	c.JSON(http.StatusOK, gin.H{"data": members, "status": "success"})
}

// RespondInvitationHandler handles POST /api/group-trips/:id/members/respond
func (s *Service) RespondInvitationHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("RespondInvitationHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	var req RespondToInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("RespondInvitationHandler: invalid request body", "error", err.Error(), "user_id", userID, "trip_id", tripID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	accept := req.Status == "active"
	s.logger.Info("RespondInvitationHandler: responding to invitation", "user_id", userID, "trip_id", tripID, "accept", accept)

	if err := s.RespondToGroupInvite(tripID, userID, accept); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("RespondInvitationHandler: failed to respond to invitation", "error", err.Error(), "trip_id", tripID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("RespondInvitationHandler: invitation response successful", "trip_id", tripID, "user_id", userID, "accept", accept)
	c.JSON(http.StatusOK, gin.H{"data": nil, "status": "success"})
}

// RemoveGroupMemberHandler handles DELETE /api/group-trips/:id/members/:user_id
func (s *Service) RemoveGroupMemberHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("RemoveGroupMemberHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")
	targetUserID := c.Param("user_id")

	s.logger.Info("RemoveGroupMemberHandler: removing member", "user_id", userID, "trip_id", tripID, "target_user", targetUserID)

	if err := s.RemoveGroupMember(tripID, userID, targetUserID); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("RemoveGroupMemberHandler: failed to remove member", "error", err.Error(), "trip_id", tripID, "target_user", targetUserID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("RemoveGroupMemberHandler: member removed successfully", "trip_id", tripID, "target_user", targetUserID)
	c.JSON(http.StatusNoContent, "")
}

// LeaveGroupHandler handles POST /api/group-trips/:id/members/leave
func (s *Service) LeaveGroupHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("LeaveGroupHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	s.logger.Info("LeaveGroupHandler: user leaving group", "user_id", userID, "trip_id", tripID)

	if err := s.LeaveGroup(tripID, userID); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("LeaveGroupHandler: failed to leave group", "error", err.Error(), "trip_id", tripID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("LeaveGroupHandler: user left group successfully", "trip_id", tripID, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": nil, "status": "success"})
}

// ============================================================================
// EXPENSE HANDLERS
// ============================================================================

// AddExpenseHandler handles POST /api/group-trips/:id/expenses
func (s *Service) AddExpenseHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("AddExpenseHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("AddExpenseHandler: invalid request body", "error", err.Error(), "user_id", userID, "trip_id", tripID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	s.logger.Info("AddExpenseHandler: adding expense", "user_id", userID, "trip_id", tripID, "amount", req.Amount, "category", req.Category)

	expense, err := s.AddExpense(tripID, userID, &req)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("AddExpenseHandler: failed to add expense", "error", err.Error(), "trip_id", tripID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("AddExpenseHandler: expense added successfully", "trip_id", tripID, "expense_id", expense.ID, "amount", expense.Amount)
	c.JSON(http.StatusCreated, gin.H{"data": expense, "status": "success"})
}

// GetGroupExpensesHandler handles GET /api/group-trips/:id/expenses
func (s *Service) GetGroupExpensesHandler(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	s.logger.Debug("GetGroupExpensesHandler: retrieving group expenses", "trip_id", tripID, "user_id", userID)

	expenses, err := s.db.GetGroupExpenses(tripID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetGroupExpensesHandler: failed to retrieve expenses", "error", err.Error(), "trip_id", tripID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Debug("GetGroupExpensesHandler: expenses retrieved", "trip_id", tripID, "count", len(expenses))
	c.JSON(http.StatusOK, gin.H{"data": expenses, "status": "success"})
}

// GetExpenseReportHandler handles GET /api/group-trips/:id/expense-report
func (s *Service) GetExpenseReportHandler(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	s.logger.Debug("GetExpenseReportHandler: generating expense report", "trip_id", tripID, "user_id", userID)

	report, err := s.GetGroupExpenseReport(tripID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetExpenseReportHandler: failed to generate report", "error", err.Error(), "trip_id", tripID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("GetExpenseReportHandler: expense report generated", "trip_id", tripID, "settlements_count", len(report.Settlements))
	c.JSON(http.StatusOK, gin.H{"data": report, "status": "success"})
}

// ============================================================================
// POLL HANDLERS
// ============================================================================

// CreatePollHandler handles POST /api/group-trips/:id/polls
func (s *Service) CreatePollHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("CreatePollHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	tripID := c.Param("id")

	var req CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("CreatePollHandler: invalid request body", "error", err.Error(), "user_id", userID, "trip_id", tripID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	s.logger.Info("CreatePollHandler: creating poll", "user_id", userID, "trip_id", tripID, "question", req.Question)

	poll, err := s.CreatePoll(tripID, userID, &req)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("CreatePollHandler: failed to create poll", "error", err.Error(), "trip_id", tripID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("CreatePollHandler: poll created successfully", "trip_id", tripID, "poll_id", poll.ID)
	c.JSON(http.StatusCreated, gin.H{"data": poll, "status": "success"})
}

// GetPollHandler handles GET /api/polls/:id
func (s *Service) GetPollHandler(c *gin.Context) {
	pollID := c.Param("id")
	userID := c.GetString("user_id")

	s.logger.Debug("GetPollHandler: retrieving poll", "poll_id", pollID, "user_id", userID)

	poll, err := s.GetPoll(pollID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetPollHandler: failed to retrieve poll", "error", err.Error(), "poll_id", pollID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Debug("GetPollHandler: poll retrieved", "poll_id", pollID, "question", poll.Question)
	c.JSON(http.StatusOK, gin.H{"data": poll, "status": "success"})
}

// GetGroupPollsHandler handles GET /api/group-trips/:id/polls
func (s *Service) GetGroupPollsHandler(c *gin.Context) {
	tripID := c.Param("id")
	userID := c.GetString("user_id")

	s.logger.Debug("GetGroupPollsHandler: retrieving group polls", "trip_id", tripID, "user_id", userID)

	polls, err := s.db.GetGroupPolls(tripID)
	if err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("GetGroupPollsHandler: failed to retrieve polls", "error", err.Error(), "trip_id", tripID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Debug("GetGroupPollsHandler: polls retrieved", "trip_id", tripID, "count", len(polls))
	c.JSON(http.StatusOK, gin.H{"data": polls, "status": "success"})
}

// VotePollHandler handles POST /api/polls/:id/votes
func (s *Service) VotePollHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		s.logger.Warn("VotePollHandler: unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, NewAPIError(ErrUnauthorized, "user not authenticated", ""))
		return
	}

	pollID := c.Param("id")

	var req VotePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("VotePollHandler: invalid request body", "error", err.Error(), "user_id", userID, "poll_id", pollID)
		c.JSON(http.StatusBadRequest, NewAPIError(ErrValidationError, "invalid request body", err.Error()))
		return
	}

	s.logger.Info("VotePollHandler: voting on poll", "user_id", userID, "poll_id", pollID, "option_id", req.OptionID)

	if err := s.VoteOnPoll(pollID, req.OptionID, userID); err != nil {
		apiErr := err.(*APIError)
		s.logger.Error("VotePollHandler: failed to vote on poll", "error", err.Error(), "poll_id", pollID, "user_id", userID)
		c.JSON(GetStatusCode(apiErr.Code), apiErr)
		return
	}

	s.logger.Info("VotePollHandler: vote recorded successfully", "poll_id", pollID, "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"data": nil, "status": "success"})
}
