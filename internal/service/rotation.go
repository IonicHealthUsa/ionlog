package service

import (
	"context"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
)

type rotationService struct {
	ctx           context.Context
	cancel        context.CancelFunc
	serviceWg     sync.WaitGroup
	serviceStatus ServiceStatus

	rotationEngine rotationengine.IRotationEngine
}

type IRotationService interface {
	IService
	RotationEngine() rotationengine.IRotationEngine
}

func NewRotationService(folder string, maxFolderSize uint, rotation rotationengine.PeriodicRotation) IRotationService {
	rs := &rotationService{}
	rs.ctx, rs.cancel = context.WithCancel(context.Background())
	rs.rotationEngine = rotationengine.NewRotationEngine(folder, maxFolderSize, rotation)
	return rs
}

func (r *rotationService) RotationEngine() rotationengine.IRotationEngine {
	return r.rotationEngine
}

func (r *rotationService) Start(startSync *sync.WaitGroup) {
	r.serviceWg.Add(1)
	defer r.serviceWg.Done()

	r.setServiceStatus(Running)
	defer r.setServiceStatus(Stopped)

	if startSync != nil {
		startSync.Done()
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-r.ctx.Done():
			return

		case <-ticker.C:
			r.rotationEngine.AutoChecks()
		}
	}
}

func (r *rotationService) Stop() {
	r.cancel()
	r.serviceWg.Wait()
	r.rotationEngine.CloseLogFile()
}

func (r *rotationService) Status() ServiceStatus {
	return r.serviceStatus
}

func (r *rotationService) setServiceStatus(status ServiceStatus) {
	r.serviceStatus = status
}
