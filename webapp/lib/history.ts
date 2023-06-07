import { Manifest } from "./manifest"
import { getBoxInfo } from "./api"

async function loadHistory(): Promise<Array<Manifest>> {
  let p: Array<Promise<Manifest>> = []
  for (let i = 0; i < localStorage.length; i++) {
    const k = localStorage.key(i)
    if (!k) continue
    try {
      let value = localStorage.getItem(k)
      if (!value) continue
      p.push(getBoxInfo(k))
    } catch (e) {
      console.log(e)
    }
  }
  return (await Promise.all(p)).filter(i => {
    if (typeof i === "string") {
      localStorage.removeItem(i)
      return false
    }
    return true
  })
}

export {
  loadHistory,
}
