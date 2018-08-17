package kk

type IChannel interface {

	/**
	 * 发送消息
	 */
	Send(data []byte) error

	/**
	 * 关闭
	 */
	Close()
}
