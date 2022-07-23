import { useEffect, useRef, useState } from 'react'

import { ProtoHttpToWs } from '../lib/util'
import { getServer, getConfig, shareGetRoom } from '../lib/api'
import LibFgg from '../libfgg/libfgg'
import log  from 'loglevel'

import Address from './Address'
import File from './File'
import Qrcode from './QRCode'

import { DomSendFile, DomRecvFile } from '../libfgg/pool/file/dom'

const fgg = new LibFgg()
let enabled = true

function Index(props: { address: string }) {
  const address = props.address

  const [progress, setProgress] = useState<number>(0)
  const [total, setTotal] = useState<number>(10)
  const [recver, setRecver] = useState<boolean>(false)

  const refIce = useRef<RTCIceServer[]>([])

  fgg.onPreTran = (meta: any) => {
    setTotal(meta.size)
  }

  fgg.onRecvFile = () => {
    fgg.setRecv(new DomRecvFile())
    setRecver(true)
  }

  fgg.setOnProgress((c: number): void => {
    setProgress(c)
    log.debug(progress)
  })

  const getfile = async () => {

    await fgg.run()

    //fgg.useWebRTC({
    //  iceServers: refIce.current,
    //}, () => {

    //  // TODO:
    //  // Need Wait to 1s
    //  setTimeout(() => {
    //    fgg.getfile()
    //  }, 1000)
    //})
    //fgg.runWebRTC()
  }
  const handleFile = function(files: FileList) {
    //fgg.useWebRTC({
    //  iceServers: refIce.current,
    //}, () => {})
    fgg.setSend(new DomSendFile(files[0]))
  }

  useEffect(() => {
    const init = async function() {
      refIce.current = await getConfig()
    }

    init()
    return () => {
      //fgg.close()
    }
  }, [])

  useEffect(() => {
    if (enabled) {
      enabled = false

      const init = async() => {
        const addr = getServer() + shareGetRoom(address)
        await fgg.useWebsocket(ProtoHttpToWs(addr))
        fgg.clientMeta()
      }

      init()
    }
  }, [props.address])

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
