package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shresth/ledgr/util"
	"github.com/stretchr/testify/require"
)

func createRandomBudget(t *testing.T) Budget {
	user := createRandomUser(t)
	category, err := testStore.CreateCategory(context.Background(), CreateCategoryParams{
		Owner: pgtype.Text{String: user.Username, Valid: true},
		Name:  "Budget_" + util.RandomString(4),
		Type:  "expense",
	})
	require.NoError(t, err)

	now := time.Now()
	arg := CreateBudgetParams{
		Owner:        user.Username,
		CategoryID:   category.ID,
		MonthlyLimit: util.RandomMoney(),
		Month:        int32(now.Month()),
		Year:         int32(now.Year()),
	}

	budget, err := testStore.CreateBudget(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, budget)

	require.Equal(t, arg.Owner, budget.Owner)
	require.Equal(t, arg.CategoryID, budget.CategoryID)
	require.Equal(t, arg.MonthlyLimit, budget.MonthlyLimit)
	require.Equal(t, arg.Month, budget.Month)
	require.Equal(t, arg.Year, budget.Year)
	require.NotZero(t, budget.ID)

	return budget
}

func TestCreateBudget(t *testing.T) {
	createRandomBudget(t)
}

func TestGetBudget(t *testing.T) {
	budget1 := createRandomBudget(t)

	budget2, err := testStore.GetBudget(context.Background(), budget1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, budget2)

	require.Equal(t, budget1.ID, budget2.ID)
	require.Equal(t, budget1.Owner, budget2.Owner)
	require.Equal(t, budget1.MonthlyLimit, budget2.MonthlyLimit)
	require.Equal(t, budget1.Month, budget2.Month)
	require.Equal(t, budget1.Year, budget2.Year)
}

func TestGetBudgetByCategoryAndMonth(t *testing.T) {
	budget1 := createRandomBudget(t)

	budget2, err := testStore.GetBudgetByCategoryAndMonth(context.Background(), GetBudgetByCategoryAndMonthParams{
		Owner:      budget1.Owner,
		CategoryID: budget1.CategoryID,
		Month:      budget1.Month,
		Year:       budget1.Year,
	})
	require.NoError(t, err)
	require.Equal(t, budget1.ID, budget2.ID)
	require.Equal(t, budget1.MonthlyLimit, budget2.MonthlyLimit)
}

func TestListBudgets(t *testing.T) {
	budget1 := createRandomBudget(t)

	budgets, err := testStore.ListBudgets(context.Background(), ListBudgetsParams{
		Owner: budget1.Owner,
		Month: budget1.Month,
		Year:  budget1.Year,
	})
	require.NoError(t, err)
	require.NotEmpty(t, budgets)

	for _, b := range budgets {
		require.Equal(t, budget1.Owner, b.Owner)
		require.Equal(t, budget1.Month, b.Month)
		require.Equal(t, budget1.Year, b.Year)
	}
}

func TestUpdateBudget(t *testing.T) {
	budget1 := createRandomBudget(t)
	newLimit := util.RandomMoney()

	budget2, err := testStore.UpdateBudget(context.Background(), UpdateBudgetParams{
		ID:           budget1.ID,
		MonthlyLimit: newLimit,
	})
	require.NoError(t, err)
	require.Equal(t, budget1.ID, budget2.ID)
	require.Equal(t, newLimit, budget2.MonthlyLimit)
}

func TestDeleteBudget(t *testing.T) {
	budget := createRandomBudget(t)
	err := testStore.DeleteBudget(context.Background(), budget.ID)
	require.NoError(t, err)

	_, err = testStore.GetBudget(context.Background(), budget.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
