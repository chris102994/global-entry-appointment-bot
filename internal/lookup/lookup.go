package lookup

type Lookup struct {
	City  string `mapstructure:"city,omitempty"`
	State string `mapstructure:"state,omitempty"`
}
