package commands

import (
	"log"
	"time"

    "github.com/shellhub-io/shellhub/integration/utils"
)

func instanceArgs(svc string, instance int, extra []string) []string {
	return append(func() []string {
		var args []string

		envs := MapEnvs(instance)
		if len(envs) == 0 {
			return args
		}

		for _, env := range envs {
			args = append(args, []string{"--build-arg", env[9:]}...)
		}

		return args
	}(), append(extra, svc)...)
}

func ApiHasInitialized() bool {
    log.Println("Api initialization check")
    service := "api"
    commandArgs := []string{
        "logs",
        service,
    }

    message, err := ExecuteCommand("docker-compose", commandArgs...)
    if err != nil {
        log.Println(err)
        return false
    }

    log.Println("checking if the API has begun")
    log.Println(message)

    return utils.ServerRunning(message)
}

func InitializeTestAPI(version int, extra []string) {
	service := "api"

	log.Println("building...")

	buildArgs := []string{
		"-f", "../docker-compose.test.yml",
		"-f", "../docker-compose.yml",
		"build",
		"--build-arg", "DB_NAME=test",
		"--build-arg", "GOPROXY=http://localhost:3333",
		"--build-arg", "NPM_CONFIG_REGISTRY=http://localhost:4873",
	}

	extraArgs := instanceArgs(service, version, extra)

	buildArgs = append(buildArgs, extraArgs...)

	_, err := ExecuteCommand("docker-compose", buildArgs...)

	if err != nil {
		log.Println(err)
		log.Println("couldnt build the project")

		return
	}

	log.Println("the project has been built succesfully")

	_, err = ExecuteCommand("docker-compose", []string{
		"-f", "../docker-compose.test.yml",
		"-f", "../docker-compose.yml",
		"up", "-d",
		service,
	}...)

	if err != nil {
		log.Println(err)
		log.Println("couldnt up the proect")
		return
	}

	log.Println("the project is going up")

	for { // wait the container to be up and running
		if ApiHasInitialized() {
            time.Sleep(2 * time.Second)
			break
		}

		time.Sleep(3 * time.Second)
		log.Println("pinging...")
	}

	log.Println("Initialization completed")
}
