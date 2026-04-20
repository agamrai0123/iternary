package itinerary

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/itinerary-backend/itinerary/common"
)

// ============================================================================
// GROUP TRIP DATABASE OPERATIONS
// ============================================================================

// CreateGroupTrip creates a new group trip
func (db *common.Database) CreateGroupTrip(groupTrip *GroupTrip) error {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO group_trips (id, title, destination_id, owner_id, budget, duration, start_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, groupTrip.Title, groupTrip.DestinationID, groupTrip.OwnerID, groupTrip.Budget, groupTrip.Duration, groupTrip.StartDate, "draft", now, now)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to create group trip: "+err.Error(), "")
	}

	groupTrip.ID = id
	groupTrip.Status = "draft"
	groupTrip.CreatedAt = now
	groupTrip.UpdatedAt = now

	return nil
}

// GetGroupTrip retrieves a group trip by ID
func (db *common.Database) GetGroupTrip(id string) (*GroupTrip, error) {
	query := `
		SELECT id, title, destination_id, owner_id, budget, duration, start_date, status, created_at, updated_at
		FROM group_trips
		WHERE id = ?
	`

	row := db.conn.QueryRow(query, id)
	groupTrip := &GroupTrip{}

	err := row.Scan(&groupTrip.ID, &groupTrip.Title, &groupTrip.DestinationID, &groupTrip.OwnerID, &groupTrip.Budget, &groupTrip.Duration, &groupTrip.StartDate, &groupTrip.Status, &groupTrip.CreatedAt, &groupTrip.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAPIError(ErrNotFound, "group trip not found", "")
		}
		return nil, NewAPIError(ErrDatabaseError, "failed to get group trip: "+err.Error(), "")
	}

	return groupTrip, nil
}

// GetUserGroupTrips retrieves all group trips for a user (as member or owner)
func (db *common.Database) GetUserGroupTrips(userID string) ([]*GroupTrip, error) {
	query := `
		SELECT DISTINCT gt.id, gt.title, gt.destination_id, gt.owner_id, gt.budget, gt.duration, gt.start_date, gt.status, gt.created_at, gt.updated_at
		FROM group_trips gt
		JOIN group_members gm ON gt.id = gm.group_trip_id
		WHERE gm.user_id = ? AND gm.status = 'active'
		ORDER BY gt.created_at DESC
	`

	rows, err := db.conn.Query(query, userID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get group trips: "+err.Error(), "")
	}
	defer rows.Close()

	var groupTrips []*GroupTrip
	for rows.Next() {
		groupTrip := &GroupTrip{}
		err := rows.Scan(&groupTrip.ID, &groupTrip.Title, &groupTrip.DestinationID, &groupTrip.OwnerID, &groupTrip.Budget, &groupTrip.Duration, &groupTrip.StartDate, &groupTrip.Status, &groupTrip.CreatedAt, &groupTrip.UpdatedAt)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan group trip: "+err.Error(), "")
		}
		groupTrips = append(groupTrips, groupTrip)
	}

	return groupTrips, nil
}

// UpdateGroupTrip updates a group trip
func (db *common.Database) UpdateGroupTrip(id string, updates *UpdateGroupTripRequest) error {
	query := `
		UPDATE group_trips
		SET title = COALESCE(?, title),
		    budget = COALESCE(?, budget),
		    duration = COALESCE(?, duration),
		    start_date = COALESCE(?, start_date),
		    status = COALESCE(?, status),
		    updated_at = ?
		WHERE id = ?
	`

	_, err := db.conn.Exec(query, updates.Title, updates.Budget, updates.Duration, updates.StartDate, updates.Status, time.Now(), id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to update group trip: "+err.Error(), "")
	}

	return nil
}

// DeleteGroupTrip deletes a group trip
func (db *common.Database) DeleteGroupTrip(id string) error {
	query := `DELETE FROM group_trips WHERE id = ?`

	_, err := db.conn.Exec(query, id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to delete group trip", err.Error())
	}

	return nil
}

// ============================================================================
// GROUP MEMBER DATABASE OPERATIONS
// ============================================================================

// AddGroupMember adds a user to a group trip
func (db *common.Database) AddGroupMember(groupTripID string, userID string, role string) (*GroupMember, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO group_members (id, group_trip_id, user_id, role, joined_at, status)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, groupTripID, userID, role, now, GroupMemberStatusPending)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to add group member", err.Error())
	}

	return &GroupMember{
		ID:          id,
		GroupTripID: groupTripID,
		UserID:      userID,
		Role:        role,
		JoinedAt:    now,
		Status:      GroupMemberStatusPending,
	}, nil
}

// GetGroupMembers retrieves all members of a group trip
func (db *common.Database) GetGroupMembers(groupTripID string) ([]*GroupMember, error) {
	query := `
		SELECT id, group_trip_id, user_id, role, joined_at, status
		FROM group_members
		WHERE group_trip_id = ?
		ORDER BY joined_at ASC
	`

	rows, err := db.conn.Query(query, groupTripID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get group members", err.Error())
	}
	defer rows.Close()

	var members []*GroupMember
	for rows.Next() {
		member := &GroupMember{}
		err := rows.Scan(&member.ID, &member.GroupTripID, &member.UserID, &member.Role, &member.JoinedAt, &member.Status)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan group member", err.Error())
		}
		members = append(members, member)
	}

	return members, nil
}

// GetGroupMember retrieves a specific group member
func (db *common.Database) GetGroupMember(groupTripID string, userID string) (*GroupMember, error) {
	query := `
		SELECT id, group_trip_id, user_id, role, joined_at, status
		FROM group_members
		WHERE group_trip_id = ? AND user_id = ?
	`

	row := db.conn.QueryRow(query, groupTripID, userID)
	member := &GroupMember{}

	err := row.Scan(&member.ID, &member.GroupTripID, &member.UserID, &member.Role, &member.JoinedAt, &member.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAPIError(ErrNotFound, "group member not found", "")
		}
		return nil, NewAPIError(ErrDatabaseError, "failed to get group member", err.Error())
	}

	return member, nil
}

// UpdateGroupMemberRole updates a member's role
func (db *common.Database) UpdateGroupMemberRole(id string, role string) error {
	query := `UPDATE group_members SET role = ? WHERE id = ?`

	_, err := db.conn.Exec(query, role, id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to update group member role", err.Error())
	}

	return nil
}

// UpdateGroupMemberStatus updates a member's status (accept/decline/leave)
func (db *common.Database) UpdateGroupMemberStatus(id string, status string) error {
	query := `UPDATE group_members SET status = ? WHERE id = ?`

	_, err := db.conn.Exec(query, status, id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to update group member status", err.Error())
	}

	return nil
}

// RemoveGroupMember removes a user from a group trip
func (db *common.Database) RemoveGroupMember(groupTripID string, userID string) error {
	query := `DELETE FROM group_members WHERE group_trip_id = ? AND user_id = ?`

	_, err := db.conn.Exec(query, groupTripID, userID)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to remove group member", err.Error())
	}

	return nil
}

// ============================================================================
// EXPENSE DATABASE OPERATIONS
// ============================================================================

// CreateExpense creates a new expense
func (db *common.Database) CreateExpense(expense *Expense) error {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO expenses (id, group_trip_id, description, amount, paid_by, category, paid_date, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if expense.PaidDate == nil {
		today := time.Now()
		expense.PaidDate = &today
	}

	_, err := db.conn.Exec(query, id, expense.GroupTripID, expense.Description, expense.Amount, expense.PaidBy, expense.Category, expense.PaidDate, ExpenseStatusPending, now)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to create expense", err.Error())
	}

	expense.ID = id
	expense.Status = ExpenseStatusPending
	expense.CreatedAt = now

	return nil
}

// GetExpense retrieves an expense by ID
func (db *common.Database) GetExpense(id string) (*Expense, error) {
	query := `
		SELECT id, group_trip_id, description, amount, paid_by, category, paid_date, status, created_at
		FROM expenses
		WHERE id = ?
	`

	row := db.conn.QueryRow(query, id)
	expense := &Expense{}

	err := row.Scan(&expense.ID, &expense.GroupTripID, &expense.Description, &expense.Amount, &expense.PaidBy, &expense.Category, &expense.PaidDate, &expense.Status, &expense.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAPIError(ErrNotFound, "expense not found", "")
		}
		return nil, NewAPIError(ErrDatabaseError, "failed to get expense", err.Error())
	}

	return expense, nil
}

// GetGroupExpenses retrieves all expenses for a group trip
func (db *common.Database) GetGroupExpenses(groupTripID string) ([]*Expense, error) {
	query := `
		SELECT id, group_trip_id, description, amount, paid_by, category, paid_date, status, created_at
		FROM expenses
		WHERE group_trip_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query, groupTripID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get group expenses", err.Error())
	}
	defer rows.Close()

	var expenses []*Expense
	for rows.Next() {
		expense := &Expense{}
		err := rows.Scan(&expense.ID, &expense.GroupTripID, &expense.Description, &expense.Amount, &expense.PaidBy, &expense.Category, &expense.PaidDate, &expense.Status, &expense.CreatedAt)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan expense", err.Error())
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// UpdateExpenseStatus marks an expense as settled or pending
func (db *common.Database) UpdateExpenseStatus(id string, status string) error {
	query := `UPDATE expenses SET status = ? WHERE id = ?`

	_, err := db.conn.Exec(query, status, id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to update expense status", err.Error())
	}

	return nil
}

// ============================================================================
// EXPENSE SPLIT DATABASE OPERATIONS
// ============================================================================

// CreateExpenseSplit creates a split record for an expense
func (db *common.Database) CreateExpenseSplit(split *ExpenseSplit) error {
	id := uuid.New().String()

	query := `
		INSERT INTO expense_splits (id, expense_id, user_id, amount_owed)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, split.ExpenseID, split.UserID, split.AmountOwed)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to create expense split", err.Error())
	}

	split.ID = id
	return nil
}

// GetExpenseSplits retrieves all splits for an expense
func (db *common.Database) GetExpenseSplits(expenseID string) ([]*ExpenseSplit, error) {
	query := `
		SELECT id, expense_id, user_id, amount_owed
		FROM expense_splits
		WHERE expense_id = ?
	`

	rows, err := db.conn.Query(query, expenseID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get expense splits", err.Error())
	}
	defer rows.Close()

	var splits []*ExpenseSplit
	for rows.Next() {
		split := &ExpenseSplit{}
		err := rows.Scan(&split.ID, &split.ExpenseID, &split.UserID, &split.AmountOwed)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan expense split", err.Error())
		}
		splits = append(splits, split)
	}

	return splits, nil
}

// GetUserExpensesByTrip retrieves how much a user paid and owes in a group trip
func (db *common.Database) GetUserExpensesByTrip(groupTripID string, userID string) (paidAmount float64, owedAmount float64, err error) {
	// Get amount paid by user
	paidQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE group_trip_id = ? AND paid_by = ?
	`

	row := db.conn.QueryRow(paidQuery, groupTripID, userID)
	err = row.Scan(&paidAmount)
	if err != nil {
		return 0, 0, NewAPIError(ErrDatabaseError, "failed to get paid amount", err.Error())
	}

	// Get amount owed by user
	owedQuery := `
		SELECT COALESCE(SUM(es.amount_owed), 0)
		FROM expense_splits es
		JOIN expenses e ON es.expense_id = e.id
		WHERE e.group_trip_id = ? AND es.user_id = ?
	`

	row = db.conn.QueryRow(owedQuery, groupTripID, userID)
	err = row.Scan(&owedAmount)
	if err != nil {
		return 0, 0, NewAPIError(ErrDatabaseError, "failed to get owed amount", err.Error())
	}

	return paidAmount, owedAmount, nil
}

// ============================================================================
// POLL DATABASE OPERATIONS
// ============================================================================

// CreatePoll creates a new poll
func (db *common.Database) CreatePoll(poll *Poll) error {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO polls (id, group_trip_id, created_by, question, poll_type, status, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, poll.GroupTripID, poll.CreatedBy, poll.Question, poll.PollType, PollStatusActive, poll.ExpiresAt, now)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to create poll", err.Error())
	}

	poll.ID = id
	poll.Status = PollStatusActive
	poll.CreatedAt = now

	return nil
}

// GetPoll retrieves a poll by ID
func (db *common.Database) GetPoll(id string) (*Poll, error) {
	query := `
		SELECT id, group_trip_id, created_by, question, poll_type, status, expires_at, created_at
		FROM polls
		WHERE id = ?
	`

	row := db.conn.QueryRow(query, id)
	poll := &Poll{}

	err := row.Scan(&poll.ID, &poll.GroupTripID, &poll.CreatedBy, &poll.Question, &poll.PollType, &poll.Status, &poll.ExpiresAt, &poll.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAPIError(ErrNotFound, "poll not found", "")
		}
		return nil, NewAPIError(ErrDatabaseError, "failed to get poll", err.Error())
	}

	return poll, nil
}

// GetGroupPolls retrieves all polls for a group trip
func (db *common.Database) GetGroupPolls(groupTripID string) ([]*Poll, error) {
	query := `
		SELECT id, group_trip_id, created_by, question, poll_type, status, expires_at, created_at
		FROM polls
		WHERE group_trip_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query, groupTripID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get group polls", err.Error())
	}
	defer rows.Close()

	var polls []*Poll
	for rows.Next() {
		poll := &Poll{}
		err := rows.Scan(&poll.ID, &poll.GroupTripID, &poll.CreatedBy, &poll.Question, &poll.PollType, &poll.Status, &poll.ExpiresAt, &poll.CreatedAt)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan poll", err.Error())
		}
		polls = append(polls, poll)
	}

	return polls, nil
}

// UpdatePollStatus updates a poll's status (active -> locked -> resolved)
func (db *common.Database) UpdatePollStatus(id string, status string) error {
	query := `UPDATE polls SET status = ? WHERE id = ?`

	_, err := db.conn.Exec(query, status, id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to update poll status", err.Error())
	}

	return nil
}

// ============================================================================
// POLL OPTION DATABASE OPERATIONS
// ============================================================================

// CreatePollOption creates a poll option
func (db *common.Database) CreatePollOption(pollID string, optionText string, sequence int) (*PollOption, error) {
	id := uuid.New().String()

	query := `
		INSERT INTO poll_options (id, poll_id, option_text, vote_count, sequence)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, pollID, optionText, 0, sequence)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to create poll option", err.Error())
	}

	return &PollOption{
		ID:         id,
		PollID:     pollID,
		OptionText: optionText,
		VoteCount:  0,
		Sequence:   sequence,
	}, nil
}

// GetPollOptions retrieves all options for a poll
func (db *common.Database) GetPollOptions(pollID string) ([]*PollOption, error) {
	query := `
		SELECT id, poll_id, option_text, vote_count, sequence
		FROM poll_options
		WHERE poll_id = ?
		ORDER BY sequence ASC
	`

	rows, err := db.conn.Query(query, pollID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get poll options", err.Error())
	}
	defer rows.Close()

	var options []*PollOption
	for rows.Next() {
		option := &PollOption{}
		err := rows.Scan(&option.ID, &option.PollID, &option.OptionText, &option.VoteCount, &option.Sequence)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan poll option", err.Error())
		}
		options = append(options, option)
	}

	return options, nil
}

// ============================================================================
// POLL VOTE DATABASE OPERATIONS
// ============================================================================

// CreatePollVote records a vote on a poll option
func (db *common.Database) CreatePollVote(optionID string, userID string) error {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO poll_votes (id, poll_option_id, user_id, voted_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, optionID, userID, now)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to create poll vote", err.Error())
	}

	// Increment vote count
	updateQuery := `UPDATE poll_options SET vote_count = vote_count + 1 WHERE id = ?`
	_, err = db.conn.Exec(updateQuery, optionID)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to increment vote count", err.Error())
	}

	return nil
}

// GetUserPollVote checks if user already voted on this poll
func (db *common.Database) GetUserPollVote(pollID string, userID string) (*PollVote, error) {
	query := `
		SELECT pv.id, pv.poll_option_id, pv.user_id, pv.voted_at
		FROM poll_votes pv
		JOIN poll_options po ON pv.poll_option_id = po.id
		WHERE po.poll_id = ? AND pv.user_id = ?
		LIMIT 1
	`

	row := db.conn.QueryRow(query, pollID, userID)
	vote := &PollVote{}

	err := row.Scan(&vote.ID, &vote.PollOptionID, &vote.UserID, &vote.VotedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User hasn't voted yet
		}
		return nil, NewAPIError(ErrDatabaseError, "failed to get poll vote", err.Error())
	}

	return vote, nil
}

// ============================================================================
// SETTLEMENT DATABASE OPERATIONS
// ============================================================================

// CreateSettlement creates a settlement record
func (db *common.Database) CreateSettlement(groupTripID string, debtorID string, creditorID string, amount float64) (*Settlement, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO settlements (id, group_trip_id, debtor_id, creditor_id, amount, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query, id, groupTripID, debtorID, creditorID, amount, SettlementStatusPending, now)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to create settlement", err.Error())
	}

	return &Settlement{
		ID:          id,
		GroupTripID: groupTripID,
		DebtorID:    debtorID,
		CreditorID:  creditorID,
		Amount:      amount,
		Status:      SettlementStatusPending,
		CreatedAt:   now,
	}, nil
}

// GetSettlements retrieves all settlements for a group trip
func (db *common.Database) GetSettlements(groupTripID string) ([]*Settlement, error) {
	query := `
		SELECT id, group_trip_id, debtor_id, creditor_id, amount, status, settled_at, created_at
		FROM settlements
		WHERE group_trip_id = ?
		ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query, groupTripID)
	if err != nil {
		return nil, NewAPIError(ErrDatabaseError, "failed to get settlements", err.Error())
	}
	defer rows.Close()

	var settlements []*Settlement
	for rows.Next() {
		settlement := &Settlement{}
		err := rows.Scan(&settlement.ID, &settlement.GroupTripID, &settlement.DebtorID, &settlement.CreditorID, &settlement.Amount, &settlement.Status, &settlement.SettledAt, &settlement.CreatedAt)
		if err != nil {
			return nil, NewAPIError(ErrDatabaseError, "failed to scan settlement", err.Error())
		}
		settlements = append(settlements, settlement)
	}

	return settlements, nil
}

// MarkSettlementSettled marks a settlement as settled
func (db *common.Database) MarkSettlementSettled(id string) error {
	now := time.Now()
	query := `UPDATE settlements SET status = ?, settled_at = ? WHERE id = ?`

	_, err := db.conn.Exec(query, SettlementStatusSettled, now, id)
	if err != nil {
		return NewAPIError(ErrDatabaseError, "failed to mark settlement settled", err.Error())
	}

	return nil
}


