package binding

type ResFeatureFlags struct {
	Version      string       `json:"version"`
	FeatureFlags FeatureFlags `json:"featureFlags"`
}

type FeatureFlags struct {
	PublicRegistration string `json:"PUBLIC_REGISTRATION"`
	WebDav             string `json:"WEBDAV_ENABLED"`
}

// ToDo: Convert string to bool
