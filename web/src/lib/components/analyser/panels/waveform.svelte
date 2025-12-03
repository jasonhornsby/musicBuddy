<script lang="ts">
	import { Separator } from "$lib/components/ui/separator";
</script>

<div class="flex flex-col h-full">
	<header class="px-5 py-4 flex items-center justify-between">
		<h2 class="text-sm font-medium tracking-wide uppercase text-muted-foreground">Waveform</h2>
	</header>

	<Separator />

	<div class="flex-1 p-5 overflow-y-auto">
		<div class="flex items-center justify-center h-full">
			<div class="text-center space-y-3">
				<div class="size-16 mx-auto rounded-full bg-muted/50 flex items-center justify-center">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="size-8 text-muted-foreground"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.5"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M2 12h2l2-7 3 14 3-9 2 4h8" />
					</svg>
				</div>
				<p class="text-sm text-muted-foreground">Waveform visualization coming soon</p>
			</div>
		</div>
	</div>
</div>

