import { getParams } from './share'

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
  const str = getParams(window.location.href)
  if (str !== '') return str

  const response = await fetch("/s/")
  const result = await response.json()
  return result.room || ''
}

function shareGetRoom(addr: string): string {
  const u = new URL(addr)
  const arr = u.pathname.split("/")
  if (arr.length > 0) {
    const a2 = arr[arr.length - 1]
    return a2.match(/\d/) ? a2 : ""
  }
	return ""
}

export {
  getServer,
  getConfig,
  getRoom,
  getLogLevel,
  shareGetRoom,
}
