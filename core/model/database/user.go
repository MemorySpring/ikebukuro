package database

var TableName = "user"

type User struct {
	Id   int    `orm:"column(id)"`
	Name string `orm:"column(name)"`
	// MD5 hash of user password
	Password string `orm:"column(password)"`
}

func (u *User) TableName() string {
	return "user"
}

type Auth struct {
	Id int `orm:"column(id)"`
	// Access token
	Token string `orm:"column(token)"`
	// Authorized user, foreign key for User
	UserId int `orm:"column(user_id)"`
	// Unix timestamp of the expiring time. When user logs out, this value should be set
	// to negative (eg. -1). Whenever server check whether a token is valid, it will expect
	// this value to be greater than the current timestamp
	ExpireAt int64 `orm:"column(expires_at)"`
}

func (a *Auth) TableName() string {
	return "auth"
}
