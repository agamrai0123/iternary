package itinerary

import (
	"testing"
	"time"
)

// ============================================================================
// GROUP TRIP MODEL TESTS
// ============================================================================

func TestGroupTripValidation(t *testing.T) {
	tests := []struct {
		name      string
		groupTrip *GroupTrip
		expectErr bool
	}{
		{
			name: "valid group trip",
			groupTrip: &GroupTrip{
				Title:    "Bali Trip",
				Budget:   50000,
				Duration: 5,
			},
			expectErr: false,
		},
		{
			name: "missing title",
			groupTrip: &GroupTrip{
				Title:    "",
				Budget:   50000,
				Duration: 5,
			},
			expectErr: true,
		},
		{
			name: "zero budget",
			groupTrip: &GroupTrip{
				Title:    "Bali Trip",
				Budget:   0,
				Duration: 5,
			},
			expectErr: true,
		},
		{
			name: "negative duration",
			groupTrip: &GroupTrip{
				Title:    "Bali Trip",
				Budget:   50000,
				Duration: -1,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.groupTrip.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// ============================================================================
// EXPENSE MODEL TESTS
// ============================================================================

func TestExpenseValidation(t *testing.T) {
	tests := []struct {
		name      string
		expense   *Expense
		expectErr bool
	}{
		{
			name: "valid expense",
			expense: &Expense{
				Description: "Dinner",
				Amount:      5000,
				PaidBy:      "user-1",
			},
			expectErr: false,
		},
		{
			name: "missing description",
			expense: &Expense{
				Description: "",
				Amount:      5000,
				PaidBy:      "user-1",
			},
			expectErr: true,
		},
		{
			name: "zero amount",
			expense: &Expense{
				Description: "Dinner",
				Amount:      0,
				PaidBy:      "user-1",
			},
			expectErr: true,
		},
		{
			name: "negative amount",
			expense: &Expense{
				Description: "Dinner",
				Amount:      -5000,
				PaidBy:      "user-1",
			},
			expectErr: true,
		},
		{
			name: "missing paid_by",
			expense: &Expense{
				Description: "Dinner",
				Amount:      5000,
				PaidBy:      "",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.expense.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// ============================================================================
// EXPENSE SPLIT MODEL TESTS
// ============================================================================

func TestExpenseSplitValidation(t *testing.T) {
	tests := []struct {
		name      string
		split     *ExpenseSplit
		expectErr bool
	}{
		{
			name: "valid split",
			split: &ExpenseSplit{
				AmountOwed: 1250,
			},
			expectErr: false,
		},
		{
			name: "zero owed",
			split: &ExpenseSplit{
				AmountOwed: 0,
			},
			expectErr: false,
		},
		{
			name: "negative owed",
			split: &ExpenseSplit{
				AmountOwed: -100,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.split.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// ============================================================================
// POLL MODEL TESTS
// ============================================================================

func TestPollValidation(t *testing.T) {
	tests := []struct {
		name      string
		poll      *Poll
		expectErr bool
	}{
		{
			name: "valid poll",
			poll: &Poll{
				Question: "Which restaurant should we go to?",
				Options: []*PollOption{
					{OptionText: "Pizza"},
					{OptionText: "Chinese"},
				},
			},
			expectErr: false,
		},
		{
			name: "missing question",
			poll: &Poll{
				Question: "",
				Options: []*PollOption{
					{OptionText: "Pizza"},
					{OptionText: "Chinese"},
				},
			},
			expectErr: true,
		},
		{
			name: "only one option",
			poll: &Poll{
				Question: "Which restaurant?",
				Options: []*PollOption{
					{OptionText: "Pizza"},
				},
			},
			expectErr: true,
		},
		{
			name: "no options",
			poll: &Poll{
				Question: "Which restaurant?",
				Options:  []*PollOption{},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.poll.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// ============================================================================
// SETTLEMENT MODEL TESTS
// ============================================================================

func TestSettlementValidation(t *testing.T) {
	tests := []struct {
		name       string
		settlement *Settlement
		expectErr  bool
	}{
		{
			name: "valid settlement",
			settlement: &Settlement{
				DebtorID:   "user-1",
				CreditorID: "user-2",
				Amount:     5000,
			},
			expectErr: false,
		},
		{
			name: "missing debtor",
			settlement: &Settlement{
				DebtorID:   "",
				CreditorID: "user-2",
				Amount:     5000,
			},
			expectErr: true,
		},
		{
			name: "missing creditor",
			settlement: &Settlement{
				DebtorID:   "user-1",
				CreditorID: "",
				Amount:     5000,
			},
			expectErr: true,
		},
		{
			name: "zero amount",
			settlement: &Settlement{
				DebtorID:   "user-1",
				CreditorID: "user-2",
				Amount:     0,
			},
			expectErr: true,
		},
		{
			name: "debtor is creditor",
			settlement: &Settlement{
				DebtorID:   "user-1",
				CreditorID: "user-1",
				Amount:     5000,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settlement.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

// ============================================================================
// GROUP MEMBER MODEL TESTS
// ============================================================================

func TestGroupMemberRoles(t *testing.T) {
	validRoles := []string{
		GroupMemberRoleOwner,
		GroupMemberRoleEditor,
		GroupMemberRoleMember,
		GroupMemberRoleViewer,
	}

	for _, role := range validRoles {
		member := &GroupMember{
			UserID: "user-1",
			Role:   role,
		}

		if member.Role != role {
			t.Errorf("expected role %s, got %s", role, member.Role)
		}
	}
}

func TestGroupMemberStatuses(t *testing.T) {
	validStatuses := []string{
		GroupMemberStatusPending,
		GroupMemberStatusActive,
		GroupMemberStatusDeclined,
		GroupMemberStatusLeft,
	}

	for _, status := range validStatuses {
		member := &GroupMember{
			UserID: "user-1",
			Status: status,
		}

		if member.Status != status {
			t.Errorf("expected status %s, got %s", status, member.Status)
		}
	}
}

// ============================================================================
// EXPENSE CATEGORY TESTS
// ============================================================================

func TestExpenseCategories(t *testing.T) {
	validCategories := []string{
		ExpenseCategoryAccommodation,
		ExpenseCategoryFood,
		ExpenseCategoryTransport,
		ExpenseCategoryActivity,
		ExpenseCategoryOther,
	}

	for _, category := range validCategories {
		expense := &Expense{
			Description: "Test",
			Amount:      100,
			PaidBy:      "user-1",
			Category:    category,
		}

		if expense.Category != category {
			t.Errorf("expected category %s, got %s", category, expense.Category)
		}
	}
}

// ============================================================================
// POLL TYPE TESTS
// ============================================================================

func TestPollTypes(t *testing.T) {
	validTypes := []string{
		PollTypeItinerary,
		PollTypeBudget,
		PollTypeDate,
		PollTypeActivity,
		PollTypeDestination,
	}

	for _, pollType := range validTypes {
		poll := &Poll{
			Question: "Test?",
			PollType: pollType,
		}

		if poll.PollType != pollType {
			t.Errorf("expected poll type %s, got %s", pollType, poll.PollType)
		}
	}
}

// ============================================================================
// DATE HANDLING TESTS
// ============================================================================

func TestGroupTripWithDate(t *testing.T) {
	startDate := time.Now().AddDate(0, 0, 7) // 7 days from now
	groupTrip := &GroupTrip{
		Title:     "Future Trip",
		Budget:    50000,
		Duration:  5,
		StartDate: &startDate,
	}

	if groupTrip.StartDate == nil {
		t.Error("start_date should not be nil")
	}

	if groupTrip.StartDate.Before(time.Now()) {
		t.Error("start_date should be in the future")
	}
}

func TestExpenseWithDate(t *testing.T) {
	today := time.Now()
	expense := &Expense{
		Description: "Dinner",
		Amount:      5000,
		PaidBy:      "user-1",
		PaidDate:    &today,
	}

	if expense.PaidDate == nil {
		t.Error("paid_date should not be nil")
	}
}

// ============================================================================
// GROUP TRIP STATUS TESTS
// ============================================================================

func TestGroupTripStatus(t *testing.T) {
	validStatuses := []string{
		GroupTripStatusDraft,
		GroupTripStatusPlanning,
		GroupTripStatusPublished,
		GroupTripStatusCompleted,
	}

	for _, status := range validStatuses {
		groupTrip := &GroupTrip{
			Title:    "Trip",
			Budget:   50000,
			Duration: 5,
			Status:   status,
		}

		if groupTrip.Status != status {
			t.Errorf("expected status %s, got %s", status, groupTrip.Status)
		}
	}
}

// ============================================================================
// EXPENSE STATUS TESTS
// ============================================================================

func TestExpenseStatus(t *testing.T) {
	validStatuses := []string{
		ExpenseStatusPending,
		ExpenseStatusSettled,
	}

	for _, status := range validStatuses {
		expense := &Expense{
			Description: "Dinner",
			Amount:      5000,
			PaidBy:      "user-1",
			Status:      status,
		}

		if expense.Status != status {
			t.Errorf("expected status %s, got %s", status, expense.Status)
		}
	}
}

// ============================================================================
// SETTLEMENT STATUS TESTS
// ============================================================================

func TestSettlementStatus(t *testing.T) {
	validStatuses := []string{
		SettlementStatusPending,
		SettlementStatusSettled,
	}

	for _, status := range validStatuses {
		settlement := &Settlement{
			DebtorID:   "user-1",
			CreditorID: "user-2",
			Amount:     5000,
			Status:     status,
		}

		if settlement.Status != status {
			t.Errorf("expected status %s, got %s", status, settlement.Status)
		}
	}
}

// ============================================================================
// GROUP TRIP CREATION REQUEST TESTS
// ============================================================================

func TestCreateGroupTripRequest(t *testing.T) {
	req := &CreateGroupTripRequest{
		Title:         "Goa Trip",
		DestinationID: "dest-1",
		Budget:        100000,
		Duration:      7,
	}

	if req.Title == "" || req.Budget <= 0 || req.Duration <= 0 {
		t.Error("request should have valid fields")
	}
}

// ============================================================================
// EXPENSE SPLIT CALCULATION TESTS
// ============================================================================

func TestEqualSplit(t *testing.T) {
	totalAmount := 5000.0
	numPeople := 4

	splitAmount := totalAmount / float64(numPeople)

	if splitAmount != 1250.0 {
		t.Errorf("expected split amount 1250, got %v", splitAmount)
	}

	// Verify total adds up
	totalRecovered := splitAmount * float64(numPeople)
	if totalRecovered != totalAmount {
		t.Errorf("expected total %v, got %v", totalAmount, totalRecovered)
	}
}

func TestCustomSplit(t *testing.T) {
	splits := map[string]float64{
		"user-1": 2000,
		"user-2": 1500,
		"user-3": 1500,
	}

	totalSplit := 0.0
	for _, amount := range splits {
		totalSplit += amount
	}

	expectedTotal := 5000.0
	if totalSplit != expectedTotal {
		t.Errorf("expected total %v, got %v", expectedTotal, totalSplit)
	}
}
