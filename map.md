---


---

<h2 id="本章知识点">本章知识点</h2>
<ul>
<li><strong>go version:</strong> go version go1.10.1 darwin/amd64</li>
<li>map 基本结构</li>
<li>map 删除及内存回收</li>
<li>sync.Map并发安全</li>
<li>常见问题</li>
</ul>
<h3 id="一、map基本结构">一、map基本结构</h3>
<p><strong>源码地址：</strong> go1.10/src/runtime/hashmap.go</p>
<p>map的底层结构是hmap（即hashmap的缩写），核心元素是一个由若干个桶（bucket，结构为bmap）组成的数组，每个bucket可以存放若干元素（通常是8个），key通过哈希算法被归入不同的bucket中。当超过8个元素需要存入某个bucket时，hmap会使用extra中的overflow来拓展该bucket。下面是hmap的结构体。</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> hmap <span class="token keyword">struct</span> <span class="token punctuation">{</span>
	count     <span class="token builtin">int</span> <span class="token comment">// # 元素个数</span>
	flags     <span class="token builtin">uint8</span>
	B         <span class="token builtin">uint8</span>  <span class="token comment">// 说明包含2^B个bucket</span>
	noverflow <span class="token builtin">uint16</span> <span class="token comment">// 溢出的bucket的个数</span>
	hash0     <span class="token builtin">uint32</span> <span class="token comment">// hash种子</span>

	buckets    unsafe<span class="token punctuation">.</span>Pointer <span class="token comment">// buckets的数组指针</span>
	oldbuckets unsafe<span class="token punctuation">.</span>Pointer <span class="token comment">// 结构扩容的时候用于复制的buckets数组</span>
	nevacuate  <span class="token builtin">uintptr</span>        <span class="token comment">// 搬迁进度（已经搬迁的buckets数量）</span>

	extra <span class="token operator">*</span>mapextra
<span class="token punctuation">}</span>

</code></pre>
<p>在extra中不仅有overflow，还有oldoverflow（用于扩容）和nextoverflow（prealloc的地址）。</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> mapextra <span class="token keyword">struct</span> <span class="token punctuation">{</span>
     overflow    <span class="token operator">*</span><span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token operator">*</span>bmap
     oldoverflow <span class="token operator">*</span><span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token operator">*</span>bmap
     nextOverflow <span class="token operator">*</span>bmap
<span class="token punctuation">}</span>
</code></pre>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> bmap <span class="token keyword">struct</span> <span class="token punctuation">{</span>
    tophash <span class="token punctuation">[</span>bucketCnt<span class="token punctuation">]</span><span class="token builtin">uint8</span>
<span class="token punctuation">}</span>
</code></pre>
<ul>
<li>tophash用于记录8个key哈希值的高8位，这样在寻找对应key的时候可以更快，不必每次都对key做全等判断。</li>
<li><strong>注意后面几行注释，hmap并非只有一个tophash，而是后面紧跟8组kv对和一个overflow的指针，这样才能使overflow成为一个链表的结构。但是这两个结构体并不是显示定义的，而是直接通过指针运算进行访问的。</strong></li>
<li>kv的存储形式为”key0key1key2key3…key7val1val2val3…val7″，这样做的好处是：在key和value的长度不同的时候，节省padding空间。如上面的例子，在map[int64]int8中，4个相邻的int8可以存储在同一个内存单元中。如果使用kv交错存储的话，每个int8都会被padding占用单独的内存单元（为了提高寻址速度）。</li>
</ul>
<h3 id="二、map-删除及内存分配">二、map 删除及内存分配</h3>
<h3 id="三、sync.map使用">三、sync.Map使用</h3>
<h3 id="四、使用中常见问题">四、使用中常见问题</h3>
<p>Q：删除掉map中的元素是否会释放内存？</p>
<p>A：不会，删除操作仅仅将对应的tophash[i]设置为empty，并非释放内存。若要释放内存只能等待指针无引用后(即 m=nil 时)被系统gc</p>
<p>Q：如何并发地使用map？</p>
<p>A：map不是goroutine安全的，所以在有多个gorountine对map进行写操作是会panic。多gorountine读写map是应加锁（RWMutex），或使用sync.Map（1.9新增)。</p>
<p>Q：map的iterator是否安全？</p>
<p>A：map的delete并非真的delete，所以对迭代器是没有影响的，是安全的。</p>
<blockquote>
<p>reference:<br>
<a href="https://github.com/golang/go/blob/dev.boringcrypto.go1.10/src/runtime/hashmap.go">github</a><br>
<a href="https://juejin.im/entry/5a1e4bcd6fb9a045090942d8">juejin</a><br>
<a href="http://blog.cyeam.com/json/2017/11/02/go-map-delete">cyeam</a><br>
<a href="https://segmentfault.com/a/1190000014107066">segmentfault</a></p>
</blockquote>

