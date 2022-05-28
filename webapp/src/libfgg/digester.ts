
// Inspired https://github.com/RobinLinus/snapdrop/blob/724f0af576852517bbea96a5d41302719ea514ec/client/scripts/network.js#L476-L505
class FileDigester {
  private buffer: Array<ArrayBuffer>
  private bytesReceived: number
  private size: number
  private mime: string
  private name: string
  private callback: (_: any) => void

  constructor(meta: {
    name: string,
    mime?: string,
    size: number,
  }, callback: ((_: any) => void)) {
    this.buffer = [];
    this.bytesReceived = 0
    this.size = meta.size
    this.mime = meta.mime || 'application/octet-stream'
    this.name = meta.name
    this.callback = callback
  }

  write(chunk: ArrayBuffer) {
    this.buffer.push(chunk)
    this.bytesReceived += chunk.byteLength
    if (this.bytesReceived < this.size) return
    // we are done
    let blob = new Blob(this.buffer, { type: this.mime })
    this.callback({
      name: this.name,
      mime: this.mime,
      size: this.size,
      blob: blob
    })
  }
  close() {}
}

export default FileDigester
