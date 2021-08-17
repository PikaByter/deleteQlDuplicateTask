package main

import (
	"io"
	"log"
	"os"
)

var (
	Info *log.Logger
	Error * log.Logger
)

func init(){
	errFile,err:=os.OpenFile("errors.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err!=nil{
		log.Fatalln("打开error日志文件失败：",err)
	}
	infoFile,err:=os.OpenFile("info.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err!=nil{
		log.Fatalln("打开info日志文件失败：",err)
	}
	Info = log.New(io.MultiWriter(os.Stderr,infoFile),"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr,errFile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)
}



