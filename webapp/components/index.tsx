import { useEffect, useRef, useState } from 'react'

import { ProtoHttpToWs } from '../lib/util'
import { getServer, getConfig, shareGetRoom } from '../lib/api'
import LibFgg from '../libfgg/libfgg'
import log  from 'loglevel'

import Address from './Address'
import File from './File'
import Qrcode from './QRCode'

const fgg = new LibFgg()

function Index(props: { address: string }) {
  const address = props.address

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

    const addr = getServer() + shareGetRoom(address)
    fgg.useWebsocket(ProtoHttpToWs(addr))
  }

  useEffect(() => {
    init()

    return () => {
      fgg.close()
    }
  }, [])

  return (
    <>
      <Qrcode address={ address }></Qrcode>
      <Address address={ address }></Address>
      <File
        recver={ recver }
        percent={ progress / total * 100 }
        handleFile={ (files: any) => { handleFile(files) } }
        getFile={ getfile }
      ></File>
    </>
  )
}

export default Index
