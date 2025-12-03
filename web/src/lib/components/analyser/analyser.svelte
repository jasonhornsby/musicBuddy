<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Item from '$lib/components/ui/item';
	import { ChevronDown, RotateCcw } from 'lucide-svelte';
	import { getAudioContext } from '$lib/context/audio.svelte';

	type View = {
		id: string;
		name: string;
		description: string;
	};

	const views: View[] = [
		{ id: 'general', name: 'General Information', description: 'Overview of file metadata, duration, sample rate, and channels' },
		{ id: 'waveform', name: 'Waveform', description: 'Visual representation of amplitude over time' },
		{ id: 'spectral-flux', name: 'Spectral Flux', description: 'Measure of how quickly the spectrum changes, useful for onset detection' },
		{ id: 'spectrogram', name: 'Spectrogram', description: 'Frequency content visualization across time' }
	];

	let selectedView = $state<View>(views[0]);

    const audioContext = getAudioContext();
</script>

<Sidebar.Provider>
	<div class="h-full w-full">
		<main></main>
	</div>
	<Sidebar.Root side="right">
		<Sidebar.Header>
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<DropdownMenu.Root>
						<DropdownMenu.Trigger>
							{#snippet child({ props })}
								<Sidebar.MenuButton {...props}>
									{selectedView.name}
									<ChevronDown class="ms-auto" />
								</Sidebar.MenuButton>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content class="w-80">
							{#each views as view (view.id)}
								<DropdownMenu.Item class="p-0" onclick={() => (selectedView = view)}>
									<Item.Root size="sm" class="w-full p-2">
										<Item.Content class="gap-0.5">
											<Item.Title>{view.name}</Item.Title>
											<Item.Description>{view.description}</Item.Description>
										</Item.Content>
									</Item.Root>
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		</Sidebar.Header>
		<Sidebar.Content></Sidebar.Content>
		<Sidebar.Footer>
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton onclick={() => audioContext.resetAudio()}>
						<RotateCcw class="size-4" />Reset file</Sidebar.MenuButton
					>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		</Sidebar.Footer>
	</Sidebar.Root>
</Sidebar.Provider>
