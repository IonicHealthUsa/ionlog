package service

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
)

type rotationService struct {
	ctx           context.Context
	cancel        context.CancelFunc
	serviceWg     sync.WaitGroup
	serviceStatus ServiceStatus

	rotationEngine rotationengine.IRotationEngine

	serviceStatusLock sync.Mutex
}

type IRotationService interface {
	IService
	RotationEngine() rotationengine.IRotationEngine
}

func NewRotationService(folder string, maxFolderSize uint, rotation rotationengine.PeriodicRotation, opts ...rotationengine.Option) IRotationService {
	rs := &rotationService{}
	rs.ctx, rs.cancel = context.WithCancel(context.Background())
	rs.rotationEngine = rotationengine.NewRotationEngine(folder, maxFolderSize, rotation, opts...)
	return rs
}

func (r *rotationService) RotationEngine() rotationengine.IRotationEngine {
	return r.rotationEngine
}

func (r *rotationService) Start(startSync *sync.WaitGroup) {
	defer func() {
		if rec := recover(); rec != nil {
			ci := runtimeinfo.GetCallerInfo(3)
			fmt.Fprintf(os.Stderr, "[ionlog internal log] rotation service panic: '%v' [%v](%v) %v:%v\n", rec, ci.Package, ci.Function, ci.File, ci.Line)
		}
	}()

	r.serviceWg.Add(1)
	defer r.serviceWg.Done()

	r.setServiceStatus(Running)
	defer r.setServiceStatus(Stopped)

	if startSync != nil {
		startSync.Done()
	}

	<-r.ctx.Done()
}

func (r *rotationService) Stop() {
	r.cancel()
	r.serviceWg.Wait()
	r.rotationEngine.CloseLogFile()
}

func (r *rotationService) Status() ServiceStatus {
	r.serviceStatusLock.Lock()
	defer r.serviceStatusLock.Unlock()
	return r.serviceStatus
}

func (r *rotationService) setServiceStatus(status ServiceStatus) {
	r.serviceStatusLock.Lock()
	defer r.serviceStatusLock.Unlock()
	r.serviceStatus = status
}
