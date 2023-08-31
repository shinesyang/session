package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	说明:
	为什么要自己写一个session中间件
	网上使用最多的session中间件是: github.com/gin-contrib/sessions
	默认情况下，github.com/gin-contrib/sessions 验证调的是两个方法:
	返回header: set-cookie  字段
	获取header: cookie      字段,

	存在的问题:
	1. 在vue3作为前后端分离的相关项目中,无法通过:axios.interceptors.response.use 获取到set-cookie这个字段,
		那么无法判断用户登录之后的界面状态
	2. vue3作为前后端分离的相关项目中使用google调试时,不管是在开发环境还是编译之后的正式环境都无法自动传cookie这个header,
		手动指定时又会有报错提示: Refused to set unsafe header "Cookie"
*/

const (
	defaultKey = "shinesyang/session"
	//defaultSetHeader = "Set-Auth"
	//defaultHeader    = "Auth"
	defaultSetHeader = "Set-Cookie"
	defaultHeader    = "Cookie"
)

type Session interface {
	Session() *session
	Get(key string) interface{}
	Set(key string, value interface{})
	Options(Option) Option
	Save()
}

type session struct {
	name   string
	values map[string]interface{}
	Option
	Store
	request *http.Request
	writer  http.ResponseWriter
	written bool
}

func (s *session) Session() *session {
	if !s.written {
		return s.Store.Get(s)
	}
	return s
}

func (s *session) Options(opt Option) Option {
	s.Option = opt
	return opt
}

func (s *session) Get(key string) interface{} {
	return s.Session().values[key]
}

func (s *session) Set(key string, value interface{}) {
	s.Session().values[key] = value
	s.written = true
}

func (s *session) Delete(key string) {
	delete(s.Session().values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().values {
		s.Delete(key)
	}
}

func (s *session) Save() {
	if s.written {
		s.Session().Store.Save()
		s.written = false
	}
}

func Sessions(name string, store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := &session{name: name, Store: store, request: c.Request, writer: c.Writer}
		c.Set(defaultKey, s)
		c.Next()
	}
}

func Default(c *gin.Context) Session {
	return c.MustGet(defaultKey).(*session)
}
