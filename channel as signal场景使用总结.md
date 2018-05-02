---


---

<h2 id="一、状态">一、状态</h2>
<p>一个 channel 的行为直接被它当前的状态所影响。一个channel 的状态是：<strong>nil</strong>，<strong>open</strong>  或  <strong>closed</strong>。</p>
<p>下面的清单2展示了怎样声明或把一个 channel放进这三个状态。</p>
<p><strong>清单2</strong></p>
<pre><code>// ** nil channel

// A channel is in a nil state when it is declared to its zero value
var ch chan string

// A channel can be placed in a nil state by explicitly setting it to nil.
ch = nil


// ** open channel

// A channel is in a open state when it’s made using the built-in function make.
ch := make(chan string)    


// ** closed channel

// A channel is in a closed state when it’s closed using the built-in function close.
close(ch)
</code></pre>
<p>状态决定了怎样<strong>send</strong>（发送）和<strong>receive</strong>（接收）操作行为。</p>
<p>信号通过一个 channel 发送和接收。不要说读和写，因为 channels 不执行 I/O。</p>
<p><strong>图2</strong></p>
<p><img src="https://segmentfault.com/img/bV83F2?w=1666&amp;h=544" alt="clipboard.png" title="clipboard.png"></p>
<p>当一个 channel 是  <strong>nil</strong>  状态，任何试图在 channel 的发送或接收都将会被阻塞。当一个 channel 是在  <strong>open</strong>  状态，信号可以被发送和接收。当一个 channel 被置为  <strong>closed</strong>  状态，信号将不在被发送，但是依然可以接收信号。</p>
<h2 id="场景">场景</h2>
<p>有了这些特性，更进一步理解它们在实践中怎样工作的最好方式就是运行一系列的代码场景。当我在读写 channel 基础代码的时候，我喜欢把goroutines想像成人。这个形象对我非常有帮助，我将把它用作下面的辅助工具。</p>
<h3 id="有数据信号---保证---无缓冲-channels">有数据信号 - 保证 - 无缓冲 Channels</h3>
<p>当你需要知道一个被发送的信号已经被接收的时候，有两种情况需要考虑。它们是  <strong>等待任务</strong>和<strong>等待结果</strong>。</p>
<h4 id="场景1---等待任务">场景1 - 等待任务</h4>
<p>考虑一下作为一名经理，需要雇佣一名新员工。在本场景中，你想你的新员工执行一个任务，但是他们需要等待直到你准备好。这是因为在他们开始前你需要递给他们一份报告。</p>
<p><strong>清单5</strong></p>
<p><a href="https://play.golang.org/p/BnVEHRCcdh">在线演示地址</a></p>
<pre><code>01 func waitForTask() {
02     ch := make(chan string)
03
04     go func() {
05         p := &lt;-ch
06
07         // Employee performs work here.
08
09         // Employee is done and free to go.
10     }()
11
12     time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
13
14     ch &lt;- "paper"
15 }
</code></pre>
<p>在清单5的第2行，一个带有属性的无缓冲channel被创建，<code>string</code>  数据将与信号一起被发送。在第4行，一名员工被雇佣并在开始工作前，被告诉等待你的信号【在第5行】。第5行是一个 channel 接收，引起员工阻塞直到等到你发送的报告。一旦报告被员工接收，员工将执行工作并在完成的时候可以离开。</p>
<p>你作为经理正在并发的与你的员工工作。因此在第4行你雇佣员工之后，你发现你自己需要做什么来解锁并且发信号给员工（第12行）。值得注意的是，不知道要花费多长的时间来准备这份报告（paper）。</p>
<p>最终你准备好给员工发信号，在第14行，你执行一个有数据信号，数据就是那份报告。由于一个无缓冲的channel被使用，你得到一个保证就是一旦你操作完成，员工就已经接收到了这份报告。接收发生在发送之前。</p>
<p>技术上你所知道的一切就是在你的channel发送操作完成的同时员工接收到了这份报告。在两个channel操作之后，调度器可以选择执行它想要执行的任何语句。下一行被执行的代码是被你还是员工是不确定的。这意味着使用print语句会欺骗你关于事件的执行顺序。</p>
<h4 id="场景2---等待结果">场景2 - 等待结果</h4>
<p>在下一个场景中，事情是相反的。这时你想你的员工一被雇佣就立即执行他们的任务。然后你需要等待他们工作的结果。你需要等待是因为在你继续前你需要他们发来的报告。</p>
<p><strong>清单6</strong><br>
<a href="https://play.golang.org/p/VFAWHxIQTP">在线演示地址</a></p>
<pre><code>01 func waitForResult() {
02     ch := make(chan string)
03
04     go func() {
05         time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
06
07         ch &lt;- "paper"
08
09         // Employee is done and free to go.
10     }()
11
12     p := &lt;-ch
13 }
</code></pre>
<h4 id="成本收益">成本/收益</h4>
<p>无缓冲 channel 提供了信号被发送就会被接收的保证，这很好，但是没有任何东西是没有代价的。这个成本就是保证是未知的延迟。在<strong>等待任务</strong>场景中，员工不知道你要花费多长时间发送你的报告。在<strong>等待结果</strong>场景中，你不知道员工会花费多长时间把报告发送给你。</p>
<p>在以上两个场景中，未知的延迟是我们必须面对的，因为它需要保证。没有这种保证行为，逻辑就不会起作用。</p>
<h3 id="有数据信号---无保证---缓冲-channels--1">有数据信号 - 无保证 - 缓冲 Channels &gt; 1</h3>
<h4 id="场景1---扇出（fan-out）">场景1 - 扇出（Fan Out）</h4>
<p>扇出模式允许你抛出明确定义数量的员工在同时工作的问题上。由于你每个任务都有一个员工，你很明确的知道你会接收多少个报告。你可能需要确保你的盒子有适量的空间来接收所有的报告。这就是你员工的收益，不需要等待你来提交他们的报告。但是他们确实需要轮流把报告放进你的盒子，如果他们几乎同一时间到达盒子。</p>
<p>再次假设你是经理，但是这次你雇佣一个团队的员工，你有一个单独的任务，你想每个员工都执行它。作为每个单独的员工完成他们的任务，他们需要给你提供一张报告放进你桌子上的盒子里面。</p>
<p><strong>清单7</strong><br>
<a href="https://play.golang.org/p/8HIt2sabs_">演示地址</a></p>
<pre><code>01 func fanOut() {
02     emps := 20
03     ch := make(chan string, emps)
04
05     for e := 0; e &lt; emps; e++ {
06         go func() {
07             time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
08             ch &lt;- "paper"
09         }()
10     }
11
12     for emps &gt; 0 {
13         p := &lt;-ch
14         fmt.Println(p)
15         emps--
16     }
17 }
</code></pre>
<p>在清单7的第3行，一个带有属性的有缓冲channel被创建，<code>string</code>  数据将与信号一起被发送。这时，由于在第2行声明的  <code>emps</code>  变量，将创建有 20个缓冲的 channel。</p>
<p>在第5行和第10行之间，20 个员工被雇佣，并且他们立即开始工作。在第7行你不知道每个员工将花费多长时间。这时在第8行，员工发送他们的报告，但这一次发送不会阻塞等待接收。因为在盒子里为每位员工准备的空间，在 channel 上的发送仅仅与其他在同一时间想发送他们报告的员工竞争。</p>
<p>在 12 行和16行之间的代码全部是你的操作。在这里你等待20个员工来完成他们的工作并且发送报告。在12行，你在一个循环中，在 13 行你被阻塞在一个 channel 等待接收你的报告。一旦报告接收完成，报告在14被打印，并且本地的计数器变量被消耗来表明一个员工意见完成了他的工作。</p>
<h4 id="场景2---drop">场景2 - Drop</h4>
<p>Drop模式允许你在你的员工在满负荷的时候丢掉工作。这有利于继续接受客户端的工作，并且从不施加压力或者是这项工作可接受的延迟。这里的关键是知道你什么时候是满负荷的，因此你不承担或过度承诺你将尝试完成的工作量。通常集成测试或度量可以帮助你确定这个数字。</p>
<p>假设你是经理，你雇佣了单个员工来完成工作。你有一个单独的任务想员工去执行。当员工完成他们任务时，你不在乎知道他们已经完成了。最重要的是你能或不能把新工作放入盒子。如果你不能执行发送，这时你知道你的盒子满了并且员工是满负荷的。这时候，新工作需要丢弃以便让事情继续进行。</p>
<p><strong>清单8</strong><br>
<a href="https://play.golang.org/p/PhFUN5itiv">演示地址</a></p>
<pre><code>01 func selectDrop() {
02     const cap = 5
03     ch := make(chan string, cap)
04
05     go func() {
06         for p := range ch {
07             fmt.Println("employee : received :", p)
08         }
09     }()
10
11     const work = 20
12     for w := 0; w &lt; work; w++ {
13         select {
14             case ch &lt;- "paper":
15                 fmt.Println("manager : send ack")
16             default:
17                 fmt.Println("manager : drop")
18         }
19     }
20
21     close(ch)
22 }
</code></pre>
<p>在清单8的第3行，一个有属性的有缓冲 channel 被创建，<code>string</code>  数据将与信号一起被发送。由于在第2行声明的<code>cap</code>  常量，这时创建了有5个缓冲的 channel。</p>
<p>从第5行到第9行，一个单独的员工被雇佣来处理工作，一个  <code>for range</code>被用于循环处理 channel 的接收。每次一份报告被接收，在第7行被处理。</p>
<p>在第11行和19行之间，你尝试发送20分报告给你的员工。这时一个  <code>select</code>语句在第14行的第一个<code>case</code>被用于执行发送。因为<code>default</code>从句被用于第16行的<code>select</code>语句。如果发送被堵塞，是因为缓冲中没有多余的空间，通过执行第17行发送被丢弃。</p>
<p>最后在第21行，内建函数<code>close</code>被调用来关闭channel。这将发送没有数据的信号给员工表明他们已经完成，并且一旦他们完成分派给他们的工作可以立即离开。</p>
<h4 id="成本收益-1">成本/收益</h4>
<p>有缓冲的 channel 缓冲大于1提供无保证发送的信号被接收到。离开保证是有好处的，在两个goroutine之间通信可以降低或者是没有延迟。在<strong>扇出</strong>场景，这有一个有缓冲的空间用于存放员工将被发送的报告。在<strong>Drop</strong>场景，缓冲是测量能力的，如果容量满，工作被丢弃以便工作继续。</p>
<p>在两个选择中，这种缺乏保证是我们必须面对的，因为延迟降低非常重要。0到最小延迟的要求不会给系统的整体逻辑造成问题。</p>
<h3 id="有数据信号---延迟保证--缓冲1的channel">有数据信号 - 延迟保证- 缓冲1的channel</h3>
<h4 id="场景1---等待任务-1">场景1 - 等待任务</h4>
<p><strong>清单9</strong><br>
<a href="https://play.golang.org/p/4pcuKCcAK3">演示地址</a></p>
<pre><code>01 func waitForTasks() {
02     ch := make(chan string, 1)
03
04     go func() {
05         for p := range ch {
06             fmt.Println("employee : working :", p)
07         }
08     }()
09
10     const work = 10
11     for w := 0; w &lt; work; w++ {
12         ch &lt;- "paper"
13     }
14
15     close(ch)
16 }
</code></pre>
<p>在清单9的第2行，一个带有属性的一个缓冲大小的 channel 被创建，<code>string</code>  数据将与信号一起被发送。在第4行和第8行之间，一个员工被雇佣来处理工作。<code>for range</code>被用于循环处理 channel 的接收。在第6行每次一份报告被接收就被处理。</p>
<p>在第10行和13行之间，你开始发送你的任务给员工。如果你的员工可以跑的和你发送的一样快，你们之间的延迟会降低。但是每次发送你成功执行，你需要保证你提交的最后一份工作正在被进行。</p>
<p>在最后的第15行，内建函数<code>close</code> 被调用关闭channel，这将会发送无数据信号给员工告知他们工作已经完成，可以离开了。尽管如此，你提交的最后一份工作将在  <code>for range</code>中断前被接收。</p>
<h3 id="无数据信号---context">无数据信号 - Context</h3>
<p>在最后这个场景中，你将看到从  <code>Context</code>  包中使用  <code>Context</code> 值怎样取消一个正在运行的goroutine。这所有的工作是通过改变一个已经关闭的无缓冲channel来执行一个无数据信号。</p>
<p>最后一次你是经理，你雇佣了一个单独的员工来完成工作，这次你不会等待员工未知的时间完成他的工作。你分配了一个截止时间，如果你的员工没有按时完成工作，你将不会等待。</p>
<p><strong>清单10</strong><br>
<a href="https://play.golang.org/p/6GQbN5Z7vC">演示地址</a></p>
<pre><code>01 func withTimeout() {
02     duration := 50 * time.Millisecond
03
04     ctx, cancel := context.WithTimeout(context.Background(), duration)
05     defer cancel()
06
07     ch := make(chan string, 1)
08
09     go func() {
10         time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
11         ch &lt;- "paper"
12     }()
13
14     select {
15     case p := &lt;-ch:
16         fmt.Println("work complete", p)
17
18     case &lt;-ctx.Done():
19         fmt.Println("moving on")
20     }
21 }
</code></pre>
<p>在清单10的第2行，一个时间值被声明，它代表了员工将花费多长时间完成他们的工作。这个值被用在第4行来创建一个50毫秒超时的  <code>context.Context</code> 值。<code>context</code> 包的 <code>WithTimeout</code> 函数返回一个  <code>Context</code>  值和一个取消函数。</p>
<p><code>context</code>包创建一个goroutine，一旦时间值到期，将关闭与<code>Context</code> 值关联的无缓冲channels。不管事情如何发生，你需要负责调用<code>cancel</code> 函数。这将清理被<code>Context</code>创建的东西。<code>cancel</code>被调用不止一次是可以的。</p>
<p>在第5行，一旦函数中断，<code>cancel</code>函数被 deferred 执行。在第7行，1个缓冲的channels被创建，它被用于被员工发送他们工作的结果给你。在第09行和12行，员工被雇佣兵立即投入工作，你不需要指定员工花费多长时间完成他们的工作。</p>
<p>在第14行和20行之间，你使用  <code>select</code> 语句来在两个channels接收。在第15行的接收，你等待员工发送他们的结果。在第18行的接收，你等待看<code>context</code> 包是否正在发送信号50毫秒的时间到了。无论你首先收到哪个信号，都将有一个被处理。</p>
<p>这个算法的一个重要方面是使用一个缓冲的channels。如果员工没有按时完成，你将离开而不会给员工任何通知。对于员工而言，在第11行他将一直发送他的报告，你在或者不在那里接收，他都是盲目的。如果你使用一个无缓冲channels，如果你离开，员工将一直阻塞在那尝试你给发送报告。这会引起goroutine泄漏。因此一个缓冲的channels用来防止这个问题发生。</p>
<h2 id="总结">总结</h2>
<p>当使用 channels（或并发） 时，在保证，channel状态和发送过程中信号属性是非常重要的。它们将帮助你实现你并发程序需要的更好的行为以及你写的算法。它们将帮助你找出bug和闻出潜在的坏代码。</p>
<p>在本文中，我分享了一些程序示例来展示信号属性工作在不同的场景中。凡事都有例外，但是这些模式是非常良好的开端。</p>
<p>作为总结回顾下这些要点，何时，如何有效地思考和使用channels：</p>
<h3 id="语言机制">语言机制</h3>
<ul>
<li>
<p>使用 channels 来编排和协作 goroutines：</p>
<ul>
<li>
<p>关注信号属性而不是数据共享</p>
</li>
<li>
<p>有数据信号和无数据信号</p>
</li>
<li>
<p>询问它们用于同步访问共享数据的用途</p>
<ul>
<li>有些情况下，对于这个问题，通道可以更简单一些，但是最初的问题是。</li>
</ul>
</li>
</ul>
</li>
<li>
<p>无缓冲 channels：</p>
<ul>
<li>接收发生在发送之前</li>
<li>收益：100%保证信号被接收</li>
<li>成本：未知的延迟，不知道信号什么时候将被接收。</li>
</ul>
</li>
<li>
<p>有缓冲 channels：</p>
<ul>
<li>
<p>发送发生在接收之前。</p>
</li>
<li>
<p>收益：降低信号之间的阻塞延迟。</p>
</li>
<li>
<p>成本：不保证信号什么时候被接收。</p>
<ul>
<li>缓冲越大，保证越少。</li>
<li>缓冲为1可以给你一个延迟发送保证。</li>
</ul>
</li>
</ul>
</li>
<li>
<p>关闭的 channels：</p>
<ul>
<li>关闭发生在接收之前（像缓冲）。</li>
<li>无数据信号。</li>
<li>完美的信号取消或截止。</li>
</ul>
</li>
<li>
<p>nil channels：</p>
<ul>
<li>发送和接收都阻塞。</li>
<li>关闭信号。</li>
<li>完美的速度限制或短时停工。</li>
</ul>
</li>
</ul>
<h3 id="设计哲学">设计哲学</h3>
<ul>
<li>
<p>如果在 channel上任何给定的发送能引起发送 goroutine 阻塞：</p>
<ul>
<li>
<p>不允许使用大于1的缓冲channels。</p>
<ul>
<li>缓冲大于1必须有原因/测量。</li>
</ul>
</li>
<li>
<p>必须知道当发送 goroutine阻塞的时候发生了什么。</p>
</li>
</ul>
</li>
<li>
<p>如果在 channel 上任何给定的发送不会引起发送阻塞：</p>
<ul>
<li>
<p>每个发送必须有确切的缓冲数字。</p>
<ul>
<li>扇出模式。</li>
</ul>
</li>
<li>
<p>有缓冲测量最大的容量。</p>
<ul>
<li>Drop 模式。</li>
</ul>
</li>
</ul>
</li>
<li>
<p>对于缓冲而言，少即是多。</p>
<ul>
<li>
<p>当考虑缓冲的时候，不要考虑性能。</p>
</li>
<li>
<p>缓冲可以帮助降低信号之间的阻塞延迟。</p>
<ul>
<li>降低阻塞延迟到0并不一定意味着更好的吞吐量。</li>
<li>如果一个缓冲可以给你足够的吞吐量，那就保持它。</li>
<li>缓冲大于1的问题需要测量大小。</li>
<li>尽可能找到提供足够吞吐量的最小缓冲</li>
</ul>
</li>
</ul>
</li>
</ul>
<blockquote>
<p>reference:<br>
<a href="https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html">https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html</a><br>
<a href="https://segmentfault.com/a/1190000014524388">https://segmentfault.com/a/1190000014524388</a></p>
</blockquote>

