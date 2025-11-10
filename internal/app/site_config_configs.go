package app

func siteConfigGetHelper[T any](c *SiteConfig, key string) *T {
	v := new(T)
	_, ok := c.get(key, v)
	if !ok {
		return new(T)
	}
	return v
}

type CaptchaSettings struct {
	GoogleRecaptchaKey string
	Type               string
}

func (c *SiteConfig) CaptchaSettings() *CaptchaSettings {
	return siteConfigGetHelper[CaptchaSettings](c, "CaptchaSettings")
}

func (c *SiteConfig) SetCaptchaSettings(v *CaptchaSettings) {
	c.set("CaptchaSettings", v)
}

type PasswordRequirements struct {
	Digits         bool
	Symbols        bool
	DifferentCases bool
	MinLength      int
}

func (c *SiteConfig) PasswordRequirements() *PasswordRequirements {
	return siteConfigGetHelper[PasswordRequirements](c, "PasswordRequirements")
}

func (c *SiteConfig) SetPasswordRequirements(v *PasswordRequirements) {
	c.set("PasswordRequirements", v)
}

type ContentRestrictions struct {
	// if true - whole website is considered "adult"
	AdultWebsite bool
}

func (c *SiteConfig) ContentRestrictions() *ContentRestrictions {
	return siteConfigGetHelper[ContentRestrictions](c, "ContentRestrictions")
}

func (c *SiteConfig) SetContentRestrictions(v *ContentRestrictions) {
	c.set("ContentRestrictions", v)
}
