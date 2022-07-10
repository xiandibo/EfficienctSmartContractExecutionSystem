package utils

import (
	"log"
	"testing"
)

const event_test = "test"

//测试用例
func TestEvent(t *testing.T) {
	dispatcher := NewEventDispatcher()

	listener1 := NewEventListener(firstEventListener)
	listener2 := NewEventListener(secondEventListener)

	dispatcher.AddEventListener(event_test, listener1)
	dispatcher.DispatchEvent(event_test, "hello1")
	//第二个注册
	dispatcher.AddEventListener(event_test, listener2)
	dispatcher.AddEventListener(event_test, listener2)
	dispatcher.DispatchEvent(event_test, "hello2")
	//删除第一个
	dispatcher.RemoveEventListener(event_test, listener1)
	dispatcher.DispatchEvent(event_test, "hello3")
	//删除所有
	dispatcher.RemoveAllEvent()
	dispatcher.DispatchEvent(event_test, "hello4")
}

//第一个监听
func firstEventListener(event Event) {
	log.Println(event.EventType, event.Object, "第一个事件接收函数")
}

//第二个监听
func secondEventListener(event Event) {
	log.Println(event.EventType, event.Object, "第二个事件接收函数")
}
