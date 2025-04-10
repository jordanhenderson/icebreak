package icebreak

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func InitLambda() {
	wrapperEnv := os.Getenv("AWS_LAMBDA_EXEC_WRAPPER")
	if wrapperEnv == "" {
		return
	}

	selfPath, err := os.Executable()
	if err != nil {
		log.Printf("[bootstrap] Could not determine self path: %v", err)
		return
	}
	selfPath, _ = filepath.EvalSymlinks(selfPath)

	wrappers := strings.Split(wrapperEnv, ",")
	for _, wrapper := range wrappers {
		wrapper = strings.TrimSpace(wrapper)
		if wrapper == "" {
			continue
		}

		resolvedPath, err := filepath.EvalSymlinks(wrapper)
		if err != nil {
			log.Printf("[bootstrap] Failed to resolve %s: %v", wrapper, err)
			continue
		}

		if resolvedPath == selfPath {
			log.Printf("[bootstrap] Skipping self invocation: %s", wrapper)
			continue
		}

		log.Printf("[bootstrap] Running: %s", resolvedPath)

		cmd := exec.Command(resolvedPath)
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("failed to attach stdout pipe: %v", err)
		}
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		if err := cmd.Start(); err != nil {
			log.Fatalf("[bootstrap] Wrapper %s failed to start: %v", resolvedPath, err)
		}

		log.Printf("[bootstrap] Wrapper %s started in background (PID %d)", resolvedPath, cmd.Process.Pid)

		// Block until stdout contains "READY"
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("[stdout] %s", line)
			if strings.Contains(line, "READY") {
				log.Println("[bootstrap] child process signaled readiness")
				break
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading iceflake stdout: %v", err)
		}
	}
}
