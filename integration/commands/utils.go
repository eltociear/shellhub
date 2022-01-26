package commands

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	Cloud = iota
	Enterprise
	Community
)

func MapEnvs(instance int) []string {
	switch instance {
	case Cloud:
		return []string{"SHELLHUB_CLOUD=true", "SHELLHUB_BILLING=true"}
	case Enterprise:
		return []string{"SHELLHUB_CLOUD=false", "SHELLHUB_ENTERPRISE=true"}
	case Community:
		return []string{"SHELLHUB_CLOUD=false", "SHELLHUB_ENTERPRISE=false"}
	default:
		return []string{}
	}
}

func CheckAPIInstance(envLog string, instance int, spec []string) bool {
	if !strings.Contains(envLog, "MONGO_DB_NAME=test") {
		return false
	}

	envsInstance := MapEnvs(instance)

	for _, env := range append(envsInstance, spec...) {
		if !strings.Contains(envLog, env) {
			return false
		}
	}

	return true
}

func IsInstanceAlive(service string, version int, extra []string) bool {
	envs, err := ExecuteCommand("../bin/docker-compose",
		"exec",
		service,
		"env",
	)

	if err != nil {
		log.Println("error: ", err)
		return false
	}

	return CheckAPIInstance(envs, version, extra)
}

func ExecuteCommand(command string, args ...string) (string, error) {
	log.Println(fmt.Sprintf("the command %s is going to be executed", command))

	log.Println("The args are ", args)
	cmd := exec.Command(command, args...)

	bytes, err := cmd.Output()
	if err != nil { // the api container is not running
		log.Println("error: ", err)
		return "", err
	}

	log.Println("success: ", string(bytes))
	return string(bytes), err
}
