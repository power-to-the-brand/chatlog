package conf

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/sjzar/chatlog/pkg/config"
)

const (
	AppName          = "chatlog"
	ServerConfigName = "chatlog-server"
	EnvPrefix        = "CHATLOG"
	EnvConfigDir     = "CHATLOG_DIR"
)

// LoadTUIConfig 加载 TUI 配置
func LoadTUIConfig(configPath string) (*TUIConfig, *config.Manager, error) {

	if configPath == "" {
		configPath = os.Getenv(EnvConfigDir)
	}

	tcm, err := config.New(AppName, configPath, "", "", true)
	if err != nil {
		log.Error().Err(err).Msg("load tui config failed")
		return nil, nil, err
	}

	conf := &TUIConfig{}
	config.SetDefaults(tcm.Viper, conf, TUIDefaults)

	if err := tcm.Load(conf); err != nil {
		log.Error().Err(err).Msg("load tui config failed")
		return nil, nil, err
	}
	conf.ConfigDir = tcm.Path

	b, _ := json.Marshal(conf)
	log.Info().Msgf("tui config: %s", string(b))

	return conf, tcm, nil
}

// LoadServiceConfig 加载服务配置
func LoadServiceConfig(configPath string, cmdConf map[string]any) (*ServerConfig, *config.Manager, error) {

	if configPath == "" {
		configPath = os.Getenv(EnvConfigDir)
	}

	scm, err := config.New(AppName, configPath, ServerConfigName, EnvPrefix, false)
	if err != nil {
		log.Error().Err(err).Msg("load server config failed")
		return nil, nil, err
	}

	conf := &ServerConfig{}
	config.SetDefaults(scm.Viper, conf, ServerDefaults)

	// Load cmd Conf
	for key, value := range cmdConf {
		scm.SetConfig(key, value)
	}

	if err := scm.Load(conf); err != nil {
		log.Error().Err(err).Msg("load server config failed")
		return nil, nil, err
	}

	// Fallback to TUI config (chatlog.json) when server config has no data_dir
	if len(conf.DataDir) == 0 {
		if tuiConf, _, err := LoadTUIConfig(configPath); err == nil {
			history := tuiConf.ParseHistory()
			account := tuiConf.LastAccount
			if account == "" && len(tuiConf.History) > 0 {
				account = tuiConf.History[0].Account
			}
			if pc, ok := history[account]; ok && len(pc.DataDir) > 0 {
				conf.Type = pc.Type
				conf.Platform = pc.Platform
				conf.Version = pc.Version
				conf.FullVersion = pc.FullVersion
				conf.DataDir = pc.DataDir
				conf.DataKey = pc.DataKey
				conf.ImgKey = pc.ImgKey
				conf.WorkDir = pc.WorkDir
				log.Info().Msgf("using account %q from TUI config", account)
			}
		}
	}

	// Load Data Dir config
	if len(conf.DataDir) != 0 && len(conf.DataKey) == 0 {
		if b, err := os.ReadFile(filepath.Join(conf.DataDir, "chatlog.json")); err == nil {
			var pconf map[string]any
			if err := json.Unmarshal(b, &pconf); err == nil {
				for key, value := range pconf {
					if !DataDirConfigs[key] {
						continue
					}
					scm.SetConfig(key, value)
				}
			}
		}
		if err := scm.Load(conf); err != nil {
			log.Error().Err(err).Msg("reload server config failed")
			return nil, nil, err
		}
	}

	b, _ := json.Marshal(conf)
	log.Info().Msgf("server config: %s", string(b))

	return conf, scm, nil
}

var DataDirConfigs = map[string]bool{
	"type":         true,
	"platform":     true,
	"version":      true,
	"full_version": true,
	"data_key":     true,
	"img_key":      true,
}
