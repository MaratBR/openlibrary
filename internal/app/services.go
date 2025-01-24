package app

type BackgroundServices struct {
	Book      BookBackgroundService
	Favorites *FavoriteRecalculationBackgroundService
}

func NewBackgroundServices(db DB) *BackgroundServices {
	return &BackgroundServices{
		Book:      NewBookBackgroundService(db),
		Favorites: NewFavoriteRecalculationBackgroundService(db),
	}
}

func (s *BackgroundServices) Stop() {
	s.Favorites.Stop()
	s.Book.Stop()
}

func (s *BackgroundServices) Start() error {
	if err := s.Favorites.Start(); err != nil {
		return err
	}
	if err := s.Book.Start(); err != nil {
		return err
	}
	return nil
}
