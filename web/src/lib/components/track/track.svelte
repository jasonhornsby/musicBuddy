<script lang="ts">
	import { EllipsisVertical } from "lucide-svelte";
	import Visualisation from "../analyser/panels/visualisation.svelte";
	import Button from "../ui/button/button.svelte";
	import ButtonGroup from "../ui/button-group/button-group.svelte";
	import {
		DropdownMenu,
		DropdownMenuTrigger,
		DropdownMenuContent,
		DropdownMenuItem,
	} from "../ui/dropdown-menu";

	const trackName = $state("Track 1");
	let channel = $state<"left" | "right" | "mono">("mono");
</script>

<div
	class="h-[280px] w-full rounded-lg border border-slate-200 bg-white flex flex-row overflow-hidden shadow-sm"
>
	<aside
		class="w-48 max-w-[28%] border-r border-slate-200 bg-slate-50/80 text-[11px] flex flex-col"
	>
		<div class="flex items-center justify-between gap-1.5 px-2 py-2 border-b border-slate-100">
			<div class="flex flex-col gap-0.5 min-w-0">
				<h2 class="text-[12px] font-semibold text-slate-800 truncate">
					{trackName}
				</h2>
			</div>
			<DropdownMenu>
				<DropdownMenuTrigger>
					{#snippet child({ props })}
						<Button size="icon-sm" variant="ghost" class="h-7 w-7" {...props}>
							<EllipsisVertical class="w-3.5 h-3.5 text-slate-500" />
						</Button>
					{/snippet}
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem onclick={() => console.log('Remove track clicked')}>
						Remove track
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<div class="flex flex-col gap-2 px-2 py-2">
			<div class="flex flex-col gap-1">
				<span class="text-[10px] font-medium text-slate-500 uppercase tracking-[0.14em]">
					Channel
				</span>
				<ButtonGroup class="mt-0.5 w-full">
					<Button
						variant={channel === 'left' ? 'default' : 'outline'}
						size="sm"
						class="flex-1 text-[11px] px-1.5"
						onclick={() => (channel = 'left')}
					>
						Left
					</Button>
					<Button
						variant={channel === 'right' ? 'default' : 'outline'}
						size="sm"
						class="flex-1 text-[11px] px-1.5"
						onclick={() => (channel = 'right')}
					>
						Right
					</Button>
					<Button
						variant={channel === 'mono' ? 'default' : 'outline'}
						size="sm"
						class="flex-1 text-[11px] px-1.5"
						onclick={() => (channel = 'mono')}
					>
						Mono
					</Button>
				</ButtonGroup>
			</div>
		</div>
	</aside>
	<div class="flex-1">
		<Visualisation />
	</div>
</div>
