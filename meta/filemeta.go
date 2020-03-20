package meta

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var FileHashForTimes []string
var FileMetas map[string]FileMeta

func init() {
	FileHashForTimes = make([]string, 0)
	FileMetas = make(map[string]FileMeta)
}
