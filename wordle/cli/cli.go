// TODO: flags
// TODO: menu
package cli

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/johnietre/term-games/common"
	"github.com/johnietre/term-games/common/ansi"
	inputpkg "github.com/johnietre/term-games/common/input"
	utils "github.com/johnietre/utils/go"
	"github.com/spf13/cobra"
)

var (
	fiveWordsPath, sixWordsPath string
)

func init() {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		common.Fatal("error getting file path")
	}
	thisDir := filepath.Dir(thisFile)
	fiveWordsPath = filepath.Join(thisDir, "five-words.txt")
	sixWordsPath = filepath.Join(thisDir, "six-words.txt")
}

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "wordle",
		Short:                 "Play Worlde",
		Run:                   run,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

func Run() {
	if err := MakeCmd().Execute(); err != nil {
		common.Fatal("error running: ", err)
	}
}

func run(_ *cobra.Command, _ []string) {
	defer common.GlobalDeferrer.Run()

	restore, _, err := common.MakeTermRaw()
	if err != nil {
		common.Fatal("error making term raw: ", err)
	}
	common.GlobalDefer(restore)

	words, err := loadWords(fiveWordsPath, 5)

	common.Println("Press ESC or CTRL-C to exit\r")
	for {
		if err != nil {
			common.Fatal("error loading words: ", err)
		}
		runWithWords(words)
		common.Print("Exit = 0, Five = 5, Six = 6: ")
		for {
			buf, n, err := common.ReadStdinBytes()
			if err != nil {
				common.Fatal("error reading from stdin: ", err)
			} else if n == 0 {
				continue
			}
			switch input := inputpkg.FromBytes(buf[:]); input {
			case inputpkg.CtrlC, inputpkg.Escape:
				common.Exit(0)
			case inputpkg.Input('0'):
				common.Println("0\r")
				common.Exit(0)
			case inputpkg.Input('5'):
				common.Println("5\r")
				words, err = loadWords(fiveWordsPath, 5)
			case inputpkg.Input('6'):
				common.Println("6\r")
				words, err = loadWords(sixWordsPath, 6)
			default:
				continue
			}
			break
		}
	}
}

func runWithWords(words []string) {
	wordsLen := len(words)
	if wordsLen == 0 {
		return
	}
	wordLen := len(words[0])
	wordsSet := utils.SetFromSlice(words)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	correctWord := words[rng.Intn(wordsLen)]

	tb := common.NewTermBuffer(nil)
	for guesses := 1; guesses <= 6; guesses++ {
		letters := bytes.Repeat([]byte{'_'}, wordLen)
		tb.Write(bytes.Repeat([]byte{'_'}, wordLen))
		tb.Write([]byte{'\n', '\r'})
		tb.WriteString(ansi.SetDim + "ENTER" + ansi.ResetDim)
		tb.WriteString(ansi.CurUp + "\r")
		tb.Flush()

		curCol, wordOk := 0, false
	LettersLoop:
		for {
			buf, n, err := common.ReadStdinBytes()
			if err != nil {
				common.Fatal("error reading from stdin: ", err)
			} else if n == 0 {
				continue
			}

			/*
			   tb.WriteString(ansi.CurDownN(3)+"\r")
			   tb.Printf("%x %x %x %x", buf[0], buf[1], buf[2], buf[3])
			   tb.WriteString(ansi.CurUpN(3)+"\r")
			   tb.Flush()
			*/

			input := inputpkg.FromBytes(buf[:])
			newWordOk := wordOk
			// TODO: fix cursor positioning/backspace
			switch input {
			case inputpkg.CtrlC, inputpkg.Escape:
				tb.PrintFlush("\n\r\r\n")
				common.Exit(0)
			case inputpkg.Backspace, inputpkg.Del, inputpkg.Del2:
				letters[curCol] = '_'
				tb.WriteString("_" + ansi.CurLeft)
				if input != inputpkg.Del2 && curCol != 0 {
					tb.WriteString(ansi.CurLeft)
					curCol--
				}
				newWordOk = false
			case inputpkg.ArrowRight:
				if curCol == wordLen-1 {
					continue
				}
				curCol++
				tb.WriteString(ansi.CurRight)
			case inputpkg.ArrowLeft:
				if curCol == 0 {
					continue
				}
				curCol--
				tb.WriteString(ansi.CurLeft)
			case inputpkg.Enter:
				if !wordOk {
					continue
				}
				break LettersLoop
			default:
				if !input.IsAsciiAlpha() {
					continue
				}
				b := input.ToAsiiUpper().Byte()
				letters[curCol] = b
				tb.WriteByte(b)
				if curCol != wordLen-1 {
					curCol++
				} else {
					tb.WriteString(ansi.CurLeft)
				}
				newWordOk = wordsSet.Contains(string(letters))
			}
			if newWordOk != wordOk {
				tb.Write([]byte{'\n', '\r'})
				if newWordOk {
					tb.WriteString("ENTER")
				} else {
					tb.WriteString(ansi.SetDim + "ENTER" + ansi.ResetDim)
				}
				tb.WriteString(ansi.CurUp + "\r")
				if curCol != 0 {
					tb.WriteString(ansi.CurRightN(curCol))
				}
				wordOk = newWordOk
			}
			tb.Flush()
		}

		tb.WriteByte('\r')
		correctLetters, numCorrect := []byte(correctWord), 0
		// Check greens
		tb.WriteString(ansi.ForeGreen)
		for i, b := range letters {
			if correctLetters[i] == b {
				tb.WriteByte(b)
				correctLetters[i], letters[i] = 0, 0
				numCorrect++
			} else {
				tb.WriteString(ansi.CurRight)
			}
		}
		tb.WriteByte('\r')
		if numCorrect != wordLen {
			// Check rest
			yellow := false
			//tb.WriteString(ansi.ForeBlack)
			tb.WriteString(ansi.ForeWhite)
			for i, b := range letters {
				if b == 0 {
					tb.WriteString(ansi.CurRight)
					continue
				}
				if ci := bytes.IndexByte(correctLetters, b); ci != -1 {
					if !yellow {
						yellow = true
						tb.WriteString(ansi.ForeYellow)
					}
					tb.WriteByte(b)
					correctLetters[ci], letters[i] = 0, 0
				} else {
					if yellow {
						//tb.WriteString(ansi.ForeBlack)
						tb.WriteString(ansi.ForeWhite)
						yellow = false
					}
					tb.WriteByte(b)
				}
			}
		}

		tb.WriteString(ansi.ForeDefault)
		tb.WriteString("\n\r" + ansi.ClearLine)
		tb.Flush()

		if numCorrect == wordLen {
			tb.PrintFlush("CORRECT!!!\n\r")
			return
		}
	}
	tb.PrintFlush(ansi.ForeRed + correctWord + ansi.ForeDefault + "\n\r")
}

func loadWords(path string, wordLen int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var words []string
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return words, err
		}
		word := strings.ToUpper(strings.TrimSpace(line))
		if len(word) == wordLen {
			words = append(words, word)
		}
	}
}
