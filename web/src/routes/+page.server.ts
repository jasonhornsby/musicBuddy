import { getDemoFiles } from "../remote/demo-files.remote";
import type { PageServerLoad } from "./$types";

export const ssr = false;

export const load: PageServerLoad = async (event) => {
    return {
        demoFiles: await getDemoFiles()
    }
}