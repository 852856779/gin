package app

import (
    "github.com/gin-gonic/gin"
	//"fmt"
    // "testproject/app/common/request"
    // "testproject/app/common/response"
    // "testproject/app/services"
	// "go.mongodb.org/mongo-driver/bson"
	"fmt"
	"time"
	"sync"
	"context"
	"golang.org/x/sync/errgroup"
	"errors"
)



var wg sync.WaitGroup
var rwmu sync.RWMutex
var mu sync.Mutex
func GoroutinesTest(c *gin.Context){

	// 一:创建带有cancel函数的context
	// ctx, cancel := context.WithCancel(context.Background());
	// ctx := context.WithValue(context.Background(), "key", "value")
	//(4)timerCtx 在原本的带有取消方法的基础上,新增了定时
	// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	// defer cancel()


	// data := "Hello, World!"
	// go func(msg string) {
	// 	  fmt.Println(msg)
	// }(data) //如果需要往协程里面传参这么使用


	// ch := make(chan int, 1) // 创建一个带缓冲的channel
	// // ch := make(chan int, 0) // 创建一个无缓冲的channel
	// go func() {
    // 	// 异步任务逻辑
	// 	testData := 40
    // 	ch <- testData // 将结果发送到channel
    // 	// 异步任务逻辑
   	//  	close(ch) // 关闭channel，表示任务完成
	// }()
	// // 在需要的时候从channel接收结果
	// result := <-ch
	// fmt.Println(result)

	// for i := 0; i < 5; i++ {
	// 	wg.Add(1)
	// 	go func(index int) {
	// 		defer wg.Done()
	// 		fmt.Println(index)
	// 		// 异步任务逻辑
	// 	}(i)
	// }
	// // 等待所有协程完成
	// wg.Wait()

	// ch := make(chan int)

	// go func() {
	// 	for i := 0; i < 5; i++ {
	// 		ch <- i // 发送值到通道
	// 	}
	// 	close(ch) // 关闭通道
	// }()
	
	// // 使用range迭代接收通道的值 因为是通道所以range只返回值
	// for val := range ch {
	// 	// 处理接收到的值
	// 	fmt.Println(val)
	// }


	// ch1 := make(chan int)
	// ch2 := make(chan string)
	
	// go func() {
	// 	// 异步任务1逻辑
	// 	result1 := 1;
	// 	ch1 <- result1
	// }()
	
	// go func() {
	// 	// 异步任务2逻辑
	// 	//time.Sleep(2 * time.Second) 
	// 	result2 := "字符";
	// 	ch2 <- result2
	// }()
	
	// // 在主goroutine中等待多个异步任务完成 select如果有多个满足条件的则随机选择一个执行
	// select {
	// case res1 := <-ch1:
	// 	// 处理结果1
	// 	fmt.Println("处理结果1")
	// 	fmt.Println(res1)
	// case res2 := <-ch2:
	// 	// 处理结果2
		
	// 	fmt.Println("处理结果2")
	// 	fmt.Println(res2)
	// }


	// ch := make(chan int)

	// go func() {
	// 	// 异步任务逻辑
	// 	result := 1
	// 	time.Sleep(5 * time.Second)
	// 	ch <- result
	// }()
	
	// // 设置超时时间
	// select {
	// case res := <-ch:
	// 	// 处理结果
	// 	fmt.Println(res)
	// case <-time.After(3 * time.Second):
	// 	// 超时处理
	// 	fmt.Println("协程超时了")
	// }

	//golang 有自己的定时器 php什么垃圾啊
	//tick := time.Tick(1 * time.Second) // 这个定时器手动关不了(单独在一个协程里面一直跑的意思)每秒执行一次操作 time.tick 返回的是一个通道如果循环他没有设置结束条件那将会是一个死循环因为通道每秒都有值返回
// 	tick := time.NewTicker(1 * time.Second) //NewTicker  提供了stop的方法可以手动关闭定时器
// 	fmt.Println(tick)
// 	timeOut := time.After(5 * time.Second);
// 	defer tick.Stop()
// 	for {
// 		select {
// 		case <-timeOut:
// 			// 在5秒后执行操作
// 			// res :=<-tick
// 			// fmt.Println(res)
// 			// tick.Stop()
// 			fmt.Println("停止循环")
// 			return;
		
// 		case res2:=<-tick.C:
// 			// 执行定时操作
// 			fmt.Println(res2)
// 			fmt.Println("定时操作")
// 	}
// }

//互斥锁
	// for i := 0; i < 5; i++ {
	// 	wg.Add(1)
	// 	go func(index int) {
	// 		mu.Lock()
	// 		defer wg.Done()
	// 		fmt.Println(index)
	// 		mu.Unlock()
	// 		// 异步任务逻辑
	// 	}(i)
	// }
	// 等待所有协程完成
	//wg.Wait()
// fmt.Println("结束定时")

	// _, cancel := context.WithCancel(context.Background())
	// wg.Add(1)
	// go f(ctx)
	// fmt.Println("11")
	// time.Sleep(time.Second*5)
	// fmt.Println("22")
	// cancel()
	// wg.Wait()

	//
	// var cond = sync.NewCond(&rwmu) //使用sync.NewCond 需要先传入一个索的实例（例如sync.Mutex、sync.RWMutex等
	// ready := false // cond所等待的条件

	// go func() {
	// 	fmt.Println("doing some work...")
	// 	time.Sleep(time.Second*5)
	// 	ready = true // 不一定加锁，因为只有一个goroutine在写ready
	// 	cond.Broadcast() //()
	// }()

	// cond.L.Lock() // 检查目标条件时先加锁
	// for !ready {
	// 	fmt.Println("wait")
	// 	cond.Wait()  //阻塞等待唤醒后执行  Wait返回并跳出for循环后需要解锁 被Signal或Broadcast方法唤醒，Wait才会返回 如果是多个协程的话需要在wait前加锁(等待必加锁)
	// 	fmt.Println("done")
	// }
	// cond.L.Unlock() //
	// fmt.Println("got work done signal!")
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("doing some work...")
	// 	time.Sleep(time.Second*5)
	// 	ready = true // 不一定加锁，因为只有一个goroutine在写ready
	// 	cond.Broadcast() // 通知多个被阻塞的goroutine
	// }()
	

	// wg.Add(5)
	// for range 5 {
	// 	go func() {
	// 		defer wg.Done()
	// 		cond.L.Lock() //必须 加锁是因为 Wait()方法里面包括了UnLock如果不加锁 会报错的
	// 		for !ready { //循环监听cond.Broadcast()/cond.Signal()
	// 			fmt.Println("wait")
	// 			cond.Wait()
	// 			fmt.Println("done")
	// 		}
	// 		cond.L.Unlock()
	// 		fmt.Println("got work done signal!")
	// 	}()
	// }

	objectPool := sync.Pool{
		New: func() interface{} {
			// 创建新对象
			return &MyObject{}
		},
	}
	// 从对象池获取对象
	obj := objectPool.Get().(*MyObject)
	// 使用对象
	fmt.Println(obj)
	obj.Name = "test";
	fmt.Println(obj);
	// 将对象放回对象池
	objectPool.Put(obj)
}
type MyObject struct {
    // 对象结构
	Name string
}
//上下文
func ContextTest(c *gin.Context){
	// ctx, cancel := context.WithCancel(context.Background());
	// someCondition := false;
	// go func() {
	// 	// 异步任务逻辑
	// 	time.Sleep(time.Second*5)
	// 	someCondition = true;
	// 	if someCondition {
	// 		fmt.Println("取消任务");
	// 		cancel() // 取消任务
	// 	}
	// }()
	
	// // 等待任务完成或取消
	// select {
	// case <-ctx.Done():
	// 	// 任务被取消或超时
	// 	fmt.Println("任务已经被取消");
	// }
	//WithDeadline和WithTimeout 设置的上下文的过期时间 可以不需要设置cancel()
    // ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, _ := context.WithDeadline(context.Background(),  time.Now().Add(time.Second*10))
    // defer cancel()


	select {
    case <-time.After(5 * time.Second): //返回一个在指定的时间间隔后触发的通道
        // 超时处理
		fmt.Println("超时了")
    case <-ctx.Done():
        // 上下文取消处理
		fmt.Println("取消了")
    }
}

func f(ctx context.Context) {
	defer wg.Done() //使用 defer wg.Done() 来保证即使协程中途发生错误或提前退出，wg.Done() 方法也会被调用，从而防止死锁的形象发生
LOOP:
	for {
		fmt.Println("hello")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}

}

//使用errorgroup 可以处理协程如果出错后的逻辑
func errorgo(){
	var eg errgroup.Group
	for i := 0; i < 5; i++ {
		eg.Go(func() error {
			return errors.New("error")
		})
	
		eg.Go(func() error {
			return nil
		})
	}
	
	if err := eg.Wait(); err != nil {
		// 处理错误
		fmt.Println(err);
	}
}