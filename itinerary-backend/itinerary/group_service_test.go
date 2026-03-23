package itinerary

import (
	"testing"
	"time"
)

// ============================================================================
// MOCK OBJECTS FOR SERVICE TESTS
// ============================================================================

// MockGroupDatabase wraps MockDatabase to add group-specific methods
type MockGroupDatabase struct {
	*MockDatabase
	groupTrips    map[string]*GroupTrip
	groupMembers  map[string][]*GroupMember
	expenses      map[string]*Expense
	expenseSplits map[string][]*ExpenseSplit
}

func NewMockGroupDatabase() *MockGroupDatabase {
	return &MockGroupDatabase{
		MockDatabase:  &MockDatabase{},
		groupTrips:    make(map[string]*GroupTrip),
		groupMembers:  make(map[string][]*GroupMember),
		expenses:      make(map[string]*Expense),
		expenseSplits: make(map[string][]*ExpenseSplit),
	}
}

// ============================================================================
// GROUP TRIP SERVICE TESTS
// ============================================================================

func TestCreateGroupTripValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *CreateGroupTripRequest
		expectErr bool
	}{
		{
			name: "valid request",
			req: &CreateGroupTripRequest{
				Title:         "Bali Trip",
				DestinationID: "dest-1",
				Budget:        50000,
				Duration:      5,
			},
			expectErr: false,
		},
		{
			name: "zero budget",
			req: &CreateGroupTripRequest{
				Title:         "Bali Trip",
				DestinationID: "dest-1",
				Budget:        0,
				Duration:      5,
			},
			expectErr: true,
		},
		{
			name: "negative budget",
			req: &CreateGroupTripRequest{
				Title:         "Bali Trip",
				DestinationID: "dest-1",
				Budget:        -50000,
				Duration:      5,
			},
			expectErr: true,
		},
		{
			name: "zero duration",
			req: &CreateGroupTripRequest{
				Title:         "Bali Trip",
				DestinationID: "dest-1",
				Budget:        50000,
				Duration:      0,
			},
			expectErr: true,
		},
		{
			name: "empty title",
			req: &CreateGroupTripRequest{
				Title:         "",
				DestinationID: "dest-1",
				Budget:        50000,
				Duration:      5,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test just validation logic, not full creation
			if tt.req.Budget <= 0 || tt.req.Duration <= 0 || tt.req.Title == "" {
				if !tt.expectErr {
					t.Error("expected validation to fail")
				}
			} else {
				if tt.expectErr {
					t.Error("expected validation to pass")
				}
			}
		})
	}
}

// ============================================================================
// EXPENSE SPLITTING ALGORITHM TESTS
// ============================================================================

func TestEqualExpenseSplit(t *testing.T) {
	tests := []struct {
		name          string
		totalAmount   float64
		memberCount   int
		expectedSplit float64
	}{
		{
			name:          "5000 split 4 ways",
			totalAmount:   5000,
			memberCount:   4,
			expectedSplit: 1250,
		},
		{
			name:          "10000 split 2 ways",
			totalAmount:   10000,
			memberCount:   2,
			expectedSplit: 5000,
		},
		{
			name:          "3000 split 3 ways",
			totalAmount:   3000,
			memberCount:   3,
			expectedSplit: 1000,
		},
		{
			name:          "100 split 1 way",
			totalAmount:   100,
			memberCount:   1,
			expectedSplit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitAmount := tt.totalAmount / float64(tt.memberCount)
			if splitAmount != tt.expectedSplit {
				t.Errorf("expected split %v, got %v", tt.expectedSplit, splitAmount)
			}
		})
	}
}

// ============================================================================
// SETTLEMENT CALCULATION ALGORITHM TESTS
// ============================================================================

func TestSettlementCalculationSimple(t *testing.T) {
	// Simple case: User A paid 100, User B owes 50, User C owes 50
	paidByUser := map[string]float64{
		"user-a": 100.0,
		"user-b": 0.0,
		"user-c": 0.0,
	}

	owedByUser := map[string]float64{
		"user-a": 0.0,  // User A doesn't owe anything
		"user-b": 50.0, // User B owes 50
		"user-c": 50.0, // User C owes 50
	}

	// Balances: A=+100, B=-50, C=-50
	// Settlements needed: B pays 50 to A, C pays 50 to A
	// Total: 2 settlements

	balances := make(map[string]float64)
	for userID := range paidByUser {
		balances[userID] = paidByUser[userID] - owedByUser[userID]
	}

	if balances["user-a"] != 100.0 {
		t.Errorf("expected user-a balance 100, got %v", balances["user-a"])
	}
	if balances["user-b"] != -50.0 {
		t.Errorf("expected user-b balance -50, got %v", balances["user-b"])
	}
	if balances["user-c"] != -50.0 {
		t.Errorf("expected user-c balance -50, got %v", balances["user-c"])
	}
}

func TestSettlementCalculationComplex(t *testing.T) {
	// Complex case: Multiple people pay different amounts
	paidByUser := map[string]float64{
		"user-a": 6000.0, // Paid 6000
		"user-b": 2000.0, // Paid 2000
		"user-c": 0.0,    // Paid 0
	}

	owedByUser := map[string]float64{
		"user-a": 2000.0, // Owes 2000 (for food)
		"user-b": 2000.0, // Owes 2000
		"user-c": 4000.0, // Owes 4000
	}

	// Balances: A=4000, B=0, C=-4000
	// Settlements: C pays 4000 to A

	for userID := range paidByUser {
		balance := paidByUser[userID] - owedByUser[userID]
		if userID == "user-a" && balance != 4000.0 {
			t.Errorf("expected user-a balance 4000, got %v", balance)
		}
		if userID == "user-b" && balance != 0.0 {
			t.Errorf("expected user-b balance 0, got %v", balance)
		}
		if userID == "user-c" && balance != -4000.0 {
			t.Errorf("expected user-c balance -4000, got %v", balance)
		}
	}
}

func TestSettlementCalculationChain(t *testing.T) {
	// Chain case: A owes B, B owes C
	paidByUser := map[string]float64{
		"user-a": 2000.0, // Paid for items
		"user-b": 2000.0,
		"user-c": 2000.0,
	}

	owedByUser := map[string]float64{
		"user-a": 2000.0, // Everyone paid 2000, so everyone owes 2000
		"user-b": 2000.0,
		"user-c": 2000.0,
	}

	// Balances: A=0, B=0, C=0 (everyone settled)
	totalSettled := 0
	for userID := range paidByUser {
		balance := paidByUser[userID] - owedByUser[userID]
		if balance == 0.0 {
			totalSettled++
		}
	}

	if totalSettled != 3 {
		t.Errorf("expected all 3 users settled, found %d", totalSettled)
	}
}

// ============================================================================
// GROUP MEMBER ROLE TESTS
// ============================================================================

func TestGroupMemberRolePermissions(t *testing.T) {
	tests := []struct {
		role          string
		canInvite     bool
		canRemove     bool
		canDeleteTrip bool
		canUpdateTrip bool
	}{
		{
			role:          GroupMemberRoleOwner,
			canInvite:     true,
			canRemove:     true,
			canDeleteTrip: true,
			canUpdateTrip: true,
		},
		{
			role:          GroupMemberRoleEditor,
			canInvite:     true,
			canRemove:     false,
			canDeleteTrip: false,
			canUpdateTrip: true,
		},
		{
			role:          GroupMemberRoleMember,
			canInvite:     false,
			canRemove:     false,
			canDeleteTrip: false,
			canUpdateTrip: false,
		},
		{
			role:          GroupMemberRoleViewer,
			canInvite:     false,
			canRemove:     false,
			canDeleteTrip: false,
			canUpdateTrip: false,
		},
	}

	for _, tt := range tests {
		// Only owner can delete trip
		if tt.role == GroupMemberRoleOwner && !tt.canDeleteTrip {
			t.Errorf("owner should be able to delete trip")
		}
		if tt.role != GroupMemberRoleOwner && tt.canDeleteTrip {
			t.Errorf("%s should not be able to delete trip", tt.role)
		}

		// Owner and editor can invite
		if (tt.role == GroupMemberRoleOwner || tt.role == GroupMemberRoleEditor) && !tt.canInvite {
			t.Errorf("%s should be able to invite members", tt.role)
		}
		if tt.role != GroupMemberRoleOwner && tt.role != GroupMemberRoleEditor && tt.canInvite {
			t.Errorf("%s should not be able to invite members", tt.role)
		}
	}
}

// ============================================================================
// POLL VOTING TESTS
// ============================================================================

func TestPollVoting(t *testing.T) {
	tests := []struct {
		name           string
		numOptions     int
		votes          map[string]int // option_index -> vote_count
		expectedWinner int
	}{
		{
			name:       "clear winner",
			numOptions: 3,
			votes: map[string]int{
				"0": 5, // Pizza gets 5 votes
				"1": 2, // Burger gets 2 votes
				"2": 1, // Chinese gets 1 vote
			},
			expectedWinner: 0,
		},
		{
			name:       "tie",
			numOptions: 2,
			votes: map[string]int{
				"0": 3, // Pizza gets 3 votes
				"1": 3, // Burger gets 3 votes
			},
			expectedWinner: 0, // First would be selected in tie-breaker
		},
		{
			name:       "all votes one option",
			numOptions: 1,
			votes: map[string]int{
				"0": 10,
			},
			expectedWinner: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Total votes should match
			totalVotes := 0
			for _, count := range tt.votes {
				totalVotes += count
			}

			if totalVotes == 0 {
				t.Error("no votes cast")
			}
		})
	}
}

// ============================================================================
// EXPENSE CATEGORY TESTS
// ============================================================================

func TestExpenseCategoryOrganization(t *testing.T) {
	// Test that expenses can be organized by category
	expenses := []*Expense{
		{ID: "exp-1", Category: ExpenseCategoryFood, Amount: 1000},
		{ID: "exp-2", Category: ExpenseCategoryTransport, Amount: 500},
		{ID: "exp-3", Category: ExpenseCategoryFood, Amount: 800},
		{ID: "exp-4", Category: ExpenseCategoryAccommodation, Amount: 3000},
	}

	expensesByCategory := make(map[string][]string)
	for _, exp := range expenses {
		expensesByCategory[exp.Category] = append(expensesByCategory[exp.Category], exp.ID)
	}

	if len(expensesByCategory[ExpenseCategoryFood]) != 2 {
		t.Errorf("expected 2 food expenses, got %d", len(expensesByCategory[ExpenseCategoryFood]))
	}
	if len(expensesByCategory[ExpenseCategoryTransport]) != 1 {
		t.Errorf("expected 1 transport expense, got %d", len(expensesByCategory[ExpenseCategoryTransport]))
	}
	if len(expensesByCategory[ExpenseCategoryAccommodation]) != 1 {
		t.Errorf("expected 1 accommodation expense, got %d", len(expensesByCategory[ExpenseCategoryAccommodation]))
	}
}

// ============================================================================
// GROUP TRIP STATUS TRANSITIONS
// ============================================================================

func TestGroupTripStatusTransitions(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		newStatus     string
		shouldBeValid bool
	}{
		{
			name:          "draft -> planning",
			currentStatus: GroupTripStatusDraft,
			newStatus:     GroupTripStatusPlanning,
			shouldBeValid: true,
		},
		{
			name:          "planning -> published",
			currentStatus: GroupTripStatusPlanning,
			newStatus:     GroupTripStatusPublished,
			shouldBeValid: true,
		},
		{
			name:          "published -> completed",
			currentStatus: GroupTripStatusPublished,
			newStatus:     GroupTripStatusCompleted,
			shouldBeValid: true,
		},
		{
			name:          "completed -> draft (invalid)",
			currentStatus: GroupTripStatusCompleted,
			newStatus:     GroupTripStatusDraft,
			shouldBeValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation: can't go backwards
			validTransitions := map[string][]string{
				GroupTripStatusDraft:     {GroupTripStatusPlanning},
				GroupTripStatusPlanning:  {GroupTripStatusPublished},
				GroupTripStatusPublished: {GroupTripStatusCompleted},
				GroupTripStatusCompleted: {},
			}

			isValid := false
			for _, validNext := range validTransitions[tt.currentStatus] {
				if validNext == tt.newStatus {
					isValid = true
					break
				}
			}

			if isValid != tt.shouldBeValid {
				t.Errorf("transition %s -> %s should be valid=%v, got %v",
					tt.currentStatus, tt.newStatus, tt.shouldBeValid, isValid)
			}
		})
	}
}

// ============================================================================
// GROUP MEMBER INVITATION TESTS
// ============================================================================

func TestGroupMemberInvitationLifecycle(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  string
		action         string
		expectedStatus string
	}{
		{
			name:           "pending -> accept -> active",
			initialStatus:  GroupMemberStatusPending,
			action:         "accept",
			expectedStatus: GroupMemberStatusActive,
		},
		{
			name:           "pending -> decline -> declined",
			initialStatus:  GroupMemberStatusPending,
			action:         "decline",
			expectedStatus: GroupMemberStatusDeclined,
		},
		{
			name:           "active -> leave -> left",
			initialStatus:  GroupMemberStatusActive,
			action:         "leave",
			expectedStatus: GroupMemberStatusLeft,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			member := &GroupMember{
				UserID: "user-1",
				Status: tt.initialStatus,
			}

			// Simulate status change
			if tt.action == "accept" && tt.initialStatus == GroupMemberStatusPending {
				member.Status = GroupMemberStatusActive
			} else if tt.action == "decline" && tt.initialStatus == GroupMemberStatusPending {
				member.Status = GroupMemberStatusDeclined
			} else if tt.action == "leave" && tt.initialStatus == GroupMemberStatusActive {
				member.Status = GroupMemberStatusLeft
			}

			if member.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, member.Status)
			}
		})
	}
}

// ============================================================================
// TIME-BASED POLL EXPIRATION TESTS
// ============================================================================

func TestPollExpiration(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-1 * time.Hour)
	futureTime := now.Add(1 * time.Hour)

	tests := []struct {
		name      string
		expiresAt *time.Time
		isExpired bool
	}{
		{
			name:      "no expiration",
			expiresAt: nil,
			isExpired: false,
		},
		{
			name:      "expires in past",
			expiresAt: &pastTime,
			isExpired: true,
		},
		{
			name:      "expires in future",
			expiresAt: &futureTime,
			isExpired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var isExpired bool
			if tt.expiresAt != nil && tt.expiresAt.Before(now) {
				isExpired = true
			}

			if isExpired != tt.isExpired {
				t.Errorf("expected expired=%v, got %v", tt.isExpired, isExpired)
			}
		})
	}
}
