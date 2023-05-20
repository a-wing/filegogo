import { useRef, useState, ChangeEvent } from "react"
import { Item } from "../lib/manifest"
import FileItem from "./file-item"

export default () => {
  const [files, setFiles] = useState<Array<Item>>([{
    name: "aaa",
    type: "xxx",
    size: 123,
    files: [],
  }, {
    name: "aaa2",
    type: "xxx2",
    size: 1234,
    files: [
      {
        name: "a2-1",
        type: "xxx",
        size: 123,
      }, {
        name: "a2-2",
        type: "xxx",
        size: 123,
      },
    ],
  }, {
    name: "aaa3",
    type: "xxx3",
    size: 12345,
    files: [
      {
        name: "a3-1",
        type: "xxx",
        size: 123,
      }, {
        name: "a3-2",
        type: "xxx",
        size: 123,
      },
    ],
  }])

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
            <p className="p-4 cursor-pointer" onClick={ () => toggleClose(index) }>X</p>

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
              <div>Copy Link</div>
            </div>

          </li>
        )}
      </ul>
    </>
  )
}
