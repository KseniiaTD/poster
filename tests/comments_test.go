package tests

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/KseniiaTD/poster/graph"
	"github.com/KseniiaTD/poster/graph/model"
	"github.com/KseniiaTD/poster/internal/database"
	"github.com/KseniiaTD/poster/internal/database/common"
	"github.com/stretchr/testify/require"
)

func TestCreateCommentOk(t *testing.T) {
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

		idStr, err := resolver.Mutation().CreateComment(context.Background(), model.NewCommentInput{
			PostID:   postId,
			AuthorID: userId,
			Body:     "Test Body",
		})
		require.NoError(t, err)
		require.NotNil(t, idStr)

		id, err := strconv.Atoi(*idStr)
		require.NoError(t, err)
		require.Greater(t, id, 0)

		childIdStr, err := resolver.Mutation().CreateComment(context.Background(), model.NewCommentInput{
			PostID:   postId,
			AuthorID: userId,
			Body:     "Test Body",
			ParentID: idStr,
		})

		require.NoError(t, err)
		require.NotNil(t, childIdStr)

		childId, err := strconv.Atoi(*childIdStr)
		require.NoError(t, err)
		require.Greater(t, childId, 0)
	}
}

func TestCreateCommentError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	inputs := []struct {
		Body       string
		WithUserId bool
		WithPostId bool
	}{
		{
			Body:       "",
			WithUserId: true,
			WithPostId: true,
		},
		{
			Body:       strings.Repeat("s", 2001),
			WithUserId: true,
			WithPostId: true,
		},
		{
			Body:       "test body",
			WithUserId: true,
			WithPostId: false,
		},
		{
			Body:       "test body",
			WithUserId: false,
			WithPostId: true,
		},
		{
			Body:       "test body",
			WithUserId: false,
			WithPostId: false,
		},
	}

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, userId, err := createTestPost(resolver)
		require.NoError(t, err)

		for _, input := range inputs {
			if !input.WithUserId {
				userId = "0"
			}
			if !input.WithPostId {
				postId = "0"
			}

			idStr, err := resolver.Mutation().CreateComment(context.Background(), model.NewCommentInput{
				Body:     input.Body,
				AuthorID: userId,
				PostID:   postId,
			})

			require.Error(t, err)
			require.Nil(t, idStr)
		}
	}
}

func TestUpdateCommentOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		createdId, _, _, err := createTestComment(resolver)
		require.NoError(t, err)

		idStr, err := resolver.Mutation().UpdateComment(context.Background(), model.UpdCommentInput{
			ID:   createdId,
			Body: "test body",
		})
		require.NoError(t, err)
		require.NotNil(t, idStr)

		id, err := strconv.Atoi(*idStr)
		require.NoError(t, err)
		require.Greater(t, id, 0)

	}
}

func TestUpdateCommentError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for j, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		id, _, _, err := createTestComment(resolver)
		require.NoError(t, err)

		inputs := []model.UpdCommentInput{
			{
				ID:   "0",
				Body: "test body",
			},
			{
				ID:   id,
				Body: "",
			},
			{
				ID:   id,
				Body: strings.Repeat("s", 2001),
			},
		}

		for _, input := range inputs {

			idStr, err := resolver.Mutation().UpdateComment(context.Background(), input)

			require.Error(t, err, j)
			require.Nil(t, idStr)
		}
	}
}

func TestDeleteCommentOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		commentId, _, _, err := createTestComment(resolver)
		require.NoError(t, err)

		idStr, err := resolver.Mutation().DeleteComment(context.Background(), commentId)
		require.NoError(t, err)
		require.NotNil(t, idStr)

		id, err := strconv.Atoi(*idStr)
		require.NoError(t, err)
		require.Greater(t, id, 0)
	}
}

func TestDeleteCommentError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	inputs := []struct {
		ID string
	}{
		{
			ID: "0",
		},
		{
			ID: "",
		},
	}

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		for _, input := range inputs {
			idStr, err := resolver.Mutation().DeleteComment(context.Background(), input.ID)
			require.Error(t, err, input)
			require.Nil(t, idStr)
		}
	}
}

func TestGetCommentOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		commentId, postId, _, err := createTestComment(resolver)
		require.NoError(t, err)

		postIdInt, err := strconv.Atoi(postId)
		require.NoError(t, err)

		page := 0
		perPage := 1000000000

		comments, err := resolver.Query().Comments(context.Background(), postIdInt, nil, &page, &perPage)
		require.NoError(t, err)
		require.NotEmpty(t, comments)

		found := false
		for _, comment := range comments {
			if comment.ID == commentId {
				found = true
				break
			}
		}

		require.True(t, found)
	}
}
