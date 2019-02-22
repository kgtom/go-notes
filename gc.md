
## 大纲
* [一.gc 源码解析](#1)
* [二.研发中gc问题排查和定位](#2)



## <span id="1">一.gc 源码解析</span>



## <span id="1">二.研发中gc问题排查和定位</span>

### 排查
* GODEBUG=gctrace=1 代表只针对当前进程开启gc追踪

~~~

# tom @ tom-pc in ~/goprojects/src/example [23:20:26] 
$ GODEBUG=gctrace=1  go run main.go
gc 1 @0.039s 0%: 0.012+0.48+0.039 ms clock, 0.050+0.47/0.32/0.77+0.15 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 2 @0.082s 0%: 0.005+0.42+0.025 ms clock, 0.022+0.20/0.32/0.59+0.10 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 3 @0.139s 0%: 0.024+1.1+0.13 ms clock, 0.097+0.24/0.79/1.7+0.55 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
gc 4 @0.147s 0%: 0.045+0.54+0.10 ms clock, 0.18+0.29/0.49/0.62+0.42 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
gc 5 @0.159s 0%: 0.009+1.3+0.041 ms clock, 0.037+0/0.12/1.4+0.16 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
gc 6 @0.174s 0%: 0.005+0.76+0.14 ms clock, 0.021+0.16/0.58/0.72+0.58 ms cpu, 4->4->1 MB, 5 MB goal, 4 P
# command-line-arguments
gc 1 @0.005s 0%: 0.011+2.4+0.40 ms clock, 0.044+1.2/2.1/3.4+1.6 ms cpu, 4->4->3 MB, 5 MB goal, 4 P
# command-line-arguments
gc 1 @0.001s 0%: 0.017+6.4+0.045 ms clock, 0.070+0.057/6.1/0.54+0.18 ms cpu, 4->5->4 MB, 5 MB goal, 4 P
gc 2 @0.017s 0%: 0.008+4.8+0.036 ms clock, 0.033+0.091/3.8/5.3+0.14 ms cpu, 8->8->7 MB, 9 MB goal, 4 P
gc 3 @0.042s 0%: 0.013+7.0+0.034 ms clock, 0.052+0.51/6.7/0.32+0.13 ms cpu, 14->15->14 MB, 15 MB goal, 4 P
gc 4 @0.084s 0%: 0.008+13+0.039 ms clock, 0.033+0.065/13/0.27+0.15 ms cpu, 26->29->26 MB, 29 MB goal, 4 P

~~~

* 参数解释

~~~
gc 1 @0.039s 0%: 0.012+0.48+0.039 ms clock, 0.050+0.47/0.32/0.77+0.15 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
~~~
 - gc 1：第一次执行
 - @0.039s: 程序执行总时间
 - 0%：gc时占总运行时间的百分比
 - 0.012+0.48+0.039 ms clock：gc时间
 - 0.050+0.47/0.32/0.77+0.15 ms cpu：gc占用cpu时间
 - 4->4->0 MB：堆的大小->gc后堆的大小->存活数据堆的大小
 - 5 MB goal：全局堆大小
 - 4 P：使用处理器数量

* 使用GCTRACE跟踪了gc信息，但具体哪个函数消耗大量内存，需要使用 pprof进行定位

### pprof 定位问题

[详见](http://39.106.173.209:88/go-web-issue/#3)
