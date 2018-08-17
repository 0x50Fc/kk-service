package kk

type ICenter interface {
	GetContainer(name string, c chan bool) IContainer
	Exit()
}

type Center struct {
	containers map[string]IContainer
	c          chan func()
	exit       bool
}

func NewCenter() *Center {
	v := Center{}
	v.containers = map[string]IContainer{}
	v.c = make(chan func(), 20480)
	v.exit = false

	go func() {

		for !v.exit {
			fn := <-v.c
			fn()
		}

		for _, container := range v.containers {
			container.Exit()
		}

		close(v.c)

	}()

	return &v
}

func (C *Center) GetContainer(name string, c chan bool) IContainer {

	if C.exit {
		return nil
	}

	var v IContainer = nil

	C.c <- func() {

		v = C.containers[name]

		if v == nil {
			v = NewContainer(C, name)
			C.containers[name] = v
		}

		c <- true
	}

	<-c

	return v
}

func (C *Center) Exit() {

	if C.exit {
		return
	}

	C.exit = true
	C.c <- func() {}
}
