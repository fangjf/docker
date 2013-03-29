package docker

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestCmdStreamLargeStderr(t *testing.T) {
	// This test checks for deadlock; thus, the main failure mode of this test is deadlocking.
	// FIXME implement a timeout to avoid blocking the whole test suite when this test fails
	cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero bs=1k count=1000 of=/dev/stderr; echo hello")
	out, err := CmdStream(cmd)
	if err != nil {
		t.Fatalf("Failed to start command: " + err.Error())
	}
	_, err = io.Copy(ioutil.Discard, out)
	if err != nil {
		t.Fatalf("Command should not have failed (err=%s...)", err.Error()[:100])
	}
}

func TestCmdStreamBad(t *testing.T) {
	badCmd := exec.Command("/bin/sh", "-c", "echo hello; echo >&2 error couldn\\'t reverse the phase pulser; exit 1")
	out, err := CmdStream(badCmd)
	if err != nil {
		t.Fatalf("Failed to start command: " + err.Error())
	}
	if output, err := ioutil.ReadAll(out); err == nil {
		t.Fatalf("Command should have failed")
	} else if err.Error() != "exit status 1: error couldn't reverse the phase pulser\n" {
		t.Fatalf("Wrong error value (%s)", err.Error())
	} else if s := string(output); s != "hello\n" {
		t.Fatalf("Command output should be '%s', not '%s'", "hello\\n", output)
	}
}

func TestCmdStreamGood(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "echo hello; exit 0")
	out, err := CmdStream(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if output, err := ioutil.ReadAll(out); err != nil {
		t.Fatalf("Command should not have failed (err=%s)", err)
	} else if s := string(output); s != "hello\n" {
		t.Fatalf("Command output should be '%s', not '%s'", "hello\\n", output)
	}
}

func TestTarUntar(t *testing.T) {
	archive, err := Tar(".", Uncompressed)
	if err != nil {
		t.Fatal(err)
	}
	tmp, err := ioutil.TempDir("", "docker-test-untar")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	if err := Untar(archive, tmp); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(tmp); err != nil {
		t.Fatalf("Error stating %s: %s", tmp, err.Error())
	}
}