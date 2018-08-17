package kk

type IContainer interface {
	GetCenter() ICenter
	GetName() string
	Add(service IService)
	Remove(service IService)
	Get(c chan bool) IService
	Exit()
}

type Container struct {
	center   ICenter
	name     string
	c        chan func()
	exit     bool
	id       int64
	services map[int64]IService
}

func NewContainer(center ICenter, name string) *Container {
	v := Container{}
	v.name = name
	v.center = center
	v.c = make(chan func(), 20480)
	v.exit = false
	v.id = 0
	v.services = map[int64]IService{}

	go func() {

		for !v.exit {
			fn := <-v.c
			fn()
		}

		for _, service := range v.services {
			service.Exit()
		}

		close(v.c)
	}()

	return &v
}

func (C *Container) GetCenter() ICenter {
	return C.center
}

func (C *Container) GetName() string {
	return C.name
}

func (C *Container) Add(service IService) {

	if C.exit {
		return
	}

	C.c <- func() {

		C.id = C.id + 1

		service.SetId(C.id)

		C.services[C.id] = service

	}
}

func (C *Container) Remove(service IService) {

	if C.exit {
		return
	}

	C.c <- func() {
		delete(C.services, service.GetId())
	}
}

func (C *Container) Get(c chan bool) IService {

	if C.exit {
		return nil
	}

	var v IService = nil

	C.c <- func() {

		var priority int64 = 0

		for _, s := range C.services {

			if v == nil || s.GetPriority() > priority {
				v = s
				priority = s.GetPriority()
			}

		}
	}

	<-c

	return v
}

func (C *Container) Exit() {

	if C.exit {
		return
	}

	C.exit = true
	C.c <- func() {}
}
