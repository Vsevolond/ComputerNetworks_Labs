package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/skorobogatov/input"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/smtp"
)

const (
	password string = "Je2dTYr6"
	login    string = "iu9networkslabs"
	host     string = "students.yss.su"
	dbname   string = "iu9networkslabs"
)

func get_file() []byte {
	var file []byte
	var err error
	for {
		fmt.Print("Path = ")
		path := input.Gets()
		file, err = ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("Wrong path")
			continue
		}
		break
	}
	return file
}

func interact() ([]string, string, string, string, [][]byte) {
	log.Println("### Send Mail ###")
	fmt.Print("to = ")
	to := input.Gets()
	fmt.Print("name = ")
	name := input.Gets()
	fmt.Print("subject = ")
	subject := input.Gets()
	fmt.Print("message = ")
	message := input.Gets()
	files := make([][]byte, 0)
	fmt.Println("Add some files? (Y/N)")
	for {
		ans := input.Gets()
		if ans == "Y" || ans == "y" {
			file := get_file()
			files = append(files, file)
			fmt.Println("More? (Y/N)")
		} else if ans == "N" || ans == "n" {
			break
		} else {
			fmt.Println("Wrong answer")
		}
	}
	return []string{to}, name, subject, message, files
}

func updateDB(db *sql.DB, to string, name string, subject string, message string, code string, description string) {
	_, err := db.Exec("insert into iu9networkslabs.smtpDonchenko (address_to, recipient, subject, message, code, description) values (?, ?, ?, ?, ?, ?)",
		to, name, subject, message, code, description)
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", login+":"+password+"@tcp("+host+")/"+dbname)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS iu9networkslabs.smtpDonchenko(address_to text, 
				recipient text,
				subject text,
				message text,
				code text, 
				description text
				)`,
	)
	if err != nil {
		panic(err)
	}
	to, name, subj, mess, files := interact()
	auth := smtp.PlainAuth("", "vsevolond@gmail.com", "lmytvuyidgwaulqs", "smtp.gmail.com")
	conf := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", "smtp.gmail.com:465", conf)
	if err != nil {
		panic(err)
	}
	client, err := smtp.NewClient(conn, "smtp.gmail.com")
	if err != nil {
		panic(err)
	}
	if err = client.Auth(auth); err != nil {
		sterr := err.Error()
		updateDB(db, to[0], name, subj, mess, sterr[0:3],
			sterr[4:])
		return
	}
	if err = client.Mail("vsevolond@gmail.com"); err != nil {
		sterr := err.Error()
		updateDB(db, to[0], name, subj, mess, sterr[0:3], sterr[4:])
		return
	}
	if err = client.Rcpt(to[0]); err != nil {
		fmt.Println("Entered email is invalid\nTry again?(Y/N)")
		for {
			ans := input.Gets()
			if ans == "Y" || ans == "y" {
				fmt.Print("Enter email address: ")
				temp := input.Gets()
				to = []string{temp}
				break
			} else if ans == "N" || ans == "n" {
				sterr := err.Error()
				updateDB(db, to[0], name, subj, mess, sterr[0:3],
					sterr[4:])
				return
			} else {
				fmt.Println("Wrong answer")
			}
		}
	}
	writer, err := client.Data()
	if err != nil {
		sterr := err.Error()
		updateDB(db, to[0], name, subj, mess, sterr[0:3],
			sterr[4:])
		return
	}
	mail := createTemplate("vsevolond@gmail.com", to[0], name, subj, mess, files)
	_, err = writer.Write(mail)
	if err != nil {
		sterr := err.Error()
		updateDB(db, to[0], name, subj, mess, sterr[0:3], sterr[4:])
		return
	}
	err = writer.Close()
	if err != nil {
		sterr := err.Error()
		updateDB(db, to[0], name, subj, mess, sterr[0:3], sterr[4:])
		return
	}
	err = client.Quit()
	if err != nil {
		sterr := err.Error()
		updateDB(db, to[0], name, subj, mess, sterr[0:3], sterr[4:])
		return
	}
	fmt.Println("Sended")
	updateDB(db, to[0], name, subj, mess, "250", "success")
}

func createTemplate(from string, to string, name string, subject string, message string, files [][]byte) []byte {
	buf := bytes.NewBuffer(nil)
	hasFiles := len(files) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", to))
	mime := "MIME-version: 1.0\n"
	buf.WriteString(mime)
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if hasFiles {
		buf.WriteString(fmt.Sprintf("Content-Type:multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}
	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\n")
	buf.WriteString(fmt.Sprintf(`
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" 
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
	<html> 
	<head>
	<meta http-equiv="Content-Type" content="text/html; 
	charset=utf-8" >
	<title>Message to %s</title>
	<table cellpadding="0" cellspacing="0" border="0"
	width="100%%" style="background: #f5f5f5; min-width: 320px; font- 
	size: 1px; line-height: normal;">
	<tr>
	<td align="center" valign="top">
	<table cellpadding="0" cellspacing="0" border="0" width="700" 
	class="table700" style="max-width: 700px; min-width: 320px; 
	background: #FFFF00;">
	<tr>
	<td align="center" valign="top">
	<span style="font-family: Arial, Helvetica, sans-serif; font- 
	size: 14px; line-height: 16px;">
	<b>Hello, %s!</b>. This is a message from %s. Read it:
	</span> 
	</td> 
	</tr> 
	<tr>
	<td align="left" valign="top">
	<span style="font-family: Arial, Helvetica, sans-serif; font- 
	size: 14px; line-height: 16px;">
	<i>%s</i>
	</span>
	</td>
	</tr>
	<tr>
	</tr>
	<tr>
	<td align="right" valign="top">
	<span style="font-family: Arial, Helvetica, sans-serif; font- 
	size: 11px; line-height: 16px;">
	<i>With respect, Vsevolod</i>
	</span>
	</td>
	</table>
	</td>
	</tr>
	</table>
	</head>
	<body>
	</body>
	</html>`, to, name, from, message) + "\n")
	if hasFiles {
		for k, v := range files {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType((v))))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition:attachment; filename=%d\n", k))
			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
		buf.WriteString("--")
	}
	return buf.Bytes()
}
