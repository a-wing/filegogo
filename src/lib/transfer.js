
// Safari
import 'blob-polyfill'

// Firefox, Safari
import 'web-streams-polyfill/dist/polyfill.min.js'

import SparkMD5 from 'spark-md5'

export default class Transfer {
  constructor(file) {
    // htmlDOMfile
    this.file = file
    this.fileStream = {}

    // safari default
    this.step = 1024 * 64,
    // chrome, firefox max-message-size
    // step: 1024 * 256,

    this.dataChannel = {}
    this.onComplete = () => {}
    this.onProgress = () => {}
    this.pointer = 0
    this.spark = new SparkMD5.ArrayBuffer()
  }
  sendBlob() {
    this.file.slice(this.pointer, this.pointer + this.step).arrayBuffer().then(buffer => {
      // Md5
      this.spark.append(buffer)
      this.progress(buffer.byteLength)

      this.dataChannel.send(buffer)
    })
  }
  onData(buffer) {
    // Md5
    this.spark.append(buffer)
    this.progress(buffer.byteLength)

    this.fileStream.write(new Uint8Array(buffer)).then(this.next())
  }
  next() {
    this.signChannel.send('req')
  }
  verify(checksum) {
    this.onComplete()
    return this.spark.end() === this.checks
  }
  progress(step) {
    if (this.pointer >= this.file.size) {
      this.onComplete(this.spark.end())
      console.log("onComplete")
    }
    // computed progress
    this.pointer = this.pointer + step

    this.onProgress((this.pointer / this.file.size) * 100)
  }
}
