package config

func getDefaultConfig() Config {
	return Config{
		Theme: ThemeConfig{
			ShowIcons: true,
		},
	}
}
