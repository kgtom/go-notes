

<h2 id="本章知识点">本章知识点</h2>
<ul>
<li><strong>go version:</strong> go version go1.10.1 darwin/amd64</li>
<li>map 基本结构</li>
<li>map创建、读取</li>
<li>map 删除及内存释放</li>
<li>sync.Map并发安全</li>
<li>常见问题</li>
</ul>
<h3 id="一、map基本结构">一、map基本结构</h3>
<p><strong>源码地址：</strong> go1.10/src/runtime/hashmap.go</p>
<p>map的底层结构是hmap（即hashmap的缩写），核心元素是一个由若干个桶（bucket，结构为bmap）组成的数组，每个bucket可以存放若干元素（通常是8个），key通过哈希算法被归入不同的bucket中。当超过8个元素需要存入某个bucket时，hmap会使用extra中的overflow来拓展该bucket。下面是hmap的结构体。</p>

~~~go

type hmap struct {
	count     int // # 元素个数
	flags     uint8
	B         uint8  // 说明包含2^B个bucket
	noverflow uint16 // 溢出的bucket的个数
	hash0     uint32 // hash种子

	buckets    unsafe.Pointer // buckets的数组指针
	oldbuckets unsafe.Pointer // 结构扩容的时候用于复制的buckets数组
	nevacuate  uintptr        // 搬迁进度（已经搬迁的buckets数量）

	extra *mapextra
}

~~~

<p>在extra中不仅有overflow，还有oldoverflow（用于扩容）和nextoverflow（prealloc的地址）。</p>

~~~go

type mapextra struct {
     overflow    *[]*bmap
     oldoverflow *[]*bmap
     nextOverflow *bmap
}
type bmap struct {
    tophash [bucketCnt]uint8
}
~~~

<ul>
<li>tophash用于记录8个key哈希值的高8位，这样在寻找对应key的时候可以更快，不必每次都对key做全等判断。</li>
<li><strong>注意后面几行注释，hmap并非只有一个tophash，而是后面紧跟8组kv对和一个overflow的指针，这样才能使overflow成为一个链表的结构。但是这两个结构体并不是显示定义的，而是直接通过指针运算进行访问的。</strong></li>
<li>kv的存储形式为”key0key1key2key3…key7val1val2val3…val7″，这样做的好处是：在key和value的长度不同的时候，节省padding空间。如上面的例子，在map[int64]int8中，4个相邻的int8可以存储在同一个内存单元中。如果使用kv交错存储的话，每个int8都会被padding占用单独的内存单元（为了提高寻址速度）。</li>
</ul>
<h3 id="二、map-创建、读取">二、map 创建、读取</h3>
<h4 id="创建">创建</h4>
<p>创建其实就是完成内存的申请，以及一些初始值的设定。那么这里假设创建的空间较大，也就是说将 overflow 区域的初始化，也一并放在这里记录。</p>

~~~go

// hint 代表的 capacity
func makemap(t *maptype, hint int64) *hmap  {
    // 条件检查
    t.keysize = sys.PtrSize = t.key.size
    t.valuesize = sys.PtrSize = t.elem.size

    // 通过 hint 确定 hmap 中最小的 B 应该是多大。
    // B 与后面的内存空间申请，以及未来可能的扩容都有关。B 是一个基数。
    // overLoadFactor 考虑了装载因子。golang 将其初始设置为 0.65
    B := uint8(0)
    for ; overLoadFactor(hint, B); B++ {}

    // golang 是 lazy 形式申请内存
        if B != 0 {
        var nextOverflow *bmap
        buckets, nextOverflow = makeBucketArray(t, B)
        if nextOverflow != nil {
            extra = new(mapextra)
            extra.nextOverflow = nextOverflow
        }
    }
    // 后面就是将内存地址关联到 hmap 结构，并返回实例
    h.count = 0  // 记录存储的 k/v pari 数量。扩容时候会用到
    h.B = B  // 记录基数
    h.flags = 0 // 与状态有关。包含并发控制，以及扩容。

    ...
}

~~~


~~~go

// makeBucketArray 会根据情况判断是否要申请 nextOverflow 。
func makeBucketArray(t *maptype, b uint8) (buckets unsafe.Pointer, nextOverflow *bmap) {
    base := uintptr(1 << b)
    nbuckets := base
    if b >= 4 {
        // 向上调整 nbuckets
    }

    // 注意，是按照 nbuckets 申请内存的
    buckets = newarray(t.bucket, int(nbuckets))

    // 处理 overflow 情况，
    if base != nbuckets {
        // 移动到 数据段 的末尾
        nextOverflow = (*bmap)(add(buckets, base*uintptr(t.bucketsize)))

        // 设置末尾地址，用来做末尾检测
        last := (*bmap)(add(buckets, (nbuckets-1)*uintptr(t.bucketsize)))
        last.setoverflow(t, (*bmap)(buckets))
    }
    return buckets, nextOverflow
}


~~~

<h4 id="读取">读取</h4>
<p>读取有  <code>mapaccess1</code>  和  <code>mapaccess2</code>  两个，前者返回指针，后者返回指针和一个  <code>bool</code>，用于判断 key 是否存在。这里只说  <code>mapaccess1</code>。 指针是  <code>value field</code>中存储的地址</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span> <span class="token function">mapaccess1</span><span class="token punctuation">(</span>t <span class="token operator">*</span>maptype<span class="token punctuation">,</span> h <span class="token operator">*</span>hmap<span class="token punctuation">,</span> key unsafe<span class="token punctuation">.</span>Pointer<span class="token punctuation">)</span> unsafe<span class="token punctuation">.</span>Pointer <span class="token punctuation">{</span>
<span class="token comment">// 如果为空或者长度为 0，那么就返回一个 0 值</span>

<span class="token comment">// 如果正在被写入，那么抛出异常</span>

<span class="token comment">// 获取 key 的 hash 值</span>

<span class="token comment">// 确认该 key 所在的 bucket 位置 （可能是在 buckets 也有可能在 oldbuckets 中）</span>
<span class="token comment">// 使用模计算，先计算出如果在 buckets 中，则是在哪个 bucket</span>
<span class="token comment">// 检测 oldbucket 是否为空，如果不为空，则用上面同样的方式得出在 oldbuckets 的位置</span>
<span class="token comment">// 并检测该 bucket 是否已经被 evacuate ，如果已经被 evacuate 则使用 buckets， 否则使用 oldbuckets 中的位置</span>
    m <span class="token operator">:=</span> <span class="token function">uintptr</span><span class="token punctuation">(</span><span class="token number">1</span><span class="token punctuation">)</span><span class="token operator">&lt;&lt;</span>h<span class="token punctuation">.</span>B <span class="token operator">-</span> <span class="token number">1</span>
    b <span class="token operator">:=</span> <span class="token punctuation">(</span><span class="token operator">*</span>bmap<span class="token punctuation">)</span><span class="token punctuation">(</span><span class="token function">add</span><span class="token punctuation">(</span>h<span class="token punctuation">.</span>buckets<span class="token punctuation">,</span> <span class="token punctuation">(</span>hash<span class="token operator">&amp;</span>m<span class="token punctuation">)</span><span class="token operator">*</span><span class="token function">uintptr</span><span class="token punctuation">(</span>t<span class="token punctuation">.</span>bucketsize<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  <span class="token comment">// buckets 结构</span>
    <span class="token keyword">if</span> c <span class="token operator">:=</span> h<span class="token punctuation">.</span>oldbuckets<span class="token punctuation">;</span> c <span class="token operator">!=</span> <span class="token boolean">nil</span> <span class="token punctuation">{</span>
        <span class="token keyword">if</span> <span class="token operator">!</span>h<span class="token punctuation">.</span><span class="token function">sameSizeGrow</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  <span class="token comment">// 上次扩容是等量还是双倍扩容， 会有影响</span>
            <span class="token comment">// There used to be half as many buckets; mask down one more power of two.</span>
            m <span class="token operator">&gt;&gt;=</span> <span class="token number">1</span>
        <span class="token punctuation">}</span>
        oldb <span class="token operator">:=</span> <span class="token punctuation">(</span><span class="token operator">*</span>bmap<span class="token punctuation">)</span><span class="token punctuation">(</span><span class="token function">add</span><span class="token punctuation">(</span>c<span class="token punctuation">,</span> <span class="token punctuation">(</span>hash<span class="token operator">&amp;</span>m<span class="token punctuation">)</span><span class="token operator">*</span><span class="token function">uintptr</span><span class="token punctuation">(</span>t<span class="token punctuation">.</span>bucketsize<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token punctuation">)</span>
        <span class="token keyword">if</span> <span class="token operator">!</span><span class="token function">evacuated</span><span class="token punctuation">(</span>oldb<span class="token punctuation">)</span> <span class="token punctuation">{</span>
            b <span class="token operator">=</span> oldb
        <span class="token punctuation">}</span>
    <span class="token punctuation">}</span>

<span class="token comment">// 得到 bucket 以后，通过 tophash 来再次定位，如果定位不到，则递归该 bucket 的 overflow 域，循环查找。</span>
<span class="token comment">// 以上步骤有两个结果。</span>
<span class="token comment">// 1 遍历到最后，都没有找到命中的 tophash ，此时则返回一个零值。</span>
<span class="token comment">// 2 命中 tophash 则进行 key 比对，相同则返回对应的 val 位置，不同则通过 overflow 继续获</span>

</code></pre>
<h3 id="三、map-删除及内存分配">三、map 删除及内存分配</h3>
<h4 id="删除">删除</h4>
<p>mapdelete 是负责删除的的主体函数</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span> <span class="token function">mapdelete</span><span class="token punctuation">(</span>t <span class="token operator">*</span>maptype<span class="token punctuation">,</span> h <span class="token operator">*</span>hmap<span class="token punctuation">,</span> key unsafe<span class="token punctuation">.</span>Pointer<span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> <span class="token punctuation">;</span> b <span class="token operator">!=</span> <span class="token boolean">nil</span><span class="token punctuation">;</span> b <span class="token operator">=</span> b<span class="token punctuation">.</span><span class="token function">overflow</span><span class="token punctuation">(</span>t<span class="token punctuation">)</span> <span class="token punctuation">{</span>
		<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token function">uintptr</span><span class="token punctuation">(</span><span class="token number">0</span><span class="token punctuation">)</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> bucketCnt<span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
			b<span class="token punctuation">.</span>tophash<span class="token punctuation">[</span>i<span class="token punctuation">]</span> <span class="token operator">=</span> empty
			h<span class="token punctuation">.</span>count<span class="token operator">--</span>
		<span class="token punctuation">}</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre>
<p>外层的循环就是在遍历整个 map，删除的核心就在那个empty。它修改了当前 key 的标记，而不是直接删除了内存里面的数据。</p>
<p><code>empty = 0 // cell is empty</code></p>
<h4 id="内存释放">内存释放</h4>
<pre class=" language-go"><code class="prism  language-go">m <span class="token operator">=</span> <span class="token boolean">nil</span>
</code></pre>
<p>这之后坐等垃圾回收器回收就好了。</p>
<p><strong><em>如果你用 map 做缓存，而每次更新只是部分更新，更新的 key 如果偏差比较大，有可能会有内存逐渐增长而不释放的问题</em></strong>。要注意。</p>
<h3 id="四、sync.map使用">四、sync.Map使用</h3>
<h3 id="sync.map提供的方法">1. sync.Map提供的方法</h3>
<ul>
<li>存储数据,存入key以及value可以为任意类型.</li>
</ul>
<pre><code>func (m *Map) Store(key, value interface{})
</code></pre>
<ul>
<li>删除对应key</li>
</ul>
<pre><code>func (m *Map) Delete(key interface{})
</code></pre>
<ul>
<li>读取对应key的值,ok表示是否在map中查询到key</li>
</ul>
<pre><code>func (m *Map) Load(key interface{}) (value interface{}, ok bool)
</code></pre>
<ul>
<li>针对某个key的存在读取不存在就存储,loaded为true表示存在值,false表示不存在值.</li>
</ul>
<pre><code>func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)
</code></pre>
<ul>
<li>表示对所有key进行遍历,并将遍历出的key,value传入回调函数进行函数调用,回调函数返回false时遍历结束,否则遍历完所有key.</li>
</ul>
<pre><code>func (m *Map) Range(f func(key, value interface{}) bool)
</code></pre>
<h3 id="原理">2. 原理</h3>
<p>通过引入两个map,将读写分离到不同的map,其中read map只提供读,而dirty map则负责写.<br>
这样read map就可以在不加锁的情况下进行并发读取,当read map中没有读取到值时,再加锁进行后续读取,并累加未命中数,当未命中数到达一定数量后,将dirty map上升为read map.</p>
<p>另外，虽然引入了两个map，但是底层数据存储的是指针，指向的是同一份值．</p>
<p>具体流程：<br>
如插入key 1,2,3时均插入了dirty map中，此时read map没有key值，读取时从dirty map中读取，并记录miss数</p>
<p><img src="https://segmentfault.com/img/bV7l2O?w=626&amp;h=246" alt="图片描述" title="图片描述"></p>
<p>当miss数大于等于dirty map的长度时，将dirty map直接升级为read map,这里直接 对dirty map进行地址拷贝．</p>
<p><img src="https://segmentfault.com/img/bV7l2T?w=625&amp;h=246" alt="图片描述" title="图片描述"></p>
<p>当有新的key 4插入时，将read map中的key值拷贝到dirty map中，这样dirty map就含有所有的值，下次升级为read map时直接进行地址拷贝．</p>
<h3 id="局限性">3. 局限性</h3>
<p>从以上的源码可知,sync.map并不适合同时存在大量读写的场景,大量的写会导致read map读取不到数据从而加锁进行进一步读取,同时dirty map不断升级为read map.<br>
从而导致整体性能较低,特别是针对cache场景.针对append-only以及大量读,少量写场景使用sync.map则相对比较合适.</p>
<p>对于map,还有一种基于hash的实现思路,具体就是对map加读写锁,但是分配n个map,根据对key做hash运算确定是分配到哪个map中.<br>
这样锁的消耗就降到了1/n(理论值).具体实现可见:<a href="https://github.com/orcaman/concurrent-map">https://github.com/orcaman/co…</a></p>
<p>相比之下, 基于hash的方式更容易理解,整体性能较稳定. sync.map在某些场景性能可能差一些,但某些场景却能取得更好的效果.<br>
所以还是要根据具体的业务场景进行取舍.</p>
<h3 id="五、使用中常见问题">五、使用中常见问题</h3>
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

