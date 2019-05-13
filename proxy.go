package zcmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"time"

	"github.com/cybozu-go/well"
	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

// Proxy is client for proxy
type Proxy struct {
	sshCfg []*sshClientConfig
}

type sshClientConfig struct {
	forwardType ProxyForwardType
	user        string
	key         ssh.Signer
	sshAddr     string
	localAddr   string
	remoteAddr  string
}

// NewProxy returns Proxy
func NewProxy(cfgs []ProxyConfig) (*Proxy, error) {
	proxy := new(Proxy)
	for _, cfg := range cfgs {
		key, err := parsePrivateKey(cfg.PrivateKey)
		if err != nil {
			return nil, err
		}
		for _, fwdCfg := range cfg.Forward {
			sshCfg := &sshClientConfig{
				forwardType: fwdCfg.Type,
				user:        cfg.User,
				key:         key,
				sshAddr:     fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
				localAddr:   fmt.Sprintf("%s:%d", fwdCfg.BindAddress, fwdCfg.BindPort),
				remoteAddr:  fmt.Sprintf("%s:%d", fwdCfg.RemoteAddress, fwdCfg.RemotePort),
			}

			proxy.sshCfg = append(proxy.sshCfg, sshCfg)
		}
	}
	return proxy, nil
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	buff, err := ioutil.ReadFile(keyPath)
	if err != err {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(buff)
	if err != nil {
		var passphrase string
		err := Ask(&passphrase, "ssh passphrase for "+keyPath, true)
		if err != nil {
			return nil, err
		}
		return ssh.ParsePrivateKeyWithPassphrase(buff, []byte(passphrase))
	}
	return signer, nil
}

func makeSshConfig(cfg sshClientConfig) (*ssh.ClientConfig, error) {
	return &ssh.ClientConfig{
		User: cfg.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(cfg.key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}, nil
}

func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.WithError(err)
		}
		chDone <- true
	}()

	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.WithError(err)
		}
		chDone <- true
	}()

	<-chDone
}

func (p *Proxy) Run(ctx context.Context) error {
	env := well.NewEnvironment(ctx)
	for _, cfg := range p.sshCfg {
		sshCfg, err := makeSshConfig(*cfg)
		if err != nil {
			return err
		}

		fmt.Printf("%#v\n", sshCfg)
		fmt.Printf("%#v\n", cfg)
		conn, err := ssh.Dial("tcp", cfg.sshAddr, sshCfg)
		if err != nil {
			return err
		}
		defer conn.Close()

		remote, err := conn.Dial("tcp", cfg.remoteAddr)
		if err != nil {
			fmt.Println(err)
			return err
		}

		local, err := net.Listen("tcp", cfg.localAddr)
		if err != nil {
			return err
		}
		defer local.Close()

		env.Go(func(ctx context.Context) error {
			for {
				client, err := local.Accept()
				if err != nil {
					return err
				}

				handleClient(client, remote)
			}
		})
	}
	env.Stop()
	return env.Wait()
}
