import { shareGetRoom } from './share'

const ws = '/signal/'

function getPrefix(): string {
  const url = window.location
  return url.pathname.slice(0, url.pathname.lastIndexOf('/')) + '/api'
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

  const response = await fetch(`${getPrefix()}/${ws}/`)
  const result = await response.json()
  return result.room || ''
}

async function putBoxFile(room: string, f: File, remain: number, expire: string): Promise<void> {
  let formData = new FormData()
  formData.append('f', f, f.name)
  await fetch(`${getPrefix()}/file/${room}?remain=${remain}&expire=${expire}`, {
    method: "post",
    body: formData,
  })
  return
}

async function getBoxFile(room: string): Promise<void> {
  window.open(`${getPrefix()}/file/${room}`)
}

async function delBoxFile(room: string): Promise<void> {
  await fetch(`${getPrefix()}/file/${room}`, {
    method: "delete",
  })
}

async function getBoxInfo(room: string): Promise<any> {
  const response = await fetch(`${getPrefix()}/info/${room}`)
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
