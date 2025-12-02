<script lang="ts">
    import * as Sidebar from '$lib/components/ui/sidebar';
</script>

<Sidebar.Provider>
	<div class="h-full w-full">
		<main></main>
	</div>
	<Sidebar.Root side="right">
		<Sidebar.Header>
			<h1>Analyser</h1>
		</Sidebar.Header>
		<Sidebar.Content>
			<Sidebar.Group>
				<Sidebar.GroupLabel>Views</Sidebar.GroupLabel>
				<Sidebar.GroupContent>
					<Sidebar.Menu>
						<Sidebar.MenuButton>General Information</Sidebar.MenuButton>
						<Sidebar.MenuButton>Waveform</Sidebar.MenuButton>
						<Sidebar.MenuButton>Spectral Flux</Sidebar.MenuButton>
						<Sidebar.MenuButton>Spectrogram</Sidebar.MenuButton>
					</Sidebar.Menu>
				</Sidebar.GroupContent>
			</Sidebar.Group>
		</Sidebar.Content>
	</Sidebar.Root>
</Sidebar.Provider>
