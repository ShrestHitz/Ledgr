package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shresth/ledgr/util"
	"github.com/stretchr/testify/require"
)

func createRandomSavingsGoal(t *testing.T) SavingsGoal {
	account := createRandomAccount(t)
	targetDate := util.RandomFutureDate()

	arg := CreateSavingsGoalParams{
		Owner:           account.Owner,
		Title:           "Goal_" + util.RandomString(6),
		TargetAmount:    util.RandomMoney() * 100, // larger target
		CurrentAmount:   0,
		TargetDate:      pgtype.Date{Time: targetDate, Valid: true},
		LinkedAccountID: account.ID,
	}

	goal, err := testStore.CreateSavingsGoal(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, goal)

	require.Equal(t, arg.Owner, goal.Owner)
	require.Equal(t, arg.Title, goal.Title)
	require.Equal(t, arg.TargetAmount, goal.TargetAmount)
	require.Equal(t, int64(0), goal.CurrentAmount)
	require.True(t, goal.TargetDate.Valid)
	require.Equal(t, arg.LinkedAccountID, goal.LinkedAccountID)
	require.NotZero(t, goal.ID)

	return goal
}

func TestCreateSavingsGoal(t *testing.T) {
	createRandomSavingsGoal(t)
}

func TestGetSavingsGoal(t *testing.T) {
	goal1 := createRandomSavingsGoal(t)

	goal2, err := testStore.GetSavingsGoal(context.Background(), goal1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, goal2)

	require.Equal(t, goal1.ID, goal2.ID)
	require.Equal(t, goal1.Title, goal2.Title)
	require.Equal(t, goal1.TargetAmount, goal2.TargetAmount)
	require.Equal(t, goal1.Owner, goal2.Owner)
}

func TestAddToSavingsGoal(t *testing.T) {
	goal1 := createRandomSavingsGoal(t)
	deposit := util.RandomMoney()

	goal2, err := testStore.AddToSavingsGoal(context.Background(), AddToSavingsGoalParams{
		ID:     goal1.ID,
		Amount: deposit,
	})
	require.NoError(t, err)
	require.Equal(t, goal1.ID, goal2.ID)
	require.Equal(t, goal1.CurrentAmount+deposit, goal2.CurrentAmount)
	require.Equal(t, goal1.TargetAmount, goal2.TargetAmount)
}

func TestListSavingsGoals(t *testing.T) {
	goal1 := createRandomSavingsGoal(t)

	goals, err := testStore.ListSavingsGoals(context.Background(), goal1.Owner)
	require.NoError(t, err)
	require.NotEmpty(t, goals)

	for _, g := range goals {
		require.Equal(t, goal1.Owner, g.Owner)
	}
}

func TestDeleteSavingsGoal(t *testing.T) {
	goal := createRandomSavingsGoal(t)

	err := testStore.DeleteSavingsGoal(context.Background(), goal.ID)
	require.NoError(t, err)

	_, err = testStore.GetSavingsGoal(context.Background(), goal.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
