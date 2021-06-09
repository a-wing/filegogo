import React from 'react';
import logo from './logo.svg';
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
  constructor(props: any) {
    super(props)

    this.fgg = new LibFgg()
    this.qrcode = React.createRef()
    this.address = document.location.href
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

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <canvas ref={ this.qrcode }></canvas>
          <button onClick={ () => { this.handleCopy() } } >COPY</button>
          <p>{ this.address }</p>
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.tsx</code> and save to reload.
          </p>
          <input type="file" onChange={ (ev: any) => { this.handleFile(ev.target.files) } } />
          <a
            className="App-link"
            href="https://github.com/a-wing/filegogo/"
            target="_blank"
            rel="noopener noreferrer"
          >
            filegogo
          </a>
        </header>
      </div>
    )
  }
}

export default App;
