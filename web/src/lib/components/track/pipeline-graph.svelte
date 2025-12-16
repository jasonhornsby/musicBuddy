<script lang="ts">
	import {
		SvelteFlow,
		Background,
		Controls,
		Position,
		type Node,
		type Edge,
	} from '@xyflow/svelte';
	import '@xyflow/svelte/dist/style.css';
	import dagre from 'dagre';


	export type RenderTree = {
		nodes: {
			id: string;
			type: string; // e.g. "STFTNode"
			category: string; // "compute" | "visualizer"
			label: string;
			duration: number;
		}[];
		edges: {
			from: string;
			to: string;
		}[];
	};

	// --- PROPS ---

	let { data }: { data: RenderTree } = $props();

	// --- STATE ---

	// We use standard $state for the arrays. Svelte Flow handles the internal reactivity.
	let nodes = $state<Node[]>([]);
	let edges = $state<Edge[]>([]);

	// --- LAYOUT LOGIC (DAGRE) ---

	const nodeWidth = 180;
	const nodeHeight = 50;

	const getLayoutedElements = (rawNodes: Node[], rawEdges: Edge[]) => {
		const dagreGraph = new dagre.graphlib.Graph();
		dagreGraph.setDefaultEdgeLabel(() => ({}));

		// Set layout direction: 'LR' = Left to Right, 'TB' = Top to Bottom
		dagreGraph.setGraph({ rankdir: 'LR' });

		rawNodes.forEach((node) => {
			dagreGraph.setNode(node.id, { width: nodeWidth, height: nodeHeight });
		});

		rawEdges.forEach((edge) => {
			dagreGraph.setEdge(edge.source, edge.target);
		});

		dagre.layout(dagreGraph);

		const layoutedNodes = rawNodes.map((node) => {
			const nodeWithPosition = dagreGraph.node(node.id);
			return {
				...node,
				targetPosition: Position.Left,
				sourcePosition: Position.Right,
				position: {
					x: nodeWithPosition.x - nodeWidth / 2,
					y: nodeWithPosition.y - nodeHeight / 2
				}
			};
		});

		return { nodes: layoutedNodes, edges: rawEdges };
	};

	// Watch for changes in the input data and recalculate layout
	const formatDuration = (ms: number): string => {
		if (ms < 1) return '<1ms';
		return `${ms}ms`;
	};

	$effect(() => {
		if (!data) return;

		const initialNodes: Node[] = data.nodes.map((n) => {
			const baseLabel = n.category === 'visualizer' ? `Visualizer: ${n.label}` : n.label;
			const durationLabel = n.duration > 0 ? `\n${formatDuration(n.duration)}` : '';
			
			return {
				id: n.id,
				// We can pass data to be rendered in the node
				data: {
					label: baseLabel + durationLabel,
					type: n.type,
					duration: n.duration
				},
				position: { x: 0, y: 0 }, // Placeholder, Dagre will fix this
				// Add a class based on category for styling
				class: n.category === 'visualizer' ? 'node-viz' : 'node-compute'
			};
		});

		const initialEdges: Edge[] = data.edges.map((e) => ({
			id: `${e.from}-${e.to}`,
			source: e.from,
			target: e.to,
			animated: true, // Animating audio paths looks nice
			style: 'stroke: #b1b1b7;'
		}));

		const layout = getLayoutedElements(initialNodes, initialEdges);

		nodes = layout.nodes;
		edges = layout.edges;
	});
</script>

<div class="graph-container">
	<SvelteFlow {nodes} {edges} fitView minZoom={0.5} maxZoom={2}>
		<Background bgColor="#444" gap={20} size={1} />
		<Controls />
	</SvelteFlow>
</div>

<style>
	/* The container must have a height for Svelte Flow to render */
	.graph-container {
		width: 100%;
		height: 300px; /* Adjust as needed */
		border: 1px solid #333;
		border-radius: 8px;
		background: #1a1a1a;
	}

	/* --- NODE STYLING --- */
	
	/* Global styles for nodes are applied via the :global selector because 
       Svelte Flow renders them outside the component scope */

	:global(.svelte-flow__node) {
		border-radius: 6px;
		font-size: 12px;
		font-family: monospace;
		padding: 10px;
		width: 180px;
		text-align: center;
		color: #fff;
		border: 1px solid #555;
		box-shadow: 0 4px 6px rgba(0,0,0,0.3);
		white-space: pre-line; /* Support multi-line labels */
	}

	/* Compute Nodes (Backend processing) */
	:global(.svelte-flow__node.node-compute) {
		background: #2b2b2b;
		border-color: #646cff; /* Go Blue-ish */
	}

	/* Visualizer Nodes (The output) */
	:global(.svelte-flow__node.node-viz) {
		background: #2e1a47; /* Dark Purple */
		border-color: #ff3e00; /* Svelte Orange */
		font-weight: bold;
	}

    /* Selected state */
    :global(.svelte-flow__node.selected) {
        border-color: #fff;
        box-shadow: 0 0 10px rgba(255, 255, 255, 0.2);
    }
</style>
