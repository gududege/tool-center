package process

import (
	"os/exec"
	"testing"
	"time"

	domainProcess "github.com/cli-tool-center/tool-center/internal/domain/process"
)

func TestRunner_StartAndWait(t *testing.T) {
	if _, err := exec.LookPath("cmd.exe"); err != nil {
		t.Skip("cmd.exe not found, skipping Windows-specific test")
	}

	runner := NewRunner()
	output := make([]string, 0)

	runner.SetOutputHandler(func(taskID, line, source string) {
		output = append(output, line)
	})

	err := runner.Start("test-1", domainProcess.StartRequest{
		Executable: "cmd.exe",
		Arguments:  []string{"/C", "echo hello world"},
	})
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	exitCode, err := runner.Wait("test-1")
	if err != nil {
		t.Fatalf("Wait failed: %v", err)
	}
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}

	// 给输出处理一点时间
	time.Sleep(100 * time.Millisecond)

	if len(output) == 0 {
		t.Fatal("expected output from echo command")
	}
	if output[0] != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", output[0])
	}

	runner.Remove("test-1")
	if runner.IsRunning("test-1") {
		t.Error("expected process to be removed")
	}
}

func TestRunner_Stop(t *testing.T) {
	if _, err := exec.LookPath("cmd.exe"); err != nil {
		t.Skip("cmd.exe not found")
	}

	runner := NewRunner()
	done := make(chan bool)

	err := runner.Start("test-stop", domainProcess.StartRequest{
		Executable: "cmd.exe",
		Arguments:  []string{"/C", "ping -n 10 127.0.0.1 > nul"},
	})
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// 在 goroutine 中等待，因为我们将会取消
	go func() {
		runner.Wait("test-stop")
		done <- true
	}()

	// 取消进程
	if err := runner.Stop("test-stop"); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	select {
	case <-done:
		// 成功停止
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for process to stop")
	}

	runner.Remove("test-stop")
}

func TestRunner_StopNonExistent(t *testing.T) {
	runner := NewRunner()
	err := runner.Stop("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent process")
	}
}

func TestRunner_WaitNonExistent(t *testing.T) {
	runner := NewRunner()
	_, err := runner.Wait("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent process")
	}
}

func TestRunner_IsRunning(t *testing.T) {
	if _, err := exec.LookPath("cmd.exe"); err != nil {
		t.Skip("cmd.exe not found")
	}

	runner := NewRunner()

	if runner.IsRunning("no-such-process") {
		t.Error("expected false for non-existent process")
	}

	runner.Start("test-run", domainProcess.StartRequest{
		Executable: "cmd.exe",
		Arguments:  []string{"/C", "echo running > nul"},
	})

	if !runner.IsRunning("test-run") {
		t.Error("expected true for running process after Start")
	}

	runner.Wait("test-run")
	runner.Remove("test-run")
}

func TestRunner_ExitCode(t *testing.T) {
	if _, err := exec.LookPath("cmd.exe"); err != nil {
		t.Skip("cmd.exe not found")
	}

	runner := NewRunner()

	// 执行一个会失败的命令
	err := runner.Start("test-fail", domainProcess.StartRequest{
		Executable: "cmd.exe",
		Arguments:  []string{"/C", "exit 42"},
	})
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	exitCode, waitErr := runner.Wait("test-fail")
	if waitErr != nil {
		// Wait returns the error for non-zero exit
		// The exit code should still be in exitCode
		if exitCode != 42 {
			t.Errorf("expected exit code 42, got %d", exitCode)
		}
	} else {
		if exitCode != 42 {
			t.Errorf("expected exit code 42, got %d", exitCode)
		}
	}

	runner.Remove("test-fail")
}

func TestRunner_EmptyExecutable(t *testing.T) {
	runner := NewRunner()
	err := runner.Start("test-empty", domainProcess.StartRequest{
		Executable: "",
		Arguments:  []string{},
	})
	if err == nil {
		t.Error("expected error for empty executable")
	}
}

func TestRunner_MultipleProcesses(t *testing.T) {
	if _, err := exec.LookPath("cmd.exe"); err != nil {
		t.Skip("cmd.exe not found")
	}

	runner := NewRunner()
	outputs := make(map[string][]string)

	runner.SetOutputHandler(func(taskID, line, source string) {
		outputs[taskID] = append(outputs[taskID], line)
	})

	// 启动多个进程
	runner.Start("p1", domainProcess.StartRequest{
		Executable: "cmd.exe",
		Arguments:  []string{"/C", "echo process one"},
	})
	runner.Start("p2", domainProcess.StartRequest{
		Executable: "cmd.exe",
		Arguments:  []string{"/C", "echo process two"},
	})

	runner.Wait("p1")
	runner.Wait("p2")

	time.Sleep(100 * time.Millisecond)

	if !runner.IsRunning("p1") && !runner.IsRunning("p2") {
		// both done
	}

	runner.Remove("p1")
	runner.Remove("p2")
}
