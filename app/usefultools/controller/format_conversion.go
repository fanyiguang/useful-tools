package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"useful-tools/app/usefultools/adapter"
	"useful-tools/app/usefultools/i18n"

	"github.com/BurntSushi/toml"
	"github.com/magiconair/properties"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var _ adapter.Controller = (*FormatConversion)(nil)

type FormatConversion struct {
	Base
	viewText       string
	data           string
	targetFormat   string
}

func NewFormatConversion() *FormatConversion {
	return &FormatConversion{}
}

func (c *FormatConversion) Data() string {
	return c.data
}

func (c *FormatConversion) SetData(data string) {
	c.data = data
}

func (c *FormatConversion) ViewText() string {
	return c.viewText
}

func (c *FormatConversion) SetViewText(viewText string) {
	c.viewText = viewText
}

func (c *FormatConversion) TargetFormat() string {
	return c.targetFormat
}

func (c *FormatConversion) SetTargetFormat(format string) {
	c.targetFormat = format
}

func (c *FormatConversion) GetFormats() []string {
	return []string{
		i18n.T(i18n.KeyFormatJson),
		i18n.T(i18n.KeyFormatYaml),
		i18n.T(i18n.KeyFormatToml),
		i18n.T(i18n.KeyFormatProperties),
	}
}

func (c *FormatConversion) Convert(targetFormat, content string) (string, error) {
	logrus.Infof("format conversion target: %s", targetFormat)

	content = strings.TrimSpace(content)
	if content == "" {
		return "", errors.New(i18n.T(i18n.KeyHintRequired))
	}

	// 1. Detect Source Format
	var data interface{}
	var err error
	sourceFormat := c.detectFormat(content)
	
	if sourceFormat == "" {
		return "", errors.New(i18n.T(i18n.KeyFormatNotSupportedError))
	}

	// 2. Unmarshal based on detected format
	switch sourceFormat {
	case "json":
		err = json.Unmarshal([]byte(content), &data)
	case "yaml":
		err = yaml.Unmarshal([]byte(content), &data)
	case "toml":
		err = toml.Unmarshal([]byte(content), &data)
	case "properties":
		p, pErr := properties.LoadString(content)
		if pErr != nil {
			err = pErr
		} else {
			data = p.Map()
		}
	default:
		return "", errors.New(i18n.T(i18n.KeyFormatNotSupportedError))
	}

	if err != nil {
		logrus.Errorf("unmarshal error for format %s: %v", sourceFormat, err)
		return "", errors.New(i18n.T(i18n.KeyFormatError))
	}

	// 3. Marshal to target format
	var result []byte
	var resultStr string

	switch targetFormat {
	case i18n.T(i18n.KeyFormatJson):
		result, err = json.MarshalIndent(data, "", "    ")
		resultStr = string(result)
	case i18n.T(i18n.KeyFormatYaml):
		result, err = yaml.Marshal(data)
		resultStr = string(result)
	case i18n.T(i18n.KeyFormatToml):
		var buf strings.Builder
		enc := toml.NewEncoder(&buf)
		err = enc.Encode(data)
		resultStr = buf.String()
	case i18n.T(i18n.KeyFormatProperties):
		// Properties encoding is tricky because data is interface{}, but properties expects map[string]string
		// We'll do a best-effort conversion to map[string]string via JSON intermediate
		// or direct traversal if possible.
		// A simple way is to flatten the map.
		flat, flatErr := flattenMap(data)
		if flatErr != nil {
			err = flatErr
		} else {
			// Sort keys for deterministic output
			// properties library writes sorted keys by default? 
			// No, LoadString -> Map gives map[string]string.
			// Write method exists.
			// But creating a properties object from map[string]string is easier.
			// properties.LoadMap(flat)
			p := properties.LoadMap(flat)
			var buf strings.Builder
			// Write with sorted keys
			_, err = p.Write(&buf, properties.UTF8)
			resultStr = buf.String()
		}
	default:
		return "", errors.New(i18n.T(i18n.KeyFormatNotSupportedError))
	}

	if err != nil {
		logrus.Errorf("marshal error to format %s: %v", targetFormat, err)
		return "", err
	}

	return resultStr, nil
}

func (c *FormatConversion) detectFormat(content string) string {
	// Simple heuristic detection
	
	// Try JSON
	var jsonData interface{}
	if json.Unmarshal([]byte(content), &jsonData) == nil {
		return "json"
	}

	// Try YAML
	// YAML is very permissible, so we need to be careful.
	// Often anything can be YAML. 
	// Check for common YAML indicators or structure?
	// Actually, strict yaml unmarshal might be okay if JSON failed.
	var yamlData interface{}
	if yaml.Unmarshal([]byte(content), &yamlData) == nil {
		// However, simple strings are also valid YAML.
		// Maybe prioritizing TOML or Properties before YAML if structure matches?
		
		// If it's a map or slice, it's likely YAML.
		switch yamlData.(type) {
		case map[string]interface{}, []interface{}, map[interface{}]interface{}:
             // Valid structural YAML
		default:
             // Just a scalar, might be property file line or just text
		}
	}

	// Try TOML
	var tomlData interface{}
	if _, err := toml.Decode(content, &tomlData); err == nil {
		return "toml"
	}

	// Try Properties
	// Properties are key=value lines.
	if p, err := properties.LoadString(content); err == nil && len(p.Map()) > 0 {
		// Check if it looks like actual properties (has =)
		if strings.Contains(content, "=") || strings.Contains(content, ":") {
			return "properties"
		}
	}
	
	// Fallback logic
	// If it parsed as YAML structure, return YAML
	if yaml.Unmarshal([]byte(content), &yamlData) == nil {
		switch yamlData.(type) {
		case map[string]interface{}, []interface{}, map[interface{}]interface{}:
			return "yaml"
		}
	}

	return ""
}

func flattenMap(input interface{}) (map[string]string, error) {
	result := make(map[string]string)
	
	var traverse func(prefix string, v interface{})
	traverse = func(prefix string, v interface{}) {
		switch val := v.(type) {
		case map[string]interface{}:
			for k, sub := range val {
				newKey := k
				if prefix != "" {
					newKey = prefix + "." + k
				}
				traverse(newKey, sub)
			}
		case map[interface{}]interface{}:
			for k, sub := range val {
				keyStr := fmt.Sprintf("%v", k)
				newKey := keyStr
				if prefix != "" {
					newKey = prefix + "." + keyStr
				}
				traverse(newKey, sub)
			}
		case []interface{}:
			for i, sub := range val {
				newKey := fmt.Sprintf("%s[%d]", prefix, i)
				traverse(newKey, sub)
			}
		default:
			result[prefix] = fmt.Sprintf("%v", val)
		}
	}
	
	traverse("", input)
	return result, nil
}

func (c *FormatConversion) ClearCache() {
	c.data = ""
	c.viewText = ""
	c.targetFormat = ""
}
