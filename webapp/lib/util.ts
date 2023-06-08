
function ProtoHttpToWs(s: string): string {
  return s.replace(/^http/, 'ws')
}

function ProtoWsToHttp(s: string): string {
  return s.replace(/^ws/, 'http')
}

function ExpiresAtHumanTime(timestamp: string): string {
  const now = new Date()
  const diff = new Date(timestamp).getTime() - now.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) {
    return `${days} Day After`
  } else if (hours > 0) {
    return `${hours} Hour After`
  } else if (minutes > 0) {
    return `${minutes} Minute After`
  } else {
    return `${seconds} Seconds After`
  }
}

export {
  ProtoHttpToWs,
  ProtoWsToHttp,
  ExpiresAtHumanTime,
}
