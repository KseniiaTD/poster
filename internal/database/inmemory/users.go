package inmemory

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *inMemoryDB) CreateUser(ctx context.Context, newUser model.NewUserInput) (int, error) {

	_, ok := db.loginUniq[strings.ToLower(newUser.Login)]

	if ok {
		return 0, errors.New("login is not unique")
	}

	_, ok = db.phoneUniq[strings.ToLower(newUser.Phone)]

	if ok {
		return 0, errors.New("login is not unique")
	}

	_, ok = db.emailUniq[strings.ToLower(newUser.Email)]

	if ok {
		return 0, errors.New("login is not unique")
	}

	u := user{
		login:      newUser.Login,
		createDate: time.Now(),
		isDeleted:  time.Now(),
		name:       newUser.Name,
		surname:    newUser.Surname,
		phone:      newUser.Phone,
		email:      newUser.Email,
		posts:      make(map[int]struct{}),
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	id := db.userId
	u.id = id
	db.users[id] = u
	db.userId++

	db.loginUniq[strings.ToLower(newUser.Login)] = struct{}{}
	db.phoneUniq[strings.ToLower(newUser.Phone)] = struct{}{}
	db.emailUniq[strings.ToLower(newUser.Email)] = struct{}{}

	return id, nil
}
