// The real auth credential is an httpOnly cookie the JS can't read. This flag is
// only a UX hint so the SPA can route immediately without an extra round-trip;
// the server still enforces auth on every /api call (a 401 clears the flag).
const KEY = 'pw_logged_in'

export function setLoggedIn() {
  localStorage.setItem(KEY, '1')
}

export function clearLoggedIn() {
  localStorage.removeItem(KEY)
}

export function isLoggedIn(): boolean {
  return localStorage.getItem(KEY) === '1'
}
