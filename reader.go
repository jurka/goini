package goini

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

func LoadConfig(fname string) (*Config, error) {
	cfg := &Config{}
	cfg.data = make(map[section]map[key]rawValue)
	cfg.inheritance = make(map[section]section)
	return _read(fname, cfg)
}

func _read(fname string, c *Config) (*Config, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	if err = c.read(bufio.NewReader(file)); err != nil {
		return nil, err
	}

	if err = file.Close(); err != nil {
		return nil, err
	}

	return c, nil
}

func (cfg *Config) read(buf *bufio.Reader) error {
	var sec, option string
	var current_section section

	for {
		l, err := buf.ReadString('\n') // parse line-by-line
		if err != nil {
			if err != io.EOF {
				return err
			} else if len(l) == 0 {
				break
			}
		}

		l = strings.TrimSpace(l)

		// Switch written for readability (not performance)
		switch {
		// Empty line and comments
		case len(l) == 0, l[0] == uint8(COMMENT):
			continue

		// New section
		case l[0] == '[' && l[len(l)-1] == ']':
			option = "" // reset multi-line value ???
			sec = strings.TrimSpace(l[1 : len(l)-1])
			pos := strings.Index(sec, ":")
			if pos != -1 {
				current_section = section(strings.TrimSpace(sec[0 : pos-1]))
				cfg.inheritance[current_section] = section(strings.TrimSpace(sec[pos+1:]))
			} else {
				current_section = section(sec)
			}
			cfg.data[current_section] = make(map[key]rawValue)

		// No new section and no section defined so
		// case sec == "":
		// 	return errors.NewError("no section defined")

		// Other alternatives
		default:
			i := strings.Index(l, string(SEPARATOR))

			switch {
			// Option and value
			case i > 0:
				option = strings.TrimSpace(l[0:i])
				value := strings.TrimSpace(stripComments(l[i+1:]))
				cfg.addValue(current_section, key(option), rawValue(value))

			// Continuation of multi-line value
			case current_section != "" && option != "":
				prev := cfg.getRawString(current_section, key(option))
				value := strings.TrimSpace(stripComments(l))
				cfg.addValue(current_section, key(option), rawValue(prev+"\n"+value))

			default:
				return errors.New("could not parse line: " + l)
			}
		}
	}
	return nil
}

func stripComments(l string) string {
	for _, c := range []string{" " + string(COMMENT), "\t" + string(COMMENT)} {
		if i := strings.Index(l, c); i != -1 {
			l = l[0:i]
		}
	}
	return l
}

func (cfg *Config) addValue(s section, k key, v rawValue) {
	cfg.data[s][k] = v
}

func (cfg *Config) getRawString(s section, k key) string {
	rs, ok := cfg.data[s][k]
	if !ok {
		return ""
	}
	return string(rs)
}
