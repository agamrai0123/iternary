package itinerary

import "time"

// ============================================================================
// GROUP TRIP MODELS
// ============================================================================

// GroupTrip represents a trip planned by multiple users
type GroupTrip struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	DestinationID string    `json:"destination_id"`
	OwnerID       string    `json:"owner_id"`
	Budget        float64   `json:"budget"`
	Duration      int       `json:"duration"`
	StartDate     *time.Time `json:"start_date"`
	Status        string    `json:"status"` // draft, planning, published, completed
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Joined data (optional, for API responses)
	Owner        *User           `json:"owner,omitempty"`
	Destination  *Destination    `json:"destination,omitempty"`
	Members      []*GroupMember  `json:"members,omitempty"`
	Expenses     []*Expense      `json:"expenses,omitempty"`
	MemberCount  int             `json:"member_count,omitempty"`
	TotalExpense float64         `json:"total_expense,omitempty"`
}

// CreateGroupTripRequest is the request body for creating a group trip
type CreateGroupTripRequest struct {
	Title          string     `json:"title" binding:"required,min=3,max=255"`
	DestinationID  string     `json:"destination_id" binding:"required"`
	Budget         float64    `json:"budget" binding:"required,gt=0"`
	Duration       int        `json:"duration" binding:"required,gt=0"`
	StartDate      *time.Time `json:"start_date"`
	InitialMembers []string   `json:"initial_members"` // User IDs to invite
}

// UpdateGroupTripRequest is the request body for updating a group trip
type UpdateGroupTripRequest struct {
	Title      string     `json:"title" binding:"min=3,max=255"`
	Budget     float64    `json:"budget" binding:"gt=0"`
	Duration   int        `json:"duration" binding:"gt=0"`
	StartDate  *time.Time `json:"start_date"`
	Status     string     `json:"status" binding:"omitempty,oneof=draft planning published completed"`
}

// ============================================================================
// GROUP MEMBER MODELS
// ============================================================================

// GroupMember represents a user's membership in a group trip
type GroupMember struct {
	ID          string    `json:"id"`
	GroupTripID string    `json:"group_trip_id"`
	UserID      string    `json:"user_id"`
	Role        string    `json:"role"` // owner, editor, member, viewer
	JoinedAt    time.Time `json:"joined_at"`
	Status      string    `json:"status"` // pending, active, declined, left

	// Joined data
	User *User `json:"user,omitempty"`
}

// AddGroupMemberRequest is the request body for adding a member to a group
type AddGroupMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=owner editor member viewer"`
}

// UpdateGroupMemberRequest is the request body for updating a member's role
type UpdateGroupMemberRequest struct {
	Role string `json:"role" binding:"required,oneof=owner editor member viewer"`
}

// RespondToInvitationRequest is the request body for responding to group invitation
type RespondToInvitationRequest struct {
	GroupTripID string `json:"group_trip_id" binding:"required"`
	Status      string `json:"status" binding:"required,oneof=active declined"`
}

// ============================================================================
// EXPENSE MODELS
// ============================================================================

// Expense represents a single expense in a group trip
type Expense struct {
	ID          string    `json:"id"`
	GroupTripID string    `json:"group_trip_id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	PaidBy      string    `json:"paid_by"`
	Category    string    `json:"category"` // accommodation, food, transport, activity, other
	PaidDate    *time.Time `json:"paid_date"`
	Status      string    `json:"status"` // pending, settled
	CreatedAt   time.Time `json:"created_at"`

	// Joined data
	PaidByUser *User          `json:"paid_by_user,omitempty"`
	Splits     []*ExpenseSplit `json:"splits,omitempty"`
}

// CreateExpenseRequest is the request body for creating an expense
type CreateExpenseRequest struct {
	Description string          `json:"description" binding:"required,min=3,max=255"`
	Amount      float64         `json:"amount" binding:"required,gt=0"`
	Category    string          `json:"category" binding:"required,oneof=accommodation food transport activity other"`
	SplitType   string          `json:"split_type" binding:"required,oneof=equal custom"` // How to split
	SplitAmong  []string        `json:"split_among" binding:"required,min=1"`            // User IDs to split among
	CustomSplit map[string]float64 `json:"custom_split,omitempty"` // user_id -> amount for custom split
}

// ExpenseSplit represents how an expense is split among members
type ExpenseSplit struct {
	ID         string  `json:"id"`
	ExpenseID  string  `json:"expense_id"`
	UserID     string  `json:"user_id"`
	AmountOwed float64 `json:"amount_owed"`

	// Joined data
	User *User `json:"user,omitempty"`
}

// ExpenseReport contains summary of expenses for a group trip
type ExpenseReport struct {
	TotalExpense    float64            `json:"total_expense"`
	ExpenseByUser   map[string]float64 `json:"expense_by_user"` // user_id -> amount paid
	OwedByUser      map[string]float64 `json:"owed_by_user"`   // user_id -> amount owed
	Settlements     []*Settlement      `json:"settlements"`
	ClearingMessage string             `json:"clearing_message"` // Explanation of who owes whom
}

// ============================================================================
// POLL MODELS
// ============================================================================

// Poll represents a poll/vote for group trip decisions
type Poll struct {
	ID          string        `json:"id"`
	GroupTripID string        `json:"group_trip_id"`
	CreatedBy   string        `json:"created_by"`
	Question    string        `json:"question"`
	PollType    string        `json:"poll_type"` // itinerary, budget, date, activity, destination
	Status      string        `json:"status"`   // active, locked, resolved
	ExpiresAt   *time.Time    `json:"expires_at"`
	CreatedAt   time.Time     `json:"created_at"`

	// Joined data
	Creator *User         `json:"creator,omitempty"`
	Options []*PollOption `json:"options,omitempty"`
}

// PollOption represents a single option in a poll
type PollOption struct {
	ID        string      `json:"id"`
	PollID    string      `json:"poll_id"`
	OptionText string     `json:"option_text"`
	VoteCount int         `json:"vote_count"`
	Sequence  int         `json:"sequence"`

	// Joined data
	Votes []*PollVote `json:"votes,omitempty"`
}

// PollVote represents a single vote on a poll option
type PollVote struct {
	ID          string    `json:"id"`
	PollOptionID string   `json:"poll_option_id"`
	UserID      string    `json:"user_id"`
	VotedAt     time.Time `json:"voted_at"`

	// Joined data
	User *User `json:"user,omitempty"`
}

// CreatePollRequest is the request body for creating a poll
type CreatePollRequest struct {
	Question  string   `json:"question" binding:"required,min=5,max=500"`
	PollType  string   `json:"poll_type" binding:"required,oneof=itinerary budget date activity destination"`
	Options   []string `json:"options" binding:"required,min=2,max=10"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// VotePollRequest is the request body for voting on a poll
type VotePollRequest struct {
	OptionID string `json:"option_id" binding:"required"`
}

// ============================================================================
// SETTLEMENT MODELS
// ============================================================================

// Settlement represents who owes money to whom in a group trip
type Settlement struct {
	ID          string    `json:"id"`
	GroupTripID string    `json:"group_trip_id"`
	DebtorID    string    `json:"debtor_id"`
	CreditorID  string    `json:"creditor_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // pending, settled
	SettledAt   *time.Time `json:"settled_at"`
	CreatedAt   time.Time `json:"created_at"`

	// Joined data
	Debtor   *User `json:"debtor,omitempty"`
	Creditor *User `json:"creditor,omitempty"`
}

// SettleExpenseRequest is the request body for marking a settlement as settled
type SettleExpenseRequest struct {
	SettlementID string `json:"settlement_id" binding:"required"`
}

// ============================================================================
// GROUP CONSTANTS
// ============================================================================

const (
	// Group trip statuses
	GroupTripStatusDraft      = "draft"
	GroupTripStatusPlanning   = "planning"
	GroupTripStatusPublished  = "published"
	GroupTripStatusCompleted  = "completed"

	// Group member roles
	GroupMemberRoleOwner  = "owner"
	GroupMemberRoleEditor = "editor"
	GroupMemberRoleMember = "member"
	GroupMemberRoleViewer = "viewer"

	// Group member statuses
	GroupMemberStatusPending  = "pending"
	GroupMemberStatusActive   = "active"
	GroupMemberStatusDeclined = "declined"
	GroupMemberStatusLeft     = "left"

	// Expense categories
	ExpenseCategoryAccommodation = "accommodation"
	ExpenseCategoryFood          = "food"
	ExpenseCategoryTransport     = "transport"
	ExpenseCategoryActivity      = "activity"
	ExpenseCategoryOther         = "other"

	// Expense statuses
	ExpenseStatusPending = "pending"
	ExpenseStatusSettled = "settled"

	// Poll types
	PollTypeItinerary   = "itinerary"
	PollTypeBudget      = "budget"
	PollTypeDate        = "date"
	PollTypeActivity    = "activity"
	PollTypeDestination = "destination"

	// Poll statuses
	PollStatusActive   = "active"
	PollStatusLocked   = "locked"
	PollStatusResolved = "resolved"

	// Settlement statuses
	SettlementStatusPending = "pending"
	SettlementStatusSettled = "settled"
)

// ============================================================================
// VALIDATION FUNCTIONS
// ============================================================================

// Validate checks if GroupTrip is valid
func (gt *GroupTrip) Validate() error {
	if gt.Title == "" {
		return NewAPIError(ErrValidationError, "title is required", "")
	}
	if gt.Budget <= 0 {
		return NewAPIError(ErrValidationError, "budget must be greater than 0", "")
	}
	if gt.Duration <= 0 {
		return NewAPIError(ErrValidationError, "duration must be greater than 0", "")
	}
	return nil
}

// Validate checks if Expense is valid
func (e *Expense) Validate() error {
	if e.Description == "" {
		return NewAPIError(ErrValidationError, "description is required", "")
	}
	if e.Amount <= 0 {
		return NewAPIError(ErrValidationError, "amount must be greater than 0", "")
	}
	if e.PaidBy == "" {
		return NewAPIError(ErrValidationError, "paid_by is required", "")
	}
	return nil
}

// Validate checks if ExpenseSplit is valid
func (es *ExpenseSplit) Validate() error {
	if es.AmountOwed < 0 {
		return NewAPIError(ErrValidationError, "amount_owed cannot be negative", "")
	}
	return nil
}

// Validate checks if Poll is valid
func (p *Poll) Validate() error {
	if p.Question == "" {
		return NewAPIError(ErrValidationError, "question is required", "")
	}
	if len(p.Options) < 2 {
		return NewAPIError(ErrValidationError, "poll must have at least 2 options", "")
	}
	return nil
}

// Validate checks if Settlement is valid
func (s *Settlement) Validate() error {
	if s.DebtorID == "" {
		return NewAPIError(ErrValidationError, "debtor_id is required", "")
	}
	if s.CreditorID == "" {
		return NewAPIError(ErrValidationError, "creditor_id is required", "")
	}
	if s.Amount <= 0 {
		return NewAPIError(ErrValidationError, "amount must be greater than 0", "")
	}
	if s.DebtorID == s.CreditorID {
		return NewAPIError(ErrValidationError, "debtor and creditor cannot be the same person", "")
	}
	return nil
}
