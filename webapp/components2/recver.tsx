import { useRef, useState, ChangeEvent } from "react"
import { Item } from "../lib/manifest"
import FileItem from "./file-item"

//export default (props: {
//  file: Item
//}) => {
export default () => {
  const [file] = useState<Item>({
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
  })

  return (
    <div className="w-full flex flex-col">
      <div className="flex flex-row justify-center">
        <h1 className="font-bold text-4xl">Download files</h1>
      </div>
      <div className="flex flex-row justify-center">
        <p>This file was shared via Send with end-to-end encryption and a link that automatically expires.</p>
      </div>
      <ul className="p-3 w-full flex flex-row justify-center">
        <li className="m-2 p-4 border-1 border-green-300 rounded-md bg-green-100 shadow-md w-full">
          <div className="flex flex-row justify-between">
            <FileItem file={file}></FileItem>
          </div>
          { file.files.length > 1
            ? <details className="cursor-pointer">
                <summary>{ file.files.length.toString() + " files" }</summary>
            { file.files.map((file, index) =>
              <div key={index} className="mx-8">
                <FileItem file={file}></FileItem>
              </div>)
            } </details>
            : null
          }
        </li>
      </ul>

      <button className="p-3 w-full block border-1 rounded-md bg-blue-500 text-white font-bold">Download</button>
    </div>
  )
}
