package phttp

import (
    "errors"
    "sync"
    "time"
    "container/list"
    "crypto/rand"
    "encoding/hex"
    "net/http"
    "net/url"

    "base/log"
)

const (
    SessionIDLength     = 32

    CookieName          = "sessid"
)

type SessionStore struct {
    lock  sync.RWMutex

    sid    string
    expire time.Time
    values  map[interface{}]interface{}
}

func (s *SessionStore) SessionID() string {
    return s.sid
}

func (s *SessionStore) Expire() time.Time {
    s.lock.RLock()
    defer s.lock.RUnlock()

    return s.expire
}

func (s *SessionStore) SetExpire(lifeTime int) {
    s.lock.Lock()
    defer s.lock.Unlock()
    s.expire = time.Now().Add(time.Duration(lifeTime) * time.Second)
}

func (s *SessionStore) Set(key interface{}, value interface{}) {
    s.lock.Lock()
    defer s.lock.Unlock()

    s.values[key] = value
}

func(s *SessionStore) Get(key interface{}) interface{} {
    s.lock.RLock()
    defer s.lock.RUnlock()

    v, ok := s.values[key]
    if ok {
        return v
    }
    return nil
}

func (s *SessionStore) Del(key interface{}) {
    s.lock.Lock()
    defer s.lock.Unlock()

    delete(s.values, key)
}

type SessionMgr struct {
    lock  sync.RWMutex

    Sessions map[string]*list.Element
    List *list.List

    lifeTime    int
}

func NewManager(lifeTime int) *SessionMgr {
    sm := &SessionMgr{
        Sessions: make(map[string]*list.Element),
        List: list.New(),
        lifeTime: lifeTime,
    }
    go sm.GC()
    return sm
}

func (sm *SessionMgr) readSid(r *http.Request) (string, error) {
    cookie, err := r.Cookie(CookieName)
    if err != nil {
        if err == http.ErrNoCookie {
            return "", nil
        }
        return "", err
    }
    if cookie.Value == "" {
        return "", nil
    }

    return url.QueryUnescape(cookie.Value)
}

func (sm *SessionMgr) writeSid(w http.ResponseWriter, sid string) {
    cookie := &http.Cookie{
        Name: CookieName,
        Value: url.QueryEscape(sid),
        Path: "/",
        HttpOnly: false,    //客户端需要通过document.cookie这种非http方式来获取
    }
    http.SetCookie(w, cookie)
}

func (sm *SessionMgr) SessionAttach(w http.ResponseWriter, r *http.Request) (*SessionStore, error) {
    sid , err := sm.readSid(r)
    if err != nil {
        return nil, err
    }
    log.Debug("SessionAttach: sid <%v>", sid)
    //sm.lock.RLock()
    if element, ok := sm.Sessions[sid]; ok {
        //sm.lock.RUnlock()
        log.Debug("session exist")
        go sm.update(sid)
        return element.Value.(*SessionStore), nil
    }

    log.Debug("session not exist, try create one")
    sid, err = sm.newSessionID()
    if err != nil {
        return nil, err
    }
    session := &SessionStore{
        sid: sid,
        expire: time.Now().Add(time.Duration(sm.lifeTime) * time.Second),
        values: make(map[interface{}]interface{}),
    }
    element := sm.List.PushFront(session)
    sm.Sessions[sid] = element
    //sid写入cookie，后面客户端的请求会自动带上这个字段
    sm.writeSid(w, sid)

    log.Debug("create session %v", sid)
    return session, nil
}


func (sm *SessionMgr) SessionDestroy(sid string) {
    sm.lock.Lock()
    defer sm.lock.Unlock()

    if element, ok := sm.Sessions[sid]; ok {
        sm.List.Remove(element)
        delete(sm.Sessions, sid)
    }
}

func (sm *SessionMgr) GC() {
    sm.doGC()

    time.AfterFunc(time.Duration(sm.lifeTime) * time.Second, func() { sm.GC() })
}

func (sm *SessionMgr) doGC() {
    for {
        sm.lock.RLock()
        element := sm.List.Back()
        sm.lock.RUnlock()
        if element != nil {
            if time.Now().After(element.Value.(*SessionStore).Expire()) {
                log.Debug("gc session %+v", element.Value.(*SessionStore))
                sm.lock.Lock()
                sm.List.Remove(element)
                delete(sm.Sessions, element.Value.(*SessionStore).SessionID())
                sm.lock.Unlock()
            } else {
                //队列再往前取Back()，都是没有过期的
                break
            }
        } else {
            break
        }
    }
}

func (sm *SessionMgr) update(sid string) {
    sm.lock.Lock()
    defer sm.lock.Unlock()
    element, ok := sm.Sessions[sid]
    if !ok {
        return
    }
    element.Value.(*SessionStore).SetExpire(sm.lifeTime)
    //最新的放到队列头
    sm.List.MoveToFront(element)
}

func (sm *SessionMgr) newSessionID() (string, error) {
    b := make([]byte, SessionIDLength)
    n, err := rand.Read(b)
    if err != nil || n != len(b) {
        return "", errors.New("error rand.Read")
    }
    return hex.EncodeToString(b), nil
}


