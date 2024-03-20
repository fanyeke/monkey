package queue

import (
	"container/list"
	"fmt"
)

// 定义一个队列结构体
type Queue struct {
	list *list.List // 使用 list 包的 List 结构作为队列的底层存储
}

// 创建一个新的队列
func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

// 向队尾添加元素
func (q *Queue) Enqueue(item interface{}) {
	q.list.PushBack(item)
}

// 从队头移除并返回元素，如果队列为空则返回 nil
func (q *Queue) Dequeue() interface{} {
	if q.list.Len() == 0 {
		return nil
	}
	firstItem := q.list.Front().Value
	q.list.Remove(q.list.Front())
	return firstItem
}

// 检查队列是否为空
func (q *Queue) IsEmpty() bool {
	return q.list.Len() == 0
}

func main() {
	queue := NewQueue()

	// 添加元素到队列
	queue.Enqueue("Apple")
	queue.Enqueue("Banana")
	queue.Enqueue("Cherry")

	fmt.Println("Before dequeue:", queue.list)

	// 从队头移除元素
	item := queue.Dequeue()
	fmt.Println("Dequeued item:", item)

	fmt.Println("After dequeue:", queue.list)
}
