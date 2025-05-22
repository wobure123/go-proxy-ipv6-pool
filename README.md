# Go Proxy IPV6 Pool

Random ipv6 egress proxy server (support http/socks5, support username/password authentication)

The Go language implementation of [zu1k/http-proxy-ipv6-pool](https://github.com/zu1k/http-proxy-ipv6-pool)

## Usage

```bash
    go run . --port <port> [--socks5-user <username>] [--socks5-pass <password>]
    # 只会使用本机已分配到网卡的IPv6地址作为出口
    # 默认用户名/密码为 user/pass，如未指定则使用默认
```

### Use as a proxy server

```bash
    curl -x http://<your-ip>:52122 http://6.ipw.cn/ # 2001:399:8205:ae00:456a:ab12 (random ipv6 address)
```

```bash
    curl -x socks5h://<username>:<password>@<your-ip>:52123 http://6.ipw.cn/ # 2001:399:8205:ae00:456a:ab12 (random ipv6 address)
    # 例如: curl -x socks5h://user:pass@127.0.0.1:52123 http://6.ipw.cn/
```

- socks5 代理支持用户名密码认证，推荐使用 socks5h 协议以确保 DNS 解析也走代理。
- 只会从本机网卡已分配的 IPv6 地址池中随机选取出口地址。

## License

MIT License (see [LICENSE](LICENSE))
