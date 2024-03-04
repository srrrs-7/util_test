import { defineConfig } from 'npm:vite@^3.2.3'
import { svelte } from 'npm:@sveltejs/vite-plugin-svelte@^1.1.0'

import 'npm:svelte@^3.52.0'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    lib: {
      entry: 'src/index.ts',
      name: 'index'
    }
  }
})

