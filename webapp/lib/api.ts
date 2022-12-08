import { shareGetRoom } from './share'

const prefix = '/s/'

function getServer(): string {
  return `${import.meta.env.VITE_APP_SERVER || window.location.origin}${prefix}`
}

function getLogLevel(): string {
  const loglevel = import.meta.env.VITE_APP_LOG_LEVEL
  return loglevel ? loglevel : 'info'
}

async function getConfig(): Promise<RTCIceServer[]> {
  const response = await fetch("/config")
  const result = await response.json()
  return result.iceServers || []
}

async function getRoom(): Promise<string> {
  const str = shareGetRoom(window.location.href)
  if (str !== '') return str

  const response = await fetch("/s/")
  const result = await response.json()
  return result.room || ''
}

async function putRawFile(f: File): Promise<void> {
  const room = shareGetRoom(window.location.href)
  if (room === '') throw "not room"

  let formData = new FormData()
  formData.append('f', f, f.name)
  await fetch(`/raw/${room}`, {
    method: "post",
    body: formData,
  })
  return
}

async function getRawFile(): Promise<void> {
  const room = shareGetRoom(window.location.href)
  window.open(`/raw/${room}`)
}

async function getRawInfo(): Promise<any> {
  const room = shareGetRoom(window.location.href)
  const response = await fetch(`/info/${room}`)
  if (response.status == 200) {
    return await response.json()
  }
}

export {
  getServer,
  getConfig,
  getRoom,
  getLogLevel,
  putRawFile,
  getRawFile,
  getRawInfo,
  shareGetRoom,
}
