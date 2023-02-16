# 简易http服务反向代理器
可以用作网关，转发各类http流量到指定的微服务端口

## 使用方法
编写一个配置文件config.json，通过命令行启动
```bash
./httpProxy --path "config path"
```
## 参数
1. --path 用作指定配置文件config的地址
    ```json
    {
        "port": ":8082", // 指定监听端口
        "proxyConfig": [ // 反向代理地址，以及前缀
            {
                "urlStr":"http://bing.com",
                "prefix":"/bing/"
            }
        ]
    }
    ```