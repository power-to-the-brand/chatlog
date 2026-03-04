package conf

type TUIConfig struct {
	ConfigDir        string            `mapstructure:"-" json:"config_dir"`
	LastAccount      string            `mapstructure:"last_account" json:"last_account"`
	History          []ProcessConfig   `mapstructure:"history" json:"history"`
	Webhook          *Webhook          `mapstructure:"webhook" json:"webhook"`
	Postgres         *PostgresConfig   `mapstructure:"postgres" json:"postgres"`
	SupplierMappings []SupplierMapping `mapstructure:"supplier_mappings" json:"supplier_mappings"`
}

// SupplierMapping associates an account with a set of talker → supplier_id mappings.
type SupplierMapping struct {
	Account  string            `mapstructure:"account" json:"account"`
	Mappings map[string]string `mapstructure:"mappings" json:"mappings"`
}

// PostgresConfig holds PostgreSQL connection settings for sync.
type PostgresConfig struct {
	URL string `mapstructure:"url" json:"url"`
}

var TUIDefaults = map[string]any{}

type ProcessConfig struct {
	Type        string `mapstructure:"type" json:"type"`
	Account     string `mapstructure:"account" json:"account"`
	Platform    string `mapstructure:"platform" json:"platform"`
	Version     int    `mapstructure:"version" json:"version"`
	FullVersion string `mapstructure:"full_version" json:"full_version"`
	DataDir     string `mapstructure:"data_dir" json:"data_dir"`
	DataKey     string `mapstructure:"data_key" json:"data_key"`
	ImgKey      string `mapstructure:"img_key" json:"img_key"`
	WorkDir     string `mapstructure:"work_dir" json:"work_dir"`
	HTTPEnabled bool   `mapstructure:"http_enabled" json:"http_enabled"`
	HTTPAddr    string `mapstructure:"http_addr" json:"http_addr"`
	LastTime    int64  `mapstructure:"last_time" json:"last_time"`
	Files       []File `mapstructure:"files" json:"files"`
}

type File struct {
	Path         string `mapstructure:"path" json:"path"`
	ModifiedTime int64  `mapstructure:"modified_time" json:"modified_time"`
	Size         int64  `mapstructure:"size" json:"size"`
}

func (c *TUIConfig) ParseHistory() map[string]ProcessConfig {
	m := make(map[string]ProcessConfig)
	for _, v := range c.History {
		m[v.Account] = v
	}
	return m
}

// GetSupplierMappings returns the talker → supplier_id mappings for the given account.
func (c *TUIConfig) GetSupplierMappings(account string) map[string]string {
	for _, sm := range c.SupplierMappings {
		if sm.Account == account {
			return sm.Mappings
		}
	}
	return nil
}
