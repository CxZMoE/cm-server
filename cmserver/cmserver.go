package cmserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	// MethodNotExists 请求方法不存在
	MethodNotExists = 1000
	// RouterNotExists 路由不存在
	RouterNotExists = 1001
)

// AppServer 服务器结构体
type AppServer struct {
	RouterList []*Router
	NowRoute   string
}

// Router 路由结构体
type Router struct {
	Name   string
	Method string
	Index  int
	Fun    func(w http.ResponseWriter, r *http.Request)
}

// NewServer 创建新的服务器结构体
func NewServer() *AppServer {
	return &AppServer{}
}

// RegisterRouter 注册路由
func (server *AppServer) RegisterRouter(router string, method string, fun func(w http.ResponseWriter, r *http.Request)) (*Router, bool) {
	r := &Router{Name: router, Index: len(server.RouterList), Fun: fun, Method: method}
	list := append(server.RouterList, r)
	server.RouterList = list
	return r, true
}

// Get 注册Get路由
func (server *AppServer) Get(router string, fun func(c *Context) error, middlewareFunc ...func(w http.ResponseWriter, r *http.Request)) *Router {
	var newFun = func(w http.ResponseWriter, r *http.Request) {
		// Middleware Runs
		for i := len(middlewareFunc) - 1; i >= 0; i-- {
			middlewareFunc[i](w, r)
		}

		// Base Run
		var context = NewContext(r)
		err := fun(context)

		if err != nil {
			log.Println(err.Error())
		}

		// Not status ok
		if context.StatusCode != http.StatusOK {
			http.Error(w, "", int(context.StatusCode))
			return
		}

		// Status ok
		w.Write(context.data)
	}

	r, b := server.RegisterRouter(router, "GET", newFun)
	if b {
		return r
	}
	return nil
}

// Post 注册Post路由
func (server *AppServer) Post(router string, fun func(c *Context) error, middlewareFunc ...func(w http.ResponseWriter, r *http.Request)) *Router {
	var newFun = func(w http.ResponseWriter, r *http.Request) {
		// Middleware Runs
		for i := len(middlewareFunc) - 1; i >= 0; i-- {
			middlewareFunc[i](w, r)
		}

		// Base Run
		var context = NewContext(r)
		err := fun(context)

		if err != nil {
			log.Println(err.Error())
		}

		// Not status ok
		if context.StatusCode != http.StatusOK {
			http.Error(w, "", int(context.StatusCode))
			return
		}

		// Status ok
		w.Write(context.data)
	}

	r, b := server.RegisterRouter(router, "POST", newFun)
	if b {
		return r
	}
	return nil
}

// ListenAndServe 开始服务器并监听端口
func (server *AppServer) ListenAndServe(addr string) string {
	log.Println("[SER] Server listening at:", addr)
	err := http.ListenAndServe(addr, server)
	if err != nil {
		return err.Error()
	}

	return "true"
}

// ListenAndServeTLS 开始服务器并启用证书：请使用443端口
func (server *AppServer) ListenAndServeTLS(addr, certFile, keyFile string) string {
	log.Println("[SER] Server listening at:", addr)
	err := http.ListenAndServeTLS(addr, certFile, keyFile, server)
	if err != nil {
		return err.Error()
	}

	return "true"
}

// Real access after routers
func (server *AppServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routeName := r.URL
	method := r.Method

	index := server.FindRouter(routeName.Path, method)
	switch index {
	case RouterNotExists: // Router not exists
		http.Error(w, fmt.Sprintf("Router %s does not exist.", routeName), http.StatusNotFound)
		return
	case MethodNotExists:
		http.Error(w, fmt.Sprintf("Method %s for router %s is not allowed.", method, routeName), http.StatusMethodNotAllowed)
		return
	default:
		break
	}
	server.RouterList[index].Fun(w, r)
}

// FindRouter 检查路由是否存在在路由表内
func (server *AppServer) FindRouter(name, method string) int {
	var findList []*Router
	for _, v := range server.RouterList {
		if v.Name == name {
			findList = append(findList, v)
		}
	}
	if len(findList) == 0 {
		return RouterNotExists
	}

	for _, v := range findList {
		if method == v.Method {
			return v.Index
		}
	}
	return MethodNotExists
}

// Context 路由上下文结构体
type Context struct {
	data           []byte
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	StatusCode     uint
}

// NewContext 新的路由上下文结构体
func NewContext(r *http.Request) *Context {
	return &Context{data: nil, Request: r, StatusCode: 200}
}

// GetParam 获取Url参数
func (c *Context) GetParam(key string) string {
	requestURL := c.Request.URL.String()
	splits := strings.Split(requestURL, "/")
	endStr := splits[len(splits)-1]
	/* base?key=123&key2=456 */
	l1 := strings.Split(endStr, "?")

	parameters := l1[1]
	/* parameters: key=123&key2=456 */

	pairs := strings.Split(parameters, "&")
	/* ["key=123","key2=456"] */
	//pairsMap := make(map[string]string, len(pairs))
	for _, v := range pairs {
		sp := strings.Split(v, "=")
		//pairsMap[sp[0]] = sp[1]

		if sp[0] == key {
			return sp[1]
		}
	}
	return ""
}

// JSON 发送JSON
func (c *Context) JSON(obj interface{}, statusCode uint) error {
	j, err := json.Marshal(obj)
	if err != nil {
		//log.Println(err)
		return err
	}

	c.data = j
	return nil
}

// String 发送字符串
func (c *Context) String(str string, statusCode uint) error {
	// Assign string to context data
	c.data = []byte(str)
	c.StatusCode = statusCode
	return nil
}

// File 发送文件
func (c *Context) File(f *os.File, statusCode uint) error {
	// Assign file to context data
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	c.data = buf
	c.StatusCode = statusCode

	return nil
}
