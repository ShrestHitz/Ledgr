package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shresth/ledgr/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {
	amount := util.RandomInt(-50000, 50000)
	if amount == 0 {
		amount = 100
	}

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	}

	entry, err := testStore.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func createRandomEntryWithCategory(t *testing.T, account Account, category Category) Entry {
	amount := -util.RandomMoney()

	arg := CreateEntryParams{
		AccountID:  account.ID,
		CategoryID: pgtype.Int8{Int64: category.ID, Valid: true},
		Amount:     amount,
		Note:       pgtype.Text{String: "test expense", Valid: true},
	}

	entry, err := testStore.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.True(t, entry.CategoryID.Valid)
	require.Equal(t, category.ID, entry.CategoryID.Int64)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	entry1 := createRandomEntry(t, account)

	entry2, err := testStore.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}

func TestGetAccountSummary(t *testing.T) {
	account := createRandomAccount(t)

	// create 3 income entries and 2 expense entries
	for i := 0; i < 3; i++ {
		_, err := testStore.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: account.ID,
			Amount:    10000, // +100.00
		})
		require.NoError(t, err)
	}
	for i := 0; i < 2; i++ {
		_, err := testStore.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: account.ID,
			Amount:    -5000, // -50.00
		})
		require.NoError(t, err)
	}

	summary, err := testStore.GetAccountSummary(context.Background(), GetAccountSummaryParams{
		AccountID: account.ID,
		FromTime:  time.Now().AddDate(0, 0, -1),
		ToTime:    time.Now().AddDate(0, 0, 1),
	})
	require.NoError(t, err)
	require.Equal(t, int64(30000), summary.TotalIncome)
	require.Equal(t, int64(-10000), summary.TotalExpense)
	require.Equal(t, int64(20000), summary.NetBalance)
}
