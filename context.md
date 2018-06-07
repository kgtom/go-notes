

<h2 id="一、背景">一、背景</h2>
<p>在后端api开发中，每一个请求在都有一个对应的 goroutine 去处理。比如数据库和RPC服务。用来处理一个请求的 goroutine 通常需要访问一些与请求特定的数据，比如终端用户的身份认证信息、验证相关的token、请求的截止时间。 当一个请求被取消或超时时，所有用来处理该请求的 goroutine 都应该迅速退出，然后系统才能释放这些 goroutine 占用的资源。</p>
<p>在Google 内部，我们开发了  <code>Context</code>  包，专门用来简化 对于处理单个请求的多个 goroutine 之间与请求域的数据、取消信号、截止时间等相关操作，这些操作可能涉及多个 API 调用。你可以通过  <code>go get golang.org/x/net/context</code>  命令获取这个包。</p>
<h2 id="二、源码分析">二、源码分析</h2>
~~~go
// +build !go1.9

	// A Context carries a deadline, a cancelation signal, and other values across
	//包含 超时、取消信号及传值

	// Context's methods may be called by multiple goroutines simultaneously.
	//被多个goroutine同时使用

	type Context interface {
		// set. Successive calls to Deadline return the same results.
		//Deadline返回context何时会超时。
		Deadline() (deadline time.Time, ok bool)

		// See http://blog.golang.org/pipelines for more examples of how to use

		// a Done channel for cancelation.
		//Done 方法在Context被取消或超时时返回一个close的channel,close的channel可以作为广播通知，告诉给context相关的函数要停止当前工作然后返回。
		Done() <-chan struct{}

		// Err returns a non-nil error value after Done is closed. Err returns

		//Err方法返回context为什么被取消。

		Err() error

		// Use context values only for request-scoped data that transits

		//Value方法对请求作用域传值使用，key-value结构，线程安全的

		Value(key interface{}) interface{}
~~~

<h2 id="三、使用场景">三、使用场景</h2>
<p>使用context实现上下文功能约定需要在你的方法的传入参数的第一个传入一个context.Context类型的变量。<br>
比如:</p>
<ul>
<li>上层需要指定超时的情况: ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)</li>
<li>上层需要主动取消的情况：ctx, cancel := context.WithCancel(ctx)；需要的地方调用cancel()</li>
</ul>
<p>***制定超时：</p>

```go
ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("time out")
				return
			default:
				fmt.Println("default...")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(6 * time.Second)//10s 时间
	fmt.Println("todo....")
	cancel()
	//看一下 "timeout"输出
	time.Sleep(3 * time.Second)
```

<p>***主动取消：</p>

~~~ go

	ctx, cancel := context.WithCancel(context.Background())
	valueCtx:=context.WithValue(ctx,"name","a")//context.WithValue传值
     go DoWork(valueCtx)

	valueCtx=context.WithValue(ctx,"name","b")
	go DoWork(valueCtx)

	valueCtx=context.WithValue(ctx,"name","c")
	go DoWork(valueCtx)

	time.Sleep(6 * time.Second)

	fmt.Println("工作时间到了，通知a、b、c 停止工作！")

	//主动取消
	cancel()

	// 为了看一下DoWork输出，实质检查goroutine是否停止

	time.Sleep(3 * time.Second)
	
	
	func DoWork(ctx context.Context) {

	for {

		select {

		case <-ctx.Done():

			fmt.Println(ctx.Value("name"), "主动停止 work")

			return

		default:

			fmt.Println(ctx.Value("name"), "work time 2s:...")

			time.Sleep(2 * time.Second)

		}

	}

}
~~~
<h2 id="四、继承-context">四、继承 context</h2>
<p>context 包提供了一些函数，协助用户从现有的  <code>Context</code>  对象创建新的  <code>Context</code>  对象。<br>
这些  <code>Context</code>  对象形成一棵树：当一个  <code>Context</code>  对象被取消时，继承自它的所有  <code>Context</code>  都会被取消。</p>
<p><code>Background</code>  是所有  <code>Context</code>  对象树的根，它不能被取消。它的声明如下：</p>
<pre><code>// Background returns an empty Context. It is never canceled, has no deadline,
// and has no values. Background is typically used in main, init, and tests,
// and as the top-level `Context` for incoming requests.
func Background() Context
</code></pre>
<p><code>WithCancel</code>  和  <code>WithTimeout</code>  函数 会返回继承的  <code>Context</code>  对象， 这些对象可以比它们的父  <code>Context</code>  更早地取消。</p>
<p>当请求处理函数返回时，与该请求关联的  <code>Context</code>  会被取消。 当使用多个副本发送请求时，可以使用  <code>WithCancel</code>取消多余的请求。  <code>WithTimeout</code>  在设置对后端服务器请求截止时间时非常有用。 下面是这三个函数的声明：</p>
<pre><code>// WithCancel returns a copy of parent whose Done channel is closed as soon as
// parent.Done is closed or cancel is called.
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)

// A CancelFunc cancels a Context.
type CancelFunc func()

// WithTimeout returns a copy of parent whose Done channel is closed as soon as
// parent.Done is closed, cancel is called, or timeout elapses. The new
// Context's Deadline is the sooner of now+timeout and the parent's deadline, if
// any. If the timer is still running, the cancel function releases its
// resources.
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
</code></pre>
<p><code>WithValue</code>  函数能够将请求作用域的数据与  <code>Context</code>  对象建立关系。声明如下：</p>
<pre><code>// WithValue returns a copy of parent whose Value method returns val for key.
func WithValue(parent Context, key interface{}, val interface{}) Context
</code></pre>
<p>当然，想要知道  <code>Context</code>  包是如何工作的，最好的方法是看一个栗子。</p>
<h2 id="五、-context-使用原则">五、 Context 使用原则</h2>
<ol>
<li>不要把Context放在结构体中，要以参数的方式传递</li>
<li>以Context作为参数的函数方法，应该把Context作为第一个参数，放在第一位。</li>
<li>给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO</li>
<li>Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递</li>
<li>Context是线程安全的，可以放心的在多个goroutine中传递</li>
</ol>
<blockquote>
<p>reference:<br>
<a href="https://github.com/golang/net/blob/master/context/pre_go19.go">https://github.com/golang/net/blob/master/context/pre_go19.go</a><br>
<a href="https://blog.csdn.net/liangzhiyang/article/details/65633960">https://blog.csdn.net/liangzhiyang/article/details/65633960</a><br>
<a href="https://segmentfault.com/a/1190000006744213">https://segmentfault.com/a/1190000006744213</a><br>
<a href="http://www.flysnow.org/2017/05/12/go-in-action-go-context.html">http://www.flysnow.org/2017/05/12/go-in-action-go-context.html</a></p>
</blockquote>

