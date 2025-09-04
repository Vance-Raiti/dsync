package main

import (
	"fmt"
	"strings"
	"os"
	"os/user"
	"os/exec"
	"path/filepath"
	"encoding/json"
)

type Record struct {
	User string
	Url string
}

type DsyncConfig map[string]Record

func perror(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", msg)
	os.Exit(1)
}

func pinfo(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("INFO: %v\n", msg)
}

func check(err error, format string, args ...any) {
	if err != nil {
		msg := fmt.Sprintf(format, args...)
		perror("%s: %v\n", msg, err)
	}
}

func getConfigDir(username string, initializeIfNone bool) string {
	var err error
	var u *user.User

	if username == "" {
		u, err = user.Current()
		check(err, "could not get current user")
	} else {
		u, err = user.Lookup(username)
		check(err, "could not lookup user %v", username)
	}

	configDir := fmt.Sprintf("%s/.config/dsync", u.HomeDir)

	if initializeIfNone {
		_, err = os.Stat(configDir)
		if os.IsNotExist(err) {
			os.MkdirAll(configDir, 0775)
		} else {
			check(err, "could not stat %s", configDir)
		}
	}

	_, err = os.Stat(configDir)
	check(err, "could not stat %s", configDir)

	return configDir
}

func getConfig(username string, initalizeIfNone bool) DsyncConfig {
	var err error

	configPath := fmt.Sprintf("%s/config.json", getConfigDir(username, initalizeIfNone))

	if initalizeIfNone {
		_, err = os.Stat(configPath)
		if os.IsNotExist(err) {
			f, err := os.Create(configPath)
			defer f.Close()
			check(err, "could not create config")

			_, err = f.WriteString("{}")
			check(err, "could not write to empty config")
		} else { 
			check(err, "could not stat config")
		}
	}

	var config DsyncConfig

	f, err := os.Open(configPath)
	check(err, "could not open config")
	defer f.Close()

	err = json.NewDecoder(f).Decode(&config)
	check(err, "could not decode config json")

	return config
}

func setConfig(config DsyncConfig) {
	var err error

	configPath := fmt.Sprintf("%s/config.json", getConfigDir("", false))

	file, err := os.Create(configPath)
	check(err, "could not open %v", configPath)
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	check(err, "could not marshal config")
}

func getAuthPath(username string) string {
	return fmt.Sprintf("%s/auth.json", getConfigDir(username, false))
}

func add(args []string) {
	var err error

	path, err := filepath.Abs(args[0])
	check(err, "could not interpret %v", args[0])

	url := args[1]

	config := getConfig("", true)
	authPath := getAuthPath("")	

	parsedUrl := strings.Split(url,":")
	cmd := exec.Command("podman", "login", fmt.Sprintf("--authfile=%s",authPath), parsedUrl[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	check(err, "error running `podman login %v`", url)

	u, err := user.Current()
	check(err, "could not get current user")

	config[path] = Record{Url: url, User: u.Username}
	setConfig(config)
}

func sync(args []string) {
	perror("nyi")
}

func remove(args []string) {
	var err error

	path, err := filepath.Abs(args[0])
	check(err, "could not interpret %v", args[0])

	config := getConfig("", true)

	if _, ok := config[path]; !ok {
		perror("path %v is not registered with dsync", path)
	}

	delete(config, path)
}

func main() {
	if len(os.Args) < 2 {
		perror("Usage: %v [command]", os.Args[0])
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	if cmd == "add" {
		add(args)
	} else if cmd == "sync" {
		sync(args)
	} else if cmd == "remove" {
		remove(args)
	} else {
		perror("%v is not a recognized command\n", os.Args[0])
	}
}
