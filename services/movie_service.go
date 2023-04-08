package services

import (
	"fmt"
	"pro-iris/repositories"
)

type MovieService interface {
	ShowMovieName() string
}

type MovieServiceManager struct {
	repo repositories.MovieRepository
}

func NewMovieServiceManager(repo repositories.MovieRepository) MovieService {
	return &MovieServiceManager{repo: repo}
}

func (m *MovieServiceManager) ShowMovieName() string {
	fmt.Println("Name: " + m.repo.GetMovieName())
	return "Name: " + m.repo.GetMovieName()
}
