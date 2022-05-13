package models

import (
	"bytes"
	"encoding/base64"
	"gfWeb/library/utils"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"net/smtp"
	"path"
	"strings"
	"time"
)

type sendMail struct {
	user     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}
type MailSmtpData struct {
	User  string    `json:"user" gorm:"primary_key"`
	Pass  string    `json:"pass"`
	Host  string    `json:"host"`
	Port  int       `json:"port"`
	State int       `json:"state"`
	IsAdd int       `json:"isAdd" gorm:"-"`
	auth  smtp.Auth `gorm:"-"`
}

type attachmentData struct {
	name        string
	contentType string
	withFile    bool
}

type messageData struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	contentType string
	attachment  attachmentData
}

// 获取单个设置数据
func GetMailDataOne(user string) (*MailSmtpData, error) {
	data := &MailSmtpData{
		User: user,
	}
	err := Db.Where(data).First(data).Error
	return data, err
}

// 获取已启用的邮件数据
func GetMailStartData() (*MailSmtpData, error) {
	mail := &MailSmtpData{
		State: 1,
	}
	err := Db.Where(mail).First(mail).Error
	return mail, err
}

// 获取邮件数据列表
func GetMailDataList(params *BaseQueryParam) ([]*MailSmtpData, int64) {
	data := make([]*MailSmtpData, 0)
	sortOrder := utils.StrHumpToUnderlineDefault(params.Sort, "state desc,user")
	if params.Order == "descending" {
		sortOrder = sortOrder + " desc"
	}
	var count int64
	err := Db.Model(&MailSmtpData{}).Offset(params.Offset).Limit(params.Limit).Order(sortOrder).Find(&data).Offset(0).Count(&count).Error
	utils.CheckError(err)
	for _, e := range data {
		e.Pass = ""
	}
	return data, count
}

// 更新邮件数据
func UpdateMailData(mail *MailSmtpData) error {
	err := Db.Save(mail).Error
	return err
}

// 删除设置数据
func DeleteMailData(users []string) error {
	err := Db.Where(users).Delete(&MailSmtpData{}).Error
	return err
}

// 发送邮件处理
func SendMailHandle(Title, bodyStr string, to []string) error {
	IsOpenMail := IsSettingOpenDefault(SETTING_DATA_IS_OPEN_MAIL, false)
	if IsOpenMail == false {
		g.Log().Infof("邮件通知关闭!! %s <bodyStr>:%s", Title, bodyStr)
		return gerror.New("邮件通知关闭")
	}
	mail, err := GetMailStartData()
	cc := []string{}
	bcc := []string{}
	if err != nil {
		g.Log().Errorf("发送邮件处理未找到配置数据")
		return gerror.New("发送邮件处理未找到配置数据")
	}
	//mailUser := "3078928990@qq.com"
	//mailPass := "xvbwdcihuncbdeae"
	//mailHost := "smtp.qq.com"
	//formUser := "天合游众<" + mailUser + ">"
	//formUser := mailUser
	//mail := &sendMail{user: mailUser, password: mailPass, host: mailHost, port: "587"}
	isFile := false
	fileName := ""
	//if fileName != "" {
	//	if utils.IsFileExists(fileName) {
	//		g.Log().Errorf("附件文件未找到:%s", fileName)
	//		return gerror.New("附件文件未找到")
	//	}
	//	isFile = true
	//}
	contentType := "text/plain;charset=utf-8"
	//if bodyContentType == "html" {
	//	contentType = "text/" + bodyContentType + ";charset=UTF-8"
	//}
	mailPrefixTitle := g.Cfg().GetString("game.mailPrefixTitle", "")
	Title = mailPrefixTitle + Title
	message := messageData{from: mail.User,
		to:          to,
		cc:          cc,
		bcc:         bcc,
		subject:     Title,
		body:        bodyStr,
		contentType: contentType,
		attachment: attachmentData{
			name:        fileName,
			contentType: getFileType(fileName),
			withFile:    isFile,
		},
	}
	err = mail.Send(message)
	if err != nil {
		g.Log().Errorf("邮件发送失败: %v", err)
		return err
	}
	g.Log().Infof("邮件发送成功:%s", message.to)
	return nil
}

func (mail *MailSmtpData) Auth() {
	mail.auth = smtp.PlainAuth("", mail.User, mail.Pass, mail.Host)
}

func (mail MailSmtpData) Send(message messageData) error {
	mail.Auth()
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	Header["From"] = message.from
	Header["To"] = strings.Join(message.to, ";")
	Header["Cc"] = strings.Join(message.cc, ";")
	Header["Bcc"] = strings.Join(message.bcc, ";")
	Header["Subject"] = message.subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	mail.writeHeader(buffer, Header)

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + message.contentType + "\r\n"
	body += "\r\n" + message.body + "\r\n"
	buffer.WriteString(body)

	if message.attachment.withFile {
		attachment := "\r\n--" + boundary + "\r\n"
		attachment += "Content-Transfer-Encoding:base64\r\n"
		attachment += "Content-Disposition:attachment\r\n"
		attachment += "Content-Type:" + message.attachment.contentType + ";name=\"" + message.attachment.name + "\"\r\n"
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				g.Log().Error("文件错误:%+v", err)
			}
		}()
		mail.writeFile(buffer, message.attachment.name)
	}

	buffer.WriteString("\r\n--" + boundary + "--")
	smtp.SendMail(mail.Host+":"+gconv.String(mail.Port), mail.auth, message.from, message.to, buffer.Bytes())
	return nil
}
func (mail MailSmtpData) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

// read and write the file to buffer
func (mail MailSmtpData) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}

// 文件后缀处理
func getFileType(fileName string) string {
	if fileName == "" {
		return "text/plain"
	}
	fileSuffix := path.Ext(fileName) //获取文件后缀
	switch fileSuffix {
	case ".jpg":
		return "image/jpg"
	case ".html":
		return "text/html"
	case "text":
		return "text/plain"
	default:
		g.Log().Infof("文件后缀处理未处理%s:%s", fileSuffix, fileName)
		return "text/plain"
	}
}
