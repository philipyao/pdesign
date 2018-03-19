package phttp

import (
    "fmt"
    "net/http"
    "encoding/json"
    "encoding/xml"
)

type Response struct {
    code int
    errmsg string

    data []byte

    file string

    cookies []*http.Cookie

    writer http.ResponseWriter
    r *http.Request
}

func makeResponse(w http.ResponseWriter, r *http.Request) *Response {
    return &Response{
        writer: w,
        r: r,
    }
}

func (rsp *Response) Text(text string) {
    rsp.data = []byte(text)
}

func (rsp *Response) Json(pkg interface{}) error {
    buff, err := json.Marshal(pkg)
    if err != nil {
        return err
    }
    rsp.data = buff
    rsp.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
    return nil
}

func (rsp *Response) Xml(pkg interface{}) error {
    buff, err := xml.Marshal(pkg)
    if err != nil {
        return err
    }
    rsp.data = buff
    rsp.writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
    return nil
}

func (rsp *Response) Error(code int, errmsg string) {
    rsp.code = code
    rsp.errmsg = errmsg
}

func (rsp *Response) File(path string) {
    rsp.file = path
}

func (rsp *Response) Cookie(key, value string) {
    cookie := &http.Cookie{
        Name:  key,
        Value: value,
        Path:  "/",
    }
    rsp.cookies = append(rsp.cookies, cookie)
}

func (rsp *Response) flush() {
    fmt.Println("start flushing response")

    // set all cookies to response object
    for _, v := range rsp.cookies {
        fmt.Printf("flush cookie %v\n", v)
        http.SetCookie(rsp.writer, v)
    }
    //none 200 status
    if rsp.code > 0 {
        rsp.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
        rsp.writer.Header().Set("X-Content-Type-Options", "nosniff")
        rsp.writer.WriteHeader(rsp.code)
        fmt.Fprintln(rsp.writer, rsp.errmsg)
        fmt.Printf("flush code %v %v\n", rsp.code, rsp.errmsg)
        return
    }
    //file
    if rsp.file != "" {
        fmt.Printf("flush file %v\n", rsp.file)
        http.ServeFile(rsp.writer, rsp.r, rsp.file)
        return
    }

    //write body
    fmt.Printf("flush data, len %v\n", len(rsp.data))
    rsp.writer.Write(rsp.data)
}