package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main(){
	const defaultDatPath = "C:\\Users\\wu869\\Documents\\WeChat Files\\WU869022133\\FileStorage\\Image\\2019-11"
	const defaultJpgPath = "d:\\decodeImage"

	fmt.Printf("\n\n\t\t## 批量转换微信 dat 文件为 jpg ##\n")
	fmt.Printf("\n\t\t代码改编自 github@Seraphli/wic\n")

	fmt.Printf("\n操作步骤:\n")
	fmt.Printf("\n1. 把目录中的 color_sheet.jpg 图片发给一个微信好友;\n")
	fmt.Printf("\n2. 在 WeChat Files\\你的微信id\\Data 目录下按创建时间倒序排列文件, 会看到产生了一个刚刚创建的 dat 文件, 把这个文件复制到当前目录下, 重命名为 color_trans.dat;\n")
	fmt.Printf("\n3. 完成上面两步之后, 按以下提示操作.\n")

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n输入地址 或者 拖动微信 Data 文件目录到窗口:\n( 直接回车则使用当前目录下的 dat 目录 )\n")
	datPath, err := reader.ReadString('\n')
	CheckErr(err)
	datPath, err = CheckInputIsPath(datPath, defaultDatPath, false)

	fmt.Printf("\n输入地址 或者 拖动图片保存目录到窗口:\n( 直接回车则使用当前目录下的 jpg 目录 )\n")
	jpgPath, err := reader.ReadString('\n')
	CheckErr(err)
	jpgPath, err = CheckInputIsPath(jpgPath, defaultJpgPath, true)

	Dat2Jpg(datPath, jpgPath)

	fmt.Printf("\n按返回键退出...\n")
	fmt.Scanln()
}
/**
	根据文件头两个字节判断文件类型和加密byte数据
 */
func guessEncoding(data []byte) (string,byte){
	headers := map[string][]byte{"jpg":[]byte{0xff,0xd8},"png":[]byte{0x89,0x50},"gif":[]byte{0x47,0x49}}
	for k,v := range headers{
		headerCode,checkCode := v[0],v[1]
		magic := doMagic(headerCode,data)
		code := decode(magic,data[:2])
		if checkCode == code[1]{
			return k,magic
		}
	}
	fmt.Println("guess encoding err")
	return "",0x00
}

func doMagic(headerCode byte,data []byte) byte{
	if(len(data) > 0){
		return headerCode ^ data[0]
	}
	return 0x00
}

func decode(magic byte,data []byte)[]byte{
	newData := []byte{}
	for _,b := range data{
		newData = append(newData,magic ^ b)
	}
	return newData
}
func Dat2Jpg(datPath string, jpgPath string){
	files,err := ioutil.ReadDir(datPath)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	for _,curFile  := range files{
		if !curFile.IsDir(){
			decodeFile(datPath+string(filepath.Separator)+curFile.Name(),jpgPath)
		}
	}
}

func decodeFile(datPath string, jpgPath string){
	data,err:= ioutil.ReadFile(datPath)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	ext,magic := guessEncoding(data)
	newFile,err := os.Create(jpgPath+string(filepath.Separator)+filepath.Base(datPath)+"."+ext)
	defer newFile.Close()
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	_,err = newFile.Write(decode(magic,data))
	if err != nil{
		fmt.Println(err.Error())
		return
	}
}
func CheckErr(err error){
	if err != nil {
		log.Printf("\n%T\n%s\n%#v\n", err, err, err)
	}
}

func CheckInputIsPath(input string, defaultValue string, createPath bool)(path string, err error){
	input = strings.Trim(input, "\n")
	input = strings.TrimSpace(input)
	if len(input)>2 {
		if c:=input[len(input)-1]; input[0]==c && (c=='"'||c=='\'') {
			input = input[1:len(input)-1]
		}
	}
	if len(input)>0 {
		fmt.Printf("路径符合要求, 使用路径 %v\n", input)
		path = input
		err = nil
	} else {
		fmt.Printf("路径不符合要求, 使用默认路径 %v\n", defaultValue)
		path = defaultValue
		err = errors.New("路径为空")
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("目录不存在, ")
		if createPath {
			err = os.MkdirAll(path, 0644)
			CheckErr(err)
			if err==nil {
				fmt.Printf("创建目录 %v 成功\n\n", path)
			} else {
				fmt.Printf("创建目录 %v 失败\n\n", path)
				err = errors.New("创建目录 %v 失败")
			}
		} else {
			fmt.Printf("使用默认路径: %v\n\n", defaultValue)
			return defaultValue, errors.New("路径不存在")
		}
	}

	return path, err
}
