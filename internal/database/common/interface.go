package common

import (
	"context"

	"github.com/KseniiaTD/poster/graph/model"
)

type Database interface {
	CreateUser(ctx context.Context, newUser model.NewUserInput) (int, error)

	CreatePost(ctx context.Context, newPost model.NewPostInput) (int, error)
	UpdatePost(ctx context.Context, post model.UpdPostInput) (int, error)
	DeletePost(ctx context.Context, id int) (int, error)

	CreateSubscription(ctx context.Context, subscr model.SubscrInput) (*string, error)
	DeleteSubscription(ctx context.Context, subscr model.SubscrInput) (*string, error)
	CheckSubscription(ctx context.Context, subscr model.Subscr) error
	GetCntNewCommentsForSubscriber(ctx context.Context, subscr model.Subscr) (int, error)

	CreateComment(ctx context.Context, newComment model.NewCommentInput) (int, error)
	UpdateComment(ctx context.Context, comment model.UpdCommentInput) (int, error)
	DeleteComment(ctx context.Context, id int) (int, error)

	GetPosts(ctx context.Context, userId int) ([]*model.Post, error)
	GetComments(ctx context.Context, postId int, commentId int, page int, perPage int) ([]*model.Comment, error)

	UpdPostLikes(ctx context.Context, authorId int, postId int, isLike bool) error
	UpdCommentLikes(ctx context.Context, authorId int, commentId int, isLike bool) error

	GetPostLikes(ctx context.Context, postId int) (int, int, error)
	GetCommentLikes(ctx context.Context, commentId int) (int, int, error)

	CloseDB()
}
