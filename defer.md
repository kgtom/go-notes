---


---

<h2 id="总结defer常见四种用法">总结defer常见四种用法</h2>
<h3 id="一.先进后出原则">一.先进后出原则</h3>
<pre class=" language-go"><code class="prism  language-go">
<span class="token keyword">func</span> <span class="token function">a</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> <span class="token number">3</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
		<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span> <span class="token comment">//333 延迟函数执行时机，a()函数执行结束前调用，执行结束前i=3</span>
		<span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">b</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> <span class="token number">3</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
		<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span> <span class="token comment">//210  FILO，值传递</span>
		<span class="token punctuation">}</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre>
<h3 id="二.defer-函数内部所使用的变量的值在这个函数运行时才确定">二.defer 函数内部所使用的变量的值在这个函数运行时才确定</h3>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">func</span> <span class="token function">c</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> <span class="token number">3</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
		<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span> <span class="token comment">//420  FILO，</span>
			<span class="token comment">//defer 函数内部所使用的变量的值在这个函数运行时才确定</span>
		<span class="token punctuation">}</span><span class="token punctuation">(</span>i <span class="token operator">*</span> <span class="token number">2</span><span class="token punctuation">)</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>

</code></pre>
<h3 id="三.defer函数中的参数值在定义时即计算">三.defer函数中的参数值在定义时即计算</h3>
<pre class=" language-go"><code class="prism  language-go">  <span class="token keyword">func</span> <span class="token function">d</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token comment">//defer 调用的函数参数的值，在defer被定义时就确定了</span>
	i <span class="token operator">=</span> <span class="token number">1</span>
	<span class="token keyword">defer</span> fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"d() defer i:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//1</span>
	i<span class="token operator">++</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"d()....:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//2</span>
	<span class="token keyword">return</span> i                   <span class="token comment">//2</span>
<span class="token punctuation">}</span>
</code></pre>
<h3 id="四、先执行完defer，后执行return，最后函数返回值">四、先执行完defer，后执行return，最后函数返回值</h3>
<pre class=" language-go"><code class="prism  language-go">  <span class="token keyword">func</span> <span class="token function">e</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token comment">//先执行完defer后 执行return，最后函数返回值</span>
	<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
		fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"e() defer i1:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//0</span>
		i<span class="token operator">++</span>
		fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"e() defer i2:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//1</span>
	<span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	<span class="token keyword">return</span> <span class="token number">0</span> <span class="token comment">//1</span>
<span class="token punctuation">}</span>
</code></pre>
<p><strong>全部代码：</strong></p>
<pre class=" language-go"><code class="prism  language-go"><span class="token keyword">package</span> main

<span class="token keyword">import</span> <span class="token punctuation">(</span>
	<span class="token string">"fmt"</span>
<span class="token punctuation">)</span>

<span class="token keyword">func</span> <span class="token function">main</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token function">a</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	<span class="token function">b</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	<span class="token function">c</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	r <span class="token operator">:=</span> <span class="token function">d</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"d() ret i:"</span><span class="token punctuation">,</span> r<span class="token punctuation">)</span> <span class="token comment">//2</span>
	r2 <span class="token operator">:=</span> <span class="token function">e</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"e() ret i:"</span><span class="token punctuation">,</span> r2<span class="token punctuation">)</span> <span class="token comment">//1</span>

<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">a</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> <span class="token number">3</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
		<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span> <span class="token comment">//333 延迟函数执行时机，a()函数执行结束前调用，执行结束前i=3</span>
		<span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">b</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> <span class="token number">3</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
		<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span> <span class="token comment">//210  FILO，值传递</span>
		<span class="token punctuation">}</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">c</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token keyword">for</span> i <span class="token operator">:=</span> <span class="token number">0</span><span class="token punctuation">;</span> i <span class="token operator">&lt;</span> <span class="token number">3</span><span class="token punctuation">;</span> i<span class="token operator">++</span> <span class="token punctuation">{</span>
		<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
			fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span>i<span class="token punctuation">)</span> <span class="token comment">//420  FILO，</span>
			<span class="token comment">//defer 函数内部所使用的变量的值需要在这个函数运行时才确定</span>
		<span class="token punctuation">}</span><span class="token punctuation">(</span>i <span class="token operator">*</span> <span class="token number">2</span><span class="token punctuation">)</span>
	<span class="token punctuation">}</span>
<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">d</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token comment">//defer 调用的函数参数的值,在 defer 被定义时就确定了</span>
	i <span class="token operator">=</span> <span class="token number">1</span>
	<span class="token keyword">defer</span> fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"d() defer i:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//1</span>
	i<span class="token operator">++</span>
	fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"d()....:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//2</span>
	<span class="token keyword">return</span> i                   <span class="token comment">//2</span>
<span class="token punctuation">}</span>
<span class="token keyword">func</span> <span class="token function">e</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">(</span>i <span class="token builtin">int</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
	<span class="token comment">//先执行完defer后 执行return，最后函数返回值</span>
	<span class="token keyword">defer</span> <span class="token keyword">func</span><span class="token punctuation">(</span><span class="token punctuation">)</span> <span class="token punctuation">{</span>
		fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"e() defer i1:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//0</span>
		i<span class="token operator">++</span>
		fmt<span class="token punctuation">.</span><span class="token function">Println</span><span class="token punctuation">(</span><span class="token string">"e() defer i2:"</span><span class="token punctuation">,</span> i<span class="token punctuation">)</span> <span class="token comment">//1</span>
	<span class="token punctuation">}</span><span class="token punctuation">(</span><span class="token punctuation">)</span>
	<span class="token keyword">return</span> <span class="token number">0</span> <span class="token comment">//1</span>
<span class="token punctuation">}</span>

</code></pre>

