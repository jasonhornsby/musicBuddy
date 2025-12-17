package core

import "encoding/json"

const (
	ParamInt    = "int"
	ParamFloat  = "float"
	ParamSelect = "select"
)

type ParamDef struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Type    string   `json:"type"`
	Default any      `json:"default"`
	Min     float64  `json:"min,omitempty"`
	Max     float64  `json:"max,omitempty"`
	Step    float64  `json:"step,omitempty"`
	Options []string `json:"options,omitempty"`
}

type ConfigMap map[string]interface{}

func (c *ConfigMap) GetInt(key string, fallback int) int {
	if val, ok := (*c)[key]; ok {
		println("[Go] Getting int for key: ", key, " value: ", val)
		switch val := val.(type) {
		case float64:
			return int(val)
		case int:
			return val
		default:
			return fallback
		}
	}
	return fallback
}

func (c *ConfigMap) GetString(key string, fallback string) string {
	if val, ok := (*c)[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return fallback
}

func (c *ConfigMap) String() string {
	jsonStr, _ := json.Marshal(*c)
	return string(jsonStr)
}
