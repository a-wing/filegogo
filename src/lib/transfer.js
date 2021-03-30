
// Safari
import 'blob-polyfill'

// Firefox, Safari
import 'web-streams-polyfill/dist/polyfill.min.js'

import streamSaver from 'streamsaver'

import SparkMD5 from 'spark-md5'

class Transfer {
  constructor(file, channel) {
    // htmlDOMfile
    this.file = file
    this.fileStream = {}

    // safari default
    this.step = 1024 * 64,
    // chrome, firefox max-message-size
    // step: 1024 * 256,

    //channel.onmessage = ev => ev.target.label === 'dataChannel' ? this.onData(ev.data) : null
    channel.onmessage = ev => this.onData(ev.data)
    this.channel = channel

    this.onComplete = () => {}
    this.onProgress = () => {}
    this.pointer = 0
    this.spark = new SparkMD5.ArrayBuffer()
  }
  start() {
    this.channel.send(JSON.stringify({ event: 'req' }))
  }
  sendBlob() {
    this.file.slice(this.pointer, this.pointer + this.step).arrayBuffer().then(buffer => {
      // Md5
      this.spark.append(buffer)
      this.progress(buffer.byteLength)

      this.channel.send(buffer)
    })
  }
  onData(buffer) {
    if (this.isComplete()) {
      if (this.verify(JSON.parse(buffer)["checksum"])) {
        console.log("checksum success")
      }

      this.fileStream.close()
    } else {

    // Md5
    this.spark.append(buffer)
    this.progress(buffer.byteLength)

    this.fileStream.write(new Uint8Array(buffer)).then(this.next())
    }
  }
  next() {
    this.channel.send(JSON.stringify({ event: 'req' }))
  }
  verify(checksum) {
    return this.spark.end() === this.checks
  }
  progress(step) {
    // computed progress
    this.pointer = this.pointer + step

    this.onProgress((this.pointer / this.file.size) * 100)
    if (this.isComplete()) {
      this.onComplete(this.spark.end())
      console.log("onComplete")
    }
  }
  isComplete() {
    return this.pointer >= this.file.size
  }
}

export class Sender extends Transfer {
  constructor(file, channel) {
    super(file, channel)
  }
  onData(data) {
    if (JSON.parse(data)["event"] == "req") {
      this.sendBlob()
    } else {
      this.channel.send(JSON.stringify({ checksum: this.checksum }))
    }
  }
}

export class Recver extends Transfer {
  constructor(file, channel) {
    super(file, channel)

    this.fileStream = streamSaver.createWriteStream(file.name, {
      size: file.size,
      //mitm: this.file.type
    }).getWriter()

  }
}

