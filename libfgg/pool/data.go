package pool

type DataChunk struct {
	Offset int64 `json:"offset"`
	Length int64 `json:"length"`
}

func (p *Pool) SendData(c *DataChunk) ([]byte, error) {
	data := make([]byte, c.Length)
	n, err := p.sender.ReadAt(data, int64(c.Offset))

	p.fileHash.onData(c, data[:n])
	p.OnProgress(p.fileHash.offset)
	return data[:n], err
}

func (p *Pool) RecvData(c *DataChunk, data []byte) error {
	p.mu.Lock()
	p.currentSize += c.Length
	p.mu.Unlock()
	p.fileHash.onData(c, data)
	p.OnProgress(p.fileHash.offset)
	_, err := p.recver.WriteAt(data, c.Offset)
	return err
}

func (p *Pool) Next() *DataChunk {
	p.mu.Lock()
	currentSize := p.currentSize
	p.mu.Unlock()

	if currentSize >= p.meta.Size {
		p.OnFinish()
		return nil
	}

	if p.pendingSize >= p.meta.Size {
		return nil
	}

	length := p.chunkSize
	next := currentSize + p.chunkSize
	if next > p.meta.Size {
		n := next - p.meta.Size
		length = p.chunkSize - n
	}

	offset := p.pendingSize

	p.pendingSize += p.chunkSize
	return &DataChunk{
		Offset: offset,
		Length: length,
	}
}
