lisp

====

a simple lisp made by go

一个简单的go-lisp实现，基本遵循lisp标准语法，包含一些（可能）独有的函数或规则

支持输入四种基本类型：整数、浮点数、字符、字符串

内部支持三种基本类型：整数、浮点数、字符（字符作为整数保存）

支持四则运算、比较运算、逻辑运算（逻辑运算和cons是懒惰执行的）

以下是所有内置函数的简介：

	none 返回一个空值，如果出现错误，一定要返回空值

	eval 为将一个列表在当前作用域下执行

	quote 返回参数本身，主要用于保护一个列表

	atom 如果参数不是列表或者是空列表则返回真，否则返回假

	eq 判断两个元素是否相等（对内部注册函数和default或omission返回的函数会出错）

	car 返回列表的第一个元素

	cdr 返回除去第一个元素后的列表

	cons 返回将一个元素加入到列表的头部的列表

	each 为顺序执行多个语句，最后一个语句的返回值作为返回值
	
	block 类似 each 但会产生一个内层环境

	if 三个参数，根据第一个的结果决定执行第二个还是第三个（可以看成是cond的包装）

	cond 参数为一系列的二元列表，依次执行列表的第一个元素，直到返回为真时执行第二个元素并退出
	
	while 二个参数，循环执行第二个直到第一个判断为假
	
	until 二个参数，循环执行第二个直到第一个判断为真

	loop 三个参数，第一个初始化，循环执行第三个直到第二个判断为假
	
	for 三个参数，第一个必须是标签，第二个必须是生成器，循环执行第三个（生成器返回值赋给第一个）
	
		each、block 返回最后被执行的语句的值

		if cond 也一样（if如果第二选项省略，判断又为假，返回none）

		while、until、loop、for 如果被正确执行，永远返回none
		
	present 返回当前环境中的所有标签

	context 返回当前环境可以查找到的所有标签，亦即所有可用的标签

	builtin 内部注册函数无法被删除只能被内层环境覆盖，builtin直接查找内部注册函数

	remove 从当前环境中删除一个标签，如果不存在会向外查找，试图删除内部注册函数会导致错误

	clear 删除当前环境中的所有标签（不包括其父环境或更上层环境的标签）

	lambda 产生一个匿名函数

	macro 产生一个匿名宏，macro宏与define、update不同之处在于可以添加额外的替换对象

	define 声明一个标签为变量或者函数，标签建立在当前环境下
		
		(define f 1)			声明一个变量
		
		(define (f n) (* n 2))	声明一个函数
		
		如果要声明一个变量为宏，需使用macro才行

	update 用来更新一个标签的值，用法和lambda相同，但它只能更新已有的标签的值

		update 会顺序查找环境（注意闭包的生成时环境设定为嵌入在运行时环境和外部环境之间）

		update 如果查找不到，或者查找到的是内部注册函数，会报错
		
	pretreat 用于宏的参数预处理，第一个参数必须是宏，且后跟元素个数等于宏参数个数

	default 用于给函数或宏绑定默认值，返回一个绑定了默认值（因而可以省略后面参数）的函数

	omission 用于产生一个可变参数函数，提供的参数必须是一个函数，该函数的最后一个参数应该为列表

	scan 从控制台获取字符串数据，并将该数据作为lisp语句执行后返回

	load 读取一个lisp文件，并执行其内容后返回

	print 输出数据

	println 输出数据并回车

	raise 将字符串转化为错误并释放

	error 打印可能的错误，并向外传递错误

	catch 捕获错误并转化为字符串，否则返回一个空表
	
	chan 生成一个channel
	
	close 关闭一个channel
	
	go 在一个新的线程里执行语句
	
	channel可以像函数一样执行，如不提供参数为读取，可接受一个参数来写入

可以通过Add方法添加自定义函数：

	func (l *Lisp)Add(name string, func([]Token,*Lisp)(Token,error))

	事实上，func([]Token,*Lisp)(Token,error)被定义为lisp.Gfac类型，是注册函数的类型

	使用者需要自行检验参数的数量，对每个参数进行Exec运算获取实际参数，验证参数类型

Token为内部用来表示元素的类型

	type Token struct{
		Kind
		Text interface{}
	}

Text只可能装入如下类型：[]Token、int64、float64、string、Name、Hong、Lfac、Gfac

对应的Kind值分别为如下：List、Int、Float、String、Macro、Label、Front、Back

交互模式和lsp文件中都支持注释，只有一种注释形式：‘;’及该行剩余部分被忽略

注意的是为了实现惰性求值，你添加的函数接收到的切片，每个元素都是未运算的，需要你进行运算或解包

为了帮助lambda实现递归调用，内置标识符 self 代表本函数，命名函数也应尽量使用 self 来递归

命名函数使用函数名递归是可行的，但如果你将函数赋给一个值，并将原函数名另外赋值的话，会导致异常

self 可以避免这种异常，因此应尽量使用 self 标识符进行递归

具体使用参见example
