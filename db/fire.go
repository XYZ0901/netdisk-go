package db

import (
	"database/sql"
	"errors"
	"log"
	"netdisk-go/db/mysql"
	"netdisk-go/meta"
	"time"
)

type tablefile struct {
	file_sha1 string
	file_name string
	file_size int64
	file_addr string
	UploadAt  time.Time
}

// 文件上传后进行本地持久化（保存到mysql）
func SaveFile(fileMeta meta.FileMeta) error {

	tb := tablefile{}
	tb, err := queryFile(fileMeta.FileSha1)
	if err != nil || err != sql.ErrNoRows {
		return err
	}
	if tb.file_sha1 == "" {
		err = insertFile(fileMeta)
		if err != nil {
			return err
		}
		return nil
	}

	// 暂定直接更新文件表中的文件，后期改为更新用户文件表的文件
	err = UpdateFile(fileMeta)
	if err != nil {
		return err
	}
	return nil
}

// 利用hash查数据库中文件信息
func queryFile(filehash string) (tablefile, error) {
	stmt, err := mysql.MysqlDb.Prepare(
		"SELECT file_sha1,file_name,file_size,file_addr,create_at FROM tbl_file WHERE file_sha1 = ?")
	if err != nil {
		return tablefile{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(filehash)
	tf := tablefile{}
	err = row.Scan(&tf.file_sha1, &tf.file_name, &tf.file_size, &tf.file_addr, &tf.UploadAt)
	return tf, nil
}
func queryFiles(limit string) (tablefile, error) {
	stmt, err := mysql.MysqlDb.Prepare(
		"SELECT file_sha1,file_name,file_size,file_addr,create_at FROM tbl_file ORDER BY id LIMIT ?")
	if err != nil {
		return tablefile{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(limit)
	tf := tablefile{}
	err = row.Scan(&tf.file_sha1, &tf.file_name, &tf.file_size, &tf.file_addr, &tf.UploadAt)
	return tf, nil
}

// 添加入数据库
func insertFile(meta meta.FileMeta) error {
	tx, err := mysql.MysqlDb.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(
		"INSERT INTO tbl_file(`file_sha1`,`file_name`,`file_size`,`file_addr`) VALUES (?,?,?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(meta.FileSha1, meta.FileName, meta.FileSize, meta.Location)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if err = tx.Commit(); err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		err = errors.New("insert fail")
		return err
	}
	return nil
}

// 修改数据库内存储信息
func UpdateFile(fileMeta meta.FileMeta) error {
	tx, err := mysql.MysqlDb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("UPDATE tbl_file SET file_name = ?,WHERE file_sha1 = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(fileMeta.FileName, fileMeta.FileSha1)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		return errors.New("update fail")
	}
	return nil
}

// 删除文件 在数据库内打标签为（禁用） 数据不删 以备多用户时
// 多个用户存储同一个hash文件
func DeleteFile(filehash string) error {
	tx, err := mysql.MysqlDb.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("UPDATE tbl_file SET status = 2,WHERE file_sha1 = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(filehash)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
