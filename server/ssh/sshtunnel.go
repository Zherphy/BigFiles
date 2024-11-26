package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"math"
	"net"
	"strings"
	"syscall"
	"time"
)

type SSHTunnel struct {
	client *ssh.Client
}

type Tunnel struct {
	Remote string `json:"remote"`
	Local  string `json:"local"`
}

func NewSSHTunnel() *SSHTunnel {
	st := new(SSHTunnel)
	return st
}

func (st *SSHTunnel) Start() {

	st.initSSHClient()
	t := Tunnel{
		Remote: "gitee.com:22",
		Local:  "127.0.0.1:23231",
	}
	go st.connect(t)
}

func (st *SSHTunnel) Close() {
	if nil != st.client {
		st.client.Close()
	}
}

func (st *SSHTunnel) GetSSHClient() (*ssh.Client, error) {
	if st.client != nil {
		return st.client, nil
	}
	var auth []ssh.AuthMethod
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(""))

	sc := &ssh.ClientConfig{
		User: "root",
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	var err error
	st.client, err = ssh.Dial("tcp", "localhost:23231", sc)
	if err != nil {
		return nil, err
	}
	log.Printf("连接到服务器成功: %s", "localhost:23231")
	return st.client, err
}

func (st *SSHTunnel) connect(t Tunnel) {
	tid := fmt.Sprintf("%s-%s", t.Local, t.Remote)
	ll, err := net.Listen("tcp", t.Local)
	if err != nil {
		log.Printf("隧道[%s]接收开启失败, 错误: %v", tid, err)
		return
	}
	defer func() {
		ll.Close()
		log.Printf("隧道[%s]接收关闭!", tid)
	}()
	log.Printf("隧道[%s]接收开启!", tid)
	sno := int64(0)
	for {
		lc, err := ll.Accept()
		if err != nil {
			log.Printf("隧道[%s]接收连接失败, 错误: %v", tid, err)
			return
		}
		sc, err := st.GetSSHClient()
		if err != nil {
			log.Printf("隧道[%s]接入服务失败, 错误: %v", tid, err)
			lc.Close()
			continue
		}
		rc, err := sc.Dial("tcp", t.Remote)
		if err != nil {
			log.Printf("隧道[%s]接入获取连接失败, 错误: %v", err)
			sc.Close()
			lc.Close()
			continue
		}
		if sno >= math.MaxInt64 {
			sno = 0
		}
		sno += 1
		cid := fmt.Sprintf("%s:%d", tid, sno)
		go st.transfer(cid, lc, rc)
	}
}

func (st *SSHTunnel) setPass() {
	fmt.Printf("请输入登陆密码[%s@%s]:", "root", "localhost:23231")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	pass := string(bytePassword)
	fmt.Println(pass)
}

func (st *SSHTunnel) initSSHClient() {
	var err error
	for {
		st.client, err = st.GetSSHClient()
		if nil != err {
			error := err.Error()
			log.Printf("连接到服务器[%s]失败, 错误: %s", "", error)
			if strings.Contains(error, "unable to authenticate") {
				st.setPass()
				continue
			}
			if strings.Contains(error, "i/o timeout") {
				log.Printf("连接到服务器[%s]超时!", "")
				time.Sleep(2 * time.Second)
				continue
			}
		}
		return
	}
}

func (st *SSHTunnel) transfer(cid string, lc net.Conn, rc net.Conn) {
	defer rc.Close()
	defer lc.Close()
	go func() {
		defer lc.Close()
		defer rc.Close()
		io.Copy(rc, lc)
	}()
	log.Printf("通道[%s]已连接!", cid)
	io.Copy(lc, rc)
	log.Printf("通道[%s]已断开!", cid)
}
