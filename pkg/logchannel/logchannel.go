package logchannel

import (
	"encoding/binary"
	"io"
	"os"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
)

type logChannel struct {
	channel ssh.Channel
	file    *os.File
}

func writeTTYRecHeader(fd io.Writer, length int) {
	t := time.Now()

	tv := syscall.NsecToTimeval(t.UnixNano())

	binary.Write(fd, binary.LittleEndian, int32(tv.Sec))
	binary.Write(fd, binary.LittleEndian, int32(tv.Usec))
	binary.Write(fd, binary.LittleEndian, int32(length))
}

func New(channel ssh.Channel) *logChannel {
	f, err := os.OpenFile("session.ttyrec", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		panic(err)
	}

	return &logChannel{
		channel: channel,
		file:    f,
	}
}

func (l *logChannel) Read(data []byte) (int, error) {
	return l.Read(data)
}

func (l *logChannel) Write(data []byte) (int, error) {
	writeTTYRecHeader(l.file, len(data))
	l.file.Write(data)

	return l.channel.Write(data)
}

func (l *logChannel) Close() error {
	l.file.Close()

	return l.channel.Close()
}
