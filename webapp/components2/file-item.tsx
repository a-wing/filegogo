import { Meta } from "../libfgg"
import logo from "/logo192.png"
import { filesize } from "filesize"

export default (props: {
  file: Meta
}) => {
  return (
    <div className="flex flex-row p-1">
      <img className="h-12" src={ logo } alt="logo" />
      <div className="flex flex-col p-1">
        <h1 className="font-medium">{ props.file.name }</h1>
        <p>{ filesize(props.file.size).toString() }</p>
      </div>
    </div>
  )
}
