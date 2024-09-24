package common

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	jtutils "github.com/johnietre/utils/go"
	"golang.org/x/term"
)

func ReadStdin(b []byte) (int, error) {
	return os.Stdin.Read(b)
}

func ReadStdinBytes() ([4]byte, int, error) {
	buf := [4]byte{}
	n, err := os.Stdin.Read(buf[:])
	return buf, n, err
}

type SyncDeferredFunc struct {
	df *jtutils.Mutex[*jtutils.DeferredFunc]
}

func NewSyncDeferredFunc(runPtr *bool) *SyncDeferredFunc {
	return &SyncDeferredFunc{
		df: jtutils.NewMutex(jtutils.NewDeferredFunc(runPtr)),
	}
}

func (sdf *SyncDeferredFunc) Add(funcs ...func()) {
	df := *sdf.df.Lock()
	df.Add(funcs...)
	sdf.df.Unlock()
}

func (sdf *SyncDeferredFunc) Run() bool {
	df := *sdf.df.Lock()
	ran := df.Run()
	sdf.df.Unlock()
	return ran
}

func (sdf *SyncDeferredFunc) Ran() bool {
	df := *sdf.df.Lock()
	ran := df.Ran()
	sdf.df.Unlock()
	return ran
}

var (
	GlobalDeferrer = NewSyncDeferredFunc(jtutils.NewT(true))
	exiting        atomic.Bool
)

func GlobalDefer(funcs ...func()) {
	GlobalDeferrer.Add(funcs...)
}

func Fatal(args ...any) {
	if exiting.CompareAndSwap(false, true) {
		GlobalDeferrer.Run()
	}
	log.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	if exiting.CompareAndSwap(false, true) {
		GlobalDeferrer.Run()
	}
	log.Fatalf(format, args...)
}

func Fatalln(args ...any) {
	if exiting.CompareAndSwap(false, true) {
		GlobalDeferrer.Run()
	}
	log.Fatalln(args...)
}

func Exit(code int) {
	if exiting.CompareAndSwap(false, true) {
		GlobalDeferrer.Run()
	}
	os.Exit(code)
}

var (
	stdinFd = os.Stdin.Fd()
)

func StdinFd() int {
	return int(stdinFd)
}

func MakeTermRaw() (restore func(), state *term.State, err error) {
	state, err = term.MakeRaw(StdinFd())
	if err == nil {
		restore = func() {
			term.Restore(StdinFd(), state)
		}
	}
	return
}

func GetTermSize() (width, height int, err error) {
	return term.GetSize(StdinFd())
}

func Print(args ...any) {
	fmt.Print(args...)
}

func Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func Println(args ...any) {
	fmt.Println(args...)
}

func PrintFlush(args ...any) {
	fmt.Print(args...)
	os.Stdout.Sync()
}

func PrintfFlush(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Stdout.Sync()
}

func PrintlnFlush(args ...any) {
	fmt.Println(args...)
	os.Stdout.Sync()
}

type TermBuffer struct {
	buf *bytes.Buffer
}

func NewTermBuffer(buf []byte) *TermBuffer {
	return &TermBuffer{
		buf: bytes.NewBuffer(buf),
	}
}

func (tb *TermBuffer) Write(b []byte) (int, error) {
	return tb.buf.Write(b)
}

func (tb *TermBuffer) WriteString(s string) (int, error) {
	return tb.buf.WriteString(s)
}

func (tb *TermBuffer) WriteByte(b byte) error {
	return tb.buf.WriteByte(b)
}

func (tb *TermBuffer) WriteRune(r rune) (int, error) {
	return tb.buf.WriteRune(r)
}

func (tb *TermBuffer) Print(args ...any) {
	fmt.Fprint(tb.buf, args...)
}

func (tb *TermBuffer) Printf(format string, args ...any) {
	fmt.Fprintf(tb.buf, format, args...)
}

func (tb *TermBuffer) Println(args ...any) {
	fmt.Fprintln(tb.buf, args...)
}

func (tb *TermBuffer) PrintFlush(args ...any) (int, error) {
	tb.Print(args...)
	return tb.Flush()
}

func (tb *TermBuffer) PrintfFlush(format string, args ...any) (int, error) {
	tb.Printf(format, args...)
	return tb.Flush()
}

func (tb *TermBuffer) PrintlnFlush(args ...any) (int, error) {
	tb.Println(args...)
	return tb.Flush()
}

func (tb *TermBuffer) Flush() (int, error) {
	n, err := tb.buf.WriteTo(os.Stdout)
	// TODO: is this necessary?
	if err == nil {
		os.Stdout.Sync()
	}
	return int(n), err
}

func (tb *TermBuffer) Buffer() *bytes.Buffer {
	return tb.buf
}
