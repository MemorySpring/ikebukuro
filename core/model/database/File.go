package database

type Files struct {
	Id int `orm:"column(id)"`
	// Stored file name on actual server storage
	StorageName string `orm:"column(storage_name)"`
	// MD5 hash for the specific file
	Hash string `orm:"column(hash)"`
	// To show how many user has this file in his/her own DSSD
	DepCount int `orm:"column(dep_count)"`
	// Save the breakpoint for client to resume upload from.
	UploadPosition int `orm:"column(upload_pos)"`
	// To show whether the file is completely uploaded.
	Complete bool `orm:"column(complete)"`
	// File is now uploading by someone
	Uploading bool `orm:"column(uploading)"`
}

func (f *Files) TableName() string {
	return "file"
}

type Directories struct {
	Id int `orm:"column(id)"`
	// foreign key for Users
	UserId int `orm:"column(user_id)"`
	// Directory Strucure Serialized Data(DSSD), contains one serialized dir_data object
	// for one specific user
	Data string `orm:"column(data)"`
	// server-side-generated DSSD version.
	Version string `orm:"column(version)"`
}

func (d *Directories) TableName() string {
	return "directory"
}
