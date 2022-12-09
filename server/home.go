package server

import (
	"chatroom/global"
	"fmt"
	"html/template"
	"net/http"
)

func homeHandleFunc(w http.ResponseWriter, req *http.Request) {
	tpl, err := template.ParseFiles(global.RootDir + "/template/home.html")
	if err != nil {
		fmt.Fprint(w, "模板解析錯誤！")
		return
	}

	err = tpl.Execute(w, nil)
	if err != nil {
		fmt.Fprint(w, "模板执行錯誤！")
		return
	}
}