import React from 'react';
//import logo from './logo.svg';
import './App.css';

import QRCode from 'qrcode'

import LibFgg from './libfgg/libfgg'
import log from 'loglevel'
import history from 'history/browser'

import copy from 'copy-to-clipboard'

class App extends React.Component {
  fgg: any
  qrcode: any
  address: string
  progress: number
  total:    number
  constructor(props: any) {
    super(props)

    this.fgg = new LibFgg()
    this.qrcode = React.createRef()
    this.address = document.location.href
    this.progress = 0
    this.total = 10
  }

  componentDidMount() {
    log.setLevel("debug")
    //const ws = new WebSocket('ws://localhost:8033/share/4553')
    const id = document.location.pathname.split("/")[2]
    //const fgg = new LibFgg()
    const fgg = this.fgg
    fgg.onShare = ((addr: any) => {
      const url = new URL(addr)
      const path = '/s/' + url.pathname.split('/')[2]
      history.push(path)

      const address = document.location.origin + path
      //this.address = "xxx"
      this.address = address
      //this.setState((state, props) => {
      this.setState(() => {
        return "address"
      })
      //this.setState(()=>{})
      QRCode.toCanvas(this.qrcode.current, address, {
        width: 400
      }, error => {
        if (error) console.error(error)
        console.log('Create QRCode:', address)
      })

      fgg.onPreTran = (meta: any) => {
        this.total += meta.size

        this.setState(()=>{
          return "total"
        })

        //if (!fgg.tran.file) {
        fgg.useWebRTC({
          iceServers: [
            {
              urls: "stun:stun.l.google.com:19302",
            }
          ]
        })
        //}
      }

      fgg.tran.onProgress = (c: number) => {
        this.progress += c
        this.setState(()=>{
          return "progress"
        })
      }

    })
    fgg.useWebsocket('ws://localhost:8033/share/' + id)

    //this.fgg = fgg

    //const ws = new WebSocket()
    //ws.onmessage = ({data}) => {
    //  log.warn(data)
    //}
    //log.debug(ws.server)
  }
  handleFile(files: FileList) {
    this.fgg.sendFile(files[0])
  }

  handleCopy() {
      copy(this.address)
  }

  // <img src={logo} className="App-logo" alt="logo" />
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <canvas className="qrcode" ref={ this.qrcode }></canvas>
          <div className="App-address">
            <p className="App-address-text" >{ this.address }</p>
            <button className="App-address-button" onClick={ () => { this.handleCopy() } } >COPY</button>
          </div>
          <progress max={ this.total } value={ this.progress } ></progress>
          <input className="App-address-button" type="file" onChange={ (ev: any) => { this.handleFile(ev.target.files) } } />
        </header>
      </div>
    )
  }
}

export default App;
