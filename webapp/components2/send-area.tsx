import { useRef, useState, ChangeEvent, useEffect } from "react"
import { useAtom } from "jotai"
import { filesize } from "filesize"

import Archive from "../lib/archive"
import LibFgg, { Meta } from "../libfgg"
import { getIceServers, getServer, putBox } from "../lib/api"
import { loadHistory } from "../lib/history"
import FileItem from "./file-item"
import { ItemsAtom, DetailAtom } from "../store"
import SendFile from "./send-file"

import { ProtoHttpToWs } from "../lib/util"
import { DomSendFile } from "../libfgg/pool/file/dom"
import Loading from "./loading"

let archive = new Archive()

export default () => {
  const hiddenFileInput = useRef<HTMLInputElement>(null)
  const [progress, setProgress] = useState<number>(0)
  const [loading, setLoading] = useState<boolean>(false)
  const [detail, setDetail] = useAtom(DetailAtom)
  const [remain, setRemain] = useState<number>(1)
  const [expire, setExpire] = useState<string>("5m")
  const [relay, setRelay] = useState<boolean>(true)
  const [total, setTotal] = useState<number>(0)
  const [files, setFiles] = useState<Array<Meta>>([])
  const [items, setItems] = useAtom(ItemsAtom)
  const toggleButton = () => {
    hiddenFileInput.current?.click?.()
  }

  const syncLoad = async () => {
    setItems(await loadHistory())
  }

  useEffect(() => {
    syncLoad()
  }, [])

  const handleFile = (filelist: FileList) => {
    let files = new Array<File>(filelist.length)
    for (let i = 0; i < filelist.length; i++) {
      files[i] = filelist[i]
    }
    archive.addFiles(files)
    setTotal(archive.size)
    setFiles([...archive.manifest])
  }

  const toggleCommit = async (store: boolean) => {
    const count = archive.files.length
    if (count === 0) {
      return
    }
    setLoading(true)
    let file = await archive.exportFile()

    let action = store ? "relay" : "p2p"
    let manifest = await putBox(file, remain, expire, action)
    setItems([manifest, ...items])

    if (store) {
      localStorage.setItem(manifest.uxid, JSON.stringify(manifest))
    } else {
      const fgg = new LibFgg()
      await fgg.useWebsocket(ProtoHttpToWs(getServer() + manifest.uxid))
      await fgg.useWebRTC({
        iceServers: await getIceServers(),
      })

      fgg.setOnProgress((c: number): void => {
        setProgress(c)
      })

      fgg.setSend(new DomSendFile(file))
    }

    setLoading(false)
    setDetail(manifest)
  }

  const toggleClose = (i: number) => {
    archive.files.splice(i, 1)
    setFiles([...archive.files])
  }

  return (
    <>
      { detail
        ? <SendFile file={ detail } />
    : <>
      <input
          style={{ display: "none" }}
          // This id e2e test need
          id="upload"
          type="file"
          multiple
          ref={ hiddenFileInput }
          onChange={ (ev: ChangeEvent<HTMLInputElement>) => !!ev.target.files && handleFile(ev.target.files) }
        />

      { !files.length ?
      <div className="flex flex-col items-center rounded-3xl border-5 border-green-500 border-dashed" onClick={toggleButton}>
        <p className="font-bold m-18" >You can click here</p>
        <button className="m-4 px-8 py-2 font-bold text-white rounded-xl bg-purple-600 border border-purple-200">Select files to share</button>
      </div>
    : <>
        <ul className="p-3 bg-gray-100">
          { files.map((file, index) =>
            <li key={ index } className="m-2 border-1 border-green-300 rounded-md bg-green-100 shadow-md flex flex-row justify-between">
              <FileItem file={file}></FileItem>
              <p className="p-4 cursor-pointer" onClick={ () => toggleClose(index) }>x</p>
            </li>
          )}

          <div className="p-2 flex flex-row justify-between">
            <button className="font-medium" onClick={toggleButton}>Add File</button>
            <p>Total size: { filesize(total).toString() }</p>
          </div>
        </ul>

        <div className="p-2">

          <label>Expires after </label>
          <select className="rounded-md cursor-pointer pl-1 pr-8 py-2 border-1"
            value={remain}
            onChange={e => setRemain(Number(e.target.value))}
          >
            <option value="1">1 Download</option>
            <option value="3">3 Downloads</option>
            <option value="5">5 Downloads</option>
            <option value="7">7 Downloads</option>
            <option value="11">11 Downloads</option>
          </select>

          <label> Or </label>
          <select className="rounded-md cursor-pointer pl-1 pr-8 py-2 border-1"
            value={expire}
            onChange={e => setExpire(e.target.value)}
          >
            <option value="5m">5m</option>
            <option value="30m">30m</option>
            <option value="1h">1h</option>
            <option value="24h">24h</option>
          </select>

        </div>
        <progress className="w-full h-1" value={ progress } max={ archive.size } ></progress>


          <div className="p-2 flex flex-row justify-between">
            <div>
              <input className="mr-1" type="checkbox" id="relay" name="scales" checked={ relay } onChange={ e => setRelay(e.target.checked) } />
              <label>Server Relay</label>
            </div>

          </div>

          <button className="p-3 w-full block border-1 rounded-md bg-blue-500 text-white font-bold flex flex-row justify-center" disabled={ loading } onClick={ () => toggleCommit(relay) }>
            { loading
              ? <Loading></Loading>
              : null
            }
            Commit
            { loading
              ? "..."
              : null
            }
          </button>
      </>
    }</>
    }</>
  )
}
