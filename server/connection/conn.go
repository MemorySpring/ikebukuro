package connection

import (
	"Ikebukuro/core/model"
	"Ikebukuro/core/model/database"
	"Ikebukuro/core/util"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" // 导入数据库驱动
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	Day = 60 * 60 * 24
)

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.DebugLevel)
}

func DataBaseInit(source string) {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = orm.RegisterDataBase("default", "mysql", source, 30)
	if err != nil {
		fmt.Println(err)
		return
	}
	orm.RegisterModel(new(database.User))
	orm.RegisterModel(new(database.Auth))
	orm.RegisterModel(new(database.Files))
	orm.RegisterModel(new(database.Directories))
}

type DiskServer struct {
	Port  string
	Ormer orm.Ormer
}

func (server *DiskServer) Start() {
	fmt.Println("服务器启动！")
	l, err := net.Listen("tcp4", server.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())
	for {
		c, err := l.Accept()
		fmt.Println("收到客户端连接")
		if err != nil {
			fmt.Println(err)
			return
		}
		go server.handleConnection(c)
	}
}

func (server *DiskServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	controlMessage := make([]byte, 1)
	_, err := conn.Read(controlMessage)
	if err != nil {
		log.Error(err)
		return
	}
	switch controlMessage[0] {
	case model.SignUp:
		fmt.Println("收到注册请求")
		server.signUp(conn)
	case model.Login:
		fmt.Println("收到登录请求")
		server.login(conn)
	case model.Logout:
		fmt.Println("收到登出请求")
		server.logout(conn)
	case model.ChangeDSSD:
	case model.GetDSSD:
	case model.Upload:
		server.uploadFile(conn)
	case model.DownloadFile:
	}
}

func (server *DiskServer) login(conn net.Conn) {
	username, err := util.GetString(conn)
	password, err := util.GetString(conn)
	user := &database.User{Name: username, Password: password}
	err = server.Ormer.Read(user, "Name", "Password")
	var result []byte
	if err != nil {
		log.WithFields(log.Fields{
			"username": user.Name,
			"password": user.Password,
		}).Debug("用户登录失败")
		fmt.Println(err)
		result, _ = util.WrapControlMessage(model.LoginFail, "")
	} else {
		token := util.RandStringRunes(32)
		auth := &database.Auth{
			Token:    token,
			UserId:   user.Id,
			ExpireAt: time.Now().Unix() + Day,
		}
		_, err = server.Ormer.Insert(auth)
		if err != nil {
			log.WithFields(log.Fields{
				"username": user.Name,
				"password": user.Password,
			}).Debug("保存用户Token失败")
			result, _ = util.WrapControlMessage(model.LoginFail, "")
		} else {
			log.WithFields(log.Fields{
				"username": user.Name,
				"password": user.Password,
			}).Debug("用户登录成功")
			result, _ = util.WrapControlMessage(model.LoginSuccess, token)
		}
	}
	conn.Write(result)
}

func (server *DiskServer) signUp(conn net.Conn) {
	username, err := util.GetString(conn)
	password, err := util.GetString(conn)
	user := &database.User{Name: username, Password: password}
	err = server.Ormer.Read(user, "Name")
	var result []byte
	if err != nil {
		log.WithFields(log.Fields{
			"username": user.Name,
			"password": user.Password,
		}).Debug("用户注册成功")
		_, err = server.Ormer.Insert(user)
		result, _ = util.WrapControlMessage(model.SignUpSuccess)
	} else {
		log.WithFields(log.Fields{
			"username": user.Name,
			"password": user.Password,
		}).Debug("用户注册失败")
		result, _ = util.WrapControlMessage(model.SignUpFailForUniqueName)
	}
	conn.Write(result)
}

func (server *DiskServer) logout(conn net.Conn) {
	var result []byte
	token, err := util.GetString(conn)
	auth := &database.Auth{Token: token}
	err2 := server.Ormer.Read(auth, "Token")
	auth.ExpireAt = -1
	_, err = server.Ormer.Update(auth)
	if err != nil || err2 != nil {
		log.WithFields(log.Fields{
			"token": token,
		}).Debug("token更新失败")
		result, _ = util.WrapControlMessage(model.LogoutFail)
	} else {
		result, _ = util.WrapControlMessage(model.LogoutSuccess)
	}
	conn.Write(result)
}

func (server *DiskServer) uploadFile(conn net.Conn) {
	var result []byte
	token, _ := util.GetString(conn)
	fileMD5, _ := util.GetString(conn)
	fileLength, _ := util.GetInt(conn)
	if !server.checkAuth(token) {
		result, _ = util.WrapControlMessage(model.AuthCheckFail)
	} else {
		file := &database.Files{Hash: fileMD5}
		err := server.Ormer.Read(file, "Hash")
		if err != nil {
			result, _ = util.WrapControlMessage(model.FileNeedUpload, 0, 0)
			conn.Write(result)
			file.Uploading = true
			file.StorageName = file.Hash
			file.DepCount = 0
			file.Complete = false
			file.UploadPosition = 0
			server.Ormer.Insert(file)
			server.recvFile(conn, fileLength, file)
		} else if file.Uploading {
			result, _ = util.WrapControlMessage(model.FileIsUploading, 0, file.Id)
			conn.Write(result)
		} else if !file.Complete {
			result, _ = util.WrapControlMessage(model.FileNeedUpload, file.UploadPosition, file.Id)
			conn.Write(result)
			file.Uploading = true
			server.Ormer.Update(file)
			server.recvFile(conn, fileLength-file.UploadPosition, file)
		} else {
			result, _ = util.WrapControlMessage(model.FileExist, 0, file.Id)
			conn.Write(result)
		}
	}

}

func (server *DiskServer) checkAuth(token string) bool {
	auth := &database.Auth{Token: token}
	err2 := server.Ormer.Read(auth, "Token")
	if err2 != nil || auth.ExpireAt < time.Now().Unix() {
		log.WithFields(log.Fields{
			"token": token,
		}).Debug("token已失效")
		return false
	}
	return true
}

func (server *DiskServer) recvFile(conn net.Conn, length int, file *database.Files) {
	fileData := make([]byte, length)
	buffer := make([]byte, util.MAX_INT)
	recvSum := 0
	recvCount := 0
	for recvCount < length {
		recvCount, err := conn.Read(buffer)
		if err != nil {
			// 本次上传异常中断 保存进度等待下次上传
			file.UploadPosition += recvSum
			server.Ormer.Update(file)
			file.UploadPosition = recvSum
			file.Complete = false
			file.DepCount = 1
			file.StorageName = file.Hash
			file.Uploading = false
			server.Ormer.Insert(file)
		} else {
			recvSum += recvCount
			recvCount = 0
		}
	}
}
