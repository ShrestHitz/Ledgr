package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shresth/ledgr/util"
	"github.com/stretchr/testify/require"
)

// createRandomUserCategory creates a category owned by a specific user.
func createRandomUserCategory(t *testing.T, owner string) Category {
	arg := CreateCategoryParams{
		Owner: pgtype.Text{String: owner, Valid: true},
		Name:  util.RandomCategoryName() + "_" + util.RandomString(4),
		Type:  util.RandomCategoryType(),
	}

	category, err := testStore.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)

	require.True(t, category.Owner.Valid)
	require.Equal(t, owner, category.Owner.String)
	require.Equal(t, arg.Name, category.Name)
	require.Equal(t, arg.Type, category.Type)
	require.NotZero(t, category.ID)

	return category
}

// createGlobalCategory creates a system-level category with no owner.
func createGlobalCategory(t *testing.T) Category {
	arg := CreateCategoryParams{
		Owner: pgtype.Text{Valid: false}, // NULL owner = global
		Name:  util.RandomCategoryName() + "_global_" + util.RandomString(4),
		Type:  util.RandomCategoryType(),
	}

	category, err := testStore.CreateCategory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, category)
	require.False(t, category.Owner.Valid)

	return category
}

func TestCreateCategory(t *testing.T) {
	user := createRandomUser(t)
	createRandomUserCategory(t, user.Username)
}

func TestCreateGlobalCategory(t *testing.T) {
	createGlobalCategory(t)
}

func TestGetCategory(t *testing.T) {
	user := createRandomUser(t)
	cat1 := createRandomUserCategory(t, user.Username)

	cat2, err := testStore.GetCategory(context.Background(), cat1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, cat2)

	require.Equal(t, cat1.ID, cat2.ID)
	require.Equal(t, cat1.Name, cat2.Name)
	require.Equal(t, cat1.Type, cat2.Type)
	require.Equal(t, cat1.Owner, cat2.Owner)
}

func TestListCategories(t *testing.T) {
	user := createRandomUser(t)

	// create 3 user-specific and 2 global categories
	for i := 0; i < 3; i++ {
		createRandomUserCategory(t, user.Username)
	}
	for i := 0; i < 2; i++ {
		createGlobalCategory(t)
	}

	// listing for this user should return both user-owned and global categories
	categories, err := testStore.ListCategories(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, categories)

	for _, cat := range categories {
		require.True(t, !cat.Owner.Valid || cat.Owner.String == user.Username)
	}
}

func TestListCategoriesByType(t *testing.T) {
	user := createRandomUser(t)

	// create known income and expense categories
	_, err := testStore.CreateCategory(context.Background(), CreateCategoryParams{
		Owner: pgtype.Text{String: user.Username, Valid: true},
		Name:  "Salary_" + util.RandomString(4),
		Type:  "income",
	})
	require.NoError(t, err)

	_, err = testStore.CreateCategory(context.Background(), CreateCategoryParams{
		Owner: pgtype.Text{String: user.Username, Valid: true},
		Name:  "Groceries_" + util.RandomString(4),
		Type:  "expense",
	})
	require.NoError(t, err)

	expenseCategories, err := testStore.ListCategoriesByType(context.Background(), ListCategoriesByTypeParams{
		Owner: user.Username,
		Type:  "expense",
	})
	require.NoError(t, err)

	for _, cat := range expenseCategories {
		require.Equal(t, "expense", cat.Type)
	}
}

func TestDeleteCategory(t *testing.T) {
	user := createRandomUser(t)
	cat := createRandomUserCategory(t, user.Username)

	err := testStore.DeleteCategory(context.Background(), cat.ID)
	require.NoError(t, err)

	_, err = testStore.GetCategory(context.Background(), cat.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}
