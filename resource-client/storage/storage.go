package storage

import (
	"fmt"
)

type repo struct {
	Name      string
	CreatedBy string
}

type Storage struct {
	set   map[string]struct{}
	repos []*repo
}

var instance *Storage // single instance

func New() *Storage {
	if instance == nil {
		fmt.Println("create a new storage")
		instance = &Storage{
			set:   make(map[string]struct{}),
			repos: make([]*repo, 0),
		}
		return instance
	}

	fmt.Println("storage already exists")
	return instance
}

func (s *Storage) AddRepo(name, subject string) bool {
	if _, ok := s.set[name]; ok {
		return false
	}

	r := &repo{
		Name:      name,
		CreatedBy: subject,
	}
	s.repos = append(s.repos, r)
	s.set[name] = struct{}{}
	return true
}

func (s *Storage) GetRepoBySubject(subject string) []*repo {
	var res []*repo
	for _, r := range s.repos {
		if r.CreatedBy == subject {
			res = append(res, r)
		}
	}
	return res
}

func (s *Storage) AllRepo() []*repo {
	return s.repos
}
