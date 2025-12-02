package httpProxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	distributed "github.com/hedykan/go-distributedutil"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ProxyConfig struct {
	UrlStr      string
	Prefix      string
	ServiceName string
}

var node *distributed.DistributedNode
var client *clientv3.Client

func init() {
	etcdEndpoints := "localhost:2379"
	var err error

	client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoints},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}
	/*
		 	node = distributed.NewNode("http-proxy")
			node.RegisterService(client, ":8083")
	*/
}

func Serve(config []ProxyConfig, port string) {
	mux := http.NewServeMux()

	for _, v := range config {
		remote, err := url.Parse(v.UrlStr)
		if err != nil {
			panic(err)
		}
		proxy := GoReverseProxy(remote, v.Prefix, v.ServiceName)
		mux.Handle(v.Prefix, proxy)
		log.Println("proxy:", v.Prefix, v.UrlStr)
	}

	log.Println("listen:", port[1:])
	err := http.ListenAndServe(port, mux)
	if err != nil {
		panic(err)
	}
}

func GoReverseProxy(remote *url.URL, prefix string, serviceName string) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(remote)

	// 配置管理器
	proxy.Director = directorFunc(remote, prefix, serviceName)
	// 修改响应头
	proxy.ModifyResponse = modifyResponseFunc()

	return proxy
}

// 修改响应头构造函数
func modifyResponseFunc() func(*http.Response) error {
	return func(response *http.Response) error {
		response.Header.Add("Access-Control-Allow-Origin", "*")
		return nil
	}
}

// 管理器构造函数
func directorFunc(remote *url.URL, prefix string, serviceName string) func(*http.Request) {
	return func(request *http.Request) {
		// path, ok := node.GetServicePath(serviceName, "")
		targetPath := remote.Host
		// if ok {
		// 	targetPath = path
		// }

		targetQuery := remote.RawQuery
		request.URL.Scheme = remote.Scheme
		request.URL.Host = targetPath
		request.Host = targetPath
		request.URL.Path, request.URL.RawPath = joinURLPath(remote, request.URL)
		// 替换前缀
		request.URL.Path = strings.Replace(request.URL.Path, prefix, "/", 1)

		if targetQuery == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = targetQuery + request.URL.RawQuery
		} else {
			request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
		}
		if _, ok := request.Header["User-Agent"]; !ok {
			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
		}
		log.Println("proxy url:", remote, "request.URL.Path:", request.URL.Path, "request.URL.RawQuery:", request.URL.RawQuery)
	}
}

// 拼接url
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
