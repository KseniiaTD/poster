package database

import (
	"context"
	"database/sql"
	"fmt"
	"poster/graph/model"

	_ "github.com/lib/pq"
)

type Database interface {
	CreatePost(ctx context.Context, newPost model.NewPostInput) (int, error)
	UpdatePost(ctx context.Context, post model.UpdPostInput) (int, error)
	DeletePost(ctx context.Context, id int) (int, error)

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

func New(isInMemory bool) (Database, error) {
	if !isInMemory {
		dsn, err := getDSN()

		if err != nil {
			return nil, err
		}
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			panic(err)
		}

		return &postgresDB{
			db: db,
		}, nil
	}

	return &inMemoryDB{users: make(map[int]user, 0),
		posts:        make(map[int]post, 0),
		comments:     make(map[int]comment, 0),
		postLikes:    make(map[postLike]bool, 0),
		commentLikes: make(map[commentLike]bool, 0),
	}, nil
}

const (
	host     = "localhost"
	port     = 5432
	userDB   = "pguser"
	password = "pgpwd"
	dbname   = "db_poster"
)

func getDSN() (string, error) {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, userDB, password, dbname), nil
}
