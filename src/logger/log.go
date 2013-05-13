package logger

import (
	"fmt"
	"log/syslog"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

var debug bool
var replacer *strings.Replacer
var random *rand.Rand

func SetDebug(d bool) {
	debug = d
}

func NOTICE(format string, arg ...interface{}) {
	fmt.Printf("[NOTIC]%s %s", time.Now().Format("2006-01-02 15:04:05"), getCallerInfo())
	fmt.Printf(format, arg...)
	fmt.Println()
}

func DEBUG(format string, arg ...interface{}) {
	if !debug {
		return
	}

	fmt.Printf("[DEBUG]%s %s", time.Now().Format("2006-01-02 15:04:05"), getCallerInfo())
	fmt.Printf(format, arg...)
	fmt.Println()
}

func ERROR(format string, arg ...interface{}) {
	fmt.Printf("[ERROR]%s %s", time.Now().Format("2006-01-02 15:04:05"), getCallerInfo())
	fmt.Printf(format, arg...)
	fmt.Println()
}

type Func struct {
	prefix string
}

func FUNC(arg ...interface{}) *Func {
	ret := new(Func)
	if !debug {
		return ret
	}

	r := random.Int()
	ptr, f, l, ok := runtime.Caller(1)
	if ok {
		files := strings.Split(f, "/")
		f = files[len(files)-1]
		func_ := runtime.FuncForPC(ptr)
		ret.prefix = fmt.Sprintf("[FCALL]%s %s(%08x@%s:%d)", time.Now().Format("2006-01-02 15:04:05"), func_.Name(), r, f, l)
	} else {
		ret.prefix = fmt.Sprintf("[FCALL]%s unknown(%08x)", time.Now().Format("2006-01-02 15:04:05"), r)
	}
	fmt.Print(ret.prefix, " enter: ")
	for i := range arg {
		a := replacer.Replace(fmt.Sprintf("%s", arg[i]))
		fmt.Printf("%s, ", a)
	}
	fmt.Println()
	return ret
}

func (f Func) Quit() {
	if !debug {
		return
	}
	fmt.Println(f.prefix, "quit")
}

func INFO(prefix string, arg ...interface{}) {
	sys, err := syslog.New(syslog.LOG_INFO, prefix)
	if err != nil {
		ERROR("can't open syslog: %s", err)
		return
	}
	defer sys.Close()

	log := fmt.Sprintf("|%s|", time.Now().Format("2006-01-02 15:04:05"))
	for i := range arg {
		log += replacer.Replace(fmt.Sprintf("%v|", arg[i]))
	}
	sys.Info(log)

	NOTICE("%s:%s", prefix, log)
}

func init() {
	debug = false
	replacer = strings.NewReplacer("\n", "", "\r", "", "\t", "")
	random = rand.New(rand.NewSource(time.Now().Unix()))
}

func getCallerInfo() string {
	_, f, l, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	files := strings.Split(f, "/")
	return fmt.Sprintf("%s(%d): ", files[len(files)-1], l)
}
