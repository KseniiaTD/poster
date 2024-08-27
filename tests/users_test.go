package tests

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/KseniiaTD/poster/graph"
	"github.com/KseniiaTD/poster/graph/model"
	"github.com/KseniiaTD/poster/internal/database"
	"github.com/KseniiaTD/poster/internal/database/common"
	"github.com/stretchr/testify/require"
)

func TestCreateUserOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	input := model.NewUserInput{
		Login:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:    "test name",
		Surname: "test surname",
		Phone:   "8800" + randomPhoneNumber(),
		Email:   fmt.Sprintf("test%d@email.com", time.Now().UnixNano()),
	}

	for i, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		userId, err := resolver.Mutation().CreateUser(context.Background(), input)
		require.NoError(t, err, i)
		require.NotNil(t, userId)

		id, err := strconv.Atoi(*userId)
		require.NoError(t, err)
		require.Greater(t, id, 0)
	}
}

func TestCreateUserError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	inputs := []model.NewUserInput{

		{
			Login:   fmt.Sprintf("test login %d", time.Now().Unix()),
			Name:    "test name",
			Surname: "test surname",
			Phone:   "88003333333",
			Email:   "incorrect",
		},
		{
			Login:   fmt.Sprintf("test login %d", time.Now().Unix()),
			Name:    "test name",
			Surname: "test surname",
			Phone:   "incorrect",
			Email:   "test@email.com",
		},
	}

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		for _, input := range inputs {
			userId, err := resolver.Mutation().CreateUser(context.Background(), input)
			require.Error(t, err)
			require.Nil(t, userId)
		}
	}
}
