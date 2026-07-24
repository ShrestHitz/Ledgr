package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shresth/ledgr/util"
	"github.com/stretchr/testify/require"
)

func createRandomRecurringPayment(t *testing.T) RecurringPayment {
	user := createRandomUser(t)
	category, err := testStore.CreateCategory(context.Background(), CreateCategoryParams{
		Owner: pgtype.Text{String: user.Username, Valid: true},
		Name:  "Recurring_" + util.RandomString(4),
		Type:  "expense",
	})
	require.NoError(t, err)

	nextDue := util.RandomFutureDate()

	arg := CreateRecurringPaymentParams{
		Owner:       user.Username,
		Title:       "Payment_" + util.RandomString(6),
		Amount:      util.RandomMoney(),
		CategoryID:  category.ID,
		Frequency:   util.RandomFrequency(),
		NextDueDate: nextDue,
	}

	payment, err := testStore.CreateRecurringPayment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, payment)

	require.Equal(t, arg.Owner, payment.Owner)
	require.Equal(t, arg.Title, payment.Title)
	require.Equal(t, arg.Amount, payment.Amount)
	require.Equal(t, arg.CategoryID, payment.CategoryID)
	require.Equal(t, arg.Frequency, payment.Frequency)
	require.NotZero(t, payment.ID)

	return payment
}

func TestCreateRecurringPayment(t *testing.T) {
	createRandomRecurringPayment(t)
}

func TestGetRecurringPayment(t *testing.T) {
	payment1 := createRandomRecurringPayment(t)

	payment2, err := testStore.GetRecurringPayment(context.Background(), payment1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, payment2)

	require.Equal(t, payment1.ID, payment2.ID)
	require.Equal(t, payment1.Title, payment2.Title)
	require.Equal(t, payment1.Amount, payment2.Amount)
	require.Equal(t, payment1.Frequency, payment2.Frequency)
}

func TestListRecurringPayments(t *testing.T) {
	payment1 := createRandomRecurringPayment(t)

	payments, err := testStore.ListRecurringPayments(context.Background(), payment1.Owner)
	require.NoError(t, err)
	require.NotEmpty(t, payments)

	for _, p := range payments {
		require.Equal(t, payment1.Owner, p.Owner)
	}
}

func TestListDueRecurringPayments(t *testing.T) {
	user := createRandomUser(t)
	category, err := testStore.CreateCategory(context.Background(), CreateCategoryParams{
		Owner: pgtype.Text{String: user.Username, Valid: true},
		Name:  "Due_" + util.RandomString(4),
		Type:  "expense",
	})
	require.NoError(t, err)

	// create one overdue payment (yesterday)
	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	_, err = testStore.CreateRecurringPayment(context.Background(), CreateRecurringPaymentParams{
		Owner:       user.Username,
		Title:       "Overdue_" + util.RandomString(4),
		Amount:      util.RandomMoney(),
		CategoryID:  category.ID,
		Frequency:   "monthly",
		NextDueDate: yesterday,
	})
	require.NoError(t, err)

	// create one future payment (next month)
	nextMonth := time.Now().AddDate(0, 1, 0).Truncate(24 * time.Hour)
	_, err = testStore.CreateRecurringPayment(context.Background(), CreateRecurringPaymentParams{
		Owner:       user.Username,
		Title:       "Future_" + util.RandomString(4),
		Amount:      util.RandomMoney(),
		CategoryID:  category.ID,
		Frequency:   "monthly",
		NextDueDate: nextMonth,
	})
	require.NoError(t, err)

	// only overdue payment should appear
	duePayments, err := testStore.ListDueRecurringPayments(context.Background(), ListDueRecurringPaymentsParams{
		Owner:       user.Username,
		NextDueDate: time.Now().Truncate(24 * time.Hour),
	})
	require.NoError(t, err)
	require.NotEmpty(t, duePayments)

	for _, p := range duePayments {
		require.Equal(t, user.Username, p.Owner)
		require.True(t, !p.NextDueDate.After(time.Now()))
	}
}

func TestUpdateRecurringPaymentDueDate(t *testing.T) {
	payment1 := createRandomRecurringPayment(t)
	newDueDate := util.RandomFutureDate()

	payment2, err := testStore.UpdateRecurringPaymentDueDate(context.Background(), UpdateRecurringPaymentDueDateParams{
		ID:          payment1.ID,
		NextDueDate: newDueDate,
	})
	require.NoError(t, err)
	require.Equal(t, payment1.ID, payment2.ID)
	require.WithinDuration(t, newDueDate, payment2.NextDueDate, 24*time.Hour)
}

func TestDeleteRecurringPayment(t *testing.T) {
	payment := createRandomRecurringPayment(t)

	err := testStore.DeleteRecurringPayment(context.Background(), payment.ID)
	require.NoError(t, err)

	_, err = testStore.GetRecurringPayment(context.Background(), payment.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
