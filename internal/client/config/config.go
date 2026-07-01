package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path"
)

type HostEntry struct {
	Addr string
	TLS  bool
}

type Config struct {
	Hosts map[string]HostEntry `json:"hosts"`
	path  string
}

type Session struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func configPath() string {
	if p := os.Getenv("SHARELOCK_CONFIG_DIR"); p != "" {
		return path.Join(p, ".hosts")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".hosts"
	}
	primary := path.Join(home, ".config", "sharelock", ".hosts")
	if _, err := os.Stat(primary); err == nil {
		return primary
	}
	return path.Join(home, ".hosts")
}

func Load() (*Config, error) {
	p := configPath()
	cfg := &Config{Hosts: make(map[string]HostEntry), path: p}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	var raw struct {
		Hosts map[string]HostEntry `json:"hosts"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	if raw.Hosts != nil {
		cfg.Hosts = raw.Hosts
	}
	return cfg, nil
}

func (c *Config) Save() error {
	dir := path.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(struct {
		Hosts map[string]HostEntry `json:"hosts"`
	}{Hosts: c.Hosts}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}

func (c *Config) Default() (string, bool) {
	h, _ := c.Hosts["default"]
	return h.Addr, h.TLS
}

func (c *Config) Get(name string) (HostEntry, bool) {
	h, ok := c.Hosts[name]
	return h, ok
}

func (c *Config) Set(name, addr string, tls bool) {
	c.Hosts[name] = HostEntry{Addr: addr, TLS: tls}
}

func (c *Config) Remove(name string) {
	delete(c.Hosts, name)
}

func (c *Config) List() map[string]HostEntry {
	return c.Hosts
}

func sessionPath() string {
	if p := os.Getenv("SHARELOCK_CONFIG_DIR"); p != "" {
		return path.Join(p, ".session")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".session"
	}
	primary := path.Join(home, ".config", "sharelock", ".session")
	if _, err := os.Stat(path.Dir(primary)); err == nil {
		return primary
	}
	return path.Join(home, ".session")
}

func SaveSession(username, password string) error {
	s := Session{
		Username: base64.StdEncoding.EncodeToString([]byte(username)),
		Password: base64.StdEncoding.EncodeToString([]byte(password)),
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	p := sessionPath()
	dir := path.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0600)
}

func LoadSession() (username, password string, ok bool) {
	p := sessionPath()
	data, err := os.ReadFile(p)
	if err != nil {
		return "", "", false
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return "", "", false
	}
	u, err := base64.StdEncoding.DecodeString(s.Username)
	if err != nil {
		return "", "", false
	}
	pw, err := base64.StdEncoding.DecodeString(s.Password)
	if err != nil {
		return "", "", false
	}
	return string(u), string(pw), true
}

func ClearSession() error {
	return os.Remove(sessionPath())
}
