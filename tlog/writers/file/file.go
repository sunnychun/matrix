package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const bufferSize = 256 * 1024
const defaultMaxBytes = 1024 * 1024 * 1024
const flushInterval = time.Minute

type Options struct {
	Dir      string
	Program  string
	MaxBytes int
}

func Open(opts Options) (*File, error) {
	if opts.Program == "" {
		opts.Program = filepath.Base(os.Args[0])
	}
	if opts.MaxBytes <= 0 {
		opts.MaxBytes = defaultMaxBytes
	}

	var f File
	if err := f.init(opts); err != nil {
		return nil, err
	}

	return &f, nil
}

type File struct {
	opts   Options
	mu     sync.Mutex
	file   *os.File
	bufw   *bufio.Writer
	nbytes int
	closed bool
	done   chan struct{}
}

func (f *File) init(opts Options) error {
	f.opts = opts
	if err := f.rotate(time.Now()); err != nil {
		return err
	}
	f.done = make(chan struct{})
	go flushing(f)
	return nil
}

func (f *File) rotate(now time.Time) (err error) {
	if f.file != nil {
		f.bufw.Flush()
		f.file.Close()
	}

	f.nbytes = 0
	f.file, _, _, err = createLogFile(f.opts.Dir, f.opts.Program, now)
	if err != nil {
		return err
	}
	f.bufw = bufio.NewWriterSize(f.file, bufferSize)

	f.nbytes, err = fmt.Fprintf(f.file, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	return err
}

func (f *File) path() string {
	return filepath.Join(f.opts.Dir, linkName(f.opts.Program))
}

func (f *File) Close() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return &os.PathError{"close", f.path(), os.ErrClosed}
	}

	close(f.done)
	if f.file != nil {
		f.bufw.Flush()
		err = f.file.Close()
	}
	f.closed = true
	return err
}

func (f *File) Flush() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return &os.PathError{"flush", f.path(), os.ErrClosed}
	}
	return f.bufw.Flush()
}

func (f *File) Sync() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return &os.PathError{"flush", f.path(), os.ErrClosed}
	}
	f.bufw.Flush()
	return f.file.Sync()
}

func (f *File) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return 0, &os.PathError{"write", f.path(), os.ErrClosed}
	}

	if f.nbytes >= f.opts.MaxBytes {
		if err := f.rotate(time.Now()); err != nil {
			return 0, err
		}
	}
	n, err = f.bufw.Write(p)
	f.nbytes += n
	return n, err
}

func flushing(f *File) {
	t := time.NewTicker(flushInterval)
	defer t.Stop()
	for {
		select {
		case <-f.done:
			return
		case <-t.C:
			f.Flush()
		}
	}
}
