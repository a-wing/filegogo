import { defineConfig } from 'windicss/helpers'
import formsPlugin from 'windicss/plugin/forms'

export default defineConfig({
  extract: {
    include: ['webapp/*.{html,jsx,tsx}', 'webapp/components2/*.{jsx,tsx}'],
    exclude: ['node_modules', '.git', 'dist'],
  },
  darkMode: 'class',
  shortcuts: {
    'btn': 'py-2 px-4 font-semibold rounded-lg shadow-md',
    'btn-green': 'text-white bg-green-500 hover:bg-green-700',
  },
  plugins: [formsPlugin],
})
