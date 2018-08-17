package kk

type IService interface {
	GetContainer() IContainer
	GetId() int64
	SetId(id int64)
	GetPriority() int64
	GetTitle() string
	Send(data []byte) error
	Exit()
}

type Service struct {
	container IContainer
	id        int64
	channel   IChannel
	priority  int64
	count     int64
	title     string
}

func NewService(container IContainer, channel IChannel, priority int64, title string) *Service {
	v := Service{}
	v.id = 0
	v.container = container
	v.channel = channel
	v.priority = priority
	v.count = 0
	v.title = title
	return &v
}

func (S *Service) GetId() int64 {
	return S.id
}

func (S *Service) SetId(id int64) {
	S.id = id
}

func (S *Service) GetTitle() string {
	return S.title
}

func (S *Service) GetPriority() int64 {
	return S.priority - S.count
}

func (S *Service) GetContainer() IContainer {
	return S.container
}

func (S *Service) Send(data []byte) error {
	return S.channel.Send(data)
}

func (S *Service) Exit() {
	if S.channel != nil {
		v := S.channel
		S.channel = nil
		v.Close()
	}
}
