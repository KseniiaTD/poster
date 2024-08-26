package database

import (
	"context"
	"errors"
	"poster/graph/model"
	"sort"
	"strconv"
	"sync"
	"time"
)

type postLike struct {
	authorId int
	postId   int
}

type commentLike struct {
	authorId  int
	commentId int
}

type comment struct {
	id            int
	parentId      int
	postId        int
	authorId      int
	body          string
	createDate    time.Time
	updDate       time.Time
	isDeleted     bool
	likesCnt      int
	dislikesCnt   int
	childComments map[int]struct{}
}

type post struct {
	id          int
	title       string
	authorId    int
	body        string
	createDate  time.Time
	updDate     time.Time
	isDeleted   bool
	isCommented bool
	likesCnt    int
	dislikesCnt int
	comments    map[int]struct{}
}

type user struct {
	id         int
	login      string
	createDate time.Time
	isDeleted  time.Time
	name       string
	surname    string
	phone      string
	email      string
	posts      map[int]struct{}
}

type inMemoryDB struct {
	users        map[int]user
	userId       int
	posts        map[int]post
	postId       int
	comments     map[int]comment
	commentId    int
	postLikes    map[postLike]bool
	commentLikes map[commentLike]bool
	mu           sync.RWMutex
}

func (db *inMemoryDB) CloseDB() {
}

func (db *inMemoryDB) CreatePost(ctx context.Context, newPost model.NewPostInput) (int, error) {
	authorIdInt, err := strconv.Atoi(newPost.AuthorID)
	if err != nil {
		return 0, err
	}

	p := post{
		title:       newPost.Title,
		authorId:    authorIdInt,
		body:        newPost.Body,
		createDate:  time.Now(),
		updDate:     time.Now(),
		isDeleted:   false,
		isCommented: newPost.IsCommented,
	}

	db.mu.Lock()
	id := db.postId
	p.id = id
	db.posts[id] = p
	db.postId++

	u := db.users[authorIdInt]
	u.posts[id] = struct{}{}
	db.users[authorIdInt] = u
	db.mu.Unlock()
	return id, nil
}
func (db *inMemoryDB) UpdatePost(ctx context.Context, post model.UpdPostInput) (int, error) {
	var err error

	postIdInt, err := strconv.Atoi(post.PostID)
	if err != nil {
		return 0, err
	}

	db.mu.Lock()
	p, ok := db.posts[postIdInt]
	if !ok {
		return 0, errors.New("post not found")
	}

	if post.Title != nil {
		p.title = *post.Title
	}

	if post.Body != nil {
		p.body = *post.Body
	}

	if post.IsCommented != nil {
		p.isCommented = *post.IsCommented
	}

	p.updDate = time.Now()

	db.posts[postIdInt] = p
	db.mu.Unlock()
	return postIdInt, nil
}
func (db *inMemoryDB) DeletePost(ctx context.Context, id int) (int, error) {

	db.mu.Lock()
	p, ok := db.posts[id]
	if !ok {
		return 0, errors.New("post not found")
	}

	p.updDate = time.Now()
	p.isDeleted = true
	authorId := p.authorId
	db.posts[id] = p

	u := db.users[authorId]
	delete(u.posts, id)
	db.users[authorId] = u
	db.mu.Unlock()
	return id, nil
}

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

func (db *inMemoryDB) GetPosts(ctx context.Context, userId int) ([]*model.Post, error) {
	u, ok := db.users[userId]
	if !ok {
		return nil, errors.New("user not found")
	}

	posts := make([]post, 0, len(u.posts))
	for k := range u.posts {
		p := db.posts[k]
		if !p.isDeleted {
			posts = append(posts, p)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].updDate.After(posts[j].updDate)
	})

	postRes := make([]*model.Post, 0, len(posts))

	for _, val := range posts {
		modelP := model.Post{
			ID:          strconv.Itoa(val.id),
			Title:       val.title,
			AuthorID:    strconv.Itoa(val.authorId),
			Body:        val.body,
			CreateDate:  val.createDate,
			UpdDate:     val.updDate,
			IsCommented: val.isCommented,
			Likes:       &val.likesCnt,
			Dislikes:    &val.dislikesCnt,
		}

		postRes = append(postRes, &modelP)
	}

	return postRes, nil
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
	db.mu.Unlock()

	if !ok {
		return 0, 0, errors.New("post not found")
	}
	return value.likesCnt, value.dislikesCnt, nil
}
func (db *inMemoryDB) GetCommentLikes(ctx context.Context, commentId int) (int, int, error) {
	db.mu.RLock()
	value, ok := db.comments[commentId]
	db.mu.Unlock()

	if !ok {
		return 0, 0, errors.New("comment not found")
	}
	return value.likesCnt, value.dislikesCnt, nil
}
