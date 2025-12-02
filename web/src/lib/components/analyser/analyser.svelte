<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { ChevronDown } from 'lucide-svelte';
	import { DropdownMenu, DropdownMenuContent, DropdownMenuItem } from '../ui/dropdown-menu';
	import DropdownMenuTrigger from '../ui/dropdown-menu/dropdown-menu-trigger.svelte';

	type View = {
		id: string;
		name: string;
	};

	const views: View[] = [
		{ id: 'general', name: 'General Information' },
		{ id: 'waveform', name: 'Waveform' },
		{ id: 'spectral-flux', name: 'Spectral Flux' },
		{ id: 'spectrogram', name: 'Spectrogram' }
	];

	let selectedView = $state<View>(views[0]);
</script>

<Sidebar.Provider>
	<div class="h-full w-full">
		<main></main>
	</div>
	<Sidebar.Root side="right">
		<Sidebar.Header>
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<DropdownMenu>
						<DropdownMenuTrigger>
							{#snippet child({ props })}
								<Sidebar.MenuButton {...props}>
									{selectedView.name}
									<ChevronDown class="ms-auto" />
								</Sidebar.MenuButton>
							{/snippet}
						</DropdownMenuTrigger>
						<DropdownMenuContent>
							{#each views as view (view.id)}
								<DropdownMenuItem onclick={() => (selectedView = view)}>
									{view.name}
								</DropdownMenuItem>
							{/each}
						</DropdownMenuContent>
					</DropdownMenu>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		</Sidebar.Header>
		<Sidebar.Content></Sidebar.Content>
	</Sidebar.Root>
</Sidebar.Provider>
