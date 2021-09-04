import { useRef, useState, ChangeEvent } from 'react'
import styles from './File.module.scss'

function File(props: { handleFile: (files: FileList | null) => void }) {
  const hiddenFileInput = useRef<HTMLInputElement>(null)
  const handleClick = () => {
    hiddenFileInput.current?.click?.()
  }

  const [filename, setFilename] = useState('Select File')

  const handleFile = (files: FileList | null) => {
    props.handleFile(files)
    const file = files?.[0]
    if (file) setFilename(file.name)
  }

  return (
    <>
      <label className={ styles.button } onClick={ handleClick } >{ filename }</label>
      <input
        className={ styles.input }
        type="file"
        ref={ hiddenFileInput }
        onChange={ (ev: ChangeEvent<HTMLInputElement>) => { handleFile(ev.target.files) } }
      />
    </>
  )
}

export default File
