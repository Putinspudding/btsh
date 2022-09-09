package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/fatih/color"
)

func main() {
	init_loop()
}

func init_loop() {
	/*
		for {
			cmd := readline()
			fmt.Println(cmd)
		}
	*/
	for {
		u, _ := user.Current()
		uname := u.Username
		hs, _ := os.Hostname()
		g := color.New(color.FgGreen, color.Bold)
		g.Printf("%s@%s", uname, hs)

		wd, _ := os.Getwd()
		wd = strings.Replace(wd, os.Getenv("HOME"), "~", 1)
		b := color.New(color.FgBlue, color.Bold)
		b.Printf(":%s", wd)
		fmt.Printf(">")

		cmd := readline()
		//fmt.Println(cmd)
		switch {
		case strings.Contains(cmd, "|"):
			cmds := split(cmd, "|")
			executePipe(cmds...)
			//fmt.Println(singlecmd)
		default:
			params := split(cmd, " ")
			//fmt.Println(params)
			execute(params...)
		}

	}
}

func readline() string {
	var msg string
	reader := bufio.NewReader(os.Stdin) // 标准输入输出
	msg, _ = reader.ReadString('\n')    // 回车结束
	msg = strings.TrimSpace(msg)        // 去除最后一个空格
	return msg
}

func execute(cmd ...string) string {
	switch cmd[0] {
	case "":
		return ""
	case "cd":
		cmd[1] = strings.Replace(cmd[1], "~", os.Getenv("HOME"), 1)
		err := os.Chdir(cmd[1])
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}

	default:
		command := exec.Command(cmd[0], cmd[1:]...)
		command.Stderr = os.Stderr
		command.Stdout = os.Stdout

		err := command.Run()
		//fmt.Println(command.Stdout)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "executable file not found in $PATH"):
				fmt.Printf("%s:command not found\n", cmd[0])
			default:
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
		}
		return ""

	}

	/*
		writer := bufio.NewWriter(os.Stdout)
		writer.WriteString(string(out))
		writer.Flush()
	*/
	//fmt.Print(string(out))
	return ""
}

func executePipe(cmds ...string) {
	cmd0 := split(cmds[0], " ")
	cmdline0 := exec.Command(cmd0[0], cmd0[1:]...)
	cmd1 := split(cmds[1], " ")
	cmdline1 := exec.Command(cmd1[0], cmd1[1:]...)
	stdout0, _ := cmdline0.StdoutPipe()
	cmdline0.Start()
	outputBuf0 := bufio.NewReader(stdout0)
	stdin1, _ := cmdline1.StdinPipe()
	outputBuf0.WriteTo(stdin1)
	var outputBuf1 bytes.Buffer
	cmdline1.Stdout = &outputBuf1

	if err := cmdline1.Start(); err != nil {
		fmt.Printf("Error:%s\n", err)
		return
	}
	stdin1.Close()
	if err := cmdline1.Wait(); err != nil {
		fmt.Printf("Error:%s", err)
	}
	fmt.Printf("%s", outputBuf1.Bytes())
}

func split(input string, separator string) []string {
	return strings.Split(input, separator)
}
