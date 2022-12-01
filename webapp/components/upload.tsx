import styles from './upload.module.scss'

function Upload(props: { name: string, file: any }) {
  const send = () => {
    const form = document.getElementById('form-upload')
    const formData = new FormData(form as HTMLFormElement)
    const xhr = new XMLHttpRequest()
    xhr.open("POST", `/raw/${props.name}`, true)
    xhr.send(formData)
  }

  return (
    <>
      <button className={ styles.button } onClick={ send } >Send Relay Server</button>
    </>
  )
}

export default Upload
