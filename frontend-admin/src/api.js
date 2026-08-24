export const API = import.meta.env.VITE_API_BASE || 'http://localhost:28392'

export function fmt(s) {
  if (!s) return '—'
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return String(s).replace('T', ' ').slice(0, 19)
  const p = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
}

export async function req(path, opts = {}) {
  const headers = { 'Content-Type': 'application/json', ...(opts.headers || {}) }
  const tok = localStorage.getItem('qxt_token')
  if (tok) headers.Authorization = `Bearer ${tok}`
  const res = await fetch(API + path, { ...opts, headers })
  const text = await res.text()
  let body = null
  try { body = text ? JSON.parse(text) : null } catch { body = { message: text } }
  if (!res.ok) {
    const err = new Error(body?.message || '请求失败')
    err.code = body?.code
    err.status = res.status
    throw err
  }
  return body?.data ?? body
}
