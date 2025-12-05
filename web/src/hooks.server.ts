import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const response = await resolve(event);

    // Set COOP/COEP headers required for SharedArrayBuffer and WebAssembly
    // These apply to all SvelteKit-handled routes
    // Static files in /static are handled by vite.config.ts server.headers in dev
    response.headers.set('Cross-Origin-Opener-Policy', 'same-origin');
    response.headers.set('Cross-Origin-Embedder-Policy', 'require-corp');
    response.headers.set('Cross-Origin-Resource-Policy', 'cross-origin');

    return response;
}