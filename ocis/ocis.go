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

func InitOcis() {
	initCmd := exec.Command(ocis, "init", "--insecure", "true")
	log.Print(initCmd.String())
	// setting env variables
	initCmd.Env = os.Environ()
	// [cleanup] not required
	initCmd.Env = append(initCmd.Env, "IDM_ADMIN_PASSWORD=admin")

	var out, err bytes.Buffer
	initCmd.Stdout = &out
	initCmd.Stderr = &err
	initCmd.Run()

	if err.String() != "" {
		log.Fatal(err.String())
	}
	fmt.Println(out.String())
}

func startOcis() {
	startCmd := exec.Command(ocis, "server")

	stderr, err := startCmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	stdout, err := startCmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = startCmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Pid: ", startCmd.Process.Pid)

	stderrScanner := bufio.NewScanner(stderr)
	for stderrScanner.Scan() {
		m := stderrScanner.Text()
		// ocisChannel <- m
		fmt.Println("[OCIS]: ", m)
	}
	stdoutScanner := bufio.NewScanner(stdout)
	for stdoutScanner.Scan() {
		m := stdoutScanner.Text()
		// ocisChannel <- m
		fmt.Println("[OCIS]: ", m)
	}
	// get pids of ocis
	// output, _ := exec.Command("ps", "-o", "pid", "-C", "ocis").Output()
	// fields := strings.Fields(string(output))
	// for i := 1; i < len(fields); i++ {
	// 	fmt.Printf("PID of ocis: %s\n", fields[i])
	// }
	// kill after timeout
	// <-time.After(5 * time.Second)
	// err := startCmd.Process.Kill()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Process killed with PID:", startCmd.Process.Pid)
}
func main() {
	var wg sync.WaitGroup
	// initOcis()
	wg.Add(1)
	go func() {
		defer wg.Done()
		startOcis()
	}()
	fmt.Println("Ready to listen requests...")
	wg.Wait()
}
