package wnet

import "winx/wiface"

type BaseRouter struct {
}

func (r *BaseRouter) PreHandle(request wiface.IRequest) {
}

func (r *BaseRouter) Handle(request wiface.IRequest) {
}

func (r *BaseRouter) PostHandle(request wiface.IRequest) {
}
