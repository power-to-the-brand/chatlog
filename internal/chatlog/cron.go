package chatlog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const cronLine = "0 16 * * * %s decrypt && %s sync"

// SetupCron adds a crontab entry for daily decrypt+sync at 4pm (macOS).
// Deduplicates: skips if a matching entry already exists.
func SetupCron() (ok bool, msg string) {
	exe, err := os.Executable()
	if err != nil {
		return false, "failed to get executable path: " + err.Error()
	}
	if abs, err := filepath.Abs(exe); err == nil {
		exe = abs
	}

	newLine := fmt.Sprintf(cronLine, exe, exe)

	// Get current crontab
	out, err := exec.Command("crontab", "-l").Output()
	current := string(out)
	if err != nil {
		// Exit 1 = no crontab; that's ok
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			current = ""
		} else {
			return false, "failed to read crontab: " + err.Error()
		}
	}

	// Deduplicate: skip if we already have a line running our chatlog for decrypt+sync
	existingLines := strings.Split(strings.TrimSpace(current), "\n")
	for _, line := range existingLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, exe) && strings.Contains(line, "decrypt") && strings.Contains(line, "sync") {
			return true, "Daily sync (4pm) already configured"
		}
	}

	// Append our line
	var buf bytes.Buffer
	if current != "" {
		buf.WriteString(current)
		if !strings.HasSuffix(current, "\n") {
			buf.WriteByte('\n')
		}
	}
	buf.WriteString(newLine)
	buf.WriteByte('\n')

	cmd := exec.Command("crontab", "-")
	cmd.Stdin = &buf
	if out, err := cmd.CombinedOutput(); err != nil {
		return false, "failed to install crontab: " + err.Error() + "\n" + string(out)
	}
	return true, "Daily sync (4pm) configured"
}
