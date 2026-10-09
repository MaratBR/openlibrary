package siteconfig

type ConfigData struct {
	FontConfiguration    FontConfiguration
	CaptchaSettings      CaptchaSettings
	PasswordRequirements PasswordRequirements
	ContentRestrictions  ContentRestrictions
}

type PasswordRequirements struct {
	Digits         bool
	Symbols        string
	SymbolsEnabled bool
	DifferentCases bool
	MinLength      int
}

type ContentRestrictions struct {
	// if true - whole website is considered "adult"
	AdultWebsite bool
}

type CaptchaSettings struct {
	GoogleRecaptchaKey string
	Type               string
}

// FontConfiguration controls chapter font limits and exempt font families.
type FontConfiguration struct {
	MaxPerChapter int
	Whitelist     []string
}
