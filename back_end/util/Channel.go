package util

type Broker struct {
	stopCh          chan struct{}
	publishlishCh   chan [2]interface{}
	subscribeChan   chan *BrokerClient
	unSubscribeChan chan *BrokerClient
}

const DEFULT_CHAN_SIZE = 1

type BrokerOptions struct {
	PublishChanSize     int
	SubscribeChanSize   int
	UnSubscribeChanSize int
}

func NewBroker(option *BrokerOptions) *Broker {
	if option == nil {
		option = &BrokerOptions{
			PublishChanSize:     DEFULT_CHAN_SIZE,
			SubscribeChanSize:   1,
			UnSubscribeChanSize: 1,
		}
	}
	return &Broker{
		stopCh:          make(chan struct{}),
		publishlishCh:   make(chan [2]interface{}, option.PublishChanSize),
		subscribeChan:   make(chan *BrokerClient, option.SubscribeChanSize),
		unSubscribeChan: make(chan *BrokerClient, option.UnSubscribeChanSize),
	}
}

const MaximumBroadcastThreads = 1000

var threadLimiter = make(chan struct{}, MaximumBroadcastThreads)

func (b *Broker) Start() {
	subs := map[*BrokerClient]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for msgCh := range subs {
				close(msgCh.C)
			}
			return
		case msgCh := <-b.subscribeChan:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unSubscribeChan:
			delete(subs, msgCh)
		case m := <-b.publishlishCh:
			msg := m[0]
			allExcept := m[1].([]*BrokerClient)
			for msgCh := range subs {
				doTransfer := func(bk *BrokerClient) {
					if allExcept != nil {
						for _, exceptMsgCh := range allExcept {
							if exceptMsgCh == msgCh {
								return
							}
						}
					}
					if msgCh.Filter == nil || msgCh.Filter(msg) {
						// msgCh is buffered, use non-blocking send to protect the broker:
						select {
						case msgCh.C <- msg:
						default:
						}
					}
				}
				threadLimiter <- struct{}{}
				go func(msgCh *BrokerClient) {
					doTransfer(msgCh)
					<-threadLimiter
				}(msgCh)
			}
		}
	}
}

func (b *Broker) Stop() {
	close(b.stopCh)
}

type BrokerClient struct {
	C      chan interface{}
	Filter func(interface{}) bool
}

func (b *Broker) Subscribe(filter func(interface{}) bool) *BrokerClient {
	msgCh := &BrokerClient{C: make(chan interface{}, 5), Filter: filter}
	b.subscribeChan <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh *BrokerClient) {
	b.unSubscribeChan <- msgCh
	close(msgCh.C)
}

func (b *Broker) Publish(msg interface{}, except ...*BrokerClient) {
	b.publishlishCh <- [2]interface{}{msg, except}
}
