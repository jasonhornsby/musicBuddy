<script>
	import { CloudUpload, Music, Plus } from "lucide-svelte";
	import Button from "../ui/button/button.svelte";
	import Track from "../track/track.svelte";
	import AddTrack from "../track/add-track.svelte";
	import { getAudioContext } from "$lib/context/audio.svelte";
	import { Menubar, MenubarContent, MenubarMenu, MenubarTrigger, MenubarItem } from "../ui/menubar";
	import {
		Empty,
		EmptyContent,
		EmptyDescription,
		EmptyHeader,
		EmptyMedia,
		EmptyTitle,
	} from "../ui/empty";

	const audioContext = getAudioContext();

	let addTrackOpen = $state(false);
</script>

<AddTrack bind:open={addTrackOpen} />

<div class="w-full h-dvh flex flex-col">
	<header class="px-2 py-1 border-b flex flex-row gap-2 items-center shrink-0">
		<Button size="icon-sm" variant="ghost">
			<Music size={18} />
		</Button>

		<Menubar>
			<MenubarMenu>
				<MenubarTrigger>Tracks</MenubarTrigger>
				<MenubarContent>
					<MenubarItem onclick={() => (addTrackOpen = true)} disabled={audioContext.isAudioLoaded}>
						<Plus /> Add track
					</MenubarItem>
				</MenubarContent>
			</MenubarMenu>
		</Menubar>
	</header>
	<main class="flex-1 overflow-hidden flex flex-row min-h-0">
		{#if audioContext.isAudioLoaded}
			<Track />
		{:else}
			<Empty class="flex-1 bg-gray-50">
				<EmptyHeader>
					<EmptyMedia variant="icon" class="bg-primary/10 text-primary">
						<Music size={22} />
					</EmptyMedia>
					<EmptyTitle>No tracks loaded</EmptyTitle>
					<EmptyDescription>Load an audio file to start analyzing.</EmptyDescription>
				</EmptyHeader>
				<EmptyContent>
					<Button onclick={() => (addTrackOpen = true)} class="gap-2">
						<CloudUpload class="h-4 w-4" />
						Load track
					</Button>
				</EmptyContent>
			</Empty>
		{/if}
	</main>
</div>
