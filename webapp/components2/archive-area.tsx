import FileItem from "./file-item"
import Copy from "copy-to-clipboard"

import { useAtom } from "jotai"
import { ItemsAtom } from "../store"

export default () => {
  const [files, setFiles] = useAtom(ItemsAtom)

  const toggleClose = (i: number) => {
    console.log(i)
    files.splice(i, 1)
    setFiles([...files])
  }

  return (
    <>
      <ul className="p-3">
        { files.map((file, index) =>
          <li key={index} className="m-2 p-4 border-1 border-green-300 rounded-md bg-green-100 shadow-md ">
            <div className="flex flex-row justify-between">

            <FileItem file={file}></FileItem>
            <p className="p-4 cursor-pointer" onClick={ () => toggleClose(index) }>x</p>

            </div>
            <div>
              Expires after 1 download or 23h 26m
            </div>
            { file.files.length > 1
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
              <div>Download</div>
              <button className="cursor-pointer" onClick={ () => Copy("Copy Link: " + file.name) }>Copy Link</button>
            </div>

          </li>
        )}
      </ul>
    </>
  )
}
