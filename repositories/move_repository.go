package repositories

type MovieRepository interface {
	GetMovieName() string
}

type MovieManager struct {
}

func NewMovieManager() MovieRepository {
	return &MovieManager{}
}

func (m *MovieManager) GetMovieName() string {
	// database operation
	//movie := &datamodels.Movie{Name: "movie test"}
	//return movie.Name
	return "movie test"
}
