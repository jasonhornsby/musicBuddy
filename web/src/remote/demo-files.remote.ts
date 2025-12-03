import { query } from "$app/server";

export type DemoFile = {
    name: string;
    src: string;
}

export const getDemoFiles = query(async (): Promise<DemoFile[]> => {
    return [
        {
            name: "Jazz",
            src: "examples/jazz.mp3"
        },
        {
            name: "Hip Hop",
            src: "examples/hiphop.mp3"
        },
        {
            name: "Error",
            src: "examples/error.mp3"
        }
    ]
})