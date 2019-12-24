package astikit

import "bytes"

// WriterAdapter represents an object that can adapt a Writer
type WriterAdapter struct {
	buffer *bytes.Buffer
	o      WriterAdapterOptions
}

// WriterAdapterOptions represents WriterAdapter options
type WriterAdapterOptions struct {
	Callback func(i []byte)
	Split    []byte
}

// NewWriterAdapter creates a new WriterAdapter
func NewWriterAdapter(o WriterAdapterOptions) *WriterAdapter {
	return &WriterAdapter{
		buffer: &bytes.Buffer{},
		o:      o,
	}
}

// Close closes the adapter properly
func (w *WriterAdapter) Close() {
	if w.buffer.Len() > 0 {
		w.write(w.buffer.Bytes())
	}
}

// Write implements the io.Writer interface
func (w *WriterAdapter) Write(i []byte) (n int, err error) {
	// Update n to avoid broken pipe error
	defer func() {
		n = len(i)
	}()

	// Split
	if len(w.o.Split) > 0 {
		// Split bytes are not present, write in buffer
		if bytes.Index(i, w.o.Split) == -1 {
			w.buffer.Write(i)
			return
		}

		// Loop in split items
		items := bytes.Split(i, w.o.Split)
		for i := 0; i < len(items)-1; i++ {
			// If this is the first item, prepend the buffer
			if i == 0 {
				items[i] = append(w.buffer.Bytes(), items[i]...)
				w.buffer.Reset()
			}

			// Write
			w.write(items[i])
		}

		// Add remaining to buffer
		w.buffer.Write(items[len(items)-1])
		return
	}

	// By default, forward the bytes
	w.write(i)
	return
}

func (w *WriterAdapter) write(i []byte) {
	if w.o.Callback != nil {
		w.o.Callback(i)
	}
}
