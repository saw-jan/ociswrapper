package ocis

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

var ocis = "/mnt/workspace/owncloud/ocis/ocis/bin/ocis"

var ocisCmd *exec.Cmd

func InitOcis() (string, string) {
	initCmd := exec.Command(ocis, "init", "--insecure", "true")
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

func StartOcis(wg *sync.WaitGroup, envMap map[string]any) {
	defer wg.Done()
	ocisCmd = exec.Command(ocis, "server")
	ocisCmd.Env = os.Environ()
	var environments []string
	if envMap != nil {
		for key, value := range(envMap){
			environments = append(environments, fmt.Sprintf("%s=%s", key, value))
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

func stopOcis(){
	err := ocisCmd.Process.Kill()
	if err != nil {
		log.Panic("Cannot kill ocis server")
	}
}

func RestartOcisServer(wg *sync.WaitGroup, envMap map[string]any){
	log.Print("Restarting ocis server...")
	stopOcis()
	wg.Add(1)
	go StartOcis(wg, envMap)
	// Todo: wait for ocis to start
}