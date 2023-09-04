package user

import (
	"sync"
)

type User struct {
	ID       int
	Username string
}

type UserManager struct {
	Users []*User
	mu    sync.Mutex
}

func NewUserManager() *UserManager {
	return &UserManager{}
}

func (um *UserManager) RegisterUser(username string) *User {
	um.mu.Lock()
	defer um.mu.Unlock()

	user := &User{
		ID:       len(um.Users) + 1,
		Username: username,
	}

	um.Users = append(um.Users, user)
	return user
}

func (um *UserManager) GetUserByID(userID int) *User {
	um.mu.Lock()
	defer um.mu.Unlock()

	for _, user := range um.Users {
		if user.ID == userID {
			return user
		}
	}

	return nil
}
