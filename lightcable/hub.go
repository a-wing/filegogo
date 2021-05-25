package lightcable

import (
	"log"
	"sync"

	"github.com/hashicorp/golang-lru"
)

type Hub struct {
	mutex sync.Mutex
	Topic map[string]*Topic
	Cache *lru.Cache
}

func NewHub() *Hub {
	cache, _ := lru.New(1024)
	return &Hub{
		Topic: make(map[string]*Topic),
		Cache: cache,
	}
}

func (this *Hub) Add(name string, topic *Topic) {
	this.mutex.Lock()
	this.Topic[name] = topic
	this.mutex.Unlock()
}

func (this *Hub) Remove(name string) {
	this.mutex.Lock()
	delete(this.Topic, name)
	this.mutex.Unlock()
}

func (this *Hub) Broadcast(name string, msg *Message) {
	this.mutex.Lock()
	if topic := this.Topic[name]; topic != nil {
		log.Println(name, "topic is: ", topic)
		topic.Broadcast(msg)
	}
	this.mutex.Unlock()
}
