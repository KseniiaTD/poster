package tests

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"time"

	"github.com/KseniiaTD/poster/config"
	"github.com/KseniiaTD/poster/graph"
	"github.com/KseniiaTD/poster/graph/model"
)

func initTestConfig() config.Config {
	return config.Config{
		DB:       "test_poster",
		User:     "pguser",
		Password: "postgres",
		DBHost:   "localhost",
		DBPort:   "5432",
	}
}

func randomPhoneNumber() string {
	b := make([]byte, 7)
	n, err := io.ReadAtLeast(rand.Reader, b, 7)
	if n != 7 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func createTestUser(resolver graph.Resolver) (string, error) {
	input := model.NewUserInput{
		Login:   fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:    "test name",
		Surname: "test surname",
		Phone:   "8800" + randomPhoneNumber(),
		Email:   fmt.Sprintf("test%d@email.com", time.Now().UnixNano()),
	}

	id, err := resolver.Mutation().CreateUser(context.Background(), input)
	if err != nil {
		return "", err
	}

	return *id, nil
}

func createTestPost(resolver graph.Resolver) (string, string, error) {
	userId, err := createTestUser(resolver)

	if err != nil {
		return "", "", err
	}

	input := model.NewPostInput{
		Title:       "test title",
		AuthorID:    userId,
		Body:        "Test Body",
		IsCommented: true,
	}
	id, err := resolver.Mutation().CreatePost(context.Background(), input)
	if err != nil {
		return "", "", err
	}

	return *id, userId, nil
}

func createTestComment(resolver graph.Resolver) (string, string, string, error) {
	postId, userId, err := createTestPost(resolver)
	if err != nil {
		return "", "", "", err
	}

	input := model.NewCommentInput{
		AuthorID: userId,
		Body:     "Test Body",
		PostID:   postId,
	}

	id, err := resolver.Mutation().CreateComment(context.Background(), input)
	if err != nil {
		return "", "", "", err
	}

	return *id, postId, userId, nil
}
