package goini

import (
	"errors"
	"strconv"
	"strings"
)

const (
	COMMENT   = ';'
	SEPARATOR = '='

	ERR_KEY_NOT_EXISTS = "Key not exist"
)

var (
	// Strings accepted as boolean.
	boolString = map[string]bool{
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
		"on":    true,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"0":     false,
	}
)

// raw value
type rawValue string

func (r rawValue) String() string {
	return string(r)
}

// section
type section string

func (s section) String() string {
	return string(s)
}

//config key
type key string

func (k key) String() string {
	return string(k)
}

// Config is the representation of configuration settings.
type Config struct {
	inheritance map[section]section
	data        map[section]map[key]rawValue
}

type OptionsMap struct {
	data map[key]rawValue
}

func (o *OptionsMap) Len() int {
	return len(o.data)
}

func (cfg *Config) GetSection(s string) *OptionsMap {
	c := &OptionsMap{}
	p, ok := cfg.inheritance[section(s)]
	if !ok {
		c.data = make(map[key]rawValue)
	} else {
		c = cfg.GetSection(string(p))
	}

	for k, v := range cfg.data[section(s)] {
		c.data[k] = v
	}
	return c
}

func (cfg *Config) GetSectionList() []string {
	a := make([]string, len(cfg.data))
	i := 0
	for k := range cfg.data {
		a[i] = string(k)
		i++
	}
	return a
}

// Get value as a string, remove quotes if value was quoted into " or '
func (o *OptionsMap) GetString(k string) (string, error) {
	value, exist := o.data[key(k)]
	if !exist {
		return "", errors.New(ERR_KEY_NOT_EXISTS)
	}
	rawString := string(value)
	if len(rawString) > 1 {
		if rawString[0] == '\'' && rawString[len(rawString)-1] == '\'' {
			rawString = rawString[1 : len(rawString)-1]
		} else if rawString[0] == '"' && rawString[len(rawString)-1] == '"' {
			rawString = rawString[1 : len(rawString)-1]
		}
	}
	return rawString, nil
}

// synonym for GetString
func (o *OptionsMap) String(k string) (string, error) {
	return o.GetString(k)
}

// Get value as a bool, uses this mapping from string to bool
// "t":     true,
// "true":  true,
// "y":     true,
// "yes":   true,
// "on":    true,
// "1":     true,
// "f":     false,
// "false": false,
// "n":     false,
// "no":    false,
// "off":   false,
// "0":     false,
func (o *OptionsMap) GetBool(k string) (bool, error) {
	sv, exists := o.data[key(k)]
	if !exists {
		return false, errors.New(ERR_KEY_NOT_EXISTS)
	}

	value, ok := boolString[strings.ToLower(string(sv))]
	if !ok {
		return false, errors.New("could not parse bool value: " + string(sv))
	}

	return value, nil
}

// synonym for GetBool
func (o *OptionsMap) Bool(k string) (bool, error) {
	return o.GetBool(k)
}

// Get value as int
func (o *OptionsMap) GetInt(k string) (int, error) {
	sv, exists := o.data[key(k)]
	if exists {
		return strconv.Atoi(string(sv))
	}

	return 0, errors.New(ERR_KEY_NOT_EXISTS)
}

// synonym for GetInt
func (o *OptionsMap) Int(k string) (int, error) {
	return o.GetInt(k)
}

// Get value as float64
func (o *OptionsMap) GetFloat(k string) (float64, error) {
	sv, exitsts := o.data[key(k)]
	if exitsts {
		return strconv.ParseFloat(string(sv), 64)
	}

	return 0, errors.New(ERR_KEY_NOT_EXISTS)
}

// synonym for GetFloat
func (o *OptionsMap) Float64(k string) (float64, error) {
	return o.GetFloat(k)
}
