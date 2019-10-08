# Selpg CLI

## flags

> CLI（Command Line Interface）实用程序是Linux下应用开发的基础。正确的编写命令行程序让应用与操作系统融为一体，通过shell或script使得应用获得最大的灵活性与开发效率。

而命令行中参数的风格有三大类，即Unix/Posix、BSD、GNU。分别有以下特征：

- Unix/Posix风格，即命令后的参数，可以分组，便必须以连字符开头，如`ps -aux`。

- BSD风格，即命令后的参数，可以分组，但不可以与连字符同用，如` ps ax`。

- GNU风格，即长选项，命令后的参数，可以分组，但必须以双横线开头，如：`ps --help`。

解析参数是CLI程序设计中的一大问题，为了实现 POSIX/GNU-风格参数处理，–flags，包括命令完成等支持，程序员们开发了无数第三方包，这些包可以在 [godoc](https://godoc.org/) 找到。善于利用这些第三方包，我们就可以不用重复造轮子，专注于所开发的CLI的运行逻辑了。

flags包就是其中之一，很好的实现了命令行参数的解析。

使用过程：

1. 使用flag.String(), Bool(), Int() 或 flag.StringVar(), BoolVar(), IntVar() 等函数定义flag
2. 在所有flag都注册之后，调用`flag.Parse()`进行解析
3. 解析后，参数对于的值就存入flag中，以供我们使用



## selpg

selpg 实用程序，名称代表 SELect PaGes。selpg 允许用户指定从输入文本抽取的页的范围，这些输入文本可以来自文件或另一个进程。selpg 是以在 Linux 中创建命令的事实上的约定为模型创建的，这些约定包括：

- 独立工作
- 在命令管道中作为组件工作（通过读取标准输入或文件名参数，以及写至标准输出和标准错误）
- 接受修改其行为的命令行选项

### Usage

所有选项都应以“-”（连字符）开头。选项可以附加参数。

selpg 要求用户用两个命令行参数“-sNumber”（例如，“-s10”表示从第 10  页开始）和“-eNumber”（例如，“-e20”表示在第 20 页结束）指定要抽取的页面范围的起始页和结束页。selpg  对所给的页号进行合理性检查；换句话说，它会检查两个数字是否为有效的正整数以及结束页是否不小于起始页。这两个选项，“-sNumber”和“-eNumber”是强制性的，而且必须是命令行上在命令名  selpg 之后的头两个参数：       

```shell
$ selpg -s10 -e20 ...
```



selpg 可以处理两种输入文本：       

*类型 1：*该类文本的页行数固定。这是缺省类型，因此不必给出选项进行说明。也就是说，如果既没有给出“-lNumber”也没有给出“-f”选项，则 selpg 会理解为页有固定的长度（每页 72 行）。       

选择 72 作为缺省值是因为在行打印机上这是很常见的页长度。这样做的意图是将最常见的命令用法作为缺省值，这样用户就不必输入多余的选项。该缺省值可以用“-lNumber”选项覆盖，如下所示：

```
$ selpg -s10 -e20 -l66 ...
```

这表明页有固定长度，每页为 66 行。

*类型 2：*该类型文本的页由 ASCII 换页字符（十进制数值为 12，在 C  中用“\f”表示）定界。该格式与“每页行数固定”格式相比的好处在于，当每页的行数有很大不同而且文件有很多页时，该格式可以节省磁盘空间。在含有文本的行后面，类型  2 的页只需要一个字符 ― 换页 ― 就可以表示该页的结束。打印机会识别换页符并自动根据在新的页开始新行所需的行数移动打印头。    

类型 2 格式由“-f”选项表示，如下所示：

```
$ selpg -s10 -e20 -f ...
```

该命令告诉 selpg 在输入中寻找换页符，并将其作为页定界符处理。

注：“-lNumber”和“-f”选项是互斥的。



### Examples

下面给出了可使用的 `selpg` 命令字符串示例：              

如何应用我们所介绍的一些原则，下面给出了可使用的 selpg 命令字符串示例：

1. `$ selpg -s1 -e1 input_file`

    该命令将把“input_file”的第 1 页写至标准输出（也就是屏幕），因为这里没有重定向或管道。
   
2. `$ selpg -s1 -e1 < input_file`

    该命令与示例 1 所做的工作相同，但在本例中，selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“input_file”而不是显式命名的文件名参数。输入的第 1 页被写至屏幕。

3. `$ other_command | selpg -s10 -e20`

    “other_command”的标准输出被 shell／内核重定向至 selpg 的标准输入。将第 10 页到第 20 页写至 selpg 的标准输出（屏幕）。

4. `$ selpg -s10 -e20 input_file >output_file`

    selpg 将第 10 页到第 20 页写至标准输出；标准输出被 shell／内核重定向至“output_file”。
    
5.  `$ selpg -s10 -e20 input_file 2>error_file`

    selpg 将第 10 页到第 20 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“error_file”。请注意：在“2”和“>”之间不能有空格；这是 shell 语法的一部分（请参阅“man bash”或“man sh”）。
    
6.  `$ selpg -s10 -e20 input_file >output_file 2>error_file`

    selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至“error_file”。当“input_file”很大时可使用这种调用；您不会想坐在那里等着 selpg 完成工作，并且您希望对输出和错误都进行保存。
  
7.  `$ selpg -s10 -e20 input_file >output_file 2>/dev/null`

    selpg 将第 10 页到第 20 页写至标准输出，标准输出被重定向至“output_file”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。
    
8.  `$ selpg -s10 -e20 input_file >/dev/null`

    selpg 将第 10 页到第 20 页写至标准输出，标准输出被丢弃；错误消息在屏幕出现。这可作为测试 selpg 的用途，此时您也许只想（对一些测试情况）检查错误消息，而不想看到正常输出。
   
9.  `$ selpg -s10 -e20 input_file | other_command`

    selpg 的标准输出透明地被 shell／内核重定向，成为“other_command”的标准输入，第 10 页到第 20 页被写至该标准输入。“other_command”的示例可以是 lp，它使输出在系统缺省打印机上打印。“other_command”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“other_command”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。
    
10. `$ selpg -s10 -e20 input_file 2>error_file | other_command`

    与上面的示例 9 相似，只有一点不同：错误消息被写至“error_file”。



## Design



## Test

`go run selpg.go --s 1 --e 2 <test.txt`


## References

[1] [开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)
[2] [标准库—命令行参数解析flag](http://blog.studygolang.com/2013/02/%E6%A0%87%E5%87%86%E5%BA%93-%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0%E8%A7%A3%E6%9E%90flag/)

