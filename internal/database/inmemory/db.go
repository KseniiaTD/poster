package inmemory

import (
	"sync"

	"time"

	"github.com/KseniiaTD/poster/internal/database/common"
)

type postLike struct {
	authorId int
	postId   int
}

type commentLike struct {
	authorId  int
	commentId int
}

/*type newPostComments struct {
	postId      int
	commentList []int
	subscribers map[int]struct{}
}*/

type subscription struct {
	userId      int
	createDate  time.Time
	updDate     time.Time
	isDeleted   bool
	commentList []int
}

type subscriptionPosts struct {
	postId     int
	subsribers map[int]subscription
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

type InMemoryDatabase interface {
	common.Database
}

type inMemoryDB struct {
	users         map[int]user
	userId        int
	posts         map[int]post
	postId        int
	comments      map[int]comment
	commentId     int
	postLikes     map[postLike]bool
	commentLikes  map[commentLike]bool
	subscriptions map[int]subscriptionPosts
	//newComments   map[int]newPostComments
	mu sync.RWMutex
}

func New() InMemoryDatabase {
	return &inMemoryDB{
		users:         make(map[int]user, 0),
		posts:         make(map[int]post, 0),
		comments:      make(map[int]comment, 0),
		postLikes:     make(map[postLike]bool, 0),
		commentLikes:  make(map[commentLike]bool, 0),
		subscriptions: make(map[int]subscriptionPosts, 0),
		//newComments:   make(map[int]newPostComments, 0),
	}
}

// реализует интерфейс Database
func (db *inMemoryDB) CloseDB() {}
