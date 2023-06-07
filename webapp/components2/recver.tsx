import { useState, useEffect } from "react"
import { Item } from "../lib/manifest"
import RecverFile from "./recver-file"
import { getBoxInfo, shareGetRoom } from "../lib/api"

export default () => {
  const [file, setFile] = useState<Item | null>(null)

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
        : <div className="flex flex-col items-center">
            <button className="px-8 py-2 text-white rounded-xl bg-purple-600 border border-purple-200"
              onClick={ () => window.location.pathname = "/"}
            >Return Home</button>
          </div>
      }
   </div>
  )
}
