package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"math"
)

type selpg struct {
	start int
	end int
	in_file string
	page_len int
	page_type bool
	out_file string
}

func get_args(args *selpg) {
	/*flag.IntVar(&bnFlag, "bn", 3, "份数") */
	flag.IntVar(&(args.start), "s", -1, "The start page number")
	flag.IntVar(&(args.end), "e", -1, "The end page number")
	flag.IntVar(&(args.page_len), "l", 72, "the length of page")
	flag.StringVar(&(args.out_file), "d", "", "the distination of the output")
	flag.BoolVar(&(args.page_type), "f", false, "read one page until it")
}

func check_for_args(args *selpg) bool {
	if args.start == -1 || args.end == -1 {
		err := "there must be both -s and -e"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	if args.start > args.end {
		err := "end must be larger than start"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	if args.start <= 0 {
		err := "start must be larger than 0"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	if args.end <= 0 {
		err := "end must be larger than 0"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	if args.start > (math.MaxInt32-1) {
		err := "start cannot be so large"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	if args.end > (math.MaxInt32-1) {
		err := "end cannot be so large"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	if args.page_len <= 0 {
		err := "page length must be larger than 0"
		e := errors.New(err)
		fmt.Fprintln(os.Stderr, "Error:", e)
		return false
	}
	return true
}

func input(args *selpg) {
	/* io.Reader
	io.WriterTo
	io.ByteScanner
	io.RuneScanner*/
	var in_path *bufio.Reader
	if args.in_file == "" {
		/*inputReader = bufio.NewReader(os.Stdin)*/
		in_path = bufio.NewReader(os.Stdin)
	} else {
		/*func OpenFile(name string, flag int, perm FileMode) (file *File, err error)*/
		input_file, err := os.OpenFile(args.in_file, os.O_RDWR, 0644)
		if err != nil {
			erro := "error occurs in inputfilapath"
			e := errors.New(erro);
			fmt.Fprintln(os.Stderr, "Error:", e)
		}
		in_path = bufio.NewReader(input_file)
	}
	
	/*Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}*/
	/*	Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
		Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
		Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")*/
	var out_path *os.File
	var err error
	if args.out_file != "" {
		/*func OpenFile(name string, flag int, perm FileMode) (file *File, err error)*/
		/*O_RDWR也就是说用读写的权限，O_CREATE然后文件存在忽略，不存在创建它，O_TRUNC文件存在截取长度为0*/
		out_path, err = os.OpenFile(args.out_file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", "-d must must a valid command that need input")
		}
	} else {
		/*struct的指针对象可以赋值为 nil 或与 nil 进行判等*/
		/*接口对象和接口对象的指针都可以赋值为 nil*/
		out_path = nil
	}
	selpg_IO(in_path, out_path, args)
}

func selpg_IO(in_path *bufio.Reader, out_path *os.File, args *selpg) {
	if args.page_type == false {
		type1(in_path, out_path, args)
	} else {
		type2(in_path, out_path, args)
	}
}

func type1(in_path *bufio.Reader, out_path *os.File, args *selpg) {
	/*type WriteCloser interface { Writer Closer }*/
	/*type Writer interface { Write(p []byte) (n int, err error)*/
	/*type Closer interface { Close() error }*/
	var std_input io.WriteCloser
	/*type Cmd　表示一个正在准备或者正在运行的外部命令*/
	var cmd *exec.Cmd

	var err error
	if args.out_file != "" {
		/*func Command(name string, arg ...string) *Cmd*/
		cmd = exec.Command(args.out_file)
		/*StdinPipe返回一个连接到command标准输入的管道pipe*/
		std_input, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", "error occurs in out_path")
		}
	}
	/*num:字符数量*/
	num := args.end - args.start + 1
	for i := 1; i < args.start; i++ {
		for j := 0; j < args.page_len; j++ {
			/* line, err := buff.ReadString('\n') //以'\n'为结束符读入一行*/
			in_path.ReadString('\n')
		}
	}
	for i := 0; i < num; i++ {
		for j := 0; j < args.page_len; j++ {
			/* line, err := buff.ReadString('\n') //以'\n'为结束符读入一行*/
			lh, err := in_path.ReadString('\n')
			if err != nil {
				/*stdin是标准输入，stdout是标准输出，stderr是标准错误输出。
				大多数的命令行程序从stdin输入，输出到stdout或stderr*/
				fmt.Fprintln(os.Stderr, "Error:", "error occurs in out_path")
				/*这里错误将导致无法执行，必须return*/
				return
			}
			if args.out_file != "" {
				_, err = std_input.Write([]byte(lh))
				/*line 的数量不够会影响结果，但是不会影响程序执行*/
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", "line is not enough")
				}					
			}
			if out_path != nil {
				/*n3,err := f.WriteString(String)*/
				out_path.WriteString(lh)
			} else {
				fmt.Print(lh)
			}
		}
	}
	/*cmd := exec.Command("tr", "a-z", "A-Z")  
    cmd.Stdin = strings.NewReader("some input")  
    var out bytes.Buffer  
    cmd.Stdout = &out  
    err := cmd.Run()  */
	if args.out_file != "" {
		std_input.Close()
		cmd.Stdout = os.Stdout
		err := cmd.Run()
	}
}

func type2(in_path *bufio.Reader, out_path *os.File, args *selpg) {
	/*type WriteCloser interface { Writer Closer }*/
	/*type Writer interface { Write(p []byte) (n int, err error)*/
	/*type Closer interface { Close() error }*/
	var std_input io.WriteCloser
	var std_input io.WriteCloser
	var cmd *exec.Cmd
	var err error
	if args.out_file != "" {
		/*func Command(name string, arg ...string) *Cmd*/
		cmd = exec.Command(args.out_file)
		/*StdinPipe返回一个连接到command标准输入的管道pipe*/
		std_input, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", "error occurs in out_path")
		}
	}
	/*num:字符数量*/
	num := args.end - args.start + 1
	for i := 1; i < args.start; i++ {
		in_path.ReadString('\f')
	}
	for i := 0; i < num; i++ {
		/* line, err := buff.ReadString('\f') //以'\f'为结束符读入一行*/
		lh, err := in_path.ReadString('\f')
		if err != nil {
			err := "page is not enough"
			e := errors.New(err)
			fmt.Fprintln(os.Stderr, "Error:", e)
			return
		}
		if args.out_file != "" {
			/*Write([]byte(String)) 它是阻塞的*/
			_, e := std_input.Write([]byte(lh))
			if e != nil {
				fmt.Fprint(os.Stderr, "Error:", e.Error())
			}
			continue
		}
		if out_path != nil {
			out_path.WriteString(lh)
		} else {
			fmt.Print(lh)
		}
	}
	/*cmd := exec.Command("tr", "a-z", "A-Z")  
    cmd.Stdin = strings.NewReader("some input")  
    var out bytes.Buffer  
    cmd.Stdout = &out  
    err := cmd.Run()  */
	if args.out_file != "" {
		std_input.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {
	args := new(selpg)
	get_args(args)
	if check_for_args(args) {
		input(args)
	}
}
