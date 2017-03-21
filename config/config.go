package config

//ACConfig struct for AccessServer start
type ACConfig struct {
	Title string `toml:"title"`
	Owner struct {
		Name  string `toml:"name"`
		Email string `toml:"Email"`
	} `toml:"owner"`
	Server struct {
		ListenIP      string `toml:"listen_ip"`
		Port          int    `toml:"port"`
		PortInternal  int    `toml:"port_internal"`
		MaxClient     int    `toml:"max_client"`
		MaxChannelBuf int    `toml:"max_channel_buf"`
	} `toml:"server"`
}
