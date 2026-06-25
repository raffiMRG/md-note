import { ref } from 'vue'

const isDark = ref(
  localStorage.getItem('theme') === 'dark' ||
  (!localStorage.getItem('theme') && window.matchMedia('(prefers-color-scheme: dark)').matches),
)

function applyTheme(dark) {
  document.documentElement.dataset.theme = dark ? 'dark' : 'light'
  localStorage.setItem('theme', dark ? 'dark' : 'light')
}

applyTheme(isDark.value)

export function useTheme() {
  function toggle() {
    isDark.value = !isDark.value
    applyTheme(isDark.value)
  }
  return { isDark, toggle }
}
