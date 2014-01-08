package goini

import (
	"testing"
)

func TestPlainConfig(t *testing.T) {
	c, err := LoadConfig("plain_test.ini")
	if err != nil {
		t.Errorf("Got an unexpected error: %v", err)
	}

	sections := c.GetSectionList()
	s := map[string]bool{
		"production": true,
		"staging":    true,
		"testing":    true,
	}
	for _, x := range sections {
		_, ok := s[x]
		if !ok {
			t.Errorf("Section not found: %s", x)
		}
	}
	v := c.GetSection("testing")
	if v.Len() == 0 {
		t.Error("No keys inside testing section")
	}

	if val, _ := v.GetString("resources.mq.host"); val != "some.super.example.com" {
		t.Error("GetString method failed")
	}
	if val, _ := v.GetString("string"); val != "foo-bar-baz" {
		t.Error("GetString method failed on " + "string")
	}
	if val, _ := v.GetString("single_string"); val != "foo-bar-baz" {
		t.Error("GetString method failed on " + "single_string")
	}
	if val, _ := v.GetString("double_string"); val != "foo-bar-baz" {
		t.Error("GetString method failed on " + "double_string")
	}
	if val, _ := v.GetString("multy_string"); val != "foo-\nbar-\nbaz" {
		t.Error("GetString method failed on " + "multy_string")
	}
	if val, _ := v.GetString("last_key"); val != "last_key_val" {
		t.Error("GetString method failed on " + "last_key")
	}
	if val, _ := v.GetInt("resources.mq.port"); val != 5672 {
		t.Error("GetInt method failed")
	}
	if val, _ := v.GetBool("resources.mq.vhost"); val {
		t.Error("GetBool does not return error")
	}
}

func TestInheritanceConfig(t *testing.T) {
	c, err := LoadConfig("inheritance_test.ini")
	if err != nil {
		t.Errorf("Got an unexpected error: %v", err)
	}

	sections := c.GetSectionList()
	s := map[string]bool{
		"production": true,
		"staging":    true,
		"testing":    true,
	}
	for _, x := range sections {
		_, ok := s[x]
		if !ok {
			t.Errorf("Section not found: %s", x)
		}
	}
	v := c.GetSection("testing")
	if v.Len() == 0 {
		t.Error("No keys inside testing section")
	}
	if val, _ := v.GetString("resources.mq.host"); val != "staging-overrided" {
		t.Error("GetString method failed")
	}
	if val, _ := v.GetString("prod_option"); val != "staging" {
		t.Error("Wrong overrided option")
	}
	if val, _ := v.GetString("resources.mq.password"); val != "blah-blah" {
		t.Error("Wrong value for overrided option")
	}
	if val, _ := v.GetString("inherited_option"); val != "foo-bar-baz" {
		t.Error("Can not get parent option")
	}

	if _, err := v.GetString("not_found_option"); err.Error() != ERR_KEY_NOT_EXISTS {
		t.Error("Can not get parent option")
	}

}
