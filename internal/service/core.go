package service

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/IonicHealthUsa/ionlog/internal/core/logengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
)

type coreService struct {
	ctx           context.Context
	cancel        context.CancelFunc
	serviceWg     sync.WaitGroup
	serviceStatus ServiceStatus

	logEngine       logengine.ILogger
	rotationService IRotationService

	serviceStatusLock sync.Mutex
}

type ICoreService interface {
	IService
	LogEngine() logengine.ILogger
	CreateRotationService(folder string, maxFolderSize uint, rotation rotationengine.PeriodicRotation)
}

func NewCoreService() ICoreService {
	cs := &coreService{}
	cs.ctx, cs.cancel = context.WithCancel(context.Background())
	cs.logEngine = logengine.NewLogger()
	cs.rotationService = nil // will be set if rotation is enabled by the user
	return cs
}

func (c *coreService) LogEngine() logengine.ILogger {
	return c.logEngine
}

func (c *coreService) CreateRotationService(folder string, maxFolderSize uint, rotation rotationengine.PeriodicRotation) {
	if c.rotationService != nil {
		c.LogEngine().Writer().DeleteWriter(c.rotationService.RotationEngine())
		c.rotationService.Stop()
	}

	c.rotationService = NewRotationService(folder, maxFolderSize, rotation)
	c.LogEngine().Writer().AddWriter(c.rotationService.RotationEngine())
}

// Start starts the logger service, it blocks until the service is stopped
func (c *coreService) Start(startSync *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			ci := runtimeinfo.GetCallerInfo(3)
			fmt.Fprintf(os.Stderr, "logger service panic: '%v' [%v](%v) %v:%v\n", r, ci.Package, ci.Function, ci.File, ci.Line)
		}
	}()

	c.serviceWg.Add(1)
	defer c.serviceWg.Done()

	c.setServiceStatus(Running)
	defer c.setServiceStatus(Stopped)

	if c.rotationService != nil {
		rotateSync := sync.WaitGroup{}
		rotateSync.Add(1)
		go c.rotationService.Start(&rotateSync)
		rotateSync.Wait()
	}

	if startSync != nil {
		startSync.Done()
	}

	c.logEngine.HandleReports(c.ctx)
}

// Stop stops the logger by canceling the context and waiting for the worker to finish
func (c *coreService) Stop() {
	c.cancel()
	c.serviceWg.Wait()
	c.logEngine.FlushReports()

	if c.rotationService != nil {
		c.rotationService.Stop()
	}
}

// Status returns the status of the logger service
func (c *coreService) Status() ServiceStatus {
	c.serviceStatusLock.Lock()
	defer c.serviceStatusLock.Unlock()
	return c.serviceStatus
}

func (c *coreService) setServiceStatus(status ServiceStatus) {
	c.serviceStatusLock.Lock()
	defer c.serviceStatusLock.Unlock()
	c.serviceStatus = status
}
