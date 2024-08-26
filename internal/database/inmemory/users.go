package inmemory

import (
	"context"
	"time"

	"github.com/KseniiaTD/poster/graph/model"
)

func (db *inMemoryDB) CreateUser(ctx context.Context, newUser model.NewUserInput) (int, error) {
	u := user{
		login:      newUser.Login,
		createDate: time.Now(),
		isDeleted:  time.Now(),
		name:       newUser.Name,
		surname:    newUser.Surname,
		phone:      newUser.Phone,
		email:      newUser.Email,
	}

	db.mu.Lock()
	id := db.userId
	u.id = id
	db.users[id] = u
	db.userId++
	db.mu.Unlock()
	return id, nil
}
