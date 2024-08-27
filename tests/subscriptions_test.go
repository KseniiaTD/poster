package tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/KseniiaTD/poster/graph"
	"github.com/KseniiaTD/poster/graph/model"
	"github.com/KseniiaTD/poster/internal/database"
	"github.com/KseniiaTD/poster/internal/database/common"
	"github.com/stretchr/testify/require"
)

func TestCreateSubscriptionOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		subId, err := resolver.Mutation().CreateSubscription(context.Background(), model.SubscrInput{
			PostID: postId,
			UserID: userId,
		})
		require.NoError(t, err)

		id, err := strconv.Atoi(*subId)
		require.NoError(t, err)
		require.Greater(t, id, 0)
	}
}

func TestCreateSubscriptionError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		inputs := []model.SubscrInput{
			{
				PostID: "0",
				UserID: userId,
			},
			{
				PostID: postId,
				UserID: "0",
			},
		}

		for _, input := range inputs {
			_, err := resolver.Mutation().CreateSubscription(context.Background(), input)
			require.Error(t, err)
		}
	}
}

func TestDeleteSubscriptionOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		subId, err := resolver.Mutation().CreateSubscription(context.Background(), model.SubscrInput{
			PostID: postId,
			UserID: userId,
		})
		require.NoError(t, err)

		id, err := strconv.Atoi(*subId)
		require.NoError(t, err)
		require.Greater(t, id, 0)

		deletedId, err := resolver.Mutation().DeleteSubscription(context.Background(), model.SubscrInput{
			PostID: postId,
			UserID: userId,
		})

		require.NoError(t, err)

		id, err = strconv.Atoi(*deletedId)
		require.NoError(t, err)
		require.Greater(t, id, 0)
	}
}

func TestDeleteSubscriptionError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		subId, err := resolver.Mutation().CreateSubscription(context.Background(), model.SubscrInput{
			PostID: postId,
			UserID: userId,
		})
		require.NoError(t, err)

		id, err := strconv.Atoi(*subId)
		require.NoError(t, err)
		require.Greater(t, id, 0)

		inputs := []model.SubscrInput{
			{
				PostID: "0",
				UserID: userId,
			},
			{
				PostID: postId,
				UserID: "0",
			},
		}

		for _, input := range inputs {
			_, err := resolver.Mutation().DeleteSubscription(context.Background(), input)
			require.Error(t, err)
		}
	}
}

func TestCheckCommentsOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		_, err = resolver.Mutation().CreateSubscription(context.Background(), model.SubscrInput{
			PostID: postId,
			UserID: userId,
		})
		require.NoError(t, err)

		ch, err := resolver.Subscription().CheckComments(context.Background(), model.Subscr{
			PostID: postId,
			UserID: userId,
		})
		require.NoError(t, err)
		require.NotNil(t, ch)

	}
}

func TestCheckCommentsError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		ch, err := resolver.Subscription().CheckComments(context.Background(), model.Subscr{
			PostID: postId,
			UserID: userId,
		})
		require.Error(t, err)
		require.Nil(t, ch)
	}
}
