import { Meta } from "../lib/archive"
import logo from "/logo192.png"

export default (props: {
  file: Meta
}) => {
  return (
    <div className="flex flex-row p-1">
      <img className="h-12" src={ logo } alt="logo" />
      <div className="flex flex-col">
        <h1 className="font-medium">{ props.file.name }</h1>
        <p>{ props.file.size }</p>
      </div>
    </div>
  )
}
