package tests

import (
	"context"
	"testing"

	"github.com/KseniiaTD/poster/graph"
	"github.com/KseniiaTD/poster/graph/model"
	"github.com/KseniiaTD/poster/internal/database"
	"github.com/KseniiaTD/poster/internal/database/common"
	"github.com/stretchr/testify/require"
)

func TestUpdPostLikesOk(t *testing.T) {
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

		_, err = resolver.Mutation().AddPostLike(context.Background(), model.NewPostLikeInput{
			Like:     true,
			PostID:   postId,
			AuthorID: userId,
		})

		require.NoError(t, err)

		_, err = resolver.Mutation().AddPostLike(context.Background(), model.NewPostLikeInput{
			Like:     true,
			PostID:   postId,
			AuthorID: userId,
		})

		require.NoError(t, err)

		_, err = resolver.Mutation().AddPostLike(context.Background(), model.NewPostLikeInput{
			Like:     false,
			PostID:   postId,
			AuthorID: userId,
		})

		require.NoError(t, err)
	}
}

func TestUpdPostLikesError(t *testing.T) {
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

		inputs := []model.NewPostLikeInput{
			{
				PostID:   "0",
				AuthorID: userId,
			},
			{
				PostID:   postId,
				AuthorID: "0",
			},
			{
				PostID:   "",
				AuthorID: "",
			},
		}

		for _, input := range inputs {
			_, err = resolver.Mutation().AddPostLike(context.Background(), input)

			require.Error(t, err)
		}
	}
}

func TestUpdCommentLikesOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		commentId, _, userId, err := createTestComment(resolver)
		require.NoError(t, err)

		_, err = resolver.Mutation().AddCommentLike(context.Background(), model.NewCommentLikeInput{
			Like:      true,
			AuthorID:  userId,
			CommentID: commentId,
		})

		require.NoError(t, err)

		_, err = resolver.Mutation().AddCommentLike(context.Background(), model.NewCommentLikeInput{
			Like:      true,
			AuthorID:  userId,
			CommentID: commentId,
		})

		require.NoError(t, err)

		_, err = resolver.Mutation().AddCommentLike(context.Background(), model.NewCommentLikeInput{
			Like:      false,
			AuthorID:  userId,
			CommentID: commentId,
		})

		require.NoError(t, err)
	}
}

func TestUpdCommentLikesError(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		commentId, _, userId, err := createTestComment(resolver)
		require.NoError(t, err)

		inputs := []model.NewCommentLikeInput{
			{
				CommentID: "0",
				AuthorID:  userId,
			},
			{
				CommentID: commentId,
				AuthorID:  "0",
			},
			{
				CommentID: "",
				AuthorID:  "",
			},
		}

		for _, input := range inputs {
			_, err = resolver.Mutation().AddCommentLike(context.Background(), input)

			require.Error(t, err)
		}
	}
}

func TestGetPostLikesOk(t *testing.T) {
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

		_, err = resolver.Mutation().AddPostLike(context.Background(), model.NewPostLikeInput{
			Like:     true,
			AuthorID: userId,
			PostID:   postId,
		})
		require.NoError(t, err)

		rating, err := resolver.Query().GetPostLikes(context.Background(), postId)
		require.NoError(t, err)
		require.Equal(t, 1, rating.Likes)
		require.Equal(t, 0, rating.Dislikes)
	}
}

func TestGetCommentLikesOk(t *testing.T) {
	inMemoryDB, err := database.New(true, initTestConfig())
	require.NoError(t, err)

	postgresDB, err := database.New(false, initTestConfig())
	require.NoError(t, err)

	for _, db := range []common.Database{inMemoryDB, postgresDB} {
		resolver := graph.Resolver{
			DB: db,
		}

		commentId, _, userId, err := createTestComment(resolver)
		require.NoError(t, err)

		_, err = resolver.Mutation().AddCommentLike(context.Background(), model.NewCommentLikeInput{
			Like:      true,
			AuthorID:  userId,
			CommentID: commentId,
		})
		require.NoError(t, err)

		rating, err := resolver.Query().GetCommentLikes(context.Background(), commentId)
		require.NoError(t, err)
		require.Equal(t, 1, rating.Likes)
		require.Equal(t, 0, rating.Dislikes)
	}
}
