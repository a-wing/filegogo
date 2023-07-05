import { useAtom } from "jotai"
import Copy from "copy-to-clipboard"

import FileItem from "./file-item"
import { ItemsAtom, DetailAtom } from "../store"
import { getRaw, delBox, generateShare } from "../lib/api"
import { ExpiresAtHumanTime } from "../lib/util"

export default () => {
  const [_, setDetail] = useAtom(DetailAtom)
  const [files, setFiles] = useAtom(ItemsAtom)

  const toggleClose = async (i: number) => {
    let item = files.splice(i, 1)[0]
    localStorage.removeItem(item.uxid)
    setFiles([...files])
    item.secret && await delBox(item.uxid, item.secret)
  }

  const toggleShow = async (i: number) => setDetail(files[i])
  const toggleDownload = async (uxid: string) => await getRaw(uxid)

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
              <button className="cursor-pointer" onClick={ () => toggleDownload(file.uxid) }>Download</button>
              <button className="cursor-pointer" onClick={ () => toggleShow(index) }>Show</button>
              <button className="cursor-pointer" onClick={ () => Copy(generateShare(file.uxid)) }>Copy Link</button>
            </div>

          </li>
        )}
      </ul>
    </>
  )
}
