package custom

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func Shell(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func BuildReleaseVersion(ver string) ReleaseVersion {
	verSplit := strings.Split(ver, ".")
	major, err := strconv.Atoi(strings.ReplaceAll(verSplit[0], "v", ""))
	CheckError(err)
	minor, err := strconv.Atoi(verSplit[1])
	CheckError(err)
	patch, err := strconv.Atoi(verSplit[2])
	CheckError(err)

	return ReleaseVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}
