package main

/*
1.启动自检
2.是否有数据
否：扫描并创建数据
有：载入数据
3.扫描并diff数据
4.启动守护进程

1.检查文件名
2.如果文件名后缀为-init，初始化版本
2.如果文件名后缀为-major，新建主版本
3.如贵文件名后缀为-minor，新建次版本
4.日过文件名后缀为-patch，新建补丁
*/
import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	Name     = "文件版本管理器"
	NameENG  = "File Version Manager"
	Version  = "0.1.0-190605"
	DataFile = "vm.json"
	LogFile  = "vm.log"
)

var (
	logger *log.Logger
	data   *Data
)

func init() { // 初始化
	os.Remove(LogFile) // 删除记录文件（如果有）

	// 设置记录文件
	logFile, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Println(err)
	}

	// 记录文件输出和控制台输出双通
	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, "", log.Lmicroseconds|log.Lshortfile)
}
func main() {
	logger.Println("[HALO]", Name, "[Ver", Version+"]")
	logger.Println("[HALO]", "欢迎使用文件版本管理测试程序！")
	loadData()
	for true {
		tick()
		time.Sleep(1 * time.Second)
	}
}
func loadData() {
	if _, err := os.Stat(DataFile); err == nil {
		// 数据文件存在
		logger.Println("[INFO]", "正在加载数据文件...")
		f, _ := os.Open(DataFile)
		reader := bufio.NewReader(f)
		writer := bufio.NewWriter(f)
		readWriter := bufio.NewReadWriter(reader, writer)
		raw := ReadString(readWriter)
		json.Unmarshal([]byte(raw), &data)
		// logger.Println(string(config.ID))
	} else if os.IsNotExist(err) {
		// 数据文件不存在
		logger.Println("[INFO]", "正在创建数据文件...")
		data = getData(".\\")
		b, _ := json.MarshalIndent(data, "", "    ")
		f, _ := os.Create(DataFile)
		reader := bufio.NewReader(f)
		writer := bufio.NewWriter(f)
		readWriter := bufio.NewReadWriter(reader, writer)
		Write(readWriter, b)
	} else {
		logger.Fatal(err)
	}
}
func tick() {
	newData := getData(".\\")
	for _, newFile := range newData.Files {
		// for _, oldFile := range data.Files {
		// 	if newFile.Path == oldFile.Path {
		// 		continue
		// 	}
		// }
		logger.Println(newFile)
	}
	b, _ := json.MarshalIndent(newData, "", "    ")
	f, _ := os.Create(DataFile)
	reader := bufio.NewReader(f)
	writer := bufio.NewWriter(f)
	readWriter := bufio.NewReadWriter(reader, writer)
	Write(readWriter, b)
	// logger.Println(newData)
	getTimeStamp()
}
func getData(path string) *Data {
	folder := &Folder{
		Path: path,
	}
	data := &Data{
		Time: time.Now(),
	}
	findFolder(folder, data)
	return data
}
func findFolder(root *Folder, data *Data) {
	dir, _ := ioutil.ReadDir(root.Path)
	for _, item := range dir {
		if item.IsDir() {
			folder := &Folder{
				Path: root.Path + item.Name() + "\\",
			}
			findFolder(folder, data)
			root.Folders = append(root.Folders, *folder)
		} else {
			d, _ := ioutil.ReadFile(root.Path + item.Name())
			file := &File{
				Name: item.Name(),
				Path: root.Path + item.Name(),
				Hash: strings.ToUpper(fmt.Sprintf("%x", md5.Sum(d))),
				Modi: item.ModTime(),
			}
			data.Files = append(data.Files, *file)
		}
	}

}
func checkNameType(name string) {}
func parseName()                {}
func getTimeStamp() {
	t, _ := time.Parse(time.Now().String(), "130203")
	logger.Println(t)
}

type Data struct {
	Time  time.Time
	Files []File
}
type Folder struct {
	Path    string
	Files   []File
	Folders []Folder
}
type File struct {
	Path      string
	Modi      time.Time
	Hash      string
	Name      string
	MajorVer  int
	MinorVer  int
	PatchVer  int
	TimeStamp string
}

func Read(readWriter *bufio.ReadWriter) (p []byte) {
	// BUG
	_, err := readWriter.Read(p)
	if err != nil {
		logger.Fatal(err)
	}
	return p
}
func Write(readWriter *bufio.ReadWriter, p []byte) {
	readWriter.Write(p)
	readWriter.Flush()
}
func ReadString(readWriter *bufio.ReadWriter) (str string) {
	raw_msg, _ := readWriter.ReadString('\n')
	msg := strings.Split(raw_msg, "\n")
	return msg[0]
}
func WriteString(readWriter *bufio.ReadWriter, str string) {
	readWriter.WriteString(str + "\n")
	readWriter.Flush()
}
func ReadMap(readWriter *bufio.ReadWriter) (m map[string]interface{}) {
	msg := ReadString(readWriter)
	return Str2Map(msg)
}
func WriteMap(readWriter *bufio.ReadWriter, m map[string]interface{}) {
	WriteString(readWriter, Map2Str(m))
}
func Str2Map(s string) (m map[string]interface{}) {
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		logger.Fatal(err)
	}
	return m
}
func Map2Str(m map[string]interface{}) (s string) {
	b, err := json.Marshal(m)
	if err != nil {
		logger.Fatal(err)
	}
	return string(b)
}
