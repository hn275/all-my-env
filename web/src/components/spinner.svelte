<svg id="spinner" viewBox="25 25 50 50">
	<circle id="spinner-circle" r="20" cy="50" cx="50" />
</svg>

<style>
	#spinner {
		width: 3em;
		transform-origin: center;
		animation: rotate4 2s linear infinite;
		transform: scale(0.8);
	}

	#spinner-circle {
		fill: none;
		stroke-width: 3px;
		stroke-dasharray: 1, 200;
		stroke-dashoffset: 0;
		stroke-linecap: round;
		animation: dash4 1.5s ease-in-out infinite;
		@apply stroke-dark-200;
	}

	@keyframes rotate4 {
		100% {
			transform: scale(0.8) rotate(360deg);
		}
	}

	@keyframes dash4 {
		0% {
			stroke-dasharray: 1, 200;
			stroke-dashoffset: 0;
		}

		50% {
			stroke-dasharray: 90, 200;
			stroke-dashoffset: -35px;
		}

		100% {
			stroke-dashoffset: -125px;
		}
	}
</style>
