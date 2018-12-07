package server

import (
	"log"
	"os"
	"syscall"
	"time"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdhandle int32, handle syscall.Handle) error {
	r0, _, e1 := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdhandle), uintptr(handle), 0)
	if r0 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return syscall.EINVAL
	}
	return nil
}

func Dump() {
	fileName := time.Now().Format("2006-01-02")
	logFilename := "c:\\var\\log\\" + fileName + ".log"
	// redirect stdout and stderr to log file
	logFile, _ := os.OpenFile(logFilename, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	//defer logFile.Close() // DO NOT defer closing the file since this will close it prior to writing the panic output - let the runtime close the file as it cleans up after the panic
	redirectStderr(logFile)
	os.Stderr.Write([]byte("\r\n" + time.Now().Format("2006-01-02 15:04:05") + "\r\n"))
}

// redirectStderr to the file passed in
func redirectStderr(f *os.File) {
	err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(f.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
	// SetStdHandle does not affect prior references to stderr
	os.Stderr = f
}
