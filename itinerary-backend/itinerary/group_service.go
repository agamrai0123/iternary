package itinerary

import (
	"fmt"
	"sort"
)

// ============================================================================
// GROUP TRIP SERVICE OPERATIONS
// ============================================================================

// CreateGroupTrip creates a new group trip and adds the owner
func (s *Service) CreateGroupTrip(userID string, req *CreateGroupTripRequest) (*GroupTrip, error) {
	// Validate input
	if req.Budget <= 0 {
		return nil, NewAPIError(ErrValidationError, "budget must be greater than 0", "")
	}
	if req.Duration <= 0 {
		return nil, NewAPIError(ErrValidationError, "duration must be greater than 0", "")
	}
	if req.Title == "" {
		return nil, NewAPIError(ErrValidationError, "title is required", "")
	}

	// Create group trip
	groupTrip := &GroupTrip{
		Title:         req.Title,
		DestinationID: req.DestinationID,
		OwnerID:       userID,
		Budget:        req.Budget,
		Duration:      req.Duration,
		StartDate:     req.StartDate,
		Status:        GroupTripStatusDraft,
	}

	if err := s.db.CreateGroupTrip(groupTrip); err != nil {
		return nil, err
	}

	// Add owner as member with owner role
	if _, err := s.db.AddGroupMember(groupTrip.ID, userID, GroupMemberRoleOwner); err != nil {
		return nil, err
	}

	// Add initial members if provided
	for _, memberID := range req.InitialMembers {
		if memberID != userID { // Don't add owner twice
			if _, err := s.db.AddGroupMember(groupTrip.ID, memberID, GroupMemberRoleMember); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to add member %s to group trip: %v", memberID, err))
				// Continue with other members
			}
		}
	}

	return groupTrip, nil
}

// GetGroupTrip retrieves a group trip with all related data
func (s *Service) GetGroupTrip(id string) (*GroupTrip, error) {
	groupTrip, err := s.db.GetGroupTrip(id)
	if err != nil {
		return nil, err
	}

	// Fetch owner - TODO: Implement GetUser method
	// owner, err := s.db.GetUser(groupTrip.OwnerID)
	// if err == nil {
	//	groupTrip.Owner = owner
	// }

	// Fetch destination - TODO: Implement GetDestination method
	// destination, err := s.db.GetDestination(groupTrip.DestinationID)
	// if err == nil {
	//	groupTrip.Destination = destination
	// }

	// Fetch members
	members, err := s.db.GetGroupMembers(groupTrip.ID)
	if err == nil {
		groupTrip.Members = members
		groupTrip.MemberCount = len(members)
	}

	// Fetch expenses
	expenses, err := s.db.GetGroupExpenses(groupTrip.ID)
	if err == nil {
		groupTrip.Expenses = expenses
		total := 0.0
		for _, e := range expenses {
			total += e.Amount
		}
		groupTrip.TotalExpense = total
	}

	return groupTrip, nil
}

// GetUserGroupTrips retrieves all group trips for a user
func (s *Service) GetUserGroupTrips(userID string) ([]*GroupTrip, error) {
	return s.db.GetUserGroupTrips(userID)
}

// UpdateGroupTrip updates a group trip
func (s *Service) UpdateGroupTrip(tripID string, userID string, req *UpdateGroupTripRequest) error {
	// Verify user is owner
	groupTrip, err := s.db.GetGroupTrip(tripID)
	if err != nil {
		return err
	}

	if groupTrip.OwnerID != userID {
		return NewAPIError(ErrForbidden, "only the owner can update the group trip", "")
	}

	return s.db.UpdateGroupTrip(tripID, req)
}

// DeleteGroupTrip deletes a group trip (owner only)
func (s *Service) DeleteGroupTrip(tripID string, userID string) error {
	// Verify user is owner
	groupTrip, err := s.db.GetGroupTrip(tripID)
	if err != nil {
		return err
	}

	if groupTrip.OwnerID != userID {
		return NewAPIError(ErrForbidden, "only the owner can delete the group trip", "")
	}

	return s.db.DeleteGroupTrip(tripID)
}

// ============================================================================
// GROUP MEMBER SERVICE OPERATIONS
// ============================================================================

// InviteMemberToGroup invites a user to a group trip
func (s *Service) InviteMemberToGroup(tripID string, userID string, targetUserID string, role string) error {
	// Verify user is editor or owner
	member, err := s.db.GetGroupMember(tripID, userID)
	if err != nil {
		return err
	}

	if member.Role != GroupMemberRoleOwner && member.Role != GroupMemberRoleEditor {
		return NewAPIError(ErrForbidden, "only owners and editors can invite members", "")
	}

	// Check if user already exists in group
	existing, _ := s.db.GetGroupMember(tripID, targetUserID)
	if existing != nil {
		return NewAPIError(ErrConflict, "user is already a member of this group trip", "")
	}

	// Add member with pending status
	_, err = s.db.AddGroupMember(tripID, targetUserID, role)
	return err
}

// RespondToGroupInvite accepts or declines a group trip invitation
func (s *Service) RespondToGroupInvite(tripID string, userID string, accept bool) error {
	member, err := s.db.GetGroupMember(tripID, userID)
	if err != nil {
		return err
	}

	if member.Status != GroupMemberStatusPending {
		return NewAPIError(ErrValidationError, "member is not pending", "")
	}

	if accept {
		return s.db.UpdateGroupMemberStatus(member.ID, GroupMemberStatusActive)
	} else {
		return s.db.UpdateGroupMemberStatus(member.ID, GroupMemberStatusDeclined)
	}
}

// RemoveGroupMember removes a user from a group trip
func (s *Service) RemoveGroupMember(tripID string, userID string, targetUserID string) error {
	// Verify user is owner
	member, err := s.db.GetGroupMember(tripID, userID)
	if err != nil {
		return err
	}

	if member.Role != GroupMemberRoleOwner {
		return NewAPIError(ErrForbidden, "only the owner can remove members", "")
	}

	// Can't remove owner
	targetMember, _ := s.db.GetGroupMember(tripID, targetUserID)
	if targetMember != nil && targetMember.Role == GroupMemberRoleOwner {
		return NewAPIError(ErrValidationError, "cannot remove the owner", "")
	}

	return s.db.RemoveGroupMember(tripID, targetUserID)
}

// LeaveGroup allows a user to leave a group trip
func (s *Service) LeaveGroup(tripID string, userID string) error {
	member, err := s.db.GetGroupMember(tripID, userID)
	if err != nil {
		return err
	}

	if member.Role == GroupMemberRoleOwner {
		return NewAPIError(ErrValidationError, "owner cannot leave the group. Transfer ownership first", "")
	}

	return s.db.UpdateGroupMemberStatus(member.ID, GroupMemberStatusLeft)
}

// ============================================================================
// EXPENSE SERVICE OPERATIONS
// ============================================================================

// AddExpense adds an expense to a group trip and creates splits
func (s *Service) AddExpense(tripID string, userID string, req *CreateExpenseRequest) (*Expense, error) {
	// Verify user is member
	_, err := s.db.GetGroupMember(tripID, userID)
	if err != nil {
		return nil, NewAPIError(ErrForbidden, "user is not a member of this group trip", "")
	}

	// Validate input
	if req.Amount <= 0 {
		return nil, NewAPIError(ErrValidationError, "amount must be greater than 0", "")
	}
	if len(req.SplitAmong) == 0 {
		return nil, NewAPIError(ErrValidationError, "must split among at least 1 person", "")
	}

	// Create expense
	expense := &Expense{
		GroupTripID: tripID,
		Description: req.Description,
		Amount:      req.Amount,
		PaidBy:      userID,
		Category:    req.Category,
	}

	if err := s.db.CreateExpense(expense); err != nil {
		return nil, err
	}

	// Create splits based on split type
	if req.SplitType == "equal" {
		// Equal split
		splitAmount := req.Amount / float64(len(req.SplitAmong))
		for _, memberID := range req.SplitAmong {
			split := &ExpenseSplit{
				ExpenseID:  expense.ID,
				UserID:     memberID,
				AmountOwed: splitAmount,
			}
			if err := s.db.CreateExpenseSplit(split); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to create split for member %s: %v", memberID, err))
			}
		}
	} else if req.SplitType == "custom" {
		// Custom split
		for memberID, amount := range req.CustomSplit {
			split := &ExpenseSplit{
				ExpenseID:  expense.ID,
				UserID:     memberID,
				AmountOwed: amount,
			}
			if err := s.db.CreateExpenseSplit(split); err != nil {
				s.logger.Error(fmt.Sprintf("Failed to create custom split for member %s: %v", memberID, err))
			}
		}
	}

	return expense, nil
}

// GetGroupExpenseReport calculates the expense report for a group trip
func (s *Service) GetGroupExpenseReport(tripID string) (*ExpenseReport, error) {
	// Get all expenses for the trip
	expenses, err := s.db.GetGroupExpenses(tripID)
	if err != nil {
		return nil, err
	}

	// Get all members
	members, err := s.db.GetGroupMembers(tripID)
	if err != nil {
		return nil, err
	}

	expenseByUser := make(map[string]float64)
	owedByUser := make(map[string]float64)

	// Calculate what each user paid
	for _, expense := range expenses {
		expenseByUser[expense.PaidBy] += expense.Amount
		if owedByUser[expense.PaidBy] == 0 {
			owedByUser[expense.PaidBy] = 0
		}
	}

	// Initialize all members with 0
	for _, member := range members {
		if expenseByUser[member.UserID] == 0 {
			expenseByUser[member.UserID] = 0
		}
		if owedByUser[member.UserID] == 0 {
			owedByUser[member.UserID] = 0
		}
	}

	// Calculate what each user owes (from expense_splits)
	for _, expense := range expenses {
		splits, err := s.db.GetExpenseSplits(expense.ID)
		if err == nil {
			for _, split := range splits {
				owedByUser[split.UserID] += split.AmountOwed
			}
		}
	}

	// Calculate settlements using algorithm
	settlements := s.calculateSettlements(expenseByUser, owedByUser)

	// Calculate total expense
	totalExpense := 0.0
	for _, expense := range expenses {
		totalExpense += expense.Amount
	}

	// Generate clearing message
	clearingMessage := s.generateClearingMessage(settlements)

	return &ExpenseReport{
		TotalExpense:    totalExpense,
		ExpenseByUser:   expenseByUser,
		OwedByUser:      owedByUser,
		Settlements:     settlements,
		ClearingMessage: clearingMessage,
	}, nil
}

// ============================================================================
// EXPENSE SETTLEMENT ALGORITHM
// ============================================================================

// calculateSettlements calculates who owes whom based on expenses paid and owed
// Algorithm: Calculate each person's net balance (paid - owed) and find minimum transactions
func (s *Service) calculateSettlements(paidByUser map[string]float64, owedByUser map[string]float64) []*Settlement {
	// Calculate net balance for each person: paid - owed
	balances := make(map[string]float64)
	for userID := range paidByUser {
		balances[userID] = paidByUser[userID] - owedByUser[userID]
	}

	// Separate into creditors (positive balance) and debtors (negative balance)
	var creditors, debtors []struct {
		userID string
		amount float64
	}

	for userID, balance := range balances {
		if balance > 0.01 { // Small epsilon for float precision
			creditors = append(creditors, struct {
				userID string
				amount float64
			}{userID, balance})
		} else if balance < -0.01 {
			debtors = append(debtors, struct {
				userID string
				amount float64
			}{userID, -balance})
		}
	}

	// Sort for consistent ordering
	sort.Slice(creditors, func(i, j int) bool { return creditors[i].userID < creditors[j].userID })
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].userID < debtors[j].userID })

	// Match creditors with debtors for settlements
	var settlements []*Settlement
	c, d := 0, 0

	for c < len(creditors) && d < len(debtors) {
		creditor := creditors[c]
		debtor := debtors[d]

		// Settle minimum of creditor's remaining and debtor's remaining
		settleAmount := creditor.amount
		if debtor.amount < settleAmount {
			settleAmount = debtor.amount
		}

		settlement := &Settlement{
			DebtorID:   debtor.userID,
			CreditorID: creditor.userID,
			Amount:     settleAmount,
		}

		settlements = append(settlements, settlement)

		// Update remaining amounts
		creditors[c].amount -= settleAmount
		debtors[d].amount -= settleAmount

		// Move to next if current is settled
		if creditors[c].amount < 0.01 {
			c++
		}
		if debtors[d].amount < 0.01 {
			d++
		}
	}

	return settlements
}

// generateClearingMessage creates a human-readable summary of settlements
func (s *Service) generateClearingMessage(settlements []*Settlement) string {
	if len(settlements) == 0 {
		return "Everyone is settled up! ✨"
	}

	message := fmt.Sprintf("💰 %d transactions needed to settle:\n", len(settlements))
	for i, settlement := range settlements {
		message += fmt.Sprintf("%d. %s owes %s ₹%.2f\n", i+1, settlement.DebtorID, settlement.CreditorID, settlement.Amount)
	}

	return message
}

// ============================================================================
// POLL SERVICE OPERATIONS
// ============================================================================

// CreatePoll creates a new poll
func (s *Service) CreatePoll(tripID string, userID string, req *CreatePollRequest) (*Poll, error) {
	// Verify user is member
	_, err := s.db.GetGroupMember(tripID, userID)
	if err != nil {
		return nil, NewAPIError(ErrForbidden, "user is not a member of this group trip", "")
	}

	// Validate options
	if len(req.Options) < 2 {
		return nil, NewAPIError(ErrValidationError, "poll must have at least 2 options", "")
	}

	// Create poll
	poll := &Poll{
		GroupTripID: tripID,
		CreatedBy:   userID,
		Question:    req.Question,
		PollType:    req.PollType,
		ExpiresAt:   req.ExpiresAt,
	}

	if err := s.db.CreatePoll(poll); err != nil {
		return nil, err
	}

	// Create poll options
	for i, optionText := range req.Options {
		_, err := s.db.CreatePollOption(poll.ID, optionText, i+1)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to create poll option: %v", err))
		}
	}

	// Fetch options
	options, err := s.db.GetPollOptions(poll.ID)
	if err == nil {
		poll.Options = options
	}

	return poll, nil
}

// GetPoll retrieves a poll with all options and votes
func (s *Service) GetPoll(pollID string) (*Poll, error) {
	poll, err := s.db.GetPoll(pollID)
	if err != nil {
		return nil, err
	}

	// Fetch creator - TODO: Implement GetUser method
	// creator, err := s.db.GetUser(poll.CreatedBy)
	// if err == nil {
	//	poll.Creator = creator
	// }

	// Fetch options
	options, err := s.db.GetPollOptions(pollID)
	if err == nil {
		poll.Options = options
	}

	return poll, nil
}

// VoteOnPoll records a vote on a poll option
func (s *Service) VoteOnPoll(pollID string, optionID string, userID string) error {
	// Check if user already voted
	existingVote, err := s.db.GetUserPollVote(pollID, userID)
	if err != nil {
		return err
	}

	if existingVote != nil {
		return NewAPIError(ErrConflict, "user has already voted on this poll", "")
	}

	// Create vote
	return s.db.CreatePollVote(optionID, userID)
}
