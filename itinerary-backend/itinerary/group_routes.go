package itinerary

// This file has been disabled - Service type is not defined
// See groups/routes.go for available group routes
	groupRoutes := router.Group("/api")
	groupRoutes.Use(authMiddleware.RequireAuth())
	{
		// ==================== Group Trip Routes ====================
		groupRoutes.POST("/group-trips", service.CreateGroupTripHandler)
		groupRoutes.GET("/group-trips/:id", service.GetGroupTripHandler)
		groupRoutes.GET("/user/group-trips", service.GetUserGroupTripsHandler)
		groupRoutes.PUT("/group-trips/:id", service.UpdateGroupTripHandler)
		groupRoutes.DELETE("/group-trips/:id", service.DeleteGroupTripHandler)

		// ==================== Group Member Routes ====================
		groupRoutes.POST("/group-trips/:id/members", service.InviteMemberHandler)
		groupRoutes.GET("/group-trips/:id/members", service.GetGroupMembersHandler)
		groupRoutes.POST("/group-trips/:id/members/respond", service.RespondInvitationHandler)
		groupRoutes.DELETE("/group-trips/:id/members/:user_id", service.RemoveGroupMemberHandler)
		groupRoutes.POST("/group-trips/:id/members/leave", service.LeaveGroupHandler)

		// ==================== Expense Routes ====================
		groupRoutes.POST("/group-trips/:id/expenses", service.AddExpenseHandler)
		groupRoutes.GET("/group-trips/:id/expenses", service.GetGroupExpensesHandler)
		groupRoutes.GET("/group-trips/:id/expense-report", service.GetExpenseReportHandler)

		// ==================== Poll Routes ====================
		groupRoutes.POST("/group-trips/:id/polls", service.CreatePollHandler)
		groupRoutes.GET("/polls/:id", service.GetPollHandler)
		groupRoutes.GET("/group-trips/:id/polls", service.GetGroupPollsHandler)
		groupRoutes.POST("/polls/:id/votes", service.VotePollHandler)
	}

	logger.Info("Registered group routes", "routes_count", 16)
}
