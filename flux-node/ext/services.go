package ext

import (
	"fmt"
	"github.com/bytepowered/flux/flux-node"
	"sync"
)

var (
	serviceNotFound flux.BackendService
	servicesMap     *sync.Map = new(sync.Map)
)

func RegisterBackendServiceById(id string, service flux.BackendService) {
	servicesMap.Store(id, service)
}

// RegisterBackendService store backend service
func RegisterBackendService(service flux.BackendService) {
	id := _ensureServiceID(&service)
	RegisterBackendServiceById(id, service)
}

// BackendServiceById load backend service by serviceId
func BackendServiceById(serviceID string) (flux.BackendService, bool) {
	v, ok := servicesMap.Load(serviceID)
	if ok {
		return v.(flux.BackendService), true
	}
	return serviceNotFound, false
}

// RemoveBackendService remove backend service by serviceId
func RemoveBackendService(serviceID string) {
	servicesMap.Delete(serviceID)
}

// HasBackendService check service exists by service id
func HasBackendService(serviceID string) bool {
	_, ok := servicesMap.Load(serviceID)
	return ok
}

func _ensureServiceID(service *flux.BackendService) string {
	id := service.ServiceId
	if "" == id {
		id = service.Interface + ":" + service.Method
	}
	if len(id) < len("a:b") {
		panic(fmt.Sprintf("BackendService must has an Id, service: %+v", service))
	}
	return id
}
