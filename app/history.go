package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type History struct {
	commands       []string
	lock           sync.RWMutex
	lastWriteIndex int
}

var (
	history *History
	once    sync.Once
)

func GetHistory() *History {

	once.Do(func() {
		history = &History{commands: make([]string, 0)}
		if path, ok := os.LookupEnv("HISTFILE"); ok {
			history.LoadHistory(path)
			history.lastWriteIndex = len(history.commands)
		}
	})

	return history
}

func (h *History) Add(command string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.commands = append(h.commands, command)
}

func (h *History) GetLast(count int) string {

	h.lock.RLock()
	defer h.lock.RUnlock()

	start := max(0, len(h.commands)-count)

	var commands strings.Builder
	for i := start; i < len(h.commands); i++ {
		commands.WriteString(fmt.Sprintf("    %d  %s\n", (i + 1), h.commands[i]))
	}
	return commands.String()

}

func (h *History) GetHistoryIndex() int {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return len(h.commands)

}

func (h *History) Prev(historyIndex *int) string {

	h.lock.Lock()
	defer h.lock.Unlock()

	if len(h.commands) == 0 {
		return ""
	}

	*historyIndex--

	if *historyIndex < 0 {
		*historyIndex = 0
	} else if *historyIndex >= len(h.commands) {
		*historyIndex = len(h.commands) - 1
	}

	lastCommand := h.commands[*historyIndex]
	return lastCommand

}

func (h *History) Next(historyIndex *int) string {
	h.lock.Lock()
	defer h.lock.Unlock()

	if len(h.commands) == 0 {
		return ""
	}

	*historyIndex++

	if *historyIndex < 0 {
		*historyIndex = 0
	} else if *historyIndex >= len(h.commands) {
		*historyIndex = len(h.commands) - 1
	}

	nextCommand := h.commands[*historyIndex]
	return nextCommand

}

func (h *History) LoadHistory(path string) {
	file, err := Open(path, os.O_RDONLY, false)

	if err != nil {
		fmt.Println(err)
		return
	}

	if file == nil {
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var lineBuf bytes.Buffer

	for {
		part, isPrefix, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		lineBuf.Write(part)

		if !isPrefix {
			h.commands = append(h.commands, lineBuf.String())
			lineBuf.Reset()
		}
	}

	if lineBuf.Len() > 0 {
		h.commands = append(h.commands, lineBuf.String())
		lineBuf.Reset()
	}

}

func (h *History) WriteHistory(path string) {

	file, err := Open(path, os.O_WRONLY, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	h.WriteFrom(file)
}

func (h *History) AppendHistory(path string) {
	file, err := Open(path, os.O_WRONLY, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	h.WriteFrom(file)

}

func (h *History) WriteFrom(file *os.File) {
	writer := bufio.NewWriter(file)

	defer func() {
		writer.Flush()
	}()

	for i := h.lastWriteIndex; i < len(h.commands); i++ {
		if _, err := writer.WriteString(h.commands[i]); err != nil {
			h.lastWriteIndex++
			fmt.Println(err)
			return
		}
		if err := writer.WriteByte('\n'); err != nil {
			h.lastWriteIndex++
			fmt.Println(err)
			return
		}
	}
}
