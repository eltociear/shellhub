package commands

func MongoIP() string {
	containerName := "shellhub-mongo-1"
	cmd := "docker"
	args := []string{"inspect", "--format", "'{{ .NetworkSettings.Networks.shellhub_network.IPAddress }}'", containerName}
	out, err := ExecuteCommand(cmd, args...)
	if err != nil {
		return ""
	}

	return out[1 : len(out)-2]
}
