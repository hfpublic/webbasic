package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var routerFuncs = make(map[string]func(*gin.RouterGroup))

var uncheckRouterFuncs = make(map[string]func(*gin.RouterGroup))

func AddRouterFunc(group string, add func(*gin.RouterGroup)) {
	_, has := routerFuncs[group]
	if has {
		logrus.Fatalf("add same group: %s", group)
	}
	routerFuncs[group] = add
}

func AddUncheckRouterFunc(group string, add func(*gin.RouterGroup)) {
	_, has := uncheckRouterFuncs[group]
	if has {
		logrus.Fatalf("add same group: %s", group)
	}
	uncheckRouterFuncs[group] = add
}

func UseRouter(g *gin.RouterGroup) {
	for key, addFun := range routerFuncs {
		addFun(g.Group(key))
	}
}

func UseUncheckRouter(g *gin.RouterGroup) {
	for key, addFun := range uncheckRouterFuncs {
		addFun(g.Group(key))
	}
}
