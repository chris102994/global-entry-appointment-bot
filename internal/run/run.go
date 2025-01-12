package run

type Run struct {
	Limit       int      `mapstructure:"limit,omitempty"`
	Minimum     int      `mapstructure:"minimum,omitempty"`
	OrderBy     string   `mapstructure:"orderby,omitempty"`
	LocationIDs []string `mapstructure:"locationids,omitempty"`
}
