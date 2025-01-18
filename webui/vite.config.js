import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig(({command, mode, ssrBuild}) => {
	const ret = {
		plugins: [vue()],
		resolve: {
			alias: {
				'@': fileURLToPath(new URL('./src', import.meta.url))
			}
		},
	};
	ret.define = {
		// Do not modify this constant, it is used in the evaluation.
		"__API_URL__": JSON.stringify("http://localhost:3000"),
	};

	return ret;
}
)

//
// // from
// import { defineConfig } from 'vite'
// import vue from '@vitejs/plugin-vue'
// import vueDevTools from 'vite-plugin-vue-devtools'
//
// export default defineConfig({
// 	plugins: [
// 		vue(),
// 		vueDevTools(),
// 	],
// 	server: {
// 		proxy: {
// 			'/webapi': {
// 				target: 'http://localhost:3000',
// 				changeOrigin: true,
// 			}
// 		},
// 		host: true,
// 		port: 5173
// 	}
// })
