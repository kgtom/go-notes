---


---

<p><strong>1.运行时系统相关模块</strong></p>
<p><img src="https://img-blog.csdn.net/20170421163340017?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>这一节只讲述调度系统。</p>
<p><strong>2. 与调度相关的有下列数据有：</strong></p>
<p><img src="https://img-blog.csdn.net/20170421163353460?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<ul>
<li>全局M\P\G列表</li>
<li>调度器空闲M\P列表</li>
<li>调度器可运行的G列表</li>
<li>调度器自由的G列表</li>
<li>本地P可运行的G列表</li>
<li>本地P自由的G列表</li>
</ul>
<p>(1)全局表就不用多解释， 调度器不是一个专门的线程，而是一种资源，每个M都可能会对其进行查询，如：M的调用的时候，会首先查找调度器中的可运行队列G，然后查找本地P运行队列; 在M不够的时候，会首先调度器空闲的M列表中查找; 在程序运行过程中，使用<strong>runtime.GOMAXPROCS</strong>接口设置P上下文数量的时候，会动态的增长P的数量，这些P就是加入到调度器的空闲列表。</p>
<p>(2)M每次调用G的时候，必定需要一个上下文P，因为当调用的时候，可能会产生另外的G，其G的取得的步骤为， 从当前P自由G列表中取-&gt; 从调度器自由G列表中取 -&gt; 全局分配。当然，当该P运行完一个G后，会首先加入到P自由G列表，然后超过数量则加入到调度器自由G列表中。</p>
<p><strong>3. M,P,G的对应关系</strong></p>
<p><img src="https://img-blog.csdn.net/20170421163409179?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>M 默认是多个， 也就是对应多个OS线程，无状态变化，唯一的变化就是是否在空闲队列中。</p>
<p>3.1 P的状态介绍</p>
<p>P可以通过runtime.GOMAXPROCS设置，一般不用设置。P具有状态：</p>
<p><img src="https://img-blog.csdn.net/20170421163423258?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>状态机如下：</p>
<p><img src="https://img-blog.csdn.net/20170421163437705?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>(1)M与P关联后，会从Pidle -&gt; Prunning</p>
<p>(2)P变成Pdead状态，只有用户设置了P的数量减少了</p>
<p>3.2 G的状态介绍</p>
<p><img src="https://img-blog.csdn.net/20170421163553534?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>状态机图：</p>
<p><img src="https://img-blog.csdn.net/20170421163617768?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>(1)Grunning -&gt; Gwaiting 这个状态是Golang自己的channel接收机制造成</p>
<p>(2)Grunning-&gt; Gsyscall 是该协程使用了某些系统调用，如IO读取，所以直接挂起，使用的是异步IO阻塞机制</p>
<p><strong>4. M,P,G运行过程：</strong></p>
<p><img src="https://img-blog.csdn.net/20170421163635440?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p>注意：上面的全力查找注释有误，为没有则 M会挂起阻塞，然后等待被其他M唤醒</p>
<p>其中，全力查找可运行G的过程会通过以下几步</p>
<p>(1)查找本地P可运行G队列</p>
<p>(2)从调度器可运行G队列查找</p>
<p>(3)从网络IO轮询器中查找就绪G</p>
<p>(4)偷取别人的本地P可运行队列G</p>
<p>(5)再次查找调度器中的可运行G</p>
<p>(6)尝试从所有P中查找可运行队列G</p>
<p>(7)再从网络IO轮询器中查找</p>
<p>实在不行，就把自己和P解除联结，放置到调度器空闲队列中，并且阻塞。</p>
<p>其中查找到多个运行G之后的流程图如下：</p>
<p><img src="https://img-blog.csdn.net/20170421163710777?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvQ3VsbWluYXRlX2lu/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/Center" alt=""></p>
<p><strong>5.G IO操作的运行过程</strong></p>
<p>(1)当G执行IO操作后， 并且如果<strong>该G和M已经绑定，那么G阻塞（阻塞只是打个标记，并不是阻塞当前线程），IO请求由IO轮询器进行，然后M继续运行</strong>，把当前的P解除关联，M阻塞等待，当IO轮询器执行完后，然后G变成可运行状态，当其他的M找到这个绑定M的G，会把自己的上下文P解除联结，并且唤醒绑定的M，然后继续运行。</p>
<p><strong>注意：当未使用G和M绑定的时候，就不会执行唤醒其他M的操作</strong></p>
<blockquote>
<p>reference : <a href="https://blog.csdn.net/Culminate_in/article/details/70326681?locationNum=2&amp;fps=1">addr</a>.</p>
</blockquote>

