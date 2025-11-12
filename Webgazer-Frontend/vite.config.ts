import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import type { Plugin } from 'vite';

export default defineConfig({
	plugins: [
		tailwindcss(),
		sveltekit(),
		// Suppress 404 warnings for missing source maps
		{
			name: 'suppress-sourcemap-404',
			configureServer(server) {
				server.middlewares.use((req, res, next) => {
					if (req.url?.endsWith('.map')) {
						// Silently handle missing source map requests
						res.statusCode = 404;
						res.end();
						return;
					}
					next();
				});
			}
		} as Plugin
	]
});
