---


---

<p><strong>源码地址</strong>：/src/time/tick.go</p>
<pre class=" language-go"><code class="prism  language-go"> 
<span class="token comment">//返回一个Ticker对象 </span>
<span class="token comment">//参数D 必须大于0，否则panic  </span>
<span class="token comment">//Stop()用完后，释放资源</span>
<span class="token keyword">func</span> <span class="token function">NewTicker</span><span class="token punctuation">(</span>d Duration<span class="token punctuation">)</span> <span class="token operator">*</span>Ticker <span class="token punctuation">{</span>  
   <span class="token keyword">if</span> d <span class="token operator">&lt;=</span> <span class="token number">0</span> <span class="token punctuation">{</span>  
      <span class="token function">panic</span><span class="token punctuation">(</span>errors<span class="token punctuation">.</span><span class="token function">New</span><span class="token punctuation">(</span><span class="token string">"non-positive interval for NewTicker"</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
   <span class="token punctuation">}</span>  
   <span class="token comment">// Give the channel a 1-element time buffer.  </span>
 <span class="token comment">// If the client falls behind while reading, we drop ticks // on the floor until the client catches up.  c := make(chan Time, 1)  </span>
   t <span class="token operator">:=</span> <span class="token operator">&amp;</span>Ticker<span class="token punctuation">{</span>  
      C<span class="token punctuation">:</span> c<span class="token punctuation">,</span>  
      r<span class="token punctuation">:</span> runtimeTimer<span class="token punctuation">{</span>  
         when<span class="token punctuation">:</span>   <span class="token function">when</span><span class="token punctuation">(</span>d<span class="token punctuation">)</span><span class="token punctuation">,</span>  
         period<span class="token punctuation">:</span> <span class="token function">int64</span><span class="token punctuation">(</span>d<span class="token punctuation">)</span><span class="token punctuation">,</span>  
         f<span class="token punctuation">:</span>      sendTime<span class="token punctuation">,</span>  
         arg<span class="token punctuation">:</span>    c<span class="token punctuation">,</span>  
      <span class="token punctuation">}</span><span class="token punctuation">,</span>  
   <span class="token punctuation">}</span>  
<span class="token comment">//将新的定时任务添加到时间堆中,runtime.time.go</span>
   <span class="token function">startTimer</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>t<span class="token punctuation">.</span>r<span class="token punctuation">)</span> 
   <span class="token keyword">return</span> t  
<span class="token punctuation">}</span>
</code></pre>
<p><strong>解析：</strong></p>
<ul>
<li>when(d):返回一个runtimeNano() + int64(d)的未来时(到期时间)<br>
runtimeNano运行时当前纳秒时间</li>
<li>period: 被唤醒的时间</li>
<li>f:时间到期后的回调函数</li>
<li>arg:text<br>
// 时间到期后的断言参数</li>
</ul>
<p><strong>源码地址</strong>：/src/runtime/time.go</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">// startTimer adds t to the timer heap.  </span>
<span class="token comment">//go:linkname startTimer time.startTimerfunc startTimer(t *timer) {  </span>
   <span class="token keyword">if</span> raceenabled <span class="token punctuation">{</span>  
      <span class="token function">racerelease</span><span class="token punctuation">(</span>unsafe<span class="token punctuation">.</span><span class="token function">Pointer</span><span class="token punctuation">(</span>t<span class="token punctuation">)</span><span class="token punctuation">)</span>  
   <span class="token punctuation">}</span>  
   <span class="token function">addtimer</span><span class="token punctuation">(</span>t<span class="token punctuation">)</span>  
<span class="token punctuation">}</span>

<span class="token comment">//addtimer</span>
<span class="token keyword">func</span> <span class="token function">addtimer</span><span class="token punctuation">(</span>t <span class="token operator">*</span>timer<span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   tb <span class="token operator">:=</span> t<span class="token punctuation">.</span><span class="token function">assignBucket</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
   <span class="token function">lock</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>tb<span class="token punctuation">.</span>lock<span class="token punctuation">)</span>  
   tb<span class="token punctuation">.</span><span class="token function">addtimerLocked</span><span class="token punctuation">(</span>t<span class="token punctuation">)</span>  
   <span class="token function">unlock</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>tb<span class="token punctuation">.</span>lock<span class="token punctuation">)</span>  
<span class="token punctuation">}</span>

</code></pre>
<p><strong>重点看addtimerLocked()</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">// Add a timer to the heap and start or kick timerproc if the new timer is  </span>
<span class="token comment">// earlier than any of the others.  </span>
<span class="token comment">// Timers are locked.  </span>
<span class="token keyword">func</span> <span class="token punctuation">(</span>tb <span class="token operator">*</span>timersBucket<span class="token punctuation">)</span> <span class="token function">addtimerLocked</span><span class="token punctuation">(</span>t <span class="token operator">*</span>timer<span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   <span class="token comment">// when must never be negative; otherwise timerproc will overflow  </span>
 <span class="token comment">// during its delta calculation and never expire other runtime timers.  if t.when &lt; 0 {  </span>
      t<span class="token punctuation">.</span>when <span class="token operator">=</span> <span class="token number">1</span><span class="token operator">&lt;&lt;</span><span class="token number">63</span> <span class="token operator">-</span> <span class="token number">1</span>  
  <span class="token punctuation">}</span>  
   t<span class="token punctuation">.</span>i <span class="token operator">=</span> <span class="token function">len</span><span class="token punctuation">(</span>tb<span class="token punctuation">.</span>t<span class="token punctuation">)</span>  
   tb<span class="token punctuation">.</span>t <span class="token operator">=</span> <span class="token function">append</span><span class="token punctuation">(</span>tb<span class="token punctuation">.</span>t<span class="token punctuation">,</span> t<span class="token punctuation">)</span>  
   <span class="token function">siftupTimer</span><span class="token punctuation">(</span>tb<span class="token punctuation">.</span>t<span class="token punctuation">,</span> t<span class="token punctuation">.</span>i<span class="token punctuation">)</span>  
   <span class="token keyword">if</span> t<span class="token punctuation">.</span>i <span class="token operator">==</span> <span class="token number">0</span> <span class="token punctuation">{</span>  
      <span class="token comment">// siftup moved to top: new earliest deadline.  </span>
  <span class="token keyword">if</span> tb<span class="token punctuation">.</span>sleeping <span class="token punctuation">{</span>  
         tb<span class="token punctuation">.</span>sleeping <span class="token operator">=</span> <span class="token boolean">false</span>  
  <span class="token function">notewakeup</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>tb<span class="token punctuation">.</span>waitnote<span class="token punctuation">)</span>  
      <span class="token punctuation">}</span>  
      <span class="token keyword">if</span> tb<span class="token punctuation">.</span>rescheduling <span class="token punctuation">{</span>  
         tb<span class="token punctuation">.</span>rescheduling <span class="token operator">=</span> <span class="token boolean">false</span>  
  <span class="token function">goready</span><span class="token punctuation">(</span>tb<span class="token punctuation">.</span>gp<span class="token punctuation">,</span> <span class="token number">0</span><span class="token punctuation">)</span>  
      <span class="token punctuation">}</span>  
   <span class="token punctuation">}</span>  
   <span class="token keyword">if</span> <span class="token operator">!</span>tb<span class="token punctuation">.</span>created <span class="token punctuation">{</span>  
      tb<span class="token punctuation">.</span>created <span class="token operator">=</span> <span class="token boolean">true</span>  
  <span class="token keyword">go</span> <span class="token function">timerproc</span><span class="token punctuation">(</span>tb<span class="token punctuation">)</span>  
   <span class="token punctuation">}</span>  
<span class="token punctuation">}</span>

</code></pre>
<p><strong>举例：</strong></p>
<pre class=" language-go"><code class="prism  language-go">
<span class="token comment">//指定时间间隔，周期性做事情(eg：定时任务)，用完后记得调用stop()释放资源  </span>
ticker <span class="token operator">:=</span> time<span class="token punctuation">.</span><span class="token function">NewTicker</span><span class="token punctuation">(</span>time<span class="token punctuation">.</span>Second <span class="token operator">*</span> <span class="token number">2</span><span class="token punctuation">)</span>  
<span class="token keyword">defer</span> ticker<span class="token punctuation">.</span><span class="token function">Stop</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
<span class="token comment">//调用 Stop() 使计时器停止,释放资源  </span>
i <span class="token operator">:=</span> <span class="token number">0</span>  
<span class="token keyword">for</span> i <span class="token operator">&lt;</span> <span class="token number">5</span> <span class="token punctuation">{</span>  
   <span class="token keyword">select</span> <span class="token punctuation">{</span>  
   <span class="token keyword">case</span> t <span class="token operator">:=</span> <span class="token operator">&lt;-</span>ticker<span class="token punctuation">.</span>C<span class="token punctuation">:</span>  
      fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>t<span class="token punctuation">,</span> <span class="token string">" received every 2s"</span><span class="token punctuation">)</span>  
   <span class="token punctuation">}</span>  
   i<span class="token operator">++</span>  
<span class="token punctuation">}</span>
</code></pre>
<p><strong>timer.After使用</strong>只触发一次，常用于取消超时等待行为，等同于NewTimer(d).C，通常使用NewTimer 代替time.stop</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">package</span> main  
  
<span class="token keyword">import</span> <span class="token punctuation">(</span>  
   <span class="token string">"fmt"</span>  
 <span class="token string">"time"</span>
 <span class="token punctuation">)</span>  
  
<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
  
   ch <span class="token operator">:=</span> <span class="token function">make</span><span class="token punctuation">(</span><span class="token keyword">chan</span> <span class="token builtin">int</span><span class="token punctuation">,</span> <span class="token number">1</span><span class="token punctuation">)</span>  
   <span class="token keyword">go</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
  
      <span class="token comment">//do sth  </span>
  time<span class="token punctuation">.</span><span class="token function">Sleep</span><span class="token punctuation">(</span>time<span class="token punctuation">.</span>Second <span class="token operator">*</span> <span class="token number">3</span><span class="token punctuation">)</span>  
      ch <span class="token operator">&lt;-</span> <span class="token number">1</span>  
  
  <span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
  
   timer <span class="token operator">:=</span> time<span class="token punctuation">.</span><span class="token function">NewTimer</span><span class="token punctuation">(</span>time<span class="token punctuation">.</span>Second<span class="token punctuation">)</span>  
  
   <span class="token keyword">for</span> <span class="token punctuation">{</span>  
      <span class="token keyword">select</span> <span class="token punctuation">{</span>  
      <span class="token keyword">case</span> <span class="token operator">&lt;-</span>ch<span class="token punctuation">:</span>  
         fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"do sth"</span><span class="token punctuation">)</span>  
         <span class="token comment">// do something  </span>
         <span class="token comment">//case &lt;-time.After(time.Second * 2): //每次进入select都会初始化一个time,浪费cup资源，</span>
         <span class="token comment">//改进如下：使用NewTimer </span>
         <span class="token comment">// fmt.Println("time out") // goto Loop  case &lt;-timer.C:  </span>
         fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"time out user NewTimer"</span><span class="token punctuation">)</span>  
         <span class="token keyword">goto</span> Loop  
  <span class="token punctuation">}</span>  
   <span class="token punctuation">}</span>  
Loop<span class="token punctuation">:</span>  
   fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"end"</span><span class="token punctuation">)</span>  
<span class="token punctuation">}</span>

</code></pre>
<ul>
<li>当case &lt;- ch 和 case &lt;- ticker.C 同时成立时,，Select会随机公平地选出一个执行，有可能选择到前者，导致超时随机行失败,改进版如下：<strong>将超时和正常业务各自放一个selcet</strong></li>
</ul>
<pre class=" language-go"><code class="prism  language-go">
<span class="token keyword">package</span> main

<span class="token keyword">import</span> <span class="token punctuation">(</span>
	<span class="token string">"fmt"</span>
	<span class="token string">"time"</span>
<span class="token punctuation">)</span>

<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>

	ch <span class="token operator">:=</span> <span class="token function">make</span><span class="token punctuation">(</span><span class="token keyword">chan</span> <span class="token builtin">int</span><span class="token punctuation">,</span> <span class="token number">1</span><span class="token punctuation">)</span>
	<span class="token keyword">go</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>

		<span class="token comment">//do sth</span>
		time<span class="token punctuation">.</span><span class="token function">Sleep</span><span class="token punctuation">(</span>time<span class="token punctuation">.</span>Second <span class="token operator">*</span> <span class="token number">3</span><span class="token punctuation">)</span>
		ch <span class="token operator">&lt;-</span> <span class="token number">1</span>

	<span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>

	timer <span class="token operator">:=</span> time<span class="token punctuation">.</span><span class="token function">NewTimer</span><span class="token punctuation">(</span>time<span class="token punctuation">.</span>Second<span class="token punctuation">)</span>
	<span class="token keyword">defer</span> timer<span class="token punctuation">.</span><span class="token function">Stop</span><span class="token punctuation">(</span><span class="token punctuation">)</span>

	<span class="token keyword">for</span> <span class="token punctuation">{</span>
		<span class="token keyword">select</span> <span class="token punctuation">{</span>
		<span class="token keyword">case</span> <span class="token operator">&lt;-</span>ch<span class="token punctuation">:</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"do sth"</span><span class="token punctuation">)</span>
		<span class="token keyword">default</span><span class="token punctuation">:</span><span class="token operator">/</span><span class="token operator">/</span>保证不阻塞

		<span class="token punctuation">}</span>
		<span class="token keyword">select</span> <span class="token punctuation">{</span>
		<span class="token keyword">case</span> <span class="token operator">&lt;-</span>timer<span class="token punctuation">.</span>C<span class="token punctuation">:</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"time out user NewTimer"</span><span class="token punctuation">)</span>
			<span class="token comment">//goto Loop</span>
			<span class="token keyword">return</span>
		<span class="token keyword">default</span><span class="token punctuation">:</span><span class="token operator">/</span><span class="token operator">/</span>保证不阻塞

		<span class="token punctuation">}</span>
	<span class="token punctuation">}</span>
	<span class="token comment">//Loop:</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"end"</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span>

</code></pre>
<p><strong><code>time.Stop</code>  停止定时器 或  <code>time.Reset</code>  重置定时器</strong></p>
<pre class=" language-go"><code class="prism  language-go">start <span class="token operator">:=</span> time<span class="token punctuation">.</span><span class="token function">Now</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
timer <span class="token operator">:=</span> time<span class="token punctuation">.</span><span class="token function">AfterFunc</span><span class="token punctuation">(</span><span class="token number">2</span><span class="token operator">*</span>time<span class="token punctuation">.</span>Second<span class="token punctuation">,</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"after func callback, elaspe:"</span><span class="token punctuation">,</span> time<span class="token punctuation">.</span><span class="token function">Now</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Sub</span><span class="token punctuation">(</span>start<span class="token punctuation">)</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span><span class="token punctuation">)</span>

time<span class="token punctuation">.</span><span class="token function">Sleep</span><span class="token punctuation">(</span><span class="token number">1</span> <span class="token operator">*</span> time<span class="token punctuation">.</span>Second<span class="token punctuation">)</span>
<span class="token comment">// time.Sleep(3*time.Second)</span>
<span class="token comment">// Reset 在 Timer 还未触发时返回 true；触发了或Stop了，返回false</span>
<span class="token keyword">if</span> timer<span class="token punctuation">.</span><span class="token function">Reset</span><span class="token punctuation">(</span><span class="token number">3</span> <span class="token operator">*</span> time<span class="token punctuation">.</span>Second<span class="token punctuation">)</span> <span class="token punctuation">{</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"timer has not trigger!"</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span> <span class="token keyword">else</span> <span class="token punctuation">{</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"timer had expired or stop!"</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span>

time<span class="token punctuation">.</span><span class="token function">Sleep</span><span class="token punctuation">(</span><span class="token number">10</span> <span class="token operator">*</span> time<span class="token punctuation">.</span>Second<span class="token punctuation">)</span>

<span class="token comment">// output:</span>
<span class="token comment">// timer has not trigger!</span>
<span class="token comment">// after func callback, elaspe: 4.00026461s</span>
</code></pre>
<p>如果定时器还未触发，<code>Stop</code>  会将其移除，并返回 true；否则返回 false；后续再对该  <code>Timer</code>  调用  <code>Stop</code>，直接返回 false。</p>
<p><code>Reset</code>  会先调用  <code>stopTimer</code>  再调用  <code>startTimer</code>，类似于废弃之前的定时器，重新启动一个定时器。返回值和  <code>Stop</code>  一样。</p>
<blockquote>
<p>Reference:<br>
<a href="https://github.com/golang/go/">github-golang</a><br>
<a href="https://zhuanlan.zhihu.com/p/31634188">zhihu</a><br>
<a href="https://github.com/polaris1119/The-Golang-Standard-Library-by-Example/blob/master/chapter04/04.4.md">github-repo</a></p>
</blockquote>

