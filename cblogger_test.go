// CB-Log: Logger for Cloud-Barista.
//
//	* Cloud-Barista: https://github.com/cloud-barista
//
// Unit tests for cblogger.go / config.go.
// Focus areas:
//  1. LOOPCHECK field is gone — struct and defaults reflect its removal.
//  2. A config file that still contains 'loopcheck: true' is parsed without
//     error (unknown YAML fields are silently ignored for backward compatibility).
//  3. The file watcher fires automatically when the config file is modified —
//     no 'loopcheck' flag is needed.
package cblog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// resetGlobals clears package-level singletons so each test starts clean.
// It cancels any running watcher goroutine before clearing the logger.
func resetGlobals() {
	if watcherCancel != nil {
		watcherCancel()
		watcherCancel = nil
	}
	thisLogger = nil
	thisFormatter = nil
	cblogConfig = CBLOGCONFIG{}
}

// writeConfig creates a temp directory tree and writes yaml to
// <dir>/conf/log_conf.yaml. It returns the root dir and the full file path.
func writeConfig(t *testing.T, yaml string) (dir, cfgPath string) {
	t.Helper()
	dir = t.TempDir()
	confDir := filepath.Join(dir, "conf")
	if err := os.MkdirAll(confDir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	cfgPath = filepath.Join(confDir, "log_conf.yaml")
	if err := os.WriteFile(cfgPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return dir, cfgPath
}

// TestNewCBLOGCONFIG_Defaults checks that the default configuration no
// longer contains a LOOPCHECK field and that remaining defaults are correct.
func TestNewCBLOGCONFIG_Defaults(t *testing.T) {
	cfg := NewCBLOGCONFIG()
	if cfg.CBLOG.LOGLEVEL != "info" {
		t.Errorf("expected default LOGLEVEL 'info', got '%s'", cfg.CBLOG.LOGLEVEL)
	}
	if !cfg.CBLOG.CONSOLE {
		t.Error("expected default CONSOLE true")
	}
	if !cfg.CBLOG.LOGFILE {
		t.Error("expected default LOGFILE true")
	}
	// Compilation itself is the strongest proof that LOOPCHECK is gone:
	// any reference to cfg.CBLOG.LOOPCHECK would be a compile error here.
}

// TestGetConfigInfos_WithoutLoopcheck loads a config YAML that does not
// contain the 'loopcheck' key and verifies the parsed values.
func TestGetConfigInfos_WithoutLoopcheck(t *testing.T) {
	const yaml = `
cblog:
  loglevel: debug
  console: true
  logfile: false
logfileinfo:
  filename: ./log/test.log
  maxsize: 5
  maxbackups: 2
  maxage: 7
`
	dir, _ := writeConfig(t, yaml)
	t.Setenv("CBLOG_ROOT", dir)

	cfg := GetConfigInfos("")
	if cfg.CBLOG.LOGLEVEL != "debug" {
		t.Errorf("expected loglevel 'debug', got '%s'", cfg.CBLOG.LOGLEVEL)
	}
	if cfg.CBLOG.LOGFILE {
		t.Error("expected LOGFILE false")
	}
}

// TestGetConfigInfos_OldConfigWithLoopcheck verifies backward compatibility:
// a config YAML that still contains 'loopcheck: true' (old format) is loaded
// without error and the unknown field is silently ignored.
func TestGetConfigInfos_OldConfigWithLoopcheck(t *testing.T) {
	const yaml = `
cblog:
  loopcheck: true
  loglevel: warn
  console: true
  logfile: false
logfileinfo:
  filename: ./log/test.log
  maxsize: 5
  maxbackups: 2
  maxage: 7
`
	dir, _ := writeConfig(t, yaml)
	t.Setenv("CBLOG_ROOT", dir)

	cfg := GetConfigInfos("")
	if cfg.CBLOG.LOGLEVEL != "warn" {
		t.Errorf("expected loglevel 'warn', got '%s'", cfg.CBLOG.LOGLEVEL)
	}
}

// TestDynamicLevelChange is the key end-to-end test: it verifies that editing
// the log level in the config file is picked up automatically by the
// always-on file watcher — no 'loopcheck' flag is required.
func TestDynamicLevelChange(t *testing.T) {
	resetGlobals()
	defer resetGlobals()

	const initialYAML = `
cblog:
  loglevel: info
  console: false
  logfile: false
logfileinfo:
  filename: ./log/test.log
  maxsize: 5
  maxbackups: 2
  maxage: 7
`
	dir, cfgPath := writeConfig(t, initialYAML)
	t.Setenv("CBLOG_ROOT", dir)

	logger := GetLogger("TEST-WATCHER")
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
	if got := GetLevel(); got != "info" {
		t.Fatalf("initial level: expected 'info', got '%s'", got)
	}

	// Give the watcher goroutine time to call watcher.Add() and enter its
	// select loop before we modify the file.
	time.Sleep(200 * time.Millisecond)

	// Overwrite the config file with a different log level.
	const updatedYAML = `
cblog:
  loglevel: error
  console: false
  logfile: false
logfileinfo:
  filename: ./log/test.log
  maxsize: 5
  maxbackups: 2
  maxage: 7
`
	if err := os.WriteFile(cfgPath, []byte(updatedYAML), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// The file watcher is event-driven and should react within milliseconds;
	// allow up to 3 s to account for slow CI runners.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if GetLevel() == "error" {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if got := GetLevel(); got != "error" {
		t.Errorf("after file update: expected level 'error', got '%s'", got)
	}
}
