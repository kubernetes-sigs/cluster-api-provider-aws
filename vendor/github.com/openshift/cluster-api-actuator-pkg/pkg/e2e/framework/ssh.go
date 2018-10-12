package framework

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func createSSHClient(sshConfig *SSHConfig) (*ssh.Client, error) {
	signer, err := LoadPrivateKey(sshConfig.Key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: sshConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", sshConfig.Host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: " + err.Error())
	}
	return client, nil

}

func LoadPrivateKey(prkey string) (ssh.Signer, error) {
	fp, err := os.Open(prkey)
	if err != nil {
		return nil, fmt.Errorf("unable to open %q: %v", prkey, err)
	}
	defer fp.Close()

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read %q: %v", prkey, err)
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q: %v", prkey, err)
	}

	return key, err
}
