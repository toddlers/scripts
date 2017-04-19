package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint
	Config *ssh.ClientConfig
}

func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dail error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func main() {
	lport := flag.Int("lport", 9000, "Local port")
	jip := flag.String("jip", "", "Jump Server IP address")
	rip := flag.String("rip", "", "Remote Server IP Address")
	rport := flag.Int("rport", 15672, "Remote Server port")
	uname := flag.String("uname", "", "Username")

	flag.Parse()
	if *jip == "" || *rip == "" || *uname == "" {
		fmt.Println("Input missing")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		return
	}

	localEndpoint := &Endpoint{
		Host: "127.0.0.1",
		Port: *lport,
	}

	serverEndpoint := &Endpoint{
		Host: *jip,
		Port: 22,
	}

	remoteEndpoint := &Endpoint{
		Host: *rip,
		Port: *rport,
	}

	sshConfig := &ssh.ClientConfig{
		User: *uname,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}

	tunnel := &SSHTunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}
	err := tunnel.Start()
	if err == nil {
		fmt.Println("Tunnel connected : 127.0.0.1", *lport)
	}
}
