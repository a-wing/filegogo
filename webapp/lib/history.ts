import { Box } from "../libfgg"
import { getBox } from "./api"

type warpResult = Promise<{ uxid: string, data: Box }>

async function loadApiInfo(uxid: string): warpResult {
  return {
    uxid: uxid,
    data: await getBox(uxid) as Box,
  }
}

async function loadHistory(): Promise<Array<Box>> {
  let promises: Array<warpResult> = []
  for (let i = 0; i < localStorage.length; i++) {
    const k = localStorage.key(i)
    if (!k) continue
    try {
      let value = localStorage.getItem(k)
      if (!value) continue
      promises.push(loadApiInfo(k))
    } catch (e) {
      console.log(e)
    }
  }
  return (await Promise.all(promises)).filter(i => i.data ? true : localStorage.removeItem(i.uxid)).map(i => i.data)
}

export {
  loadHistory,
}
