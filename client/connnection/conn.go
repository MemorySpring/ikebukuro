package connnection

import (
	"Ikebukuro/core/model"
	"Ikebukuro/core/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

type DiskClient struct {
	Addr  string
	conn  *net.TCPConn
	token []byte
}

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.DebugLevel)
}

func (client *DiskClient) connect() error {
	var err error
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", client.Addr)
	client.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	return err
}

func (client *DiskClient) disConnect() {
	client.conn.Close()
}

func (client *DiskClient) SignUp(username string, password string) error {
	err := client.connect()
	defer client.disConnect()
	if err != nil {
		log.Error(err)
		return err
	}
	data, err := util.WrapControlMessage(model.SignUp, username, password)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = client.conn.Write(data)
	signResult := make([]byte, 1)
	_, err = client.conn.Read(signResult)
	if signResult[0] == model.SignUpSuccess {
		_, err = client.conn.Read(client.token)
		return nil
	} else if signResult[0] == model.SignUpFailForWeakPassword {
		return errors.New("注册失败，密码强度不足")
	} else if signResult[0] == model.SignUpFailForUniqueName {
		return errors.New("注册失败，用户名已存在")
	} else if signResult[0] == model.SignUpFailForOther {
		return errors.New("注册失败，异常错误")
	}
	return errors.New("Fail:Unknown Control Message")
}

// TODO 完善错误处理
func (client *DiskClient) Login(username string, password string) error {
	err := client.connect()
	defer client.disConnect()
	if err != nil {
		log.Error(err)
		return err
	}
	data, err := util.WrapControlMessage(model.Login, username, password)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = client.conn.Write(data)
	loginResult := make([]byte, 1)
	client.token = make([]byte, 32)
	_, err = client.conn.Read(loginResult)
	if loginResult[0] == model.LoginSuccess {
		_, err = client.conn.Read(client.token)
		return nil
	} else if loginResult[0] == model.LoginFail {
		return errors.New("登录失败，用户名或密码错误")
	}
	return errors.New("Fail:Unknown Control Message")
}
