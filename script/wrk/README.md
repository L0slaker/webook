## 压测接口
 1. 注册：写为主的接口
 2. 登录：写为主的接口
 3. Profile：读为主的接口

## 使用命令
    wrk -t1 -d1s -c2 -s ./script/wrk/signup.lua http://localhost:8080/api/v1/user/signup
    wrk -t1 -d1s -c2 -s ./script/wrk/login.lua http://localhost:8080/api/v1/user/login
    wrk -t1 -d1s -c2 -s ./script/wrk/profile.lua http://localhost:8080/api/v1/user/info
    -t 后跟着的是线程数量
    -d 后跟着的是持续时间
    -c 后跟着的是并发数
    -s 后跟着的是测试的脚本

## 压测前准备
 1. 启动JWT来测试：相对Session来说比较好测试
 2. 修改登录态时长：保证在测试Profile接口时登录态不会过期
 3. 去除限流的限制

## 结果集数据
```c
Running 1s test @ http://localhost:8080/api/v1/user/signup
1 threads and 2 connections
Thread Stats   Avg      Stdev     Max   +/- Stdev
Latency       64.01ms   2.51ms  72.50ms   86.67%
Req/Sec       30.00     9.43    40.00     20.00%
30 requests in 1.00s, 3.87KB read
Requests/sec:     29.92
Transfer/sec:      3.86KB

# Latency 延迟
## ±2.51ms	延迟的标准差（波动范围），越小说明响应越稳定。
## 72.50ms	最大延迟。
## 86.67%	延迟分布：86.67% 的请求延迟在 (Avg ± Stdev) 范围内（64.01ms±2.51ms）。

# Req/Sec：平均每秒完成的请求数（QPS）
## ±9.43	QPS 的标准差，波动较大（可能因网络或服务不稳定）。
## 40.00	峰值 QPS（瞬时最高值）。
## 20.00%	QPS 分布：仅 20% 的时间能达到峰值 QPS。

# Requests/sec	29.92	实际平均 QPS（与 Thread Stats 中的 Req/Sec 略有差异，因统计维度不同）。
# Transfer/sec	3.86KB	平均每秒传输的数据量（带宽占用）。
```