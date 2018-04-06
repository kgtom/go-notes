---


---

<h1 id="源码解析-go1.10">源码解析 go1.10</h1>
<h3 id="hchan-是-channel-在-runtime-中的数据结构">hchan 是 channel 在 runtime 中的数据结构</h3>
<hr>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> hchan <span class="token keyword">struct</span> <span class="token punctuation">{</span>

qcount <span class="token builtin">uint</span> <span class="token comment">// total data in the queue</span>

dataqsiz <span class="token builtin">uint</span> <span class="token comment">// size of the circular queue</span>

buf unsafe<span class="token punctuation">.</span>Pointer <span class="token comment">// points to an array of dataqsiz elements</span>

elemsize <span class="token builtin">uint16</span>

closed <span class="token builtin">uint32</span>

elemtype <span class="token operator">*</span>_type <span class="token comment">// element type</span>

sendx <span class="token builtin">uint</span> <span class="token comment">// send index</span>

recvx <span class="token builtin">uint</span> <span class="token comment">// receive index</span>

recvq waitq <span class="token comment">// list of recv waiters</span>

sendq waitq <span class="token comment">// list of send waiters</span>

<span class="token comment">// lock protects all fields in hchan, as well as several</span>

<span class="token comment">// fields in sudogs blocked on this channel.</span>

<span class="token comment">//</span>

<span class="token comment">// Do not change another G's status while holding this lock</span>

<span class="token comment">// (in particular, do not ready a G), as this can deadlock</span>

<span class="token comment">// with stack shrinking.</span>

lock mutex

<span class="token punctuation">}</span>
</code></pre>
<h3 id="初始化-make-channel">初始化 make channel</h3>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span> <span class="token function">makechan64</span><span class="token punctuation">(</span>t <span class="token operator">*</span>chantype<span class="token punctuation">,</span> size <span class="token builtin">int64</span><span class="token punctuation">)</span> <span class="token operator">*</span>hchan <span class="token punctuation">{</span>

<span class="token keyword">if</span> <span class="token function">int64</span><span class="token punctuation">(</span><span class="token function">int</span><span class="token punctuation">(</span>size<span class="token punctuation">)</span><span class="token punctuation">)</span> <span class="token operator">!=</span> size <span class="token punctuation">{</span>

<span class="token function">panic</span><span class="token punctuation">(</span><span class="token function">plainError</span><span class="token punctuation">(</span><span class="token string">"makechan: size out of range"</span><span class="token punctuation">)</span><span class="token punctuation">)</span>

<span class="token punctuation">}</span>

<span class="token keyword">return</span> <span class="token function">makechan</span><span class="token punctuation">(</span>t<span class="token punctuation">,</span> <span class="token function">int</span><span class="token punctuation">(</span>size<span class="token punctuation">)</span><span class="token punctuation">)</span>

<span class="token punctuation">}</span>

<span class="token keyword">func</span> <span class="token function">makechan</span><span class="token punctuation">(</span>t <span class="token operator">*</span>chantype<span class="token punctuation">,</span> size <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token operator">*</span>hchan <span class="token punctuation">{</span>

elem <span class="token operator">:=</span> t<span class="token punctuation">.</span>elem

<span class="token comment">// compiler checks this but be safe.</span>

<span class="token keyword">if</span> elem<span class="token punctuation">.</span>size <span class="token operator">&gt;=</span> <span class="token number">1</span><span class="token operator">&lt;&lt;</span><span class="token number">16</span> <span class="token punctuation">{</span>

<span class="token function">throw</span><span class="token punctuation">(</span><span class="token string">"makechan: invalid channel element type"</span><span class="token punctuation">)</span>

<span class="token punctuation">}</span>

<span class="token keyword">if</span> hchanSize<span class="token operator">%</span>maxAlign <span class="token operator">!=</span> <span class="token number">0</span> <span class="token operator">||</span> elem<span class="token punctuation">.</span>align <span class="token operator">&gt;</span> maxAlign <span class="token punctuation">{</span>

<span class="token function">throw</span><span class="token punctuation">(</span><span class="token string">"makechan: bad alignment"</span><span class="token punctuation">)</span>

<span class="token punctuation">}</span>

<span class="token keyword">if</span> size <span class="token operator">&lt;</span> <span class="token number">0</span> <span class="token operator">||</span> <span class="token function">uintptr</span><span class="token punctuation">(</span>size<span class="token punctuation">)</span> <span class="token operator">&gt;</span> <span class="token function">maxSliceCap</span><span class="token punctuation">(</span>elem<span class="token punctuation">.</span>size<span class="token punctuation">)</span> <span class="token operator">||</span> <span class="token function">uintptr</span><span class="token punctuation">(</span>size<span class="token punctuation">)</span><span class="token operator">*</span>elem<span class="token punctuation">.</span>size <span class="token operator">&gt;</span> _MaxMem<span class="token operator">-</span>hchanSize <span class="token punctuation">{</span>

<span class="token function">panic</span><span class="token punctuation">(</span><span class="token function">plainError</span><span class="token punctuation">(</span><span class="token string">"makechan: size out of range"</span><span class="token punctuation">)</span><span class="token punctuation">)</span>

<span class="token punctuation">}</span>

<span class="token comment">// Hchan does not contain pointers interesting for GC when elements stored in buf do not contain pointers.</span>

<span class="token comment">//当存储在buf中的元素不包含指针时，Hchan不包含GC感兴趣的指针。意味着此时channel 与gc 没有关系了。</span>
<span class="token comment">// buf points into the same allocation, elemtype is persistent.</span>

<span class="token comment">// SudoG's are referenced from their owning thread so they can't be collected.</span>

<span class="token comment">// TODO(dvyukov,rlh): Rethink when collector can move allocated objects.</span>

<span class="token keyword">var</span> c <span class="token operator">*</span>hchan

<span class="token keyword">switch</span> <span class="token punctuation">{</span>

<span class="token keyword">case</span> size <span class="token operator">==</span> <span class="token number">0</span> <span class="token operator">||</span> elem<span class="token punctuation">.</span>size <span class="token operator">==</span> <span class="token number">0</span><span class="token punctuation">:</span>

<span class="token comment">// Queue or element size is zero.</span>
<span class="token comment">//元素为0：var  ch  =  make(chan  struct{})</span>
c <span class="token operator">=</span> <span class="token punctuation">(</span><span class="token operator">*</span>hchan<span class="token punctuation">)</span><span class="token punctuation">(</span><span class="token function">mallocgc</span><span class="token punctuation">(</span>hchanSize<span class="token punctuation">,</span> <span class="token boolean">nil</span><span class="token punctuation">,</span> <span class="token boolean">true</span><span class="token punctuation">)</span><span class="token punctuation">)</span>

<span class="token comment">// Race detector uses this location for synchronization.</span>

c<span class="token punctuation">.</span>buf <span class="token operator">=</span> unsafe<span class="token punctuation">.</span><span class="token function">Pointer</span><span class="token punctuation">(</span>c<span class="token punctuation">)</span>

<span class="token keyword">case</span> elem<span class="token punctuation">.</span>kind<span class="token operator">&amp;</span>kindNoPointers <span class="token operator">!=</span> <span class="token number">0</span><span class="token punctuation">:</span>

<span class="token comment">// Elements do not contain pointers.</span>

<span class="token comment">// Allocate hchan and buf in one call.</span>
<span class="token comment">//分配空间：size*elem.size</span>
c <span class="token operator">=</span> <span class="token punctuation">(</span><span class="token operator">*</span>hchan<span class="token punctuation">)</span><span class="token punctuation">(</span><span class="token function">mallocgc</span><span class="token punctuation">(</span>hchanSize<span class="token operator">+</span><span class="token function">uintptr</span><span class="token punctuation">(</span>size<span class="token punctuation">)</span><span class="token operator">*</span>elem<span class="token punctuation">.</span>size<span class="token punctuation">,</span> <span class="token boolean">nil</span><span class="token punctuation">,</span> <span class="token boolean">true</span><span class="token punctuation">)</span><span class="token punctuation">)</span>

c<span class="token punctuation">.</span>buf <span class="token operator">=</span> <span class="token function">add</span><span class="token punctuation">(</span>unsafe<span class="token punctuation">.</span><span class="token function">Pointer</span><span class="token punctuation">(</span>c<span class="token punctuation">)</span><span class="token punctuation">,</span> hchanSize<span class="token punctuation">)</span>

<span class="token keyword">default</span><span class="token punctuation">:</span>

<span class="token comment">// Elements contain pointers.</span>
<span class="token comment">//两次分配内存空间</span>
c <span class="token operator">=</span> <span class="token function">new</span><span class="token punctuation">(</span>hchan<span class="token punctuation">)</span>

c<span class="token punctuation">.</span>buf <span class="token operator">=</span> <span class="token function">mallocgc</span><span class="token punctuation">(</span><span class="token function">uintptr</span><span class="token punctuation">(</span>size<span class="token punctuation">)</span><span class="token operator">*</span>elem<span class="token punctuation">.</span>size<span class="token punctuation">,</span> elem<span class="token punctuation">,</span> <span class="token boolean">true</span><span class="token punctuation">)</span>

<span class="token punctuation">}</span>




</code></pre>

