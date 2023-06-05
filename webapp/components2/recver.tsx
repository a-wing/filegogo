import { useRef, useState, ChangeEvent, useEffect } from "react"
import { Item } from "../lib/manifest"
import RecverFile from "./recver-file"
import { getBoxFile, getBoxInfo, shareGetRoom } from "../lib/api"

//export default (props: {
//  file: Item
//}) => {
export default () => {
  const [file, setFile] = useState<Item | null>(
    //{
    //  name: "aaa2",
    //  type: "xxx2",
    //  size: 1234,
    //  files: [
    //    {
    //      name: "a2-1",
    //      type: "xxx",
    //      size: 123,
    //    }, {
    //      name: "a2-2",
    //      type: "xxx",
    //      size: 123,
    //    },
    //  ],
    //}
  )

  const loadFile = async () => {
    setFile(await getBoxInfo(shareGetRoom(window.location.href)))
  }


  useEffect(() => {
    loadFile()
  }, [])

  return (
    <div className="w-full flex flex-col">
      <div className="flex flex-row justify-center">
        <h1 className="font-bold text-4xl">Download files</h1>
      </div>
      <div className="flex flex-row justify-center">
        <p>This file was shared via Send with end-to-end encryption and a link that automatically expires.</p>
      </div>
      { file
        ? <RecverFile file={ file }></RecverFile>
        : <></>
      }
   </div>
  )
}
