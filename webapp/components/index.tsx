import { useEffect, useRef, useState } from 'react'

import { ProtoHttpToWs } from '../lib/util'
import { getServer, getConfig, shareGetRoom, getBoxInfo } from '../lib/api'
import LibFgg from '../libfgg/libfgg'
import log  from 'loglevel'
import { Meta } from '../libfgg/pool/data'

import Address from './Address'
import File from './File'
import Qrcode from './QRCode'
import Card from './card'

import { DomSendFile, DomRecvFile } from '../libfgg/pool/file/dom'

type IMeta = Meta & {
  remain?: number
  expire?: string
}

const fgg = new LibFgg()
let enabled = true

function Index(props: { address: string }) {
  const address = props.address

  const [meta, setMeta] = useState<IMeta | null>(null)
  const [progress, setProgress] = useState<number>(0)
  const [recver, setRecver] = useState<boolean>(false)
  const [isBox, setIsBox] = useState<boolean>(false)

  const refIce = useRef<RTCIceServer[]>([])

  fgg.onSendFile = (meta: IMeta) => {
    setMeta(meta)
  }

  fgg.onRecvFile = (meta: IMeta) => {
    setMeta(meta)
    fgg.setRecv(new DomRecvFile())
    setRecver(true)
  }

  fgg.setOnProgress((c: number): void => {
    setProgress(c)
    log.debug(progress)
  })

  const getfile = async () => {
    await fgg.run()
  }
  const handleFile = function(files: FileList) {
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

        fgg.useWebRTC({
          iceServers: refIce.current,
        })

        fgg.clientMeta()
      }

      init()
    }
  }, [props.address])

  const load = async () => {
    let data = await getBoxInfo()
    if (data) {
      setMeta(data)
      setIsBox(true)
      setRecver(true)
    } else {
      setIsBox(false)
      setRecver(false)
      setMeta(null)
    }
  }
  useEffect(() => {
    load()
  }, [props.address])

  return (
    <>
      <div style={{ width: '100%' }}>
      { meta
        ? <Card
            name={ meta.name }
            type={ meta.type }
            size={ meta.size }
            remain={ meta.remain }
            expire={ meta.expire }
          ></Card>
        : <>
            <div style={{ display: 'flex', justifyContent: 'center' }}>
              <Qrcode address={ address }></Qrcode>
            </div>
            <Address address={ address }></Address>
          </>
      }
        <File
          recver={ recver }
          isBox={ isBox }
          reLoad={ load }
          percent={ progress / (meta ? meta.size : 0.01) * 100 }
          handleFile={ (files: any) => { handleFile(files) } }
          getFile={ getfile }
        ></File>
      </div>
    </>
  )
}

export default Index
