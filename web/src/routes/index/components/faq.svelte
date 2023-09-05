<script lang="ts">
</script>

<section class="join join-vertical">
	<div class="faq">
		<input type="checkbox" />
		<h5 class="collapse-title">What is the purpose of EnvHub?</h5>
		<div class="collapse-content">
			<p>
				EnvHub offers an innovative solution for the efficient
				management of your critical environment variables. Our platform
				seamlessly integrates with GitHub, simplifying the task of
				sharing variables among collaborators.
			</p>
		</div>
	</div>

	<div class="faq">
		<input type="checkbox" />
		<h5 class="collapse-title">How does it work?</h5>
		<div class="collapse-content">
			<p>
				EnvHub functionalities can be summarized in <span
					class="font-bold">3 steps</span
				>:
			</p>
			<ol class="mt-2 list-decimal">
				<li>Authentication with GitHub.</li>
				<li>Create projects and their variables.</li>
				<li>Share it with your collaborators!</li>
			</ol>
		</div>
	</div>

	<div class="faq">
		<input type="checkbox" />
		<h5 class="collapse-title">Is my data secure with this app?</h5>
		<div class="collapse-content">
			<p>
				We exercise meticulouse care when handling your environment
				variables. Our server is hosted with <a
					href="https://www.cloudflare.com/learning/ssl/what-is-ssl/"
					target="_blank"
					class="link link-info">SSL Certificates</a
				>
				and your variables undergo encryption and decryption using the robust
				<a
					href="https://ieeexplore.ieee.org/document/7927078"
					target="_blank"
					class="link link-info">ChaCha20-Poly1305</a
				>. This choice reflects our commitment to maintaining a
				high-security standard while optimizing performance.
				<br class="mb-2" />
				In case you meant your user data, yes, we believe telemetry has its
				usage... elsewhere.
			</p>
		</div>
	</div>

	<div class="faq">
		<input type="checkbox" />
		<h5 class="collapse-title">Who built this awesome app??</h5>
		<div class="collapse-content">
			<p>
				EnvHub was initiated by computer science students from the
				<a
					href="https://uvic.ca"
					target="_blank"
					class="link link-info">University of Victoria</a
				>. Led by
				<a
					href="https://github.com/hn275"
					target="_blank"
					class="link link-info">Hal</a
				>, a proactive student balancing academics with professional
				engagements, the project was born from the need for simplified
				variable sharing within teams, both in professional and
				collaborative settings.
			</p>
		</div>
	</div>

	<div class="faq">
		<input type="checkbox" />
		<h5 class="collapse-title">Why "EnvHub"?</h5>
		<div class="collapse-content">
			<p>
				This question is best answered by the talented <a
					href="https://sidhu.dev/"
					target="_blank"
					class="link link-info">Jaspreet</a
				> (who also designed this UI!)
			</p>
		</div>
	</div>
</section>

<style lang="postcss">
	h5 {
		@apply font-bold text-lg;
	}

	li {
		@apply ml-5;
	}

	.faq {
		@apply join-item border border-neutral collapse collapse-arrow;
	}
</style>
