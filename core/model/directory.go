package model

type Directory struct {
	Name   string      `json:"name"`
	Files  []File      `json:"files"`
	SubDir []Directory `json:"sub_dir"`
}

type File struct {
	Id       int    `json:"file_id"`
	FileName string `json:"filename"`
}
