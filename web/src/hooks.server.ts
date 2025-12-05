import type { Handle } from "@sveltejs/kit";

export const handle: Handle = async ({ event, resolve }) => {
    const response = await resolve(event);
    if (event.url.pathname.startsWith('/workers')) {
        response.headers.set('Cross-Origin-Embedder-Policy', 'require-corp');
    }
    return response;
}