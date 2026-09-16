package term

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// DisplayWithPager outputs the rendered text. If the text is longer than the
// terminal screen and stdout is an interactive terminal, it automatically invokes `less -R`.
func DisplayWithPager(content string, disablePager bool) error {
	isTerminal := isStdoutTerminal()
	_, termHeight := GetTerminalSize()
	lineCount := strings.Count(content, "\n") + 1

	// Conditions where pager is bypassed:
	// 1. User specified --no-pager
	// 2. Stdout is not an interactive terminal (piped to file/cmd)
	// 3. Output fits comfortably within current terminal screen
	if disablePager || !isTerminal || lineCount <= termHeight {
		fmt.Print(content)
		return nil
	}

	// Try using less -R
	pagerPath, err := exec.LookPath("less")
	if err != nil {
		fmt.Print(content)
		return nil
	}

	cmd := exec.Command(pagerPath, "-R", "-F", "-X")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Print(content)
		return nil
	}

	if err := cmd.Start(); err != nil {
		fmt.Print(content)
		return nil
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, content)
	}()

	return cmd.Wait()
}

// isStdoutTerminal checks if standard output is connected to a character device.
func isStdoutTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
