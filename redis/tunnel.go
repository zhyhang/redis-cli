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
	bad      bool
}

const timeOut = time.Second * 60

func Establish(host string, port int) *Tunnel {
	addr := address(host, port)
	c, err := net.DialTimeout("tcp", addr, timeOut)
	if err != nil {
		fmt.Printf(notLinkMsg(addr)+": %v\n", err)
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
		bad:      false,
	}
}

func (tun *Tunnel) Destroy() error {
	if tun.Linked {
		err := tun.conn.Close()
		tun.Linked = false
		if err != nil {
			fmt.Printf("Disconnect error %v\n", err)
			return err
		}
	}
	return nil
}

func (tun *Tunnel) Request(cmd string) (string, error) {
	if cmd == "" {
		return "", nil
	}
	if tun.bad || !tun.Linked {
		return "", errors.New("")
	}
	_, err := tun.writer.Write(util.StringToBytes(cmd + "\r\n"))
	if err != nil {
		tun.bad = true
		return "", err
	}
	err = tun.writerBg.Flush()
	if err != nil {
		tun.bad = true
		return "", err
	}
	strReply, err := tun.reader.ReadCmdText()
	if err != nil {
		tun.bad = true
		return "", err
	}
	return strReply, nil
}

func address(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

func notLinkMsg(addr string) string {
	return "Could not connect to Redis at " + addr
}

//func Connect(flags *terminal.CmdFlags) {
//	addr := flags.Host + ":" + strconv.Itoa(flags.Port)
//	c, err := net.DialTimeout("tcp", addr, timeOut)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer func(conn net.Conn) {
//		_ = conn.Close()
//	}(c)
//	// write
//	bw := bufio.NewWriter(c)
//	w := NewWriter(bw)
//	_, err = w.Write([]byte("set key1 value2\r\n"))
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	err = bw.Flush()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	// read
//	r :=NewReader(c)
//	result, err := r.ReadString()
//	if err != nil {
//		return
//	}
//	fmt.Println(result)
//}
