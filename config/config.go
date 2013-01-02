// Package config stores global system settings, persisted in a file.

package config

import (
	"fmt"
	"github.com/kless/goconfig/config"
	"os"
	"strings"
)

var fileName string
var c *config.Config
var defaults = make(map[string]map[string]optionDefault)

type optionDefault struct {
	value       interface{}
	description string
}

// SetFile parses a config file.  If the filename that is provided does not
// exist, the file is created and the default sections and options are written
// as comments.  This function should only be called after all options are
// registered using the Register Function, so take care putting it in init
// functions.
func SetFile(fileName string) {
	var err error
	c, err = config.ReadDefault(fileName)
	if err != nil {
		WriteDefaultFile(fileName)
		c, err = config.ReadDefault(fileName)
		if err != nil {
			panic(err)
		}
	}
}

// WriteDefaultFile writes the registerd defaults, as comments, to the given file.
// This function is called automatically by SetFile if the file does not exist,
// but can also be called manually if appropriate
func WriteDefaultFile(fileName string) error {
	os.Remove(fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("# This is the config file for the data system, values below should be s\n\n")
	if _, ok := defaults["DEFAULT"]; ok {
		writeSection(file, "DEFAULT")
	}

	for section, _ := range defaults {
		if section != "DEFAULT" {
			writeSection(file, section)
		}
	}

	return nil
}

// Helper function to write a section's defaults
func writeSection(file *os.File, section string) {
	file.WriteString(fmt.Sprintf("[%v]\n", section))
	for name, value := range defaults[section] {
		descParts := strings.Split(value.description, "\n")
		for _, descPart := range descParts {
			trimmed := strings.Trim(descPart, "\r\n\t ")
			if trimmed != "" {
				file.WriteString("# ")
				file.WriteString(trimmed)
				file.WriteString("\n")
			}
		}
		file.WriteString(fmt.Sprintf("# %v = %v\n\n", name, value.value))
	}
}

// Register registers a config option and it's default.  Default values are
// stored in the package, but can be overridden by changing the config file.
// Config items should be registered in a package's init function.  It is
// assumed that all config values have been registered before the SetFile
// function is called.
func Register(section, option, defaultValue, description string) error {
	if c != nil {
		return fmt.Errorf("Config file has already been parsed.  Register config options in init function")
	}
	if _, ok := defaults[section]; !ok {
		defaults[section] = make(map[string]optionDefault)
	}
	defaults[section][option] = optionDefault{defaultValue, description}
	return nil
}

// Get returns a config option as a string, from the config file or the
// registered default value.  Errors are returned if the config option
// has not been registered or if the registered default is not a string.
func Get(section, option string) (string, error) {
	if c.HasOption(section, option) {
		return c.String(section, option)
	} else {
		if _, ok := defaults[section]; !ok {
			fmt.Println(defaults)
			return "", fmt.Errorf("Config option [%v]:%v does not exist", section, option)
		}
		if _, ok := defaults[section][option]; !ok {
			return "", fmt.Errorf("Config Option [%v]:%v does not exist", section, option)
		}
		if s, ok := defaults[section][option].value.(string); ok {
			return s, nil
		} else {
			return "", fmt.Errorf("Default config option [%v]:%v is not a string", section, option)
		}
	}
	return "", nil
}

// MustGet is a wrapper around get that panics if a value is not available
func MustGet(section, option string) string {
	s, err := Get(section, option)
	if err != nil {
		panic(err)
	}
	return s
}
