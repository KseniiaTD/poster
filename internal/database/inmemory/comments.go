package inmemory

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *inMemoryDB) CreateComment(ctx context.Context, newComment model.NewCommentInput) (int, error) {

	var parentIdInt int
	var err error
	if newComment.ParentID == nil {
		parentIdInt = 0
	} else {
		parentIdInt, err = strconv.Atoi(*newComment.ParentID)
		if err != nil {
			return 0, err
		}
	}

	authorIdInt, err := strconv.Atoi(newComment.AuthorID)
	if err != nil {
		return 0, err
	}

	postIdInt, err := strconv.Atoi(newComment.PostID)
	if err != nil {
		return 0, err
	}

	c := comment{
		parentId:   parentIdInt,
		postId:     postIdInt,
		authorId:   authorIdInt,
		body:       newComment.Body,
		createDate: time.Now(),
		updDate:    time.Now(),
		isDeleted:  false,
	}

	db.mu.Lock()
	id := db.commentId
	c.id = id
	db.comments[id] = c
	db.commentId++

	p := db.posts[postIdInt]
	p.comments[id] = struct{}{}
	db.posts[postIdInt] = p

	parent := db.comments[parentIdInt]
	parent.childComments[id] = struct{}{}
	db.comments[parentIdInt] = parent

	s, ok := db.subscriptions[postIdInt]
	if ok {
		for _, v := range s.subsribers {
			if !v.isDeleted {
				v.commentList = append(v.commentList, id)
			}
		}
		db.subscriptions[postIdInt] = s
	}

	db.mu.Unlock()
	return id, nil
}
func (db *inMemoryDB) UpdateComment(ctx context.Context, comment model.UpdCommentInput) (int, error) {
	var err error

	commentIdInt, err := strconv.Atoi(comment.ID)
	if err != nil {
		return 0, err
	}

	db.mu.Lock()
	c, ok := db.comments[commentIdInt]
	if !ok {
		return 0, errors.New("comment not found")
	}
	c.body = comment.Body
	c.updDate = time.Now()
	db.comments[commentIdInt] = c
	db.mu.Unlock()
	return commentIdInt, nil
}

func (db *inMemoryDB) DeleteComment(ctx context.Context, id int) (int, error) {
	db.mu.Lock()
	c, ok := db.comments[id]
	if !ok {
		return 0, errors.New("post not found")
	}

	c.updDate = time.Now()
	c.isDeleted = true
	postId := c.postId
	parentId := c.parentId
	db.comments[id] = c

	p := db.posts[postId]
	delete(p.comments, id)
	db.posts[postId] = p

	parent := db.comments[parentId]
	delete(parent.childComments, id)
	db.comments[parentId] = parent
	db.mu.Unlock()
	return id, nil
}

func (db *inMemoryDB) GetComments(ctx context.Context, postId int, commentId int, page int, perPage int) ([]*model.Comment, error) {
	p, ok := db.posts[postId]
	if !ok {
		return nil, errors.New("post not found")
	}

	comments := make([]comment, 0, len(p.comments))
	for k := range p.comments {
		c := db.comments[k]
		if c.parentId == commentId && !(c.isDeleted && len(c.childComments) == 0) {
			if c.isDeleted {
				c.body = "Comment was deleted"
			}
			comments = append(comments, c)
		}
	}

	sort.Slice(comments, func(i, j int) bool {
		return comments[i].createDate.Before(comments[j].createDate)
	})

	elemStart := page * perPage
	elemFinish := (page + 1) * perPage

	commentsPerPage := comments[min(len(comments), elemStart):min(len(comments), elemFinish)]

	if len(commentsPerPage) == 0 {
		return make([]*model.Comment, 0), nil
	}

	childcommentsPerPage := make(map[int]comment, len(commentsPerPage))

	for _, k := range commentsPerPage {
		var firstComment comment
		for j := range k.childComments {
			value := db.comments[j]
			if firstComment.id == 0 && !(value.isDeleted && len(value.childComments) == 0) {
				firstComment = value
				continue
			}
			if value.createDate.Before(firstComment.createDate) && !(value.isDeleted && len(value.childComments) == 0) {
				firstComment = value
			}
		}
		if firstComment.id != 0 {
			if firstComment.isDeleted {
				firstComment.body = "Comment was deleted"
			}
			childcommentsPerPage[firstComment.parentId] = firstComment
		}
	}

	/*commentsRes := make([]comment, 0, len(commentsPerPage)+len(childcommentsPerPage))
	for _, k := range commentsPerPage {
		commentsRes = append(commentsRes, k)
		value, ok := childcommentsPerPage[k.id]
		if ok {
			commentsRes = append(commentsRes, value)
		}
	}*/

	commentModelRes := make([]*model.Comment, 0, len(commentsPerPage))

	for _, val := range commentsPerPage {
		u := db.users[val.authorId]
		parentIdStr := strconv.Itoa(val.parentId)
		childCnt := len(val.childComments)

		modelC := model.Comment{
			ID:         strconv.Itoa(val.id),
			ParentID:   &parentIdStr,
			PostID:     strconv.Itoa(val.postId),
			AuthorID:   strconv.Itoa(val.authorId),
			Body:       val.body,
			Login:      u.login,
			CreateDate: val.createDate,
			UpdDate:    val.updDate,
			Likes:      &val.likesCnt,
			Dislikes:   &val.dislikesCnt,
			Cnt:        &childCnt,
		}

		commentModelRes = append(commentModelRes, &modelC)

		childComment, ok := childcommentsPerPage[val.id]
		if ok {
			u := db.users[val.authorId]
			parentIdStr := strconv.Itoa(childComment.parentId)
			childCnt := len(childComment.childComments)

			modelC := model.Comment{
				ID:         strconv.Itoa(childComment.id),
				ParentID:   &parentIdStr,
				PostID:     strconv.Itoa(childComment.postId),
				AuthorID:   strconv.Itoa(childComment.authorId),
				Body:       childComment.body,
				Login:      u.login,
				CreateDate: childComment.createDate,
				UpdDate:    childComment.updDate,
				Likes:      &childComment.likesCnt,
				Dislikes:   &childComment.dislikesCnt,
				Cnt:        &childCnt,
			}

			commentModelRes = append(commentModelRes, &modelC)

		}
	}

	return commentModelRes, nil

}
