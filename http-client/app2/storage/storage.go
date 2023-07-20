package storage

import (
	"fmt"
)

type user struct {
	Subject  string
	Name     string
	Audience string
	Email    string
}

type Storage struct {
	set   map[string]struct{}
	users []*user
}

var instance *Storage

func New() *Storage {
	if instance == nil {
		fmt.Println("create a new storage")
		instance = &Storage{
			set:   make(map[string]struct{}),
			users: make([]*user, 0),
		}
		return instance
	}

	fmt.Println("storage already exists")
	return instance
}

func (s *Storage) AddUser(subject, name, audience, email string) bool {
	if _, ok := s.set[subject]; ok {
		return false
	}

	u := &user{
		Subject:  subject,
		Name:     name,
		Audience: audience,
		Email:    email,
	}
	s.users = append(s.users, u)
	s.set[subject] = struct{}{}

	return true
}

func (s *Storage) AllUser() []*user {
	return s.users
}
