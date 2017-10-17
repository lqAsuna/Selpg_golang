# 程序说明：

## selpg设计思想
该程序主要完成了在GO语言下实现selpg的任务，主要编写了与 cat、ls、pr 和 mv 等标准命令类似的 Linux命令行实用程序。具体内容可以参考：
开发Linux命令实用程序

## 输入
通过初始化一个struct类型获得输入args

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/code_1.png)

然后实用GO语言自带的flag库函数获得输入值，注意，由于flag的特殊性，如果不输入某变量的值，则默认是初始值。

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/code_2.png)

通过函数input转化input_path和output_path，分别将其转化成bufio.Reader和os.File，目的是为了方便后续操作。
## 检查
我在函数check_for_args总共检查了未输入start或end、start过大或过小、end过大或过小、line过大或过小这七种情况。在实际操作的时候，还检查了路径错误、页数过少、行数过少等情况，具体代码不在这里详细给出。

## 操作
根据f的值不同，操作会分成两个大类，所以我创建了type1和type2两个函数，通过简单的两次两个for循环和库函数ReadString()即可截取需要输出的字符，第一次按照开始点和页截取，第二次按照结束点和页截取，在两次截取过程中都需要进行错误判定，因为可能出现也过得或者行过大的错误；第二种情况与第一种情况类似，不再赘述。

## 结果
1.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_1.png)

2.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_2.png)

3.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_3.png)

4.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_4.png)

5.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_5.png)

6.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_6.png)

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_7.png)

7.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_8.png)

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_7.png)

8.

![image](https://github.com/lqAsuna/Selpg_golang/blob/master/image/res_9.png)
