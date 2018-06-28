---


---

<h2 id="本章内容">本章内容</h2>
<ul>
<li><a href="#1">一、string</a></li>
<li><a href="#2">二、byte</a></li>
<li><a href="#3">三、rune</a></li>
<li><a href="#4">四、字符串编码及长度</a></li>
<li><a href="#5">五、字符串修改</a></li>
<li><a href="#6">六、字符串遍历</a></li>
</ul>
<h2 id="span-id1一、stringspan"><span id="1">一、string</span></h2>
<p><strong>源码地址：</strong> <code>/src/builtin/builtin.go</code></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">// string is the set of all strings of 8-bit bytes, conventionally but not</span>
<span class="token comment">// necessarily representing UTF-8-encoded text. A string may be empty, but</span>
<span class="token comment">// not nil. Values of string type are immutable.</span>
<span class="token keyword">type</span> <span class="token builtin">string</span> <span class="token builtin">string</span>
</code></pre>
<p><strong>解析</strong> 8位字节集合，不一定是UTF-8,可以为空，不能是nil，值不可以改变。</p>
<p><strong>源码地址：</strong> <code>/src/runtime/string.go</code></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">type</span> stringStruct <span class="token keyword">struct</span> <span class="token punctuation">{</span>
	str unsafe<span class="token punctuation">.</span>Pointer
	<span class="token builtin">len</span> <span class="token builtin">int</span>
<span class="token punctuation">}</span>


<span class="token keyword">func</span> <span class="token function">gostringnocopy</span><span class="token punctuation">(</span>str <span class="token operator">*</span><span class="token builtin">byte</span><span class="token punctuation">)</span> <span class="token builtin">string</span> <span class="token punctuation">{</span>
	ss <span class="token operator">:=</span> stringStruct<span class="token punctuation">{</span>str<span class="token punctuation">:</span> unsafe<span class="token punctuation">.</span><span class="token function">Pointer</span><span class="token punctuation">(</span>str<span class="token punctuation">)</span><span class="token punctuation">,</span> <span class="token builtin">len</span><span class="token punctuation">:</span> <span class="token function">findnull</span><span class="token punctuation">(</span>str<span class="token punctuation">)</span><span class="token punctuation">}</span>
	s <span class="token operator">:=</span> <span class="token operator">*</span><span class="token punctuation">(</span><span class="token operator">*</span><span class="token builtin">string</span><span class="token punctuation">)</span><span class="token punctuation">(</span>unsafe<span class="token punctuation">.</span><span class="token function">Pointer</span><span class="token punctuation">(</span><span class="token operator">&amp;</span>ss<span class="token punctuation">)</span><span class="token punctuation">)</span>
	<span class="token keyword">return</span> s
<span class="token punctuation">}</span>

</code></pre>
<p><strong>解析</strong> string 底层是一个byte数组，string也是一个struct。</p>
<h2 id="span-id2二、bytespan"><span id="2">二、byte</span></h2>
<p><strong>源码地址：</strong> <code>/src/builtin/builtin.go</code></p>
<pre class=" language-go"><code class="prism  language-go">
<span class="token comment">// byte is an alias for uint8 and is equivalent to uint8 in all ways. It is</span>
<span class="token comment">// used, by convention, to distinguish byte values from 8-bit unsigned</span>
<span class="token comment">// integer values.</span>
<span class="token keyword">type</span> <span class="token builtin">byte</span> <span class="token operator">=</span> <span class="token builtin">uint8</span>

</code></pre>
<p><strong>解析</strong> byte是uint8的别名，byte类型的底层类型是int8类型</p>
<h2 id="span-id3三、runespan"><span id="3">三、rune</span></h2>
<p><strong>源码地址：</strong> <code>/src/builtin/builtin.go</code></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token comment">// rune is an alias for int32 and is equivalent to int32 in all ways. It is</span>
<span class="token comment">// used, by convention, to distinguish character values from integer values.</span>
<span class="token keyword">type</span> <span class="token builtin">rune</span> <span class="token operator">=</span> <span class="token builtin">int32</span>

</code></pre>
<p><strong>解析</strong>rune类型的底层类型是int32类型,比byte表现数值也多。</p>
<h2 id="span-id4四、字符串编码及长度span"><span id="4">四、字符串编码及长度</span></h2>
<p>在unicode中，一个中文占两个字节，utf-8中一个中文占三个字节，golang支持编码是utf-8编码，默认中文三个字符</p>
<pre class=" language-go"><code class="prism  language-go">    str <span class="token operator">:=</span> <span class="token string">"hello 中国"</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token function">len</span><span class="token punctuation">(</span>str<span class="token punctuation">)</span><span class="token punctuation">)</span>                    <span class="token comment">//12 默认底层byte，utf-8 一个中文三个字节</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token function">len</span><span class="token punctuation">(</span><span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token function">rune</span><span class="token punctuation">(</span>str<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token punctuation">)</span>            <span class="token comment">//8 unicode编码，一个中文2个字节</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>utf8<span class="token punctuation">.</span><span class="token function">RuneCountInString</span><span class="token punctuation">(</span>str<span class="token punctuation">)</span><span class="token punctuation">)</span> <span class="token comment">//8 unicode/utf8</span>
</code></pre>
<h2 id="span-id5五、字符串修改span"><span id="5">五、字符串修改</span></h2>
<p>因为string源码中告诉,字符串的值不能修改，所以我们需要将string转换为[]byte或[]rune，实质重新申请一块内存地址进行修改。</p>
<pre class=" language-go"><code class="prism  language-go">    str <span class="token operator">:=</span> <span class="token string">"hello 中国"</span>
	str2 <span class="token operator">:=</span> <span class="token punctuation">[</span><span class="token punctuation">]</span><span class="token function">rune</span><span class="token punctuation">(</span>str<span class="token punctuation">)</span>
	t <span class="token operator">:=</span> str2
	t<span class="token punctuation">[</span><span class="token number">0</span><span class="token punctuation">]</span> <span class="token operator">=</span> <span class="token string">'w'</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"str2"</span><span class="token punctuation">,</span> <span class="token function">string</span><span class="token punctuation">(</span>str2<span class="token punctuation">)</span><span class="token punctuation">)</span> <span class="token comment">//wello 中国</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"t"</span><span class="token punctuation">,</span> <span class="token function">string</span><span class="token punctuation">(</span>t<span class="token punctuation">)</span><span class="token punctuation">)</span><span class="token comment">//wello 中国</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"str"</span><span class="token punctuation">,</span> str<span class="token punctuation">)</span><span class="token comment">//hello 中国</span>
</code></pre>
<p><strong>解析</strong><br>
str2 是一个rune类型切片的引用，相对于str重新申请一块内存地址，str2修改不影响str。</p>
<h2 id="span-id6六、字符串遍历span"><span id="6">六、字符串遍历</span></h2>

~~~go
       str := "hello 中国"

	//第一种 实质byte,utf-8编码遍历
	for i := 0; i < len(str); i++ {
		fmt.Println(str[i])
	}
	//第二种 使用range隐式转换成Unicode
	for _, v := range str {
		//fmt.Printf("%c,", v)
		fmt.Println(string(v))
	}
	//第三种 先转换成[]rune,再遍历
	str2 := []rune(str)
	for i := 0; i < len(str2); i++ {
		fmt.Printf("%c,", str2[i])
		//fmt.Println(string(str2[i]))
	}

~~~
<blockquote>
<p>reference:<br>
<a href="https://github.com/golang/go/blob/2f2e8f9c81d899d8eefb1f2f98ce5c90976c4f61/src/runtime/string.go">github/golang</a></p>
</blockquote>

