package app

// Compatibility with existing callers; new feature code belongs in reader.
import (
	"github.com/MaratBR/openlibrary/internal/app/reader"
)

const ReaderFontSerif = reader.ReaderFontSerif
const ReaderFontSans = reader.ReaderFontSans
const ReaderFontDyslexic = reader.ReaderFontDyslexic
const ReaderPageBackground = reader.ReaderPageBackground
const ReaderPageSurface = reader.ReaderPageSurface
const ReaderThemeSystem = reader.ReaderThemeSystem
const ReaderThemeLight = reader.ReaderThemeLight
const ReaderThemeDark = reader.ReaderThemeDark
const ReaderThemeOLED = reader.ReaderThemeOLED

type ReaderPreferences = reader.ReaderPreferences

func DefaultReaderPreferences() ReaderPreferences { return reader.DefaultReaderPreferences() }

type ReaderPreferencesService = reader.ReaderPreferencesService
type ReaderFont = reader.ReaderFont
type ReaderFontService = reader.ReaderFontService

func NewReaderFontService() ReaderFontService { return reader.NewReaderFontService() }
func NewReaderPreferencesService(db DB, fonts ReaderFontService) ReaderPreferencesService {
	return reader.NewReaderPreferencesService(db, fonts)
}
