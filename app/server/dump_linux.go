package server

import (
	"log"
	"os"
	"syscall"
	"time"
)

func Dump() {
	fileName := time.Now().Format("2006-01-02")
	logFilename := "/var/log/" + fileName + ".log"
	logFile, _ := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	//defer logFile.Close() // DO NOT defer closing the file since this will close it prior to writing the panic output - let the runtime close the file as it cleans up after the panic
	redirectStderr(logFile)
	os.Stderr.Write([]byte(time.Now().Format("2006-01-02 15:04:05") + "\n"))
	//
	//if crashFile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664); err == nil {
	//	crashFile.WriteString(fmt.Sprintf("%v Opened crashfile at %v", os.Getpid(), time.Now()))
	//	os.Stderr = crashFile
	//	syscall.Dup2(int(crashFile.Fd()), 2)
	//}
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	//int(os.Stderr.Fd())
	os.Stderr = f
	err := syscall.Dup2(int(f.Fd()), 2)
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
