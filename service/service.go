package service

import "sync"

// Manager holds services
type Manager struct {
	services map[string]*SparrowService
	Active   bool
}

// SparrowService interface for services
type SparrowService interface {
	Start()
	Stop()
}

// AddService add service
func (bge *Manager) AddService(name string, v SparrowService) {
	bge.services[name] = &v
}

// StartAll starts all services
func (bge *Manager) StartAll() {
	bge.Active = true

	var wg sync.WaitGroup

	for _, v := range bge.services {
		wg.Add(1)
		go func(service *SparrowService) {
			defer wg.Done()
			(*service).Start()
		}(v)
	}

	wg.Wait()
}

// StopAll stops all services
func (bge *Manager) StopAll() {
	bge.Active = false

	var wg sync.WaitGroup

	for _, v := range bge.services {
		wg.Add(1)
		go func(service *SparrowService) {
			defer wg.Done()
			(*service).Stop()
		}(v)
	}

	wg.Wait()
}

// NewManager returns new Manager
func NewManager() Manager {
	return Manager{
		services: make(map[string]*SparrowService),
		Active:   false,
	}
}
