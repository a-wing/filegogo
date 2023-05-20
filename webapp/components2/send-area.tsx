import { useRef, useState, ChangeEvent } from "react"

import Archive, { Meta } from "../lib/archive"
import FileItem from "./file-item"

let archive = new Archive()

export default () => {
  const hiddenFileInput = useRef<HTMLInputElement>(null)
  const [files, setFiles] = useState<Array<Meta>>([{
    name: "aaa",
    type: "xxx",
    size: 123,
  }, {
    name: "aaa2",
    type: "xxx2",
    size: 1234,
  }, {
    name: "aaa3",
    type: "xxx3",
    size: 12345,
  }])

  const toggleButton = () => {
    console.log("toggleButton")
    hiddenFileInput.current?.click?.()
  }

  const handleFile = (filelist: FileList) => {
    let files = new Array<File>(filelist.length)
    for (let i = 0; i < filelist.length; i++) {
      files[i] = filelist[i]
    }
    archive.addFiles(files)

    setFiles([...archive.manifest])
  }

  const toggleClose = (i: number) => {
    console.log(i)
    files.splice(i, 1)
    setFiles([...files])
  }

  return (
    <>{ !files.length ?
      <div className="flex flex-col items-center rounded-3xl border-5 border-green-500 border-dashed" onClick={toggleButton}>
        <input
          style={{ display: "none" }}
          // This id e2e test need
          id="upload"
          type="file"
          multiple
          ref={ hiddenFileInput }
          onChange={ (ev: ChangeEvent<HTMLInputElement>) => ev.target.files ? handleFile(ev.target.files) : null }
        />

        <button className="px-8 py-2 text-white rounded-xl bg-purple-600 border border-purple-200">Select Files</button>
      </div>
    : <>
        <ul className="p-3 bg-gray-100">
          { files.map((file, index) =>
            <li key={ index } className="m-2 border-1 border-green-300 rounded-md bg-green-100 shadow-md flex flex-row justify-between">
              <FileItem file={file}></FileItem>
              <p className="p-4 cursor-pointer" onClick={ () => toggleClose(index) }>X</p>
            </li>
          )}

          <div className="p-2 flex flex-row justify-between">
            <p className="font-medium">Add File</p>
            <p>Size Total</p>
          </div>
        </ul>

        <div className="p-2">

          <label>Expires after </label>
          <select className="rounded-md cursor-pointer pl-1 pr-8">
            <option value="1">1 Download</option>
            <option value="3">3 Downloads</option>
            <option value="5">5 Downloads</option>
            <option value="7">7 Downloads</option>
            <option value="11">11 Downloads</option>
          </select>

          <label> Or </label>
          <select className="rounded-md cursor-pointer pl-1 pr-8">
            <option value="5m">5m</option>
            <option value="30m">30m</option>
            <option value="1h">1h</option>
            <option value="24h">24h</option>
          </select>

        </div>
          <hr className="border-2"/>

          <div className="p-2 flex flex-row justify-between">
            <div>
              <input className="mr-1" type="checkbox" id="relay" name="scales" />
              <label>Server Relay</label>
            </div>

            <div>
              <input className="mr-1" type="checkbox" id="encryption" name="scales" />
              <label>P2P encryption</label>
            </div>
          </div>


        <button className="p-3 w-full block border-1 rounded-md bg-blue-500 text-white font-bold">Commit</button>
      </>
    }</>
  )
}
