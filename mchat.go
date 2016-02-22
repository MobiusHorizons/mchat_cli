package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/MobiusHorizons/go_mqtt_chat"
	"io"
	"os"
	"strings"
)

func main() {
	nick := flag.String("nick", "", "Nickname on the chat server")
	server := flag.String("server", "mqtt://test.mosquitto.org:1883", "server to connect to")
	flag.Parse()

	fmt.Println(flag.Args())
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s --nick <nickname> [--server mqtt(s)://server.tld:port ] recipient\n", os.Args[0])
		os.Exit(1)
	}

	recipient := flag.Arg(0)

	//	fmt.Println("nick = ", *nick, "server = ", *server)
	client := go_mqtt_chat.New(*nick, *server, nil)
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	go func() {
		m := client.Listen()
		for m != nil {
			fmt.Println(m.From + ": " + m.Message)
			m = client.Listen()
		}
		fmt.Println("channel closed")
	}()

	reader := bufio.NewReader(os.Stdin)
	for line, err := reader.ReadString('\n'); err != io.EOF; line, err = reader.ReadString('\n') {
		client.Say(recipient, strings.TrimRight(line, "\n"))
	}
}
