import React from 'react';
//import logo from './logo.svg';
import './App.css';

import QRCode from 'qrcode'

import { ProtoHttpToWs } from './lib/util'
import { getServer, getRoom } from './lib/api'
import LibFgg from './libfgg/libfgg'
import log from 'loglevel'
import history from 'history/browser'

import Address from './address'

class App extends React.Component {
  fgg: any
  qrcode: any
  address: string
  progress: number
  total:    number

  sender: boolean
  recver: boolean

  constructor(props: any) {
    super(props)

    this.fgg = new LibFgg()
    this.qrcode = React.createRef()
    this.address = document.location.href
    this.progress = 0
    this.total = 10

    this.sender = false
    this.recver = false
  }

  componentDidMount() {
    log.setLevel("debug")

    getRoom().then(room => {
      const addr = getServer() + room
      this.ShowQRcode(document.location.origin + '/' + room)
      this.historyPush(room)
      this.address = document.location.origin + '/' + room
      this.setState(() => {
        return "address"
      })
      this.wsconn(ProtoHttpToWs(addr))
    })
  }
  historyPush(path: string) {
    history.push(path)
  }
  ShowQRcode(addr: string) {
    QRCode.toCanvas(this.qrcode.current, addr, {
      width: 300
    }, error => {
      if (error) console.error(error)
      console.log('Create QRCode:', addr)
    })
  }
  wsconn(addr: string) {
    const fgg = this.fgg
    fgg.onPreTran = (meta: any) => {
      this.total = meta.size
      this.setState(()=>{
        return "total"
      })

    }

    fgg.onRecvFile = () => {
      this.recver = true
      this.setState(() => {return "recver"})
    }

    fgg.tran.onProgress = (c: number) => {
      this.progress += c
      this.setState(()=>{
        return "progress"
      })
    }

    fgg.useWebsocket(addr)
  }
  getfile() {
    this.fgg.useWebRTC({
      iceServers: [
        {
          urls: "stun:stun.l.google.com:19302",
        }
      ]
    }, () => {

      // TODO:
      // Need Wait to 1s
      setTimeout(() => {
        this.fgg.getfile()
      }, 1000)
    })
    this.fgg.runWebRTC()
  }
  handleFile(files: FileList) {
    this.sender = true

    this.fgg.useWebRTC({
      iceServers: [
        {
          urls: "stun:stun.l.google.com:19302",
        }
      ]
    }, () => {})

    this.fgg.sendFile(files[0])
  }
  // <img src={logo} className="App-logo" alt="logo" />
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <canvas className="qrcode" ref={ this.qrcode }></canvas>
          <Address address={ this.address }></Address>

          { this.recver
          ? <button className="App-address-button" onClick={ () => { this.getfile() } } >GetFile</button>
          : <input className="App-address-button" type="file" onChange={ (ev: any) => { this.handleFile(ev.target.files) } } />
          }
          <progress max={ this.total } value={ this.progress } ></progress>
        </header>
      </div>
    )
  }
}

export default App;
