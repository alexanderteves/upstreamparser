package upstreamparser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type UpstreamConfig struct {
	Name  string
	Lines [][]string
}

func splitLineIntoElements(configLineString string) (configLine []string) {
	for _, configElement := range strings.Split(configLineString, " ") {
		if configElement != "" {
			configLine = append(configLine, configElement)
		}
	}
	return
}

func getConfigName(fileContent string) (configName string, err error) {
	re := regexp.MustCompile(`upstream ([a-zA-Z0-9\-]+)`)
	match := re.FindStringSubmatch(string(fileContent))
	if match != nil {
		return match[1], nil
	} else {
		return "", errors.New("Name could not be parsed")
	}
}

func getConfigLines(fileContent string) (configLines [][]string, err error) {
	re := regexp.MustCompile(`upstream[ a-zA-Z0-9\-\n]+{(.*)}`)
	match := re.FindStringSubmatch(string(fileContent))
	if match != nil {
		for _, configLineString := range strings.Split(match[1], ";") {
			if configLineString != "" {
				configLines = append(configLines, splitLineIntoElements(configLineString))
			}
		}
		return configLines, nil
	} else {
		return nil, errors.New("Configuration could not be parsed")
	}
}

func Load(filename string) (upstreamConfig UpstreamConfig, err error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return UpstreamConfig{}, err
	}
	upstreamConfig, err = Loads(string(fileContent))
	if err != nil {
		return UpstreamConfig{}, err
	} else {
		return upstreamConfig, nil
	}
}

func Loads(fileContent string) (upstreamConfig UpstreamConfig, err error) {
	fileContent = strings.Replace(fileContent, "\n", "", -1)
	configName, err := getConfigName(fileContent)
	if err != nil {
		return UpstreamConfig{}, err
	} else {
		upstreamConfig.Name = configName
	}
	configLines, err := getConfigLines(fileContent)
	if err != nil {
		return UpstreamConfig{}, err
	} else {
		upstreamConfig.Lines = configLines
	}
	return
}

func Dump(upstreamConfig UpstreamConfig, filename string) (err error) {
	fileContent, err := Dumps(upstreamConfig)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, []byte(fileContent), 0644); err != nil {
		return err
	}
	return nil
}

func Dumps(upstreamConfig UpstreamConfig) (fileContent string, err error) {
	var configLineString string
	for _, configLine := range upstreamConfig.Lines {
		configLineString = "   "
		for _, configElement := range configLine {
			configLineString += " " + configElement
		}
		fileContent += configLineString + ";\n"
	}
	fileContent = fmt.Sprintf("upstream %s {\n%s}", upstreamConfig.Name, fileContent)
	return fileContent, nil
}
