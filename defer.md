
<h2 id="总结defer常见四种用法">总结defer常见四种用法</h2>
<h3 id="一.先进后出原则">一.先进后出原则</h3>

~~~go
func a() {
	for i := 0; i < 3; i++ {
		defer func() {
			fmt.Println(i) //333 延迟函数执行时机，a()函数执行结束前调用，执行结束前i=3
		}()
	}
}
func b() {
	for i := 0; i < 3; i++ {
		defer func(i int) {
			fmt.Println(i) //210  FILO，值传递
		}(i)
	}
}

~~~

<h3 id="二.defer-函数内部所使用的变量的值在这个函数运行时才确定">二.defer 函数内部所使用的变量的值在这个函数运行时才确定</h3>

~~~ go
func c() {
	for i := 0; i < 3; i++ {
		defer func(i int) {
			fmt.Println(i) //420  FILO，
			//defer 函数内部所使用的变量的值在这个函数运行时才确定
		}(i * 2)
	}
}
~~~
<h3 id="三.defer函数中的参数值在定义时即计算">三.defer函数中的参数值在定义时即计算</h3>

~~~go
func d() (i int) {
	//defer 调用的函数参数的值，在defer被定义时就确定了
	i = 1
	defer fmt.Println("d() defer i:", i) //1
	i++
	fmt.Println("d()....:", i) //2
	return i                   //2
}

~~~

<h3 id="四、先执行完defer，后执行return，最后函数返回值">四、先执行完defer，后执行return，最后函数返回值</h3>

~~~go
func e() (i int) {
	//先执行完defer后 执行return，最后函数返回值
	defer func() {
		fmt.Println("e() defer i1:", i) //0
		i++
		fmt.Println("e() defer i2:", i) //1
	}()
	return 0 //1
}

~~~

<p><strong>全部代码：</strong></p>

~~~go
package main

import (
	"fmt"
)

func main() {
	a()
	b()
	c()
	r := d()
	fmt.Println("d() ret i:", r) //2
	r2 := e()
	fmt.Println("e() ret i:", r2) //1

}
func a() {
	for i := 0; i < 3; i++ {
		defer func() {
			fmt.Println(i) //333 延迟函数执行时机，a()函数执行结束前调用，执行结束前i=3
		}()
	}
}
func b() {
	for i := 0; i < 3; i++ {
		defer func(i int) {
			fmt.Println(i) //210  FILO，值传递
		}(i)
	}
}
func c() {
	for i := 0; i < 3; i++ {
		defer func(i int) {
			fmt.Println(i) //420  FILO，
			//defer 函数内部所使用的变量的值需要在这个函数运行时才确定
		}(i * 2)
	}
}
func d() (i int) {
	//defer 调用的函数参数的值,在 defer 被定义时就确定了
	i = 1
	defer fmt.Println("d() defer i:", i) //1
	i++
	fmt.Println("d()....:", i) //2
	return i                   //2
}
func e() (i int) {
	//先执行完defer后 执行return，最后函数返回值
	defer func() {
		fmt.Println("e() defer i1:", i) //0
		i++
		fmt.Println("e() defer i2:", i) //1
	}()
	return 0 //1
}
~~~

