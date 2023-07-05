import FileItem from "./file-item"
import { Box } from "../libfgg/index"
import { ExpiresAtHumanTime } from "../lib/util"
import Qrcode from "./qr-code"
import { getIceServers, getRaw, getServer } from "../lib/api"
import LibFgg from "../libfgg/libfgg"
import { ProtoHttpToWs } from "../lib/util"
import { DomRecvFile } from '../libfgg/pool/file/dom'

export default (props: { file: Box }) => {
  const file = props.file

  const toggleDownload = async (uxid: string) => {
    if (file.action === "relay") {
      await getRaw(uxid)
    } else {
      const fgg = new LibFgg()
      const addr = getServer() + uxid
      await fgg.useWebsocket(ProtoHttpToWs(addr))
      await fgg.useWebRTC({
        iceServers: await getIceServers(),
      })

      let p = new Promise<void>(resolve => {
        fgg.onRecvFile = (_) => {
          fgg.setRecv(new DomRecvFile())
          resolve()
        }
      })
      await fgg.clientMeta()
      await p
      await fgg.run()
    }
  }

  return (
    <>
      <ul className="p-3 w-full flex flex-row justify-center">
        <li className="m-2 p-4 border-1 border-green-300 rounded-md bg-green-100 shadow-md w-full">
          <div className="flex flex-row justify-between">
            <FileItem file={file}></FileItem>
          </div>
          { file.files?.length > 1
            ? <details className="cursor-pointer">
                <summary>{ file.files.length.toString() + " files" }</summary>
            { file.files.map((file, index) =>
              <div key={index} className="mx-8">
                <FileItem file={file}></FileItem>
              </div>)
            } </details>
            : null
          }
          <div>
            Expires after { file.remain } download or { ExpiresAtHumanTime(file.expire) }
          </div>
        </li>
      </ul>

      <div className="w-full flex flex-row justify-center">
        <details className="cursor-pointer">
          <summary>QRCode</summary>
          <Qrcode address={ window.location.href }></Qrcode>
        </details>
      </div>

      <button className="p-3 w-full block border-1 rounded-md bg-blue-500 text-white font-bold" onClick={ () => toggleDownload(file.uxid) }>Download</button>
    </>
  )
}
