package ionservice

import "sync"

type IService interface {
	Status() ServiceStatus
	Start(startSync *sync.WaitGroup)
	Stop()
}

type ServiceStatus int

const (
	Stopped ServiceStatus = iota
	Running
)
