---


---

<h2 id="runtime-介绍">runtime 介绍</h2>
<pre><code>源码文件地址：/src/runtime/proc.go:15行
Goroutine scheduler  
//go的调度器
 The scheduler's job is to distribute ready-to-run goroutines over worker threads.  
//调度器任务：将goroutne 分配到可用的工作线程上。
 The main concepts are:  
 G - goroutine.  
 M - worker thread, or machine.  
 P - processor, a resource that is required to execute Go code.  
M must have an associated P to execute Go code, however it can be blocked or in a syscall w/o an associated P.
//M 必须关联一个P才能执行go代码。
</code></pre>
<h2 id="bootstrap-启动顺序">bootstrap 启动顺序</h2>
<pre><code>// The bootstrap sequence is:  
//  
// call osinit  //os初始化
// call schedinit  //调度器初始化
// make &amp; queue new G  //生成G，放到G队列，等待运行
// call runtime·mstart //调用mstart方法，启动M
</code></pre>
<h1 id="runtime.osinit">runtime.osinit()</h1>
<pre><code>源码文件地址：/src/runtime/os_darwin.go:50行
func osinit() {  
   // bsdthread_register delayed until end of goenvs so that we  
 // can look at the environment first.  
  ncpu = getncpu()  //获取cup核数
   physPageSize = getPageSize()  
   darwinVersion = getDarwinVersion()  
}
</code></pre>
<p>ncpu的数量也就是P最大数量，也就是可运行G的队列的数量，实际上对并发运行G规模的一种限制。</p>
<h1 id="runtime·schedinit">runtime·schedinit()</h1>
<p>部分源码：</p>
<pre><code>// raceinit must be the first call to race detector.  
// In particular, it must be done before mallocinit below calls racemapshadow.  
_g_ := getg()  
if raceenabled {  
   _g_.racectx, raceprocctx0 = raceinit()  
}  
  
sched.maxmcount = 10000  //M最大数量初始设置
  
tracebackinit()  
moduledataverify()  
stackinit()  
mallocinit()  
mcommoninit(_g_.m)  
alginit()       // maps must not be used before this call  
modulesinit()   // provides activeModules  
typelinksinit() // uses maps, activeModules  
itabsinit()     // uses activeModules  
  
msigsave(_g_.m)  
initSigmask = _g_.m.sigmask  
  
goargs()  
goenvs()  
parsedebugvars()  
gcinit() //GC 初始化工作
  
sched.lastpoll = uint64(nanotime())  
procs := ncpu  //P的最大数量
//从环境变量GOMAXPROCS中获取P的最大数量，并检查n值的有效性
if n, ok := atoi32(gogetenv("GOMAXPROCS")); ok &amp;&amp; n &gt; 0 {  
   procs = n  
}  
//重置P的最大数量
if procresize(procs) != nil {  
   throw("unknown runnable goroutine during bootstrap")  
}
</code></pre>
<ul>
<li>M最大数量初始设置1000，意味着go程序最多可以使用1000内核线程，也可以通过SetMaxThreads(i int)设置</li>
</ul>
<pre><code>func setMaxThreads(in int) (out int) {  
  lock(&amp;sched.lock)  
  out = int(sched.maxmcount)  
  if in &gt; 0x7fffffff { // MaxInt32  
 sched.maxmcount = 0x7fffffff  
 } else {  
     sched.maxmcount = int32(in)  
  }  
  checkmcount()  
  unlock(&amp;sched.lock)  
  return  
}
</code></pre>
<ul>
<li>找到当前G,然后执行一系列初始化工作</li>
<li>从环境变量GOMAXPROCS中获取P的数量进行调整</li>
</ul>
<h1 id="make--queue-new-g">make &amp; queue new G</h1>
<p>生成g，放到队列中，等待运行</p>
<pre><code>func newproc(siz int32, fn *funcval) {  
  argp := add(unsafe.Pointer(&amp;fn), sys.PtrSize)  
  pc := getcallerpc()  
  systemstack(func() {  
     newproc1(fn, (*uint8)(argp), siz, pc)  
  })  
}  
 //newproc函数调用newproc1函数
// Create a new g running fn with narg bytes of arguments starting  
// at argp. callerpc is the address of the go statement that created  
// this. The new g is put on the queue of g's waiting to run.  
func newproc1(fn *funcval, argp *uint8, narg int32, callerpc uintptr) {  
  _g_ := getg()  
 
  if fn == nil {  
     _g_.m.throwing = -1 // do not dump full stacks  
 throw("go of nil func value")  
  }  
  ....
  }
</code></pre>
<h2 id="runtime.mstart--启动m">runtime.mstart  启动M</h2>
<pre><code>  
// Called to start an M.  
//  
// This must not split the stack because we may not even have stack  
// bounds set up yet.  
//  
// May run during STW (because it doesn't have a P yet), so write  
// barriers are not allowed.

func mstart() {  
   _g_ := getg()  
  
   osStack := _g_.stack.lo == 0  
  if osStack {  
      // Initialize stack bounds from system stack.  
 // Cgo may have left stack size in stack.hi.  size := _g_.stack.hi  
      if size == 0 {  
         size = 8192 * sys.StackGuardMultiplier  
  }  
      _g_.stack.hi = uintptr(noescape(unsafe.Pointer(&amp;size)))  
      _g_.stack.lo = _g_.stack.hi - size + 1024  
  }  
   // Initialize stack guards so that we can start calling  
 // both Go and C functions with stack growth prologues.  _g_.stackguard0 = _g_.stack.lo + _StackGuard  
  _g_.stackguard1 = _g_.stackguard0  
   mstart1(0)  
  
   // Exit this thread.  
  if GOOS == "windows" || GOOS == "solaris" || GOOS == "plan9" {  
      // Window, Solaris and Plan 9 always system-allocate  
 // the stack, but put it in _g_.stack before mstart, // so the logic above hasn't set osStack yet.  osStack = true  
  }  
   mexit(osStack)  
}  
  
func mstart1(dummy int32) {
...
}
</code></pre>
<h1 id="runtime.main">runtime.main</h1>
<pre><code>// The main goroutine.  
func main() {  
   g := getg()  
  
   // Racectx of m0-&gt;g0 is used only as the parent of the main goroutine.  
 // It must not be used for anything else.  g.m.g0.racectx = 0  
  //初始化栈的大小
  // Max stack size is 1 GB on 64-bit, 250 MB on 32-bit.  
 // Using decimal instead of binary GB and MB because // they look nicer in the stack overflow failure message.  if sys.PtrSize == 8 {  
      maxstacksize = 1000000000  
  } else {  
      maxstacksize = 250000000  
  }  
  //生成新的M
   // Allow newproc to start new Ms.  
  mainStarted = true  
  
  systemstack(func() {  
      newm(sysmon, nil)  
   })  
  
   // Lock the main goroutine onto this, the main OS thread,  
 // during initialization. Most programs won't care, but a few // do require certain calls to be made by the main thread. // Those can arrange for main.main to run in the main thread // by calling runtime.LockOSThread during initialization // to preserve the lock.  lockOSThread()  
  
   if g.m != &amp;m0 {  
      throw("runtime.main not on m0")  
   }  
  
   runtime_init() // must be before defer  
  if nanotime() == 0 {  
      throw("nanotime returning zero")  
   }  
  
   // Defer unlock so that runtime.Goexit during init does the unlock too.  
  needUnlock := true  
  defer func() {  
      if needUnlock {  
         unlockOSThread()  
      }  
   }()  
  
   // Record when the world started. Must be after runtime_init  
 // because nanotime on some platforms depends on startNano.  runtimeInitTime = nanotime()  
  //执行GC
   gcenable()  
  
   main_init_done = make(chan bool)  
   if iscgo {  
      if _cgo_thread_start == nil {  
         throw("_cgo_thread_start missing")  
      }  
      if GOOS != "windows" {  
         if _cgo_setenv == nil {  
            throw("_cgo_setenv missing")  
         }  
         if _cgo_unsetenv == nil {  
            throw("_cgo_unsetenv missing")  
         }  
      }  
      if _cgo_notify_runtime_init_done == nil {  
         throw("_cgo_notify_runtime_init_done missing")  
      }  
      // Start the template thread in case we enter Go from  
 // a C-created thread and need to create a new thread.  startTemplateThread()  
      cgocall(_cgo_notify_runtime_init_done, nil)  
   }  
  //初始化工作
   fn := main_init // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime  
  fn()  
   close(main_init_done)  
  
   needUnlock = false  
   //解绑M与G
  unlockOSThread()  
  
   if isarchive || islibrary {  
      // A program compiled with -buildmode=c-archive or c-shared  
 // has a main, but it is not executed.  return  
  }  
  //执行main包下面的main，即开始执行项目的入口函数main()
   fn = main_main // make an indirect call, as the linker doesn't know the address of the main package when laying down the runtime  
  fn()  
   if raceenabled {  
      racefini()  
   }  
  
   // Make racy client program work: if panicking on  
 // another goroutine at the same time as main returns, // let the other goroutine finish printing the panic trace. // Once it does, it will exit. See issues 3934 and 20018.  if atomic.Load(&amp;runningPanicDefers) != 0 {  
      // Running deferred functions should not take long.  
  for c := 0; c &lt; 1000; c++ {  
         if atomic.Load(&amp;runningPanicDefers) == 0 {  
            break  
  }  
         Gosched()  
      }  
   }  
   if atomic.Load(&amp;panicking) != 0 {  
      gopark(nil, nil, "panicwait", traceEvGoStop, 1)  
   }  
  
   exit(0)  
   for {  
      var x *int32  
      *x = 0  
  }  
}
</code></pre>
<ul>
<li>初始化栈空间大小</li>
<li>执行main.main</li>
<li>执行GC</li>
</ul>

