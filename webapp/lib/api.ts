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

export {
  getServer,
  getConfig,
  getRoom,
  getLogLevel,
  shareGetRoom,
}
