package config

import "strconv"

type smtpConfig struct {
	domain   string
	username string
	password string
	port     int
}

func Smtp() smtpConfig {
	return appConfig.smtp
}

func (c smtpConfig) Domain() string {
	return c.domain
}

func (c smtpConfig) Port() string {
	return strconv.Itoa(c.port)
}

func (c smtpConfig) Username() string {
	return c.username
}

func (c smtpConfig) Password() string {
	return c.password
}

func newSmtpConfig() smtpConfig {
	return smtpConfig{
		domain:   readEnvString("SMTP_DOMAIN"),
		username: readEnvString("SMTP_USERNAME"),
		password: readEnvString("SMTP_PASSWORD"),
		port:     readEnvInt("SMTP_PORT"),
	}
}
