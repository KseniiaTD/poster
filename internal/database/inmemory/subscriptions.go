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

	p, ok := db.subscriptions[post]
	if !ok {

		postS := subscriptionPosts{
			postId:     post,
			subsribers: map[int]subscription{},
		}

		sub := subscription{
			createDate:  time.Now(),
			updDate:     time.Now(),
			isDeleted:   false,
			userId:      user,
			commentList: []int{},
		}
		postS.subsribers[user] = sub
		db.subscriptions[post] = postS

	} else {
		s, ok := p.subsribers[user]
		if !ok {
			sub := subscription{
				createDate:  time.Now(),
				updDate:     time.Now(),
				isDeleted:   false,
				userId:      user,
				commentList: []int{},
			}
			p.subsribers[user] = sub
		} else {
			s.updDate = time.Now()
			s.isDeleted = false
			p.subsribers[user] = s
		}
		db.subscriptions[post] = p
	}
	db.mu.Unlock()

	str := "ok"
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
			s.commentList = []int{}
			p.subsribers[user] = s
		}
		db.subscriptions[post] = p
	}
	db.mu.Unlock()

	str := "ok"
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
	db.mu.RUnlock()

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
	db.mu.Unlock()

	return commentsCnt, nil

}
