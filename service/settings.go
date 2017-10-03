package service

import "github.com/jennal/goplay/log"

type Settings struct {
	IsDisconnectOnError bool
}

var _defaultSettings = &Settings{
	IsDisconnectOnError: true,
}

func DefaultSettings() *Settings {
	return _defaultSettings
}

type SettingContainer struct {
	settings *Settings
}

func NewSettingContainer() *SettingContainer {
	return &SettingContainer{
		settings: DefaultSettings(),
	}
}

func (sc *SettingContainer) SetSettings(s *Settings) error {
	if s == nil {
		return log.NewErrorf("can't set settings to nil")
	}

	sc.settings = s
	return nil
}

func (sc *SettingContainer) Settings() *Settings {
	return sc.settings
}
