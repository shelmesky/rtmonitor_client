### RTMonitor 客户端

RTMonitor 是使用Go语言编写的，内嵌到Go代码中，支持收集进程、操作系统、Golang运行时的信息库。
收集到数据后，会将数据传送给API Server，在RTMonitor服务端中可以使用图表的形式方便查看。

目前支持收集以下信息：

+ OS的内存使用
+ OS的负载
+ 进程的内存使用
+ Golang运行时的内存使用
+ Golang运行时的GC状态 

### 安装方式

```
go get github.com/shelmesky/rtmonitor_client
```

在代码中使用：

```
rtmonitor_client.Start("http://apiserver/api/", "CLIENT_KEY")
```

`CLIENT_KEY`是在RTMonitor服务端中预先创建的KEY。

### 进程的内存占用

```
VirtualMemory   // 进程占用的虚拟内存字节数
ResidentMemory  // 虚拟内存映射到物理页面的字节数
SharedMemory    // 共享内存字节数
```


### 运行时状态信息汇总

```
// General statistics
Alloc       // 正在使用的字节数
Sys         // 向操作系统申请的总内存数量(RSS)
Mallocs     // malloc的次数
Frees       // free的次数

// 堆内存分配
HeapAlloc       // 当前堆中已经分配的字节数
HeapSys         // 进程向OS申请的堆内存字节数
HeapIdle        // span中空闲的字节数
HeapInuse       // span中非空闲的字节数
HeapObjects     // 当前对中的对象数量

// 栈内存分配
StackInuse      // 栈上使用的字节数
StackSys        // 栈申请的字节数

// GC status
GCPause             // 当前GC的暂停时间
GCPausePerSecond    // 平均每秒GC暂停的时间
GCPerSecond         // 平均每秒发生GC的次数
GCTotalPause        // 总体的GC时间

//Num of goroutines
Goroutines          // goroutine的数量
```
