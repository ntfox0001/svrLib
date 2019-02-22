package util

import (
	"math"
)

// 浮点底数整数指数求幂
func Power(x float64, n int) float64 {
	ans := 1.0

	for n != 0 {
		ans *= x
		n--
	}
	return ans
}

// golang 标准库的求幂
func Powerf(x, n float64) float64 {
	return math.Pow(x, n)
}

/*
---------------------
作者：陈鹏万里
来源：CSDN
原文：https://blog.csdn.net/QQ245671051/article/details/70342047?utm_source=copy
版权声明：本文为博主原创文章，转载请附上博文链接！
注：次幂n为整数，底数可以是整数、小数、矩阵等(只要能进行乘法运算的
*/

func PowerI(x, n int64) int64 {
	var ret int64 = 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}
