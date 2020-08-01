package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func main() {
	httpSimpleDemo()

	// customizedRouter()
}

func httpSimpleDemo() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()                      ////解析参数，默认是不会解析的
		fmt.Println(r.Form)                // map[url_long:[111 222]]
		fmt.Println("path = ", r.URL.Path) // 如 /login
		fmt.Println("scheme = ", r.URL.Scheme)
		//
		//r.Form里面包含了所有请求的参数，比如URL中query-string、POST的数据、PUT的数据，
		//所以当你在URL中的query-string字段和POST冲突时，会保存成一个slice，里面存储了多个值
		//
		// r.Form 是个 url.Values 类型, so, 还可以 r.Form.Get("xxx"), Add, Set, for迭代
		//
		// 还可以 r.FormValue("username")。调用r.FormValue时底层会自动调用r.ParseForm, 只会返回同名参数中的第一个，若参数不存在则返回空字符串
		fmt.Println(r.Form["url_long"]) //[111 222], 是个 slice
		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, " "))
		}
		fmt.Fprintf(w, "hello world")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("method = ", r.Method)
		if r.Method == "GET" {
			// 防止重复提交
			//
			//在模版里面增加了一个隐藏字段token，这个值我们通过MD5(时间戳)来获取唯一值，然后我们把这个值存储到服务器端(session来控制
			curtime := time.Now().Unix()
			m := md5.New()
			io.WriteString(m, strconv.FormatInt(curtime, 10))
			token := fmt.Sprintf("%x", m.Sum(nil))

			tpl, _ := template.ParseFiles("login.gtpl", token) // 相对路径, 并且 页面代码 支持热加载
			log.Println(tpl.Execute(w, token))
		} else if r.Method == "POST" {
			r.ParseForm()

			// 重复提交校验
			token := r.Form.Get("token")
			if token != "" {
				//验证token的合法性
			} else {
				//不存在token报错
			}

			fmt.Println(r.Form)
			//请求的是登录数据，那么执行登录的逻辑判断
			fmt.Println("username:", r.Form["username"])
			fmt.Println("password:", r.Form["password"])

			// 校验
			//
			// 这种方式不好, 如果 username 没有值, 则根本不会在r.Form中产生相应条目, 则 r.Form[username] 会报错
			// 所以仅仅适用于 username 肯定存在, 且 有多个值的场景
			if len(r.Form["username"][0]) == 0 {
				fmt.Fprintf(w, "username is invalid \n")
			}
			// 通过  r.Form.Get(xxx) 更好
			if len(r.Form.Get("password")) == 0 {
				fmt.Fprintf(w, "password is invalid \n")
			}

			// 数字
			age, err := strconv.Atoi(r.Form.Get("age"))
			if err != nil {
				fmt.Fprintf(w, "age is invalid \n")
			}
			if age < 0 || age > 150 {
				fmt.Fprintf(w, "age out of range \n")
			}
			//
			// 使用正则的方式
			if matched, _ := regexp.MatchString("^[0-9]+$", r.Form.Get("age1")); !matched {
				fmt.Fprintf(w, "age1 is invalid \n")
			}

			// 中文汉字
			//
			if m, _ := regexp.MatchString("^\\p{Han}+$", r.Form.Get("realname")); !m {
				fmt.Fprintf(w, "realname is not Han \n")
			}
			// 英文
			if m, _ := regexp.MatchString("^[a-zA-Z]+$", r.Form.Get("username")); !m {
				fmt.Fprintf(w, "Username is not in english \n")
			}

			// email
			if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,})\.([a-z]{2,4})$`, r.Form.Get("email")); !m {
				fmt.Println("no")
			} else {
				fmt.Println("yes")
			}
			// 手机号
			if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
				//
			}

			// 下拉菜单
			fruits := []string{"apple", "pear", "banana"}
			fr := r.Form.Get("fruit")
			for _, v := range fruits {
				if fr == v {
					// ...
				}
			}

			//单选框
			genders := []string{"1", "2"}
			for _, v := range genders {
				if v == r.Form.Get("gender") {
					// ...
				}
			}

			// 复选框
			intrs := []string{"football", "basketball", "tennis"}
			interests := r.Form["interests"]
			for _, v := range interests {
				for _, item := range intrs {
					if v == item {
						// ..
					}
				}
			}

			// 身份证
			//验证15位身份证，15位的是全部数字
			if m, _ := regexp.MatchString(`^(\d{15})$`, r.Form.Get("usercard")); !m {
			}

			//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
			if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, r.Form.Get("usercard")); !m {
			}

			// 防止跨站脚本攻击
			//
			// Go的html/template里面带有下面几个函数可以帮你转义
			//func HTMLEscape(w io.Writer, b []byte) //把b进行转义之后写到w
			//func HTMLEscapeString(s string) string //转义s之后返回结果字符串
			//func HTMLEscaper(args ...interface{}) string //支持多个参数一起转义，返回结果字符串
			//
			//t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
			//err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")
			//
			//t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
			// err = t.ExecuteTemplate(out, "T", template.HTML("<script>alert('you have been pwned')</script>"))
		}
	})

	http.HandleFunc("/upload", upload)

	err := http.ListenAndServe(":8080", nil) // 第二个参数是路由器, nil那么会使用默认路由器
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 处理/upload 逻辑
//
// application/x-www-form-urlencoded   表示在发送前编码所有字符（默认）
// multipart/form-data	  不对字符编码。在使用包含文件上传控件的表单时，必须使用该值。
// text/plain	  空格转换为 "+" 加号，但不对特殊字符编码。
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		//把上传的文件存储在内存和临时文件中
		//
		// 获取其他非文件字段信息的时候就不需要调用r.ParseForm，因为
		//在需要的时候Go自动会去调用。而且ParseMultipartForm调用一次之后，后面再次调用不会再有效果。
		//
		// 参数表示缓冲区大小, 上传文件存储在这块内存中, 如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中
		r.ParseMultipartForm(32 << 20)
		file, fileHeader, err := r.FormFile("uploadfile") //获取文件句柄
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", fileHeader.Header) // map[Content-Disposition:[form-data; name="uploadfile"; filename="老上海_窗台.png"] Content-Type:[image/png]]

		f, err := os.OpenFile("./"+fileHeader.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

////////////////////////////////////////////////////////////////

type myMux struct {
}

func (p *myMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		sayhelloName(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello xy!")
}

func customizedRouter() {
	mux := &myMux{}
	http.ListenAndServe(":9090", mux)
}
