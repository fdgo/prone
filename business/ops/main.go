package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"business/ops/handler"
	"business/ops/model"
	"github.com/dchest/captcha"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	conf model.Config
	addr = flag.String("addr", ":8080", "server listen address")
)

func main() {
	flag.Parse()

	path := filepath.Dir(os.Args[0])
	// 初始化模板
	handler.Start(path)
	static := http.StripPrefix("/static/",
		http.FileServer(http.Dir(filepath.Join(path, "static"))))

	r := mux.NewRouter()
	// 静态文件
	r.PathPrefix("/static/").Handler(static)
	// 404页面
	r.NotFoundHandler = handler.Limit(2, http.NotFound)
	// 加载配置文件
	if err := conf.Load(path); err == nil {
		if err = model.Open(&conf); err != nil {
			log.Panic(err)
		}
	}
	if !model.IsOpen() {
		r.Use(func(h http.Handler) http.Handler {
			if model.IsOpen() {
				return h
			}
			if reflect.ValueOf(h).Pointer() == reflect.ValueOf(static).Pointer() {
				return h
			}
			return handler.Install(path, r)
		})
	}
	// 登录相关
	r.Handle("/", handler.Check(http.HandlerFunc(handler.Home)))
	r.Handle("/login", handler.Limit(2, handler.Login)).Methods(http.MethodPost)
	r.Handle("/captcha/{png}", captcha.Server(120, 35)).Methods(http.MethodGet)

	r.HandleFunc("/login", handler.Login).Methods(http.MethodGet)
	r.HandleFunc("/logout", handler.Logout)
	r.HandleFunc("/password", handler.Password).Methods(http.MethodPost)
	// 后台主页
	s := r.PathPrefix("/").Subrouter()
	// 检查登陆状态
	s.Use(handler.Check)
	// 系统管理
	s.HandleFunc("/users", handler.Users).Methods(http.MethodGet)
	s.HandleFunc("/user/add", handler.UserAdd).Methods(http.MethodPost)
	s.HandleFunc("/user/delete/{id:[0-9]+}", handler.UserDelete)
	s.HandleFunc("/group/{id:[0-9]+}", handler.GroupEdit)
	s.HandleFunc("/group/add", handler.GroupAdd).Methods(http.MethodPost)
	s.HandleFunc("/logs", handler.Logs).Methods(http.MethodGet)
	// 个人中心
	s.HandleFunc("/profile", handler.Profile).Methods(http.MethodGet)
	// 文件上传
	s.HandleFunc("/upload", handler.Upload).Methods(http.MethodPost)

	log.Panic(http.ListenAndServe(*addr, handlers.CustomLoggingHandler(os.Stdout,
		handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(r), handler.WriteLog)))
}
