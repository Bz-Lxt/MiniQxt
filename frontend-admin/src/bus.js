import { reactive } from 'vue'

export const ui = reactive({
  toasts: [],
  modal: null,
})

let n = 0
export function toast(message, type = 'ok') {
  const id = ++n
  ui.toasts.push({ id, message, type })
  setTimeout(() => dismiss(id), 5000)
}
export function dismiss(id) {
  ui.toasts = ui.toasts.filter((t) => t.id !== id)
}
export function confirm(title, body) {
  return new Promise((resolve) => {
    ui.modal = { title, body, resolve }
  })
}
export function closeModal(ok) {
  const m = ui.modal
  ui.modal = null
  m?.resolve(ok)
}
