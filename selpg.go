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

func main() {
	args := new(selpg)
	get_args(args)
	if check_for_args(args) {
		input(args)
	}
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
	var in_path *bufio.Reader
	if args.in_file == "" {
		in_path = bufio.NewReader(os.Stdin)
	} else {
		input_file, err := os.OpenFile(args.in_file, os.O_RDWR, 0644)
		if err != nil {
			erro := "error occurs in inputfilapath"
			e := errors.New(erro);
			fmt.Fprintln(os.Stderr, "Error:", e)
		}
		in_path = bufio.NewReader(input_file)
	}
	
	var out_path *os.File
	var err error
	if args.out_file != "" {
		out_path, err = os.OpenFile(args.out_file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error=>", "-d must must a valid command that need input")
		}
	} else {
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
	var std_input io.WriteCloser
	var cmd *exec.Cmd
	var err error
	if args.out_file != "" {
		cmd = exec.Command(args.out_file)
		std_input, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", "error occurs in out_path")
		}
	}
	num := args.end - args.start + 1
	for i := 1; i < args.start; i++ {
		for j := 0; j < args.page_len; j++ {
			in_path.ReadString('\n')
		}
	}
	for i := 0; i < num; i++ {
		for j := 0; j < args.page_len; j++ {
			line, err := in_path.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", "error occurs in out_path")
				return
			}
			if args.out_file != "" {
				_, err = std_input.Write([]byte(line))
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", "line is not enough")
				}					
			}
			if out_path != nil {
				out_path.WriteString(line)
			} else {
				fmt.Print(line)
			}
		}
	}
	if args.out_file != "" {
		std_input.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func type2(in_path *bufio.Reader, out_path *os.File, args *selpg) {
	var std_input io.WriteCloser
	var cmd *exec.Cmd
	var err error
	if args.out_file != "" {
		cmd = exec.Command(args.out_file)
		std_input, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", "error occurs in out_path")
		}
	}
	num := args.end - args.start + 1
	for i := 1; i < args.start; i++ {
		in_path.ReadString('\f')
	}
	for i := 0; i < num; i++ {
		line, err := in_path.ReadString('\f')
		if err != nil {
			err := "page is not enough"
			e := errors.New(err)
			fmt.Fprintln(os.Stderr, "Error:", e)
			return
		}
		if args.out_file != "" {
			_, e := std_input.Write([]byte(line))
			if e != nil {
				fmt.Fprint(os.Stderr, "Error:", e.Error())
			}
			continue
		}
		if out_path != nil {
			out_path.WriteString(line)
		} else {
			fmt.Print(line)
		}
	}
	if args.out_file != "" {
		std_input.Close()
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
