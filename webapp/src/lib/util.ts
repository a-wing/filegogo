
function ProtoHttpToWs(s: string): string {
  return s.replace(/^http/, 'ws')
}

function ProtoWsToHttp(s: string): string {
  return s.replace(/^ws/, 'http')
}

export {
  ProtoHttpToWs,
  ProtoWsToHttp,
}
