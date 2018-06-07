---


---

<h2 id="一-源码">一 源码</h2>
<p><code>地址：https://github.com/golang/go/tree/master/src/reflect</code></p>
<h2 id="二.重要接口与类">二.重要接口与类</h2>
<p><strong>golang中reflect包</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">package</span> reflect
</code></pre>
<p><strong>重要接口</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> Type <span class="token keyword">interface</span><span class="token punctuation">{</span><span class="token punctuation">}</span>
</code></pre>
<p><strong>重要类</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> Value <span class="token keyword">struct</span><span class="token punctuation">{</span><span class="token punctuation">}</span>
</code></pre>
<p><strong>关键方法</strong>：</p>
<p><strong>与Type有关</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span> <span class="token function">TypeOf</span><span class="token punctuation">(</span>i <span class="token keyword">interface</span><span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">)</span> Type
</code></pre>
<p><strong>与Value有关</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span> <span class="token function">ValueOf</span><span class="token punctuation">(</span>i <span class="token keyword">interface</span><span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">)</span> Value
</code></pre>
<h2 id="三.-reflect-作用">三. reflect 作用</h2>
<h3 id="获取变量对象类型">获取变量对象类型</h3>
<p>golang强类型语言，函数如果想接收所有的类型,定义参数类型为interface{}，要想获取变量真实类型，通过reflect</p>
<pre class=" language-go"><code class="prism  language-go">
<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">{</span>
<span class="token function">TestGetType</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">GetType</span><span class="token punctuation">(</span>i <span class="token keyword">interface</span><span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">)</span> <span class="token builtin">string</span> <span class="token punctuation">{</span>  
   val <span class="token operator">:=</span> reflect<span class="token punctuation">.</span><span class="token function">ValueOf</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span>  
   <span class="token keyword">return</span> val<span class="token punctuation">.</span><span class="token function">Kind</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">String</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
<span class="token punctuation">}</span>  
  
<span class="token keyword">func</span> <span class="token function">TestGetType</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   i <span class="token operator">:=</span> <span class="token number">100</span>  
  s <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">{</span><span class="token string">"tom"</span><span class="token punctuation">,</span> <span class="token string">"jim"</span><span class="token punctuation">}</span>  
   m <span class="token operator">:=</span> <span class="token keyword">map</span><span class="token punctuation">[</span><span class="token builtin">string</span><span class="token punctuation">]</span><span class="token builtin">string</span><span class="token punctuation">{</span><span class="token string">"tom"</span><span class="token punctuation">:</span> <span class="token string">"beijing"</span><span class="token punctuation">,</span> <span class="token string">"jim"</span><span class="token punctuation">:</span> <span class="token string">"shanghai"</span><span class="token punctuation">}</span>  
   fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token function">GetType</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token comment">//int  </span>
  fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token function">GetType</span><span class="token punctuation">(</span>s<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token comment">//slice  </span>
  fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token function">GetType</span><span class="token punctuation">(</span>m<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token comment">//map  </span>
<span class="token punctuation">}</span>
</code></pre>
<h3 id="section"></h3>
<h3 id="获取对象方法、字段">获取对象方法、字段</h3>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">// 获取一个对象的字段和方法  </span>
<span class="token keyword">func</span> <span class="token function">TestStruct</span><span class="token punctuation">(</span>i <span class="token keyword">interface</span><span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   t <span class="token operator">:=</span> reflect<span class="token punctuation">.</span><span class="token function">TypeOf</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span>  
   val <span class="token operator">:=</span> reflect<span class="token punctuation">.</span><span class="token function">ValueOf</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span>  
   kind <span class="token operator">:=</span> val<span class="token punctuation">.</span><span class="token function">Kind</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
  
  
   <span class="token keyword">if</span> kind <span class="token operator">!=</span> reflect<span class="token punctuation">.</span>Ptr <span class="token operator">&amp;&amp;</span> val<span class="token punctuation">.</span><span class="token function">Elem</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Kind</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token operator">==</span> reflect<span class="token punctuation">.</span>Struct <span class="token punctuation">{</span>  
      fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"error struct"</span><span class="token punctuation">)</span>  
      <span class="token keyword">return</span>  
  <span class="token punctuation">}</span>  
   <span class="token comment">//获取字段数量  </span>
  fields <span class="token operator">:=</span> val<span class="token punctuation">.</span><span class="token function">Elem</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">NumField</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
   fmt<span class="token punctuation">.</span><span class="token function">Printf</span><span class="token punctuation">(</span><span class="token string">"struct 共有 %d fields\n"</span><span class="token punctuation">,</span> fields<span class="token punctuation">)</span>  
   <span class="token comment">//获取字段的类型  </span>
  <span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> fields<span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>  
  
      fmt<span class="token punctuation">.</span><span class="token function">Printf</span><span class="token punctuation">(</span><span class="token string">"第 %d 字段，类型：%v ,"</span><span class="token punctuation">,</span>i<span class="token punctuation">,</span>val<span class="token punctuation">.</span><span class="token function">Elem</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Field</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Kind</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
      <span class="token comment">//获取tag  </span>
  tag <span class="token operator">:=</span> t<span class="token punctuation">.</span><span class="token function">Elem</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Field</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span><span class="token punctuation">.</span>Tag<span class="token punctuation">.</span><span class="token function">Get</span><span class="token punctuation">(</span><span class="token string">"json"</span><span class="token punctuation">)</span>  
      fmt<span class="token punctuation">.</span><span class="token function">Printf</span><span class="token punctuation">(</span><span class="token string">"tag:%s\n"</span><span class="token punctuation">,</span> tag<span class="token punctuation">)</span>  
   <span class="token punctuation">}</span>  
   <span class="token comment">//获取方法数量  </span>
  
  fmt<span class="token punctuation">.</span><span class="token function">Printf</span><span class="token punctuation">(</span><span class="token string">"struct 共有 %d 方法:\n"</span><span class="token punctuation">,</span> t<span class="token punctuation">.</span><span class="token function">NumMethod</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
   <span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> t<span class="token punctuation">.</span><span class="token function">NumMethod</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>  
      fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"方法："</span><span class="token punctuation">,</span>i<span class="token punctuation">,</span>t<span class="token punctuation">.</span><span class="token function">Method</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span><span class="token punctuation">.</span>Name<span class="token punctuation">)</span>  
   <span class="token punctuation">}</span>  
  
<span class="token punctuation">}</span>  
  
<span class="token comment">//测试用例  </span>
<span class="token keyword">type</span> Books <span class="token keyword">struct</span> <span class="token punctuation">{</span>  
   Id   <span class="token builtin">int64</span>  <span class="token string">`json:"id"`</span>  
  Name <span class="token builtin">string</span> <span class="token string">`json:"name"`</span>  
<span class="token punctuation">}</span>  
  
<span class="token keyword">func</span> <span class="token punctuation">(</span>s <span class="token operator">*</span>Books<span class="token punctuation">)</span> <span class="token function">BorrowBooks</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"BorrowBooks:"</span><span class="token punctuation">,</span> s<span class="token punctuation">.</span>Name<span class="token punctuation">)</span>  
<span class="token punctuation">}</span>  
  
<span class="token comment">// 接收器为指针类型  </span>
<span class="token keyword">func</span> <span class="token punctuation">(</span>s <span class="token operator">*</span>Books<span class="token punctuation">)</span> <span class="token function">BackBooks</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"BackBooks:"</span><span class="token punctuation">,</span> s<span class="token punctuation">.</span>Name<span class="token punctuation">)</span>  
<span class="token punctuation">}</span>

<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">{</span>
<span class="token function">TestStruct</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>Books<span class="token punctuation">{</span><span class="token punctuation">}</span><span class="token punctuation">)</span>
<span class="token punctuation">}</span>

</code></pre>
<h3 id="类型判断条件">类型判断条件</h3>
<p><strong>只有以下类型可以进一步获取其信息</strong></p>
<pre class=" language-go"><code class="prism  language-go">  <span class="token comment">//github.com/kr/pretty</span>
<span class="token keyword">func</span> <span class="token function">canExpand</span><span class="token punctuation">(</span>t reflect<span class="token punctuation">.</span>Type<span class="token punctuation">)</span> <span class="token builtin">bool</span> <span class="token punctuation">{</span>  
   <span class="token keyword">switch</span> t<span class="token punctuation">.</span><span class="token function">Kind</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
   <span class="token keyword">case</span> reflect<span class="token punctuation">.</span>Map<span class="token punctuation">,</span> reflect<span class="token punctuation">.</span>Struct<span class="token punctuation">,</span>  
      reflect<span class="token punctuation">.</span>Interface<span class="token punctuation">,</span> reflect<span class="token punctuation">.</span>Array<span class="token punctuation">,</span> reflect<span class="token punctuation">.</span>Slice<span class="token punctuation">,</span>  
      reflect<span class="token punctuation">.</span>Ptr<span class="token punctuation">:</span>  
      <span class="token keyword">return</span> <span class="token boolean">true</span>  
  <span class="token punctuation">}</span>  
   <span class="token keyword">return</span> <span class="token boolean">false</span>  
<span class="token punctuation">}</span>
</code></pre>
<h3 id="elem.set使用">Elem().Set使用</h3>
<pre class=" language-go"><code class="prism  language-go">
<span class="token keyword">func</span> <span class="token function">test</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
  
   ch <span class="token operator">:=</span> <span class="token function">make</span><span class="token punctuation">(</span><span class="token keyword">chan</span> <span class="token builtin">bool</span><span class="token punctuation">,</span> <span class="token number">1</span><span class="token punctuation">)</span>  
   i <span class="token operator">:=</span> <span class="token number">0</span>  
  v <span class="token operator">:=</span> reflect<span class="token punctuation">.</span><span class="token function">ValueOf</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>i<span class="token punctuation">)</span>  
   <span class="token keyword">go</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>  
      v<span class="token punctuation">.</span><span class="token function">Elem</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Set</span><span class="token punctuation">(</span>reflect<span class="token punctuation">.</span><span class="token function">ValueOf</span><span class="token punctuation">(</span><span class="token number">1</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
      fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"i:"</span><span class="token punctuation">,</span>i<span class="token punctuation">)</span>  
      ch <span class="token operator">&lt;-</span> <span class="token boolean">true</span>  
  <span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>  
   <span class="token comment">//用来获取指针指向的变量  </span>
  v<span class="token punctuation">.</span><span class="token function">Elem</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">.</span><span class="token function">Set</span><span class="token punctuation">(</span>reflect<span class="token punctuation">.</span><span class="token function">ValueOf</span><span class="token punctuation">(</span><span class="token number">2</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
   <span class="token comment">//_= v.Elem().Int()  </span>
  fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"type:"</span><span class="token punctuation">,</span>v<span class="token punctuation">.</span><span class="token function">Type</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
  fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"value:"</span><span class="token punctuation">,</span>i<span class="token punctuation">)</span>  
  fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"kind:"</span><span class="token punctuation">,</span>v<span class="token punctuation">.</span><span class="token function">Kind</span><span class="token punctuation">(</span><span class="token punctuation">)</span><span class="token punctuation">)</span>  
  
   <span class="token operator">&lt;-</span>ch  
  
<span class="token punctuation">}</span>
</code></pre>
<blockquote>
<p>reference:<br>
<a href="https://github.com/golang/go/blob/master/src/reflect/">https://github.com/golang/go/blob/master/src/reflect/</a><br>
<a href="https://jimmyfrasche.github.io/go-reflection-codex/">https://jimmyfrasche.github.io/go-reflection-codex/</a></p>
</blockquote>

