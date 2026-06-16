package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/cs161-staff/project2-starter-code/internal/client/app"
	"github.com/cs161-staff/project2-starter-code/internal/client/config"
	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	// Handle host management commands
	switch os.Args[1] {
	case "host":
		hostCmd(cfg, os.Args[2:])
		return
	}

	// Connect to KV server via config
	if os.Args[1] != "read" {
		hostName := "default"
		if h := os.Getenv("SHARELOCK_HOST"); h != "" {
			hostName = h
		}
		entry, ok := cfg.Get(hostName)
		if ok && entry.Addr != "" {
			app.Connect(entry.Addr, entry.TLS)
			defer app.Disconnect()
		}
	}

	cli := &app.Client{}

	switch os.Args[1] {
	case "inituser":
		cmd := flag.NewFlagSet("inituser", flag.ExitOnError)
		username := cmd.String("username", "", "")
		password := cmd.String("password", "", "")
		cmd.Parse(os.Args[2:])
		if err := cli.InitUser(*username, *password); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "getuser":
		cmd := flag.NewFlagSet("getuser", flag.ExitOnError)
		username := cmd.String("username", "", "")
		password := cmd.String("password", "", "")
		cmd.Parse(os.Args[2:])
		if err := cli.GetUser(*username, *password); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "storefile":
		cmd := flag.NewFlagSet("storefile", flag.ExitOnError)
		filename := cmd.String("filename", "", "")
		content := cmd.String("content", "", "")
		cmd.Parse(os.Args[2:])
		if err := cli.StoreFile(*filename, []byte(*content)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "loadfile":
		cmd := flag.NewFlagSet("loadfile", flag.ExitOnError)
		filename := cmd.String("filename", "", "")
		cmd.Parse(os.Args[2:])
		data, err := cli.LoadFile(*filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(data))

	case "appendtofile":
		cmd := flag.NewFlagSet("appendtofile", flag.ExitOnError)
		filename := cmd.String("filename", "", "")
		content := cmd.String("content", "", "")
		cmd.Parse(os.Args[2:])
		if err := cli.AppendToFile(*filename, []byte(*content)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "createinvitation":
		cmd := flag.NewFlagSet("createinvitation", flag.ExitOnError)
		filename := cmd.String("filename", "", "")
		recipient := cmd.String("recipient", "", "")
		cmd.Parse(os.Args[2:])
		invite, err := cli.CreateInvitation(*filename, *recipient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(invite.String())

	case "acceptinvitation":
		cmd := flag.NewFlagSet("acceptinvitation", flag.ExitOnError)
		sender := cmd.String("sender", "", "")
		invitation := cmd.String("invitation", "", "")
		filename := cmd.String("filename", "", "")
		cmd.Parse(os.Args[2:])
		invUUID, err := uuid.Parse(*invitation)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid invitation UUID: %v\n", err)
			os.Exit(1)
		}
		if err := cli.AcceptInvitation(*sender, invUUID, *filename); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "revokeaccess":
		cmd := flag.NewFlagSet("revokeaccess", flag.ExitOnError)
		filename := cmd.String("filename", "", "")
		recipient := cmd.String("recipient", "", "")
		cmd.Parse(os.Args[2:])
		if err := cli.RevokeAccess(*filename, *recipient); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "read":
		cmd := flag.NewFlagSet("read", flag.ExitOnError)
		filename := cmd.String("filename", "", "")
		address := cmd.String("address", "localhost:8080", "")
		tlsEnabled := cmd.Bool("tls", true, "use TLS encryption")
		cmd.Parse(os.Args[2:])
		if err := cli.ReadFileTLS(*filename, *address, *tlsEnabled); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	default:
		printUsage()
		os.Exit(1)
	}
}

func hostCmd(cfg *config.Config, args []string) {
	if len(args) == 0 {
		hostUsage()
		return
	}
	switch args[0] {
	case "list":
		hosts := cfg.List()
		names := make([]string, 0, len(hosts))
		for n := range hosts {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			h := hosts[n]
			tlsStr := "tls"
			if !h.TLS {
				tlsStr = "plain"
			}
			mark := " "
			if n == "default" {
				mark = "*"
			}
			fmt.Printf("%s %-12s %s (%s)\n", mark, n, h.Addr, tlsStr)
		}

	case "add":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: host add <name> <addr> [--tls=true]")
			os.Exit(1)
		}
		name := args[1]
		addr := args[2]
		tls := true
		for i := 3; i < len(args); i++ {
			if args[i] == "--tls=false" {
				tls = false
			}
		}
		cfg.Set(name, addr, tls)
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "save error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "remove":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: host remove <name>")
			os.Exit(1)
		}
		cfg.Remove(args[1])
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "save error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	case "default":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "usage: host default <name>")
			os.Exit(1)
		}
		entry, ok := cfg.Get(args[1])
		if !ok {
			fmt.Fprintf(os.Stderr, "host %q not found\n", args[1])
			os.Exit(1)
		}
		cfg.Set("default", entry.Addr, entry.TLS)
		cfg.Remove(args[1])
		if err := cfg.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "save error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("ok")

	default:
		hostUsage()
	}
}

func hostUsage() {
	fmt.Fprintln(os.Stderr, `Usage: host <command> [args]

Commands:
  list                   list all hosts
  add <name> <addr>      add a host (--tls=false for plain TCP)
  remove <name>          remove a host
  default <name>         set an existing host as default`)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: %s <command> [flags]

Environment:
  SHARELOCK_HOST   select host from .hosts config (default: "default")

Config file:
  ~/.config/sharelock/.hosts  (fallback: ~/.hosts)

Host management:
  host list
  host add <name> <addr> [--tls=false]
  host remove <name>
  host default <name>

Commands:
  inituser -username <name> -password <pw>
  getuser -username <name> -password <pw>
  storefile -filename <name> -content <data>
  loadfile -filename <name>
  appendtofile -filename <name> -content <data>
  createinvitation -filename <name> -recipient <user>
  acceptinvitation -sender <user> -invitation <uuid> -filename <name>
  revokeaccess -filename <name> -recipient <user>
  read -filename <name> [-address <host:port>] [-tls=true]
`, os.Args[0])
}
