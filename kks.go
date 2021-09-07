package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkga/kks/cmd"
)

type KakContext struct {
	session string
	client  string
}

//go:embed init.kak
var initStr string

var session string
var client string

func main() {
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	attachCmd := flag.NewFlagSet("attach", flag.ExitOnError)
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	killCmd := flag.NewFlagSet("kill", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	envCmd := flag.NewFlagSet("env", flag.ExitOnError)

	sessionCmds := []*flag.FlagSet{
		editCmd, sendCmd, attachCmd, getCmd, killCmd,
	}
	for _, cmd := range sessionCmds {
		cmd.StringVar(&session, "s", "", "Kakoune session")
		cmd.StringVar(&client, "c", "", "Kakoune client")
	}

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "edit", "e":
		editCmd.Parse(os.Args[2:])
	case "send", "s":
		sendCmd.Parse(os.Args[2:])
	case "attach", "a":
		attachCmd.Parse(os.Args[2:])
	case "get":
		getCmd.Parse(os.Args[2:])
	case "kill":
		killCmd.Parse(os.Args[2:])
	case "list", "l", "ls":
		listCmd.Parse(os.Args[2:])
	case "env":
		envCmd.Parse(os.Args[2:])
	case "init":
		fmt.Print(initStr)
	default:
		fmt.Println("unknown command:", os.Args[1])
		os.Exit(1)
	}

	if editCmd.Parsed() {
		filename := editCmd.Arg(0)
		if filename == "" {
			printHelp()
			os.Exit(2)
		}

		context, err := getContext()
		if err != nil {
			log.Fatal(err)
		}
		cmd.Edit(filename, context.session, context.client)
	}

	if attachCmd.Parsed() {
		context, err := getContext()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Edit("", context.session, context.client); err != nil {
			log.Fatal(err)
		}
	}

	if sendCmd.Parsed() {
		args := sendCmd.Args()
		kakCommand := strings.Join(args, " ")

		context, err := getContext()
		if err != nil {
			log.Fatal(err)
		}
		cmd.Send(kakCommand, context.session, context.client)
	}

	if getCmd.Parsed() {
		arg := getCmd.Arg(0)

		context, err := getContext()
		if err != nil {
			log.Fatal(err)
		}

		out, err := cmd.Get(arg, context.session, context.client)
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(arg, "buflist") {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("cwd:", cwd)

			kakwd, err := cmd.Get("%sh{pwd}", context.session, context.client)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("kakwd:", kakwd)

			relPath, _ := filepath.Rel(cwd, kakwd[0])
			fmt.Println("rel path:", relPath)

			for i, buf := range out {
				out[i] = fmt.Sprintf("%s/%s", relPath, buf)
			}
		}

		fmt.Println(strings.Join(out, "\n"))
	}

	if killCmd.Parsed() {
		kakCommand := "kill"
		context, err := getContext()
		if err != nil {
			log.Fatal(err)
		}

		cmd.Send(kakCommand, context.session, context.client)
	}

	if listCmd.Parsed() {
		cmd.List()
	}

	if envCmd.Parsed() {
		context, err := getContext()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("session: %s\n", context.session)
		fmt.Printf("client: %s\n", context.client)
	}

}

func getContext() (*KakContext, error) {
	c := KakContext{
		session: os.Getenv("KKS_SESSION"),
		client:  os.Getenv("KKS_CLIENT"),
	}
	if session != "" {
		c.session = session
	}
	if client != "" {
		c.client = client
	}
	if c.session == "" {
		return nil, errors.New("No session in context")
	}
	return &c, nil
}

func printHelp() {
	fmt.Println("Handy Kakoune companion.")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  kks <command> [-s <session>] [-c <client>] [<args>]")
	fmt.Println()
	fmt.Println("COMMANDS")
	fmt.Println("  edit, e         edit file")
	fmt.Println("  send, s         send command")
	fmt.Println("  attach, a       attach to session")
	fmt.Println("  list, l         list sessions and clients")
	fmt.Println("  kill, k         kill session")
	fmt.Println("  get             get %{val}, %{opt} and friends")
	fmt.Println("  env             print env")
	fmt.Println("  init            print Kakoune definitions")
	fmt.Println()
	fmt.Println("ENVIRONMENT VARIABLES")
	fmt.Println("  KKS_SESSION     Kakoune session")
	fmt.Println("  KKS_CLIENT      Kakoune client")
}
