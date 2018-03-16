package phttp

import (
    "strings"
    "os"
    "fmt"
    "path/filepath"
)

type static struct {
    serves  map[string]string
}

func (s *static) serve(prefix, dir string) error {
    if s.serves == nil {
        s.serves = make(map[string]string)
    }
    path, err := filepath.Abs(dir)
    if err != nil {
        return fmt.Errorf("cannot make absolute path for " + dir)
    }
    if !existDir(path) {
        return fmt.Errorf("path %v not exist", path)
    }
    s.serves[prefix] = path
    return nil
}

func (s *static) match(url string) string {
    if s.serves == nil {
        return ""
    }
    if url == "/" {
        url = "/index.html"
    }

    for prefix, dir := range s.serves {
        if strings.HasPrefix(url, prefix) {
            //文件路径
            path := strings.Replace(url, prefix, dir, 1)
            //如果文件存在，则返回
            if existFile(path) {
                return path
            }
        }
    }
    return ""
}

//================================================
func existDir(path string) bool {
    stat, err := os.Stat(path)
    if err == nil {
        return stat.IsDir()
    }
    return false
}

func existFile(path string) bool {
    stat, err := os.Stat(path)
    if err == nil {
        return !stat.IsDir()
    }
    return false
}
