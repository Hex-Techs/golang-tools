/*
Copyright The Michaelson.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	DefaultProtocol = "tcp"
)

type SshClient struct {
	Host       string            `json:"host"`
	Port       string            `json:"port"`
	Protocol   string            `json:"protocol"` // tcp tcp4 tcp6 unix
	Command    string            `json:"cmd"`
	Config     *ssh.ClientConfig `json:"config"`
	Client     *ssh.Client       `json:"client"`
	Session    *ssh.Session      `json:"session"`
	SftpClient *sftp.Client      `json:"sftpClient"`
}

// NewInternalSshClient return new SshClient
// user: ssh user
// host: ssh remote host
// port: ssh remote port
// privateKeyPath: private key path
// protocol: ssh protocol, default tcp
func NewInternalSshClient(user, host, port, privateKeyPath, protocol string) (*SshClient, error) {
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	if protocol == "" {
		protocol = DefaultProtocol
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	sshClient := SshClient{
		Host:     host,
		Port:     port,
		Protocol: protocol,
		Config: &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}

	err = sshClient.Dail()
	if err != nil {
		return nil, err
	}

	err = sshClient.NewSession()
	if err != nil {
		return nil, err
	}

	err = sshClient.NewSftpClinet()
	if err != nil {
		return nil, err
	}

	return &sshClient, nil
}

func (s *SshClient) Dail() error {
	client, err := ssh.Dial(s.Protocol, fmt.Sprintf("%s:%s", s.Host, s.Port), s.Config)
	if err != nil {
		return err
	}

	s.Client = client
	return nil
}

func (s *SshClient) NewSession() error {
	session, err := s.Client.NewSession()
	if err != nil {
		return err
	}

	s.Session = session
	return nil
}

func (s *SshClient) NewSftpClinet() error {
	client, err := sftp.NewClient(s.Client)
	if err != nil {
		return err
	}

	s.SftpClient = client
	return nil
}

// Copy file from local to remote
// src: /go/local/src.file
// dest: /go/remote/dest.file
func (s *SshClient) Copy(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	destFile, err := s.SftpClient.Create(dest)
	if err != nil {
		return err
	}

	defer destFile.Close()
	buf, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return err
	}

	_, err = destFile.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (s *SshClient) SetCommand(cmd string) {
	s.Command = cmd
}

func (s *SshClient) Run() ([]byte, error) {
	return s.Session.CombinedOutput(s.Command)
}

func (s *SshClient) Close() {
	s.SftpClient.Close()
	s.Session.Close()
	s.Client.Close()
}
