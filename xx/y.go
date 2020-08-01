package xx

import (
	"fmt"
	"math/rand"
)

// Hello is func
// 导出的函数必须有注释
// 变量or方法都必须大写才能在包外访问
func Hello() {
	fmt.Println("hello", "world", rand.Intn(10))
}
