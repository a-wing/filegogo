import FileItem from "./file-item"
import Copy from "copy-to-clipboard"
import Qrcode from "./qr-code"

import { useAtom } from "jotai"
import { ItemsAtom } from "../store"

import { getBoxFile, delBoxFile, generateShare } from "../lib/api"
import { ExpiresAtHumanTime } from "../lib/util"

export default () => {
  const [files, setFiles] = useAtom(ItemsAtom)

  const toggleClose = async (i: number) => {
    let item = files.splice(i, 1)
    item[0] ? localStorage.removeItem(item[0].uxid) : null
    setFiles([...files])
    await delBoxFile(item[0].uxid)
  }

  const toggleDownload = async (i: number) => {
    let item = files[i]
    await getBoxFile(item.uxid)
  }

  return (
    <>
      <ul className="p-3">
        { files.map((file, index) =>
          <li key={index} className="m-2 p-4 border-1 border-green-300 rounded-md bg-green-100 shadow-md">
            <div className="flex flex-row justify-between">

            <FileItem file={file}></FileItem>
            <p className="p-4 cursor-pointer" onClick={ () => toggleClose(index) }>x</p>

            </div>
            <div>
              Expires after { file.remain } download or { ExpiresAtHumanTime(file.expire) }
            </div>
            { file.files?.length > 1
              ? <details className="cursor-pointer">
                  <summary>{ file.files.length.toString() + " files" }</summary>
                { file.files.map((file, index) =>
                  <div key={index} className="mx-8">
                    <FileItem file={file}></FileItem>
                  </div>
                )}
                </details>
              : null
            }
            <hr className="my-2" />
            <div className="flex flex-row justify-between">
              <button className="cursor-pointer" onClick={ () => { toggleDownload(index) } }>Download</button>
              <details className="cursor-pointer">
                <summary>QRCode</summary>
                <Qrcode address={ generateShare(file.uxid) }></Qrcode>
              </details>
              <button className="cursor-pointer" onClick={ () => Copy(generateShare(file.uxid)) }>Copy Link</button>
            </div>

          </li>
        )}
      </ul>
    </>
  )
}
