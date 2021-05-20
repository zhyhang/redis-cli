package redis

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/zhyhang/redis-client/util"
	"net"
	"strconv"
	"time"
)

type Tunnel struct {
	Address  string
	conn     net.Conn
	writerBg *bufio.Writer
	writer   *Writer
	reader   *Reader
	Linked   bool
}

const timeOut = time.Second * 60

func Establish(host string, port int) *Tunnel {
	addr := address(host, port)
	c, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		fmt.Printf(NotLinkMsg(addr)+": %v\n", err)
		return &Tunnel{
			Address: addr,
			Linked:  false}
	}
	wbg := bufio.NewWriter(c)
	return &Tunnel{
		Address:  addr,
		conn:     c,
		Linked:   true,
		writerBg: wbg,
		writer:   NewWriter(wbg),
		reader:   NewReader(c),
	}
}

func (tun *Tunnel) Destroy() {
	if tun.Linked {
		err := tun.conn.Close()
		tun.Linked = false
		if err != nil {
			fmt.Printf("Disconnect error %v\n", err)
		}
	}
}

func (tun *Tunnel) Request(cmd string) (string, error) {
	if cmd == "" {
		return "", nil
	}
	if !tun.Linked {
		return "", errors.New("")
	}
	_, err := tun.writer.Write(util.StringToBytes(cmd + "\r\n"))
	if err != nil {
		tun.Destroy()
		return "", err
	}
	err = tun.writerBg.Flush()
	if err != nil {
		tun.Destroy()
		return "", err
	}
	strReply, err := tun.reader.ReadCmdText()
	if err != nil {
		tun.Destroy()
		return "", err
	}
	return strReply, nil
}

func address(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func NotLinkMsg(addr string) string {
	return "Could not connect to Redis at " + addr
}
