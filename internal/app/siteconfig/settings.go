package siteconfig

var settingsDefinitions = []SettingsDefinition{
	Definition[FontConfiguration]{
		Name:  "FontConfiguration",
		Field: func(c *ConfigData) *FontConfiguration { return &c.FontConfiguration },
		Migrations: []func(*FontConfiguration) error{
			func(c *FontConfiguration) error {
				if c.MaxPerChapter == 0 && c.Whitelist == nil {
					c.MaxPerChapter = 10
				}
				return nil
			},
		},
	},
	Definition[CaptchaSettings]{Name: "CaptchaSettings", Field: func(c *ConfigData) *CaptchaSettings { return &c.CaptchaSettings }},
	Definition[PasswordRequirements]{Name: "PasswordRequirements", Field: func(c *ConfigData) *PasswordRequirements { return &c.PasswordRequirements }},
	Definition[ContentRestrictions]{Name: "ContentRestrictions", Field: func(c *ConfigData) *ContentRestrictions { return &c.ContentRestrictions }},
}
