package go_utils

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

// 执行Shell脚本，返回行解析对象数组
func ExecuteBashLiner(shellScripts string, liner func(line string) bool) error {
	cmd := exec.Command("bash", "-c", shellScripts)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Start()
	defer cmd.Process.Kill()
	defer cmd.Wait()

	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line != "" {
			if !liner(line) {
				return nil
			}
		}
	}

	return nil
}

func ExecuteBash(shellScripts string) (string, error) {
	stdout := ""

	err := ExecuteBashLiner(shellScripts, func(line string) bool {
		stdout += line
		return true
	})

	return stdout, err
}
