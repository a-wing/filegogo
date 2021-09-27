import { useEffect, useRef, useState } from 'react'
//import logo from './logo.svg';
import './App.css';

import { ProtoHttpToWs } from './lib/util'
import { getServer, getConfig, getRoom } from './lib/api'
import LibFgg from './libfgg/libfgg'
import log from 'loglevel'
import history from 'history/browser'

import Address from './components/Address'
import File from './components/File'
import Qrcode from './components/QRCode'

const fgg = new LibFgg()

function App() {
  const [address, setAddress] = useState<string>(document.location.href)

  const [progress, setProgress] = useState<number>(0)
  const [total, setTotal] = useState<number>(10)
  const [recver, setRecver] = useState<boolean>(false)

  const refIce = useRef<RTCIceServer[]>([])

  fgg.onPreTran = (meta: any) => {
    setTotal(meta.size)
  }

  fgg.onRecvFile = () => setRecver(true)

  fgg.tran.onProgress = (c: number) => {
    setProgress(progress + c)
    log.debug(progress)
  }

  const getfile = function() {
    fgg.useWebRTC({
      iceServers: refIce.current,
    }, () => {

      // TODO:
      // Need Wait to 1s
      setTimeout(() => {
        fgg.getfile()
      }, 1000)
    })
    fgg.runWebRTC()
  }
  const handleFile = function(files: FileList) {
    fgg.useWebRTC({
      iceServers: refIce.current,
    }, () => {})

    fgg.sendFile(files[0])
  }

  const init = async function() {
    refIce.current = await getConfig()

    const room = await getRoom()
    console.log(room)
    history.push(room)
    setAddress(document.location.origin + '/' + room)
    const addr = getServer() + room
    fgg.useWebsocket(ProtoHttpToWs(addr))
  }

  useEffect(() => {
    // log.setLevel('debug')
    init()
  }, [])

  // <img src={logo} className="App-logo" alt="logo" />
  return (
      <div className="App">
        <header className="App-header">
        <div className="App-card">
          <Qrcode address={ address }></Qrcode>
          <Address address={ address }></Address>
          <File
            recver={ recver }
            percent={ progress / total * 100 }
            handleFile={ (files: any) => { handleFile(files) } }
            getFile={ getfile }
          ></File>
        </div>
        </header>
      </div>
    )
}

export default App;
