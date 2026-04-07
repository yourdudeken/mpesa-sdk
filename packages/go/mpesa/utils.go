package mpesa

import (
	"encoding/base64"
	"time"
)

type Helpers struct {
	config *Config
}

func NewHelpers(config *Config) *Helpers {
	return &Helpers{config: config}
}

func (h *Helpers) PhoneValidator(phoneNo string) string {
	if len(phoneNo) > 0 && phoneNo[0] == '+' {
		phoneNo = phoneNo[1:]
	}
	if len(phoneNo) > 0 && phoneNo[0] == '0' {
		phoneNo = "254" + phoneNo[1:]
	} else if len(phoneNo) > 0 && phoneNo[0] == '7' {
		phoneNo = "254" + phoneNo
	}
	return phoneNo
}

func (h *Helpers) GetFormattedTimestamp() string {
	now := time.Now()
	return now.Format("20060102150405")
}

func (h *Helpers) LipaNaMpesaPassword() string {
	timestamp := h.GetFormattedTimestamp()
	password := h.config.Shortcode + h.config.Passkey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func (h *Helpers) GetConfig(key string) string {
	if key == "shortcode" {
		return h.config.Shortcode
	}
	if key == "passkey" {
		return h.config.Passkey
	}
	if key == "till_number" {
		return h.config.TillNumber
	}
	if key == "initiator_name" {
		return h.config.InitiatorName
	}
	return ""
}

func (h *Helpers) ResolveCallbackURL(paramURL string, configKey string) string {
	if paramURL != "" {
		return paramURL
	}
	if h.config.Callbacks != nil {
		if url, ok := h.config.Callbacks[configKey]; ok {
			return url
		}
	}
	return ""
}
