package inmemory

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/KseniiaTD/poster/graph/model"
)

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

	/*newC := newPostComments{
		postId:      id,
		commentList: []int{},
	}
	db.newComments[id] = newC*/
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
