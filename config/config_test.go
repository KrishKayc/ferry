package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJiraFinder_CreateConfigEmpty(t *testing.T) {
	r := assert.New(t)

	err, _ := New("")
	r.Errorf(err, "expected createConfig to fail")
	r.Containsf(err.Error(), "empty config not allowed", "expected 'empty config not allowed', got '%s'", err)
}

func TestJiraFinder_CreateConfigNotExisting(t *testing.T) {
	r := assert.New(t)

	err, _ := New("/nofile/here")
	r.Errorf(err, "expected createConfig to fail")
	r.Containsf(err.Error(), "no such file", "expected 'Not Exists Error', got '%s'", err)
}

func TestJiraFinder_CreateConfigSuccess(t *testing.T) {
	r := assert.New(t)
	err, c := New("../example_config/sample_config_bug_search.json")

	r.NoErrorf(err, "expected reading config succeed, got error: '%s'", err)
	r.NotNil(c, "expected to have an healthy config, got nil")
	r.NotEmpty(c.AuthToken, ".AuthToken should not be empty")
}