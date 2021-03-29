
// Safari
import 'blob-polyfill'

// Firefox, Safari
import 'web-streams-polyfill/dist/polyfill.min.js'

import SparkMD5 from 'spark-md5'

export default class Transfer {
  constructor(file) {
    // htmlDOMfile
    this.file = file
    this.metadata = {}

    // safari default
    this.step = 1024 * 64,
    // chrome, firefox max-message-size
    // step: 1024 * 256,

    this.dataChannel = {}
    this.onComplete = () => {}
    this.pointer = 0
    this.spark = new SparkMD5.ArrayBuffer()
  }
  sendBlob() {
    const p = this.pointer

    if (p >= this.file.size) {
      this.checksum = this.spark.end()
      this.onComplete(JSON.stringify({ checksum: this.checksum }))
    }

    this.file.slice(p, p + this.step).arrayBuffer().then(buffer => {
      // Md5
      this.spark.append(buffer)

      this.dataChannel.send(buffer)
    })
    this.pointer = p + this.step
  }
  onData(data) {
    // computed progress
    this.pointer = this.pointer + this.step

    // Md5
    this.spark.append(data)

    this.write(data)
  }
  write(buf) {
    console.log(buf)
    const readableStream = new Response(buf).body

    const reader = readableStream.getReader()
    const pump = () => reader.read()
      .then(res => res.done
        ? this.next()
        : this.file.write(res.value).then(pump))

    pump()
  }
  next() {
    this.signChannel.send('req')
  }
  verify(checksum) {
    this.onComplete()
    return this.spark.end() === this.checks
  }
}
