package messagequeue

/*
消息队列是典型的生产者、消费者模型。生产者不断地向消息队列中添加消息，消费者不断地从消息队列中取出消息。
因为消息的生产和消费是异步的，只需要关心消息的发送和接收，不需要关心业务逻辑，实现了生产者和消费者的解耦。
*/

// 用户数据接口
type Message interface {
	Do()
}

// 消息队列
type MessageQueue struct {
	msgChan chan Message
}

// 创建消息队列
func New(size int) *MessageQueue {
	return &MessageQueue{
		msgChan: make(chan Message, size),
	}
}

// 添加消息
func (mq *MessageQueue) Add(msg Message) {
	mq.msgChan <- msg
}

// 启动
func (mq *MessageQueue) Start() {
	go func() {
		// 从消息队列依次取出消息，创建工作协程执行任务
		for msg := range mq.msgChan {
			go msg.Do()
		}
	}()
}

// 停止
func (mq *MessageQueue) Stop() {
	close(mq.msgChan)
}
