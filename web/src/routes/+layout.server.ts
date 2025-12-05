import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async (event) => {
    event.setHeaders({
        'Cross-Origin-Opener-Policy': 'same-origin',
        'Cross-Origin-Embedder-Policy': 'require-corp',
    })
}