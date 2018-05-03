---


---

<h3 id="基础">1.基础</h3>
<pre><code>runtime/slice.go
type slice struct {
    array unsafe.Pointer
    len   int
    cap   int
}
</code></pre>
<p>slice 的底层结构定义，指向底层数组的指针，当前长度 len 和当前 容量 的 cap。</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">var</span>  s  <span class="token operator">=</span> <span class="token punctuation">[</span><span class="token number">5</span><span class="token punctuation">]</span><span class="token builtin">int</span><span class="token punctuation">{</span><span class="token number">1</span><span class="token punctuation">,</span> <span class="token number">2</span><span class="token punctuation">,</span> <span class="token number">3</span><span class="token punctuation">,</span> <span class="token number">4</span><span class="token punctuation">,</span> <span class="token number">5</span><span class="token punctuation">}</span>

<span class="token keyword">var</span>  b  <span class="token operator">=</span> s<span class="token punctuation">[</span><span class="token number">2</span><span class="token punctuation">:</span><span class="token number">3</span><span class="token punctuation">]</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"b len()"</span><span class="token punctuation">,</span> <span class="token function">len</span><span class="token punctuation">(</span>b<span class="token punctuation">)</span><span class="token punctuation">,</span> <span class="token string">"cap():"</span><span class="token punctuation">,</span> <span class="token function">cap</span><span class="token punctuation">(</span>b<span class="token punctuation">)</span><span class="token punctuation">,</span> <span class="token string">"b:"</span><span class="token punctuation">,</span> b<span class="token punctuation">)</span>

<span class="token comment">//长度：元素个数，不能超过slice的容量</span>

<span class="token comment">//容量：从slice的起始元素到底层数组的最后一个元素的个数</span>
</code></pre>
<h3 id="copy-使用">2.copy 使用</h3>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">//将第二slice元素拷贝到第一个slice</span>

s1  <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">int</span><span class="token punctuation">{</span><span class="token number">1</span><span class="token punctuation">,</span> <span class="token number">2</span><span class="token punctuation">,</span> <span class="token number">3</span><span class="token punctuation">,</span> <span class="token number">4</span><span class="token punctuation">,</span> <span class="token number">5</span><span class="token punctuation">}</span>

s2  <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">int</span><span class="token punctuation">{</span><span class="token number">7</span><span class="token punctuation">,</span> <span class="token number">8</span><span class="token punctuation">}</span>

<span class="token function">copy</span><span class="token punctuation">(</span>s1<span class="token punctuation">,</span> s2<span class="token punctuation">)</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"s1:"</span><span class="token punctuation">,</span> s1<span class="token punctuation">,</span> <span class="token string">"s1cap:"</span><span class="token punctuation">,</span> <span class="token function">cap</span><span class="token punctuation">(</span>s1<span class="token punctuation">)</span><span class="token punctuation">,</span> <span class="token string">"s2:"</span><span class="token punctuation">,</span> s2<span class="token punctuation">)</span>

<span class="token comment">//s1: [7 8 3 4 5] s1cap: 5 s2: [7 8]</span>

<span class="token function">copy</span><span class="token punctuation">(</span>s2<span class="token punctuation">,</span> s1<span class="token punctuation">)</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"s1:"</span><span class="token punctuation">,</span> s1<span class="token punctuation">,</span> <span class="token string">"s2:"</span><span class="token punctuation">,</span> s2<span class="token punctuation">)</span>

<span class="token comment">//s1: [1 2 3 4 5] s2: [1 2]</span>
</code></pre>
<h3 id="append">3.append</h3>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">//将第二slice元素拷贝到第一个slice</span>

s1  <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">int</span><span class="token punctuation">{</span><span class="token number">1</span><span class="token punctuation">,</span> <span class="token number">2</span><span class="token punctuation">,</span> <span class="token number">3</span><span class="token punctuation">,</span> <span class="token number">4</span><span class="token punctuation">,</span> <span class="token number">5</span><span class="token punctuation">}</span>

s2  <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">int</span><span class="token punctuation">{</span><span class="token number">7</span><span class="token punctuation">,</span> <span class="token number">8</span><span class="token punctuation">}</span>

s1  <span class="token operator">=</span>  <span class="token function">append</span><span class="token punctuation">(</span>s1<span class="token punctuation">,</span> <span class="token number">6</span><span class="token punctuation">)</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"t1:"</span><span class="token punctuation">,</span> s1<span class="token punctuation">)</span>

<span class="token comment">// t1: [1 2 3 4 5 6]</span>

s1  <span class="token operator">=</span>  <span class="token function">append</span><span class="token punctuation">(</span>s1<span class="token punctuation">,</span> s2<span class="token operator">...</span><span class="token punctuation">)</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"t2:"</span><span class="token punctuation">,</span> s1<span class="token punctuation">)</span>

  

<span class="token comment">// t2: [1 2 3 4 5 6 7 8]</span>

  

<span class="token keyword">var</span>  s <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">rune</span>

<span class="token keyword">for</span>  <span class="token boolean">_</span><span class="token punctuation">,</span> v  <span class="token operator">:=</span>  <span class="token keyword">range</span>  <span class="token string">"hello world"</span> <span class="token punctuation">{</span>

s  <span class="token operator">=</span>  <span class="token function">append</span><span class="token punctuation">(</span>s<span class="token punctuation">,</span> v<span class="token punctuation">)</span>

<span class="token punctuation">}</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"s:"</span><span class="token punctuation">,</span> s<span class="token punctuation">)</span>

fmt<span class="token punctuation">.</span><span class="token function">Printf</span><span class="token punctuation">(</span><span class="token string">"%q\n"</span><span class="token punctuation">,</span> s<span class="token punctuation">)</span>

  

<span class="token comment">//重点：append 扩容，当长度&lt;容量时，不需要扩容；反之，两倍长度扩容</span>
</code></pre>
<h3 id="slice比较">4.slice比较</h3>
<ul>
<li><code>reflect</code>比较的方法，代码简单，性能低</li>
<li>循环遍历比较的方法，代码多，性能高</li>
</ul>
<p>reflect比较法：</p>
<pre class=" language-go"><code class="prism  language-go"></code></pre>
<p>func StrSliceReflectEqual(a, b []string) bool {<br>
return reflect.DeepEqual(a, b)<br>
}</p>
<pre><code></code></pre>
<p>循环比较法：</p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span>  <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>

s1  <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">{</span><span class="token string">"a"</span><span class="token punctuation">,</span> <span class="token string">"b"</span><span class="token punctuation">}</span>

s2  <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">{</span><span class="token string">"a"</span><span class="token punctuation">,</span> <span class="token string">"c"</span><span class="token punctuation">}</span>

ret  <span class="token operator">:=</span>  <span class="token function">StrSliceEqual</span><span class="token punctuation">(</span>s1<span class="token punctuation">,</span> s2<span class="token punctuation">)</span>

fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>ret<span class="token punctuation">)</span>

<span class="token punctuation">}</span>

<span class="token keyword">func</span>  <span class="token function">StrSliceEqual</span><span class="token punctuation">(</span>a<span class="token punctuation">,</span> b <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>

<span class="token comment">//先比较长度是否相等,否则为false</span>

<span class="token keyword">if</span>  <span class="token function">len</span><span class="token punctuation">(</span>a<span class="token punctuation">)</span> <span class="token operator">!=</span>  <span class="token function">len</span><span class="token punctuation">(</span>b<span class="token punctuation">)</span> <span class="token punctuation">{</span>

<span class="token keyword">return</span>  <span class="token boolean">false</span>

<span class="token punctuation">}</span>

<span class="token comment">//再比较两个slice是否 都为nil或都不为nil,否则为false</span>

<span class="token keyword">if</span> <span class="token punctuation">(</span>a <span class="token operator">==</span>  <span class="token boolean">nil</span><span class="token punctuation">)</span> <span class="token operator">!=</span> <span class="token punctuation">(</span>b <span class="token operator">==</span>  <span class="token boolean">nil</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>

<span class="token keyword">return</span>  <span class="token boolean">false</span>

<span class="token punctuation">}</span>

<span class="token comment">//再比较对应索引处两个slice的元素是否相等,否则为false</span>

<span class="token keyword">for</span>  i<span class="token punctuation">,</span> v  <span class="token operator">:=</span>  <span class="token keyword">range</span> a <span class="token punctuation">{</span>

<span class="token keyword">if</span> v <span class="token operator">!=</span> b<span class="token punctuation">[</span>i<span class="token punctuation">]</span> <span class="token punctuation">{</span>

<span class="token keyword">return</span>  <span class="token boolean">false</span>

<span class="token punctuation">}</span>

<span class="token punctuation">}</span>

<span class="token keyword">return</span>  <span class="token boolean">true</span>

<span class="token punctuation">}</span>
</code></pre>
<h3 id="slice-删除">5.slice 删除</h3>
<pre><code>index  :=  3  //删除元素的索引

var  s  = []int{1, 2, 3, 4, 5, 6}

s  =  append(s[:index], s[index+1:]...)

fmt.Println(s)

//[1 2 3 5 6]
</code></pre>
<blockquote>
<p>reference:<a href="https://www.jianshu.com/p/80f5f5173fca">addr</a></p>
</blockquote>

