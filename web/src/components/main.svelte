<script lang="ts">
  import "../index.css";
  import cx from "classnames";
  import Logo from "@assets/logo.svg";
	import { AuthStore, type User } from "@lib/auth";
	import { onMount } from "svelte";
	import { makeUrl } from "@lib/url";
	import { apiFetch } from "@lib/requests";
  import type { Breadcrumbs } from "@lib/types"


  export let breadcrumbs: Array<Breadcrumbs> = [];

  let user: User | undefined;
  onMount(() => {
    user = AuthStore.user();
    if (!user) {
      AuthStore.refreshSession();
      return
    }
  })

  async function logout() {
    const url: string = makeUrl("/auth/logout");
    const rsp = await apiFetch(url, { method: "GET" });
    if (rsp.status === 200) {
      AuthStore.logout();
      window.location.replace("/");
      return;
    } else {
      const payload: EnvHub.Error = await rsp.json();
      console.error(payload.message);
    }
  }
</script>

<nav class={cx([
  "bg-dark-200 border-dark-100 h-16",
  "flex items-center justify-between border-b-2 p-2 px-4"
])}>
  <div class="breadcrumbs flex items-center gap-9">
    <ul class="text-sm">
      <li>
        <a href="/">
          <img src={Logo} alt="EnvHub" class="w-14">
        </a>
      </li>
      <li>
        <div>
          <img 
            src={user?.avatar_url ?? ""} 
            alt="" 
            role="presentation" 
            class="inline w-5 rounded-full" 
          />
          <a 
            class="link link-hover ml-1 text-sm"
            href="/dashboard"
          >{user?.login ?? ""}
          </a>
        </div>
      </li>
      {#if breadcrumbs}
        {#each breadcrumbs as {text, href}}
          <li>
            <a {href}>
              {text}
            </a>
          </li>
        {/each}
      {/if}
    </ul>
  </div>

  <div>
    <button on:click={logout} class="btn hover:bg-dark-100 border-none bg-transparent text-sm font-normal normal-case">
      <i class="fa-solid fa-arrow-right-from-bracket fa-sm mr-2" />
      Log out
    </button>
  </div>
</nav>

<main {...$$restProps}>
  <slot/>
</main>

<footer class="footer footer-center text-base-content h-12 p-4">
  <div>
    <a 
      href="https://github.com/hn275/envhub" 
      target="_blank"
      class="text-light/60 hover:text-main transition"
    >
      <span>
        <i class="fa-regular fa-star"></i>
      </span>&nbsp;
      Star us on Github
    </a>
  </div>
</footer>

<style lang="postcss">
  main {
    min-height: calc(100vh - 48px - 64px);
    @apply bg-dark-100;
  }
</style>
