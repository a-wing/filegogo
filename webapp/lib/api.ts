import { shareGetRoom } from './share'

const ws = '/s/'

function getPrefix(): string {
  const url = window.location
  return url.pathname.slice(0, url.pathname.lastIndexOf('/'))
}

function getServer(): string {
  return `${import.meta.env.VITE_APP_SERVER || window.location.origin}${getPrefix()}${ws}`
}

function getLogLevel(): string {
  const loglevel = import.meta.env.VITE_APP_LOG_LEVEL
  return loglevel ? loglevel : 'info'
}

async function getConfig(): Promise<RTCIceServer[]> {
  const response = await fetch(`${getPrefix()}/config`)
  const result = await response.json()
  return result.iceServers || []
}

async function getRoom(): Promise<string> {
  const str = shareGetRoom(window.location.href)
  if (str !== '') return str

  const response = await fetch(`${getPrefix()}/s/`)
  const result = await response.json()
  return result.room || ''
}

async function putBoxFile(f: File): Promise<void> {
  const room = shareGetRoom(window.location.href)
  if (room === '') throw "not room"

  let formData = new FormData()
  formData.append('f', f, f.name)
  await fetch(`${getPrefix()}/api/file/${room}`, {
    method: "post",
    body: formData,
  })
  return
}

async function getBoxFile(): Promise<void> {
  const room = shareGetRoom(window.location.href)
  window.open(`${getPrefix()}/api/file/${room}`)
}

async function delBoxFile(): Promise<void> {
  const room = shareGetRoom(window.location.href)
  await fetch(`${getPrefix()}/api/file/${room}`, {
    method: "delete",
  })
}

async function getBoxInfo(): Promise<any> {
  const room = shareGetRoom(window.location.href)
  const response = await fetch(`${getPrefix()}/api/info/${room}`)
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
