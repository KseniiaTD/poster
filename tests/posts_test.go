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

func TestCreatePostOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for i, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		userId, err := createTestUser(resolver)
		require.NoError(t, err, i)

		input := model.NewPostInput{
			Title:       "test title",
			AuthorID:    userId,
			Body:        "Test Body",
			IsCommented: true,
		}
		idStr, err := resolver.Mutation().CreatePost(context.Background(), input)
		require.NoError(t, err)
		require.NotNil(t, idStr)

		id, err := strconv.Atoi(*idStr)
		require.NoError(t, err)
		require.Greater(t, id, 0)
	}
}

func TestCreatePostError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	inputs := []struct {
		Title       string
		Body        string
		IsCommented bool
		WithAuthor  bool
	}{
		{
			Title:       "test title",
			Body:        "Test Body",
			IsCommented: true,
			WithAuthor:  false,
		},
		{
			Title:       "",
			Body:        "Test Body",
			IsCommented: true,
			WithAuthor:  true,
		},
	}

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		for _, input := range inputs {
			userId := "0"
			if input.WithAuthor {
				id, err := createTestUser(resolver)
				require.NoError(t, err)
				userId = id
			}

			idStr, err := resolver.Mutation().CreatePost(context.Background(), model.NewPostInput{
				Title:       input.Title,
				Body:        input.Body,
				AuthorID:    userId,
				IsCommented: input.IsCommented,
			})

			require.Error(t, err)
			require.Nil(t, idStr)
		}
	}
}

func TestUpdatePostOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, _, err := createTestPost(resolver)
		require.NoError(t, err)

		title := "test title"
		body := "test body"
		isCommented := true
		inputs := []model.UpdPostInput{
			{
				PostID: postId,
				Title:  &title,
			},
			{
				PostID: postId,
				Body:   &body,
			},
			{
				PostID:      postId,
				IsCommented: &isCommented,
			},
			{
				PostID:      postId,
				Title:       &title,
				Body:        &body,
				IsCommented: &isCommented,
			},
		}

		for _, input := range inputs {
			idStr, err := resolver.Mutation().UpdatePost(context.Background(), input)
			require.NoError(t, err)
			require.NotNil(t, idStr)

			id, err := strconv.Atoi(*idStr)
			require.NoError(t, err)
			require.Greater(t, id, 0)
		}
	}
}

func TestUpdatePostError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	inputs := []struct {
		Title       string
		Body        string
		IsCommented bool
		WithPost    bool
	}{
		{
			Title:       "test title",
			Body:        "Test Body",
			IsCommented: true,
			WithPost:    false,
		},
		{
			Title:       "",
			Body:        "Test Body",
			IsCommented: true,
			WithPost:    true,
		},
		{
			Title:       "title",
			Body:        "",
			IsCommented: true,
			WithPost:    true,
		},
	}

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		for _, input := range inputs {
			postId := "0"
			if input.WithPost {
				id, _, err := createTestPost(resolver)
				require.NoError(t, err)
				postId = id
			}

			idStr, err := resolver.Mutation().UpdatePost(context.Background(), model.UpdPostInput{
				Title:       &input.Title,
				Body:        &input.Body,
				PostID:      postId,
				IsCommented: &input.IsCommented,
			})

			require.Error(t, err)
			require.Nil(t, idStr)
		}
	}
}

func TestDeletePostOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		postId, _, err := createTestPost(resolver)
		require.NoError(t, err)

		idStr, err := resolver.Mutation().DeletePost(context.Background(), postId)
		require.NoError(t, err)
		require.NotNil(t, idStr)

		id, err := strconv.Atoi(*idStr)
		require.NoError(t, err)
		require.Greater(t, id, 0)
	}
}

func TestDeletePostError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	inputs := []struct {
		PostID string
	}{
		{
			PostID: "0",
		},
		{
			PostID: "",
		},
	}

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		for _, input := range inputs {
			idStr, err := resolver.Mutation().DeletePost(context.Background(), input.PostID)
			require.Error(t, err, input)
			require.Nil(t, idStr)
		}
	}
}

func TestGetPostOk(t *testing.T) {
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

		userIdInt, err := strconv.Atoi(userId)
		require.NoError(t, err)

		posts, err := resolver.Query().Posts(context.Background(), userIdInt)
		require.NoError(t, err)
		require.NotEmpty(t, posts)

		found := false
		for _, post := range posts {
			if post.ID == postId {
				found = true
				break
			}
		}

		require.True(t, found)
	}
}
