<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="h-dvh w-dvw flex flex-col overflow-hidden">
	<header class="p-4 border-b">
		<h1>Music Analyser</h1>
	</header>
	<main class="flex-1 overflow-y-auto">
		{@render children()}
	</main>
</div>
