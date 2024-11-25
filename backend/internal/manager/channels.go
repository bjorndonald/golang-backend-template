package manager

import (
	"github.com/dustin/go-broadcast"
)

type Message struct {
	UserId string     `json:"userId"`
	Action ActionType `json:"action"`
}

type Listener struct {
	UserId string
	Chan   chan interface{}
}

type ActionType string

var Logout ActionType = "logout"

var ActionTypeMap = map[string]ActionType{
	"logout": Logout,
}

type Manager struct {
	channels map[string]broadcast.Broadcaster
	open     chan *Listener
	close    chan *Listener
	delete   chan string
	messages chan *Message
}

func NewGameManager() *Manager {
	manager := &Manager{
		channels: make(map[string]broadcast.Broadcaster),
		open:     make(chan *Listener, 100),
		close:    make(chan *Listener, 100),
		delete:   make(chan string, 100),
		messages: make(chan *Message, 100),
	}

	go manager.run()
	return manager
}

func (m *Manager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case roomid := <-m.delete:
			m.deleteBroadcast(roomid)
		case message := <-m.messages:
			m.room(message.UserId).Submit(message)
		}
	}
}

func (m *Manager) register(listener *Listener) {
	m.room(listener.UserId).Register(listener.Chan)
}

func (m *Manager) deregister(listener *Listener) {
	m.room(listener.UserId).Unregister(listener.Chan)
	close(listener.Chan)
}

func (m *Manager) deleteBroadcast(userid string) {
	b, ok := m.channels[userid]
	if ok {
		b.Close()
		delete(m.channels, userid)
	}
}

func (m *Manager) room(gameid string) broadcast.Broadcaster {
	b, ok := m.channels[gameid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.channels[gameid] = b
	}
	return b
}

func (m *Manager) OpenListener(userid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		UserId: userid,
		Chan:   listener,
	}
	return listener
}

func (m *Manager) CloseListener(userid string, channel chan interface{}) {
	m.close <- &Listener{
		UserId: userid,
		Chan:   channel,
	}
}

func (m *Manager) DeleteBroadcast(userid string) {
	m.delete <- userid
}

func (m *Manager) Submit(userId string, action ActionType) {
	msg := &Message{
		UserId: userId,
		Action: action,
	}
	m.messages <- msg
}
