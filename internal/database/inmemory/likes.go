package inmemory

import (
	"context"
	"errors"
)

func (db *inMemoryDB) UpdPostLikes(ctx context.Context, authorId int, postId int, isLike bool) error {
	pl := postLike{authorId: authorId, postId: postId}

	db.mu.Lock()
	value, ok := db.postLikes[pl]
	if !ok || value != isLike {
		db.postLikes[pl] = isLike
		if isLike {
			p := db.posts[postId]
			p.likesCnt++
			db.posts[postId] = p
		} else {
			p := db.posts[postId]
			p.dislikesCnt++
			db.posts[postId] = p
		}
	} else {
		delete(db.postLikes, pl)
		if isLike {
			p := db.posts[postId]
			p.likesCnt--
			db.posts[postId] = p
		} else {
			p := db.posts[postId]
			p.dislikesCnt--
			db.posts[postId] = p
		}
	}
	db.mu.Unlock()

	return nil
}
func (db *inMemoryDB) UpdCommentLikes(ctx context.Context, authorId int, commentId int, isLike bool) error {
	cl := commentLike{authorId: authorId, commentId: commentId}

	db.mu.Lock()
	value, ok := db.commentLikes[cl]
	if !ok || value != isLike {
		db.commentLikes[cl] = isLike
		if isLike {
			c := db.comments[commentId]
			c.likesCnt++
			db.comments[commentId] = c
		} else {
			c := db.comments[commentId]
			c.dislikesCnt++
			db.comments[commentId] = c
		}
	} else {
		delete(db.commentLikes, cl)
		if isLike {
			c := db.comments[commentId]
			c.likesCnt--
			db.comments[commentId] = c
		} else {
			c := db.comments[commentId]
			c.dislikesCnt--
			db.comments[commentId] = c
		}
	}
	db.mu.Unlock()

	return nil
}

func (db *inMemoryDB) GetPostLikes(ctx context.Context, postId int) (int, int, error) {
	db.mu.RLock()
	value, ok := db.posts[postId]
	db.mu.RUnlock()

	if !ok {
		return 0, 0, errors.New("post not found")
	}
	return value.likesCnt, value.dislikesCnt, nil
}
func (db *inMemoryDB) GetCommentLikes(ctx context.Context, commentId int) (int, int, error) {
	db.mu.RLock()
	value, ok := db.comments[commentId]
	db.mu.RUnlock()

	if !ok {
		return 0, 0, errors.New("comment not found")
	}
	return value.likesCnt, value.dislikesCnt, nil
}
