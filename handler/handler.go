package handler

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"netdisk-go/db"
	"netdisk-go/meta"
	"netdisk-go/util"
	"os"
	"strconv"
	"time"
)

// 处理函数

// 上传
// method: GET(get upload.html) POST(upload)
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			w.Write([]byte("internal server error"))
			return
		}
		w.Write(data)
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("Failed to get data, err: " + err.Error()))
			return
		}
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			FileSize: head.Size,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:5"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		err = db.SaveFile(fileMeta)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传成功
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload finished!"))
}

// 查找文件信息
//method: GET
//// param: filehash
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fMeta := meta.FileMeta{}
	ok := false
	fMeta
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 查询最近存储的文件信息，
// method: GET
// param: limit
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitS := r.FormValue("limit")
	limit, err := strconv.Atoi(limitS)
	if err != nil {
		w.Write([]byte("please input right limit for int"))
		return
	}
	filemetas := []meta.FileMeta{}
	n := len(meta.FileHashForTimes)
	for i := n - 1; i >= n-limit && i >= 0; i-- {
		filemetas = append(filemetas,
			meta.FileMetas[meta.FileHashForTimes[i]])
	}
	data, err := json.Marshal(filemetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 下载文件
// method: GET
// param: filehash
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fhash := r.FormValue("filehash")

	fMeta := meta.FileMeta{}
	ok := false
	if fMeta, ok = meta.FileMetas[fhash]; !ok {
		w.Write([]byte("have not the hash in all file"))
		return
	}
	f, err := os.Open(fMeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 下载文件必备请求头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fMeta.FileName+`"`)
	w.Write(data)
}

// 更新文件名
// method: POST
// param: op=0 filehash filename
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.FormValue("op")
	fileSha1 := r.FormValue("filehash")
	newFileName := r.FormValue("filename")

	if opType != "0" || len(newFileName) < 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fMeta := meta.FileMeta{}
	ok := false
	if fMeta, ok = meta.FileMetas[fileSha1]; !ok {
		w.Write([]byte("have not the hash in all file"))
		return
	}

	fMeta.FileName = newFileName
	meta.FileMetas[fileSha1] = fMeta

	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
