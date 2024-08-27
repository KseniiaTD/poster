package inmemory

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *inMemoryDB) CreateSubscription(ctx context.Context, subscr model.SubscrInput) (*string, error) {

	user, err := strconv.Atoi(subscr.UserID)
	if err != nil {
		return nil, err
	}

	post, err := strconv.Atoi(subscr.PostID)
	if err != nil {
		return nil, err
	}

	db.mu.Lock()

	var id int
	p, ok := db.subscriptions[post]
	if !ok {

		postS := subscriptionPosts{
			postId:     post,
			subsribers: map[int]subscription{},
		}

		id = db.subscriptionId
		db.subscriptionId++

		sub := subscription{
			createDate:  time.Now(),
			updDate:     time.Now(),
			isDeleted:   false,
			userId:      user,
			id:          id,
			commentList: []int{},
		}
		postS.subsribers[user] = sub
		db.subscriptions[post] = postS

	} else {
		s, ok := p.subsribers[user]
		if !ok {

			id = db.subscriptionId
			db.subscriptionId++

			sub := subscription{
				createDate:  time.Now(),
				updDate:     time.Now(),
				isDeleted:   false,
				userId:      user,
				id:          id,
				commentList: []int{},
			}
			p.subsribers[user] = sub
		} else {
			s.updDate = time.Now()
			s.isDeleted = false
			id = s.id
			p.subsribers[user] = s
		}
		db.subscriptions[post] = p
	}
	db.mu.Unlock()

	str := strconv.Itoa(id)
	return &str, nil
}

func (db *inMemoryDB) DeleteSubscription(ctx context.Context, subscr model.SubscrInput) (*string, error) {

	user, err := strconv.Atoi(subscr.UserID)
	if err != nil {
		return nil, err
	}

	post, err := strconv.Atoi(subscr.PostID)
	if err != nil {
		return nil, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	var id int
	p, ok := db.subscriptions[post]
	if !ok {

		return nil, errors.New("subscription not found")

	} else {
		s, ok := p.subsribers[user]
		if !ok {
			return nil, errors.New("subscription not found")
		} else {
			s.updDate = time.Now()
			s.isDeleted = true
			id = s.id
			s.commentList = []int{}
			p.subsribers[user] = s
		}
		db.subscriptions[post] = p
	}

	str := strconv.Itoa(id)
	return &str, nil
}

func (db *inMemoryDB) CheckSubscription(ctx context.Context, subscr model.Subscr) error {
	user, err := strconv.Atoi(subscr.UserID)
	if err != nil {
		return err
	}

	post, err := strconv.Atoi(subscr.PostID)
	if err != nil {
		return err
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	p, ok := db.subscriptions[post]
	if !ok {

		return errors.New("subscription not found")

	} else {
		s, ok := p.subsribers[user]
		if !ok {
			if s.isDeleted {
				return errors.New("subscription not found")
			}
		}
	}

	return nil
}

func (db *inMemoryDB) GetCntNewCommentsForSubscriber(ctx context.Context, subscr model.Subscr) (int, error) {
	user, err := strconv.Atoi(subscr.UserID)
	if err != nil {
		return 0, err
	}

	post, err := strconv.Atoi(subscr.PostID)
	if err != nil {
		return 0, err
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	subscrPost, ok := db.subscriptions[post]
	if !ok {
		return 0, errors.New("subscription not found")
	}

	subscriber := subscrPost.subsribers[user]
	if !ok {
		return 0, errors.New("subscription not found")
	}

	commentsCnt := len(subscriber.commentList)

	subscriber.commentList = []int{}
	subscrPost.subsribers[user] = subscriber
	db.subscriptions[post] = subscrPost

	return commentsCnt, nil

}
