## 学习大纲
* [一、channel状态](#1)
* [二、无缓冲channel使用](#2)
* [三、有缓存channel使用](#3)
* [四、缓存为1，延迟保证channel](#4)
* [五、context中channel使用](#5)
* [六、goroutine泄露](#6)


## <span id="1">一、channel状态</span>
~~~go
// ** nil 发送或接收都会堵塞
// A channel is in a nil state when it is declared to its zero value
var ch chan string

// A channel can be placed in a nil state by explicitly setting it to nil.
ch = nil

// ** open 允许发送、接收

// A channel is in a open state when it’s made using the built-in function make.
ch := make(chan string)

// ** closed 可以接收，不能发送

// A channel is in a closed state when it’s closed using the built-in function close.
close(ch)
~~~

## <span id="2">二、无缓冲channel使用<span>
    
    
 *  等待任务场景：等待领导做完后，员工等着领导给部署的任务再开始做
 ~~~go
   
  func doWorkAsTask() {
	var ch = make(chan int)

	// goroutine 等待领导任务
	go func() {

		<-ch //堵塞于此，等着 send 后执行【接收发送在发送之前，堵塞于此】
		doEmpWork()
	}()
	doLeaderWork()

	ch <- 1         //领导做完后，发送信号，员工再做
	time.Sleep(1e9) //暂停一下给doyourwork一个执行的机会
	fmt.Println("end")
}
~~~

* 等待结果场景：员工先做完，领导等着员结果再去做。
~~~go

func doWorkAsResult() {
	ch := make(chan int)
	go func() {

		ch <- 1
		doEmpWork()
	}()

	<-ch //堵塞于此，等待员工做完后，领导再做。【接收发送在发送之前，堵塞于此】
	doLeaderWork()
	fmt.Println("end")
}
~~~
## <span id="3">三、有缓存channel使用<span>
    
 *  有缓存，无保证

~~~go
    
func doWorkWithChannel() {
	emp := 20
	ch := make(chan int, emp)
	for i := 0; i < emp; i++ {
		go func(who int) {
			doEmpWork()
			ch <- who
		}(i)

	}

	for emp > 0 {
		fmt.Println(<-ch) //堵塞于此，等待员工做完,领导再做
		doLeaderWork()
		emp--
		fmt.Println("............")
	}
	fmt.Println("end")
}
  ~~~
  
  
  * 有缓存、无保证、缓冲满了会被丢弃
  
  ~~~go
  func doWorkWithSelect() {
	cap := 5
	ch := make(chan int, cap) //如果容量满了，则
	go func() {
		for c := range ch {
			fmt.Println(c)
			doEmpWork()
		}

	}()
	const work = 20
	for i := 0; i < work; i++ {
		select {
		case ch <- i:
			doLeaderWork()
		default:
			fmt.Println("容量慢了，不能再发送了")

		}
	}
	time.Sleep(1e9)
	close(ch) //关闭channel，发送结束的信号给员工
}
  ~~~
## <span id="4">四、缓存为1，延迟保证channel<span>
    
    ~~~go
    func doWorkWithOneChannel() {
	ch := make(chan int, 1)
	go func() {
		for c := range ch {
			fmt.Println(c)
			doEmpWork()
			fmt.Println("........")
		}
	}()
	for i := 0; i < 20; i++ {
		doLeaderWork()
		ch <- i

	}
	close(ch)
}
    ~~~
    
## <span id="5">五、context中channel使用<span>

~~~go
func doWorkWithContext() {

	duration := 300 * time.Microsecond //超时时间设定
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	ch := make(chan int, 1)
	go func() {
		time.Sleep(500 * time.Millisecond) //模拟员工需要300ms完成
		doEmpWork()
		ch <- 1
	}()
	select {
	case c := <-ch:
		fmt.Println("work end:", c)
	case <-ctx.Done():
		fmt.Println("time out ")
	}
}
~~~

**注意**
 若此时不用缓冲channel，超时后，主gourotine退出，开启goroutine无法关闭，造成goroutine泄露。
## <span id="6">六、goroutine泄露<span>
开启goroutine 堵塞，无法退出。
    
>Reference
[ardanlabs](https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html)
