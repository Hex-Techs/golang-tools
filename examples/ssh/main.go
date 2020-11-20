package main

import (
	"flag"
	"fmt"
	"os"

	ssh "github.com/ymichaelson/golang-tools/ssh"
	"github.com/ymichaelson/klog"
)

var (
	username       = flag.String("username", "root", "remote ssh server username.")
	host           = flag.String("host", "127.0.0.1", "remote ssh server host.")
	port           = flag.String("port", "22", "remote ssh server port.")
	privateKeyPath = flag.String("keypath", "", "private key path.")
	protocol       = flag.String("protocol", "tcp", "protocol.")
	destPath       = flag.String("destpath", "", "Destination file path.")
	destFileName   = flag.String("destfile", "", "Destination file Name.")
)

func main() {
	flag.Parse()

	sshClient, err := ssh.NewInternalSshClient(*username, *host, *port, *privateKeyPath, *protocol)
	if err != nil {
		klog.Error("create new internal ssh client failed: ", err)
		os.Exit(1)
	}

	defer sshClient.Close()

	// copy file from location to remote
	err = sshClient.Copy("./examples/ssh/test.txt", fmt.Sprintf("%s/%s", *destPath, *destFileName))
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}

	// set command, and exec it.
	command := "ls -l " + *destPath
	sshClient.SetCommand(command)
	out, err := sshClient.Run()
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}

	klog.Info("\n", string(out))
}
