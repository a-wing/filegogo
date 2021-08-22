import { getParams } from './share'

const prefix = '/s/'

function getServer(): string {
  return `${process.env.REACT_APP_SERVER || window.location.origin}${prefix}`
}

async function getRoom(): Promise<string> {
  const str = getParams(window.location.href)
  if (str !== '') return str

  const response = await fetch("/s/")
  const result = await response.json()
  return result.room || ''
}

export {
  getServer,
  getRoom,
}
