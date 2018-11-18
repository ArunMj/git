package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func debug() {
	f, e := os.OpenFile("/tmp/go-get-shallow-git.debug", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()
	if e != nil {
		panic(e)
	}
	fmt.Fprintf(f, "$$$$ Running gitwrapper.  IsGoGet: %t, Args: %s\n", isThisGoGet(), os.Args)
	// fmt.Println(pid)
	c := exec.Command("pstree", "-p", "--show-parent", strconv.Itoa(os.Getpid()))
	c.Stdout = f
	c.Stderr = f
	if err := c.Run(); err != nil {
		panic(err)
	}
}

func isThisGoGet() bool {
	ppid := os.Getppid()
	cmdlineBytes, err := ioutil.ReadFile(`/proc/` + strconv.Itoa(ppid) + `/cmdline`)
	if err != nil {
		panic(err)
	}
	cmdlineArgs := strings.FieldsFunc(string(cmdlineBytes), func(r rune) bool { return r == '\u0000' })
	return len(cmdlineArgs) > 1 && cmdlineArgs[0] == "go"
}

func main() {
	os.Stderr.Write([]byte("WARN: Using patched git wrapper `" + os.Args[0] + "` for shallow 'go get'\n"))

	args := os.Args

	if isThisGoGet() {
		if len(args) > 1 {
			if args[1] == "pull" {
				args = append(args[:2], append([]string{"--depth=1"}, args[2:]...)...)
			} else if args[1] == "clone" {
				args = append(args[:2], append([]string{"--depth=1", "--shallow-submodules", "--single-branch"}, args[2:]...)...)
			}
		}
	}

	args[0] = "/usr/bin/git"
	cmd := exec.Command(args[0])
	if len(args) > 0 {
		cmd = exec.Command(args[0], args[1:]...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

}
