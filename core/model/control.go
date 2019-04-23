package model

// client Control Message
const (
	SignUp       = 0x01
	Login        = 0x02
	Logout       = 0x03
	Upload       = 0x08
	GetDSSD      = 0x0A
	ChangeDSSD   = 0x0B
	DownloadFile = 0x0C
)

// Server Control Message
const (
	SignUpSuccess             = 0xC1
	SignUpFailForWeakPassword = 0xE2
	SignUpFailForUniqueName   = 0xE3
	SignUpFailForOther        = 0xE0

	LoginSuccess = 0xC1
	LoginFail    = 0xE0

	LogoutSuccess = 0xC1
	LogoutFail    = 0xE0

	FileExist         = 0xC1
	FileIsUploading   = 0xC2
	FileNeedUpload    = 0xC3
	UploadRequestFail = 0xE0

	AuthCheckFail = 0xF0

	FileUploadSucess = 0xC1
	FileUploadFail   = 0xE0

	ChangeDSSDSucess     = 0xC1
	ChangeDSSDFail       = 0xE0
	ClientDSSDNeedUpdate = 0xE1

	DownloadSuccess = 0xC1
	DownloadFail    = 0xE0
)
