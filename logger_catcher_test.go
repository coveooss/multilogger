package multilogger

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger_Write(t *testing.T) {
	t.Parallel()
	type packet struct {
		sent, expect      string
		expectSize, delay int
	}
	tests := []struct {
		name       string
		afterClose string
		logs       string
		sequence   []packet
	}{
		{"Single", "Hello world\n", "", []packet{
			{"Hello world\n", "Hello world\n", 12, 0},
		}},
		{"Two batch", "Incomplete message\n", "", []packet{
			{"Incomplete ", "", 11, 0},
			{"message\n", "Incomplete message\n", 8, 0},
		}},
		{"Flush after close", "Not CR terminated", "", []packet{
			{"Not CR ", "", 7, 0},
			{"terminated", "", 10, 0},
		}},
		{"Test with error in the middle", "That's all",
			"[Test with error in the middle] 2018/06/24 12:34:56.789 ERROR    This is an error situation\n",
			[]packet{
				{"This is an [error] situation\n", "", 29, 0},
				{"That's all", "", 10, 0},
			}},
		{"Test with Error", "",
			"[Test with Error] 2018/06/24 12:34:56.789 ERROR    This is an error\n",
			[]packet{
				{"    [ERROR] This is an error", "", 28, 0},
			}},
		{"Test with Warn", "",
			"[Test with Warn] 2018/06/24 12:34:56.789 WARNING  This is an short warning\n",
			[]packet{
				{"    [WARN] This is an short warning", "", 35, 0},
			}},
		{"Test with cut warning", "",
			"[Test with cut warning] 2018/06/24 12:34:56.789 WARNING  This warning message is cut\n",
			[]packet{
				{"This [war", "", 9, 0},
				{"nin", "", 3, 0},
				{"g] message is cut", "", 17, 0},
			}},
		{"Test with info (including warning)", "",
			"[Test with info (including warning)] 2018/06/24 12:34:56.789 INFO     This is a [Warning] message\n",
			[]packet{
				{"[INFO] This is a [Warning] message\n", "", 35, 0},
			}},
		{"Test with info and debug (including warning)", "data ...\n",
			"[Test with info and debug (including warning)] 2018/06/24 12:34:56.789 INFO     This is an important information\n[Test with info and debug (including warning)] 2018/06/24 12:34:56.789 DEBUG    Useless trace\n",
			[]packet{
				{"[info] This is an important information", "", 39, 0},
				{"\ndata ...\n", "data ...\n", 10, 0},
				{"\t\t[debug]Useless trace", "data ...\n", 22, 0},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output, logs bytes.Buffer
			log := getTestLogger(tt.name, "Trace").SetOut(&logs).SetStdout(&output)

			for i, p := range tt.sequence {
				i++
				n, err := log.Write([]byte(p.sent))
				assert.NoError(t, err, "LogCatcher.Write() #%d error = %v", i, err)
				assert.Equal(t, p.expectSize, n, "Unexpected size")
				assert.Equal(t, p.expect, output.String())
				if i == len(tt.sequence) {
					log.Close()
				}
			}
			assert.Equal(t, tt.afterClose, output.String())
			assert.Equal(t, tt.logs, logs.String())
			assert.NoError(t, log.GetError())
		})
	}
}

func TestLogCatcher_WriteWithError(t *testing.T) {
	tests := []struct {
		name            string
		text            string
		err             error
		wantResultCount int
		wantErr         string
	}{
		{"No error", "", nil, 0, ""},
		{"Count error", "Write something\n", nil, 0, `ConsoleHook: Wrong number of bytes written (0) for "Write something\n"`},
		{"Several errors", "First line\n[debug] Hello\nNo output\n", nil, 0, `ConsoleHook: Wrong number of bytes written (0) for "First line\n"`},
		{"Disk full", "Won't be printed\n", fmt.Errorf("Disk is full"), 0, "ConsoleHook: Disk is full"},
		{"Disk full 2", "Error will be on close call", fmt.Errorf("Disk is full"), 27, "ConsoleHook: Disk is full"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := getTestLogger(tt.name, "trace").SetAllOutputs(&buggyWriter{tt.err})

			gotResultCount, err := log.Write([]byte(tt.text))
			if err == nil {
				err = log.Close()
			}
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantResultCount, gotResultCount)
		})
	}
}

type buggyWriter struct{ err error }

func (bw *buggyWriter) Write(buffer []byte) (int, error) {
	if bw.err != nil {
		return 0, bw.err
	}
	return 0, nil
}

func ExampleNew_log_catcher() {
	log := getTestLogger("catcher", logrus.TraceLevel)
	cmd := exec.Command("cat")
	stdin, _ := cmd.StdinPipe()
	cmd.Stdout, cmd.Stderr = log, log
	cmd.Start()

	lines := []string{
		"Hello,",
		"",
		"This is a text that contains:",
		"[Error] Oops! there is an error",
		"This should be considered as a [warning] message",
		"This should go directly to output",
		"Format tags like %hello in output aren't considered",
	}
	for _, line := range lines {
		io.WriteString(stdin, line+"\n")
	}
	stdin.Close()
	cmd.Wait()
	// Output:
	// Hello,
	//
	// This is a text that contains:
	// [catcher] 2018/06/24 12:34:56.789 ERROR    Oops! there is an error
	// [catcher] 2018/06/24 12:34:56.789 WARNING  This should be considered as a warning message
	// This should go directly to output
	// Format tags like %hello in output aren't considered
}

func ExampleNew_log_catcher_disabled() {
	log := getTestLogger("catcher", logrus.TraceLevel)

	// We disable the log catcher.
	log.Catcher = false
	fmt.Fprintln(log, "Hello,")
	fmt.Fprintln(log)
	fmt.Fprintln(log, "This is a text that contains:")
	fmt.Fprintln(log, "[Error] Oops! there is an error")
	fmt.Fprintln(log, "This should be considered as a [warning] message")
	fmt.Fprintln(log, "This should go directly to output")
	// Output:
	// Hello,
	//
	// This is a text that contains:
	// [Error] Oops! there is an error
	// This should be considered as a [warning] message
	// This should go directly to output
}

func ExampleNew_output_sent_to_info() {
	log := getTestLogger("catcher", logrus.TraceLevel)

	// We send all regular text to the InfoLevel.
	log.PrintLevel = logrus.InfoLevel
	fmt.Fprintln(log, "This is a text that contains:")
	fmt.Fprintln(log, "[Error] Oops! there is an error")
	fmt.Fprintln(log, "This should be considered as a [warning] message")
	fmt.Fprintln(log, "This should go directly to output")
	// Output:
	// [catcher] 2018/06/24 12:34:56.789 INFO     This is a text that contains:
	// [catcher] 2018/06/24 12:34:56.789 ERROR    Oops! there is an error
	// [catcher] 2018/06/24 12:34:56.789 WARNING  This should be considered as a warning message
	// [catcher] 2018/06/24 12:34:56.789 INFO     This should go directly to output
}
