package ocis

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"ociswrapper/common"
	"ociswrapper/ocis/config"
)

// var ocis = "/mnt/workspace/owncloud/ocis/ocis/bin/ocis"

var ocisCmd *exec.Cmd

func InitOcis() (string, string) {
	initCmd := exec.Command(config.Get("bin"), "init", "--insecure", "true")
	log.Print(initCmd.String())
	initCmd.Env = os.Environ()
	// [cleanup] not required
	initCmd.Env = append(initCmd.Env, "IDM_ADMIN_PASSWORD=admin")

	var out, err bytes.Buffer
	initCmd.Stdout = &out
	initCmd.Stderr = &err
	initCmd.Run()

	return out.String(), err.String()
}

func Start(envMap map[string]any) {
	defer common.Wg.Done()
	ocisCmd = exec.Command(config.Get("bin"), "server")
	ocisCmd.Env = os.Environ()
	var environments []string
	if envMap != nil {
		for key, value := range envMap {
			environments = append(environments, fmt.Sprintf("%s=%v", key, value))
		}
	}
	ocisCmd.Env = append(ocisCmd.Env, environments...)

	stderr, err := ocisCmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	stdout, err := ocisCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = ocisCmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		m := stderrScanner.Text()
		fmt.Println(m)
	}
	stdoutScanner := bufio.NewScanner(stdout)
	for stdoutScanner.Scan() {
		m := stdoutScanner.Text()
		fmt.Println(m)
	}
}

func Stop() {
	err := ocisCmd.Process.Kill()
	if err != nil {
		log.Panic("Cannot kill oCIS server")
	}
}

func WaitForConnection() bool {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	timeoutValue := 5 * time.Second
	client := http.Client{
		Timeout:   timeoutValue,
		Transport: transport,
	}

	timeout := time.After(timeoutValue)

	for {
		select {
		case <-timeout:
			fmt.Println(fmt.Sprintf("Timeout waiting for oCIS server [%f] seconds", timeoutValue.Seconds()))
			return false
		default:
			_, err := client.Get(config.Get("url"))
			if err != nil {
				fmt.Println("Waiting for oCIS server...")
			} else {
				fmt.Println(fmt.Sprintf("oCIS server is ready to accept requests"))
				return true
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func Restart(envMap map[string]any) bool {
	log.Print("Restarting oCIS server...")
	Stop()

	common.Wg.Add(1)
	go Start(envMap)

	return WaitForConnection()
}
