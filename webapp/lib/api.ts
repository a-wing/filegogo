import { Box } from "../libfgg/index"
import { shareGetRoom, generateShare } from "./share"

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

async function getIceServers(): Promise<RTCIceServer[]> {
  const response = await fetch(`${getPrefix()}/config`)
  const result = await response.json()
  return result.iceServers || []
}

async function putBox(f: File, remain: number, expire: string, action: string): Promise<Box> {
  let formData = new FormData()
  if (action === "relay") {
    formData.append('file', f, f.name)
  }
  return (await fetch(`${getPrefix()}/box?remain=${remain}&expire=${expire}&action=${action}`, {
    method: "post",
    body: formData,
  })).json()
}

async function delBox(uxid: string, secret: string): Promise<void> {
  await fetch(`${getPrefix()}/box/${uxid}?secret=${secret}`, {
    method: "delete",
  })
}

async function getRaw(uxid: string): Promise<void> {
  window.open(`${getPrefix()}/raw/${uxid}`)
}

async function getBox(room: string): Promise<Box | void> {
  const response = await fetch(`${getPrefix()}/box/${room}`)
  if (response.status == 200) {
    return await response.json()
  }
}

export {
  getServer,
  getIceServers,
  getLogLevel,

  putBox,
  getBox,
  delBox,

  getRaw,

  shareGetRoom,
  generateShare,
}
