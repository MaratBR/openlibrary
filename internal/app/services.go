package app

type BackgroundServices struct {
	Book BookBackgroundService
}

func NewBackgroundServices(db DB) *BackgroundServices {
	return &BackgroundServices{
		Book: NewBookBackgroundService(db),
	}
}

func (s *BackgroundServices) Stop() {
	s.Book.Stop()
}

func (s *BackgroundServices) Start() error {
	if err := s.Book.Start(); err != nil {
		return err
	}
	return nil
}
