package util

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "path/filepath"
)

func GenPidFilePath(sname string) string {
    wd, err := os.Getwd()
    if err != nil {
        return ""
    }
    return filepath.Join(wd, "pid", fmt.Sprintf("run.%v.pid", sname))
}

func PidFromFile(filepath string) (error, int) {
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        return nil, 0
    }
    f, err := os.Open(filepath)
    if err != nil {
        return err, 0
    }
    defer f.Close()
    buf := make([]byte, 64)
    n, err := f.Read(buf)
    if err != nil {
        return err, 0
    }
    str := string(buf[:n])
    str = strings.TrimSpace(str)
    pid, err := strconv.Atoi(str)
    if err != nil {
        return err, 0
    }
    return nil, pid
}

func WritePidToFile(filepath string, pid int) error {
    f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
    if err != nil {
        return err
    }
    defer f.Close()

    f.WriteString(fmt.Sprintf("%v", pid))

    return nil
}

func DeletePidFile(filepath string) error {
    return os.Remove(filepath)
}

