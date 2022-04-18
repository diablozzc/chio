package cnet

import "github.com/diablozzc/chio/ciface"

// 实现router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

func (r *BaseRouter) PreHandle(request ciface.IRequest) {

}

func (r *BaseRouter) Handle(request ciface.IRequest) {

}

func (r *BaseRouter) PostHandle(request ciface.IRequest) {

}
