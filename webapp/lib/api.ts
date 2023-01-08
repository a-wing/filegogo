import { shareGetRoom } from './share'

declare global {
  interface Window {
    sub_folder: String;
  }
}

const prefix = window.sub_folder + '/s/'

function getServer(): string {
  return `${import.meta.env.VITE_APP_SERVER || window.location.origin}${prefix}`
}

function getLogLevel(): string {
  const loglevel = import.meta.env.VITE_APP_LOG_LEVEL
  return loglevel ? loglevel : 'info'
}

async function getConfig(): Promise<RTCIceServer[]> {
  const response = await fetch(window.sub_folder + "/config")
  const result = await response.json()
  return result.iceServers || []
}

async function getRoom(): Promise<string> {
  const str = shareGetRoom(window.location.href)
  if (str !== '') return str

  const response = await fetch(window.sub_folder + "/s/")
  const result = await response.json()
  return result.room || ''
}

async function putBoxFile(f: File): Promise<void> {
  const room = shareGetRoom(window.location.href)
  if (room === '') throw "not room"

  let formData = new FormData()
  formData.append('f', f, f.name)
  await fetch(`/api/file/${room}`, {
    method: "post",
    body: formData,
  })
  return
}

async function getBoxFile(): Promise<void> {
  const room = shareGetRoom(window.location.href)
  window.open(`/api/file/${room}`)
}

async function delBoxFile(): Promise<void> {
  const room = shareGetRoom(window.location.href)
  await fetch(`/api/file/${room}`, {
    method: "delete",
  })
}

async function getBoxInfo(): Promise<any> {
  const room = shareGetRoom(window.location.href)
  const response = await fetch(`/api/info/${room}`)
  if (response.status == 200) {
    return await response.json()
  }
}

export {
  getServer,
  getConfig,
  getRoom,
  getLogLevel,
  putBoxFile,
  getBoxFile,
  delBoxFile,
  getBoxInfo,
  shareGetRoom,
}
