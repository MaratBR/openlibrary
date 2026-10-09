package app

// ReaderFont describes a bundled font available to the reader.
type ReaderFont struct {
	ID       string `json:"id"`
	LabelKey string `json:"labelKey"`
	Family   string `json:"family"`
}

type ReaderFontService interface {
	List() []ReaderFont
}

type readerFontService struct{}

func NewReaderFontService() ReaderFontService { return readerFontService{} }

func (readerFontService) List() []ReaderFont {
	return []ReaderFont{
		{ID: ReaderFontSerif, LabelKey: "reader.fontSerif", Family: "Libron, Merriweather, serif"},
		{ID: ReaderFontSans, LabelKey: "reader.fontSans", Family: "'Atkinson Hyperlegible', sans-serif"},
		{ID: ReaderFontDyslexic, LabelKey: "reader.fontDyslexic", Family: "OpenDyslexic, sans-serif"},
	}
}
