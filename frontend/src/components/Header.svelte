<script lang="ts">
    import { fly, fade } from "svelte/transition";
    import { page, router, inertia } from "@inertiajs/svelte";
    import DarkModeToggle from "./DarkModeToggle.svelte";

    // Lucide icons
import {
        LayoutDashboard,
        Settings,
        LogOut,
        Menu,
        X,
        User,
} from "lucide-svelte";

    interface User {
        id: number;
        email: string;
        name: string;
        avatar: string;
        role: string;
        email_verified: boolean;
    }

    let user = $state<User | undefined>(page.props.user);
    let isMenuOpen = $state(false);
    let isUserMenuOpen = $state(false);

    let { group = "" } = $props();

    const menuLinks = [
        {
            href: "/app",
            label: "Dashboard",
            group: "dashboard",
            show: true,
            icon: LayoutDashboard,
        },
        {
            href: "/app/profile",
            label: "Settings",
            group: "profile",
            show: true,
            icon: Settings,
        },
    ];

    function handleLogout() {
        router.post("/logout");
    }

    // Prevent body scroll when menu is open
    $effect(() => {
        if (typeof document !== "undefined") {
            document.body.style.overflow = isMenuOpen ? "hidden" : "unset";
        }
    });
</script>

<!-- Desktop Sidebar -->
<aside
    class="hidden lg:flex flex-col fixed left-0 top-0 h-full w-72 bg-white/95 dark:bg-neutral-950/95 backdrop-blur-xl border-r border-neutral-200/80 dark:border-white/[0.04] z-40 transition-all duration-300"
>
    <!-- Logo -->
    <a href="/app" use:inertia class="flex items-center gap-3 px-6 py-6 border-b border-neutral-200/80 dark:border-white/[0.04] hover:opacity-80 transition-opacity no-underline">
        <svg width="36" height="36" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
                <linearGradient id="grad1" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" style="stop-color:#14b8a6;stop-opacity:1" />
                    <stop offset="100%" style="stop-color:#22d3ee;stop-opacity:1" />
                </linearGradient>
            </defs>
            <path d="M30 10 H65 L55 50 H20 Z" fill="url(#grad1)" />
            <path d="M20 58 H85 L75 90 H10 Z" fill="url(#grad1)" />
            <rect x="70" y="58" width="20" height="32" transform="skewX(-14)" fill="white" fill-opacity="0.1" />
        </svg>
        <div>
            <h1 class="text-xl font-black italic text-neutral-900 dark:text-white">
                Laju<span class="text-brand-400">Go</span>
            </h1>
            {#if group === "dashboard"}
                <p class="text-xs text-neutral-500 dark:text-neutral-400">Dashboard</p>
            {:else if group === "profile"}
                <p class="text-xs text-neutral-500 dark:text-neutral-400">Settings</p>
            {:else}
                <p class="text-xs text-neutral-500 dark:text-neutral-400">App</p>
            {/if}
        </div>
    </a>

    <!-- Navigation -->
    <nav class="flex-1 px-4 py-6 space-y-2">
        {#each menuLinks.filter((item) => item.show) as item}
            {@const Icon = item.icon}
            <a
                href={item.href}
                use:inertia
                class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 group {item.group ===
                group
                    ? 'bg-brand-400/10 text-brand-600 dark:text-brand-400 border border-brand-400/20'
                    : 'text-neutral-600 dark:text-neutral-400 hover:text-neutral-900 dark:hover:text-white hover:bg-neutral-100 dark:hover:bg-neutral-800/50 border border-transparent'}"
            >
                <Icon
                    size="20"
                    class={item.group === group
                        ? "text-brand-400"
                        : "text-neutral-500 dark:text-neutral-400 group-hover:text-neutral-900 dark:group-hover:text-white"}
                    stroke-width="2"
                />
                {item.label}
                {#if item.group === group}
                    <div class="ml-auto w-1.5 h-1.5 rounded-full bg-brand-400"></div>
                {/if}
            </a>
        {/each}
    </nav>

    <!-- User Section -->
    {#if user && user.id}
        <div class="p-4 border-t border-neutral-200/80 dark:border-white/[0.04]">
            <div class="bg-neutral-100/50 dark:bg-neutral-800/50 rounded-xl p-4 border border-neutral-200/80 dark:border-white/[0.06]">
                <div class="flex items-center gap-3 mb-3">
                    <div class="w-10 h-10 rounded-full bg-linear-to-br from-brand-400 to-secondary-500 flex items-center justify-center text-neutral-950 font-bold text-sm ring-2 ring-neutral-300 dark:ring-neutral-700">
                        {user.name.charAt(0).toUpperCase()}
                    </div>
                    <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold text-neutral-900 dark:text-white truncate">{user.name}</p>
                        <p class="text-xs text-neutral-500 dark:text-neutral-400 truncate">{user.email || "Member"}</p>
                    </div>
                </div>
                <div class="flex items-center justify-between gap-2">
                    <DarkModeToggle />
                    <button
                        onclick={handleLogout}
                        class="inline-flex items-center gap-2 px-3 py-2 rounded-lg bg-red-500/10 hover:bg-red-500/20 text-red-500 dark:text-red-400 text-sm font-medium transition-colors"
                    >
                        <LogOut size="16" />
                        Logout
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <div class="p-4 border-t border-neutral-200/80 dark:border-white/[0.04] space-y-2">
            <a href="/login" use:inertia class="block w-full px-4 py-2.5 rounded-lg bg-neutral-200/80 dark:bg-neutral-800 hover:bg-neutral-300/80 dark:hover:bg-neutral-700 text-neutral-700 dark:text-neutral-300 text-sm font-medium transition-colors text-center">Sign In</a>
            <a href="/register" use:inertia class="block w-full px-4 py-2.5 rounded-lg bg-linear-to-r from-brand-500 to-brand-400 hover:from-brand-400 hover:to-brand-300 text-neutral-950 text-sm font-semibold transition-all text-center shadow-lg shadow-brand-500/25">Get Started</a>
        </div>
    {/if}
</aside>

<!-- Mobile Header -->
<header class="lg:hidden fixed top-0 left-0 right-0 z-50 bg-white dark:bg-neutral-950 backdrop-blur-xl border-b border-neutral-200/80 dark:border-white/[0.04]">
    <div class="flex items-center justify-between px-4 h-16">
        <!-- Logo -->
        <a href="/app" use:inertia class="flex items-center gap-2">
            <span class="text-lg font-black italic text-neutral-900 dark:text-white">Laju<span class="text-brand-400">Go</span></span>
        </a>

        <!-- Right Actions -->
        <div class="flex items-center gap-2">
            {#if user && user.id}
                <div class="relative" role="menu">
                    <button
                        onclick={() => (isUserMenuOpen = !isUserMenuOpen)}
                        class="w-9 h-9 rounded-full bg-linear-to-br from-brand-400 to-secondary-500 flex items-center justify-center text-neutral-950 font-bold text-sm ring-2 ring-neutral-300 dark:ring-neutral-700"
                    >
                        {user.name.charAt(0).toUpperCase()}
                    </button>

                    {#if isUserMenuOpen}
                        <div class="fixed inset-0 z-10" role="presentation" onclick={() => (isUserMenuOpen = false)}></div>
                        <div class="absolute right-0 mt-2 w-48 bg-white dark:bg-neutral-925 rounded-xl shadow-xl border border-neutral-200/80 dark:border-white/[0.06] overflow-hidden ring-1 ring-black/10 dark:ring-white/10" transition:fly={{ y: 10, duration: 200 }}>
                            <div class="p-3 border-b border-neutral-200/80 dark:border-white/[0.04]">
                                <p class="text-xs font-medium text-neutral-500 dark:text-neutral-400 uppercase">Signed in as</p>
                                <p class="text-sm font-semibold text-neutral-900 dark:text-white truncate">{user.name}</p>
                            </div>
                            <div class="p-2">
                                <a href="/app/profile" use:inertia class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-neutral-700 dark:text-neutral-300 hover:bg-neutral-100 dark:hover:bg-neutral-800 hover:text-neutral-900 dark:hover:text-white transition-colors">
                                    <User size="16" />
                                    Profile
                                </a>
                            </div>
                            <div class="p-2 border-t border-neutral-200/80 dark:border-white/[0.04]">
                                <button onclick={handleLogout} class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-red-500 dark:text-red-400 hover:bg-red-500/10 transition-colors">
                                    <LogOut size="16" />
                                    Logout
                                </button>
                            </div>
                        </div>
                    {/if}
                </div>
            {:else}
                <a href="/login" use:inertia class="px-4 py-2 rounded-lg bg-neutral-200/80 dark:bg-neutral-800 hover:bg-neutral-300/80 dark:hover:bg-neutral-700 text-neutral-700 dark:text-neutral-300 text-sm font-medium transition-colors">Sign In</a>
            {/if}

            <button
                onclick={() => (isMenuOpen = !isMenuOpen)}
                class="p-2 rounded-lg bg-neutral-200/80 dark:bg-neutral-800 text-neutral-600 dark:text-neutral-400 hover:text-neutral-900 dark:hover:text-white transition-colors"
                aria-label="Menu"
            >
                {#if isMenuOpen}
                    <X size="20" />
                {:else}
                    <Menu size="20" />
                {/if}
            </button>
        </div>
    </div>
</header>

<!-- Mobile Menu Drawer -->
{#if isMenuOpen}
    <div class="lg:hidden fixed inset-0 z-50">
        <button class="absolute inset-0 w-full h-full bg-neutral-900/50 backdrop-blur-sm" transition:fade={{ duration: 200 }} onclick={() => (isMenuOpen = false)} aria-label="Close menu"></button>

        <div class="absolute right-0 top-0 h-full w-[85%] max-w-[320px] bg-white dark:bg-neutral-925 shadow-2xl border-l border-neutral-200/80 dark:border-white/[0.04] flex flex-col" transition:fly={{ x: 300, duration: 400, opacity: 1 }}>
            <!-- Header -->
            <div class="flex items-center justify-between p-4 border-b border-neutral-200/80 dark:border-white/[0.04] bg-neutral-50 dark:bg-neutral-925/50">
                <span class="text-base font-bold text-neutral-900 dark:text-white">Menu</span>
                <button onclick={() => (isMenuOpen = false)} class="p-2 rounded-lg hover:bg-neutral-200/80 dark:hover:bg-neutral-800 text-neutral-600 dark:text-neutral-400 transition-colors">
                    <X size="20" />
                </button>
            </div>

            <!-- Navigation -->
            <div class="flex-1 overflow-y-auto p-4 space-y-2 bg-white dark:bg-neutral-925">
                {#each menuLinks.filter((item) => item.show) as item}
                    {@const Icon = item.icon}
                    <a
                        href={item.href}
                        use:inertia
                        class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all {item.group ===
                        group
                            ? 'bg-brand-400/10 text-brand-600 dark:text-brand-400 border border-brand-400/20'
                            : 'text-neutral-700 dark:text-neutral-400 hover:text-neutral-900 dark:hover:text-white hover:bg-neutral-100 dark:hover:bg-neutral-800/50 border border-transparent'}"
                    >
                        <Icon size="20" stroke-width="2" />
                        {item.label}
                    </a>
                {/each}
            </div>

            <!-- Footer -->
            {#if user}
                <div class="p-4 border-t border-neutral-200/80 dark:border-white/[0.04] bg-neutral-50 dark:bg-neutral-925/50">
                    <div class="bg-neutral-100/50 dark:bg-neutral-800/50 rounded-xl p-4 border border-neutral-200/80 dark:border-white/[0.06] mb-3">
                        <div class="flex items-center gap-3">
                            <div class="w-10 h-10 rounded-full bg-linear-to-br from-brand-400 to-secondary-500 flex items-center justify-center text-neutral-950 font-bold text-sm">
                                {user.name.charAt(0).toUpperCase()}
                            </div>
                            <div class="flex-1 min-w-0">
                                <p class="text-sm font-semibold text-neutral-900 dark:text-white truncate">{user.name}</p>
                                <p class="text-xs text-neutral-500 dark:text-neutral-400 truncate">{user.email || "Member"}</p>
                            </div>
                        </div>
                    </div>
                    <div class="flex items-center gap-2">
                        <DarkModeToggle />
                        <button onclick={handleLogout} class="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg bg-red-500/10 hover:bg-red-500/20 text-red-500 dark:text-red-400 font-medium transition-colors">
                            <LogOut size="18" />
                            Logout
                        </button>
                    </div>
                </div>
            {:else}
                <div class="p-4 border-t border-neutral-200/80 dark:border-white/[0.04] bg-neutral-50 dark:bg-neutral-925/50 space-y-2">
                    <a href="/login" use:inertia class="block w-full px-4 py-3 rounded-lg bg-neutral-200/80 dark:bg-neutral-800 hover:bg-neutral-300/80 dark:hover:bg-neutral-700 text-neutral-700 dark:text-neutral-300 font-medium transition-colors text-center">Sign In</a>
                    <a href="/register" use:inertia class="block w-full px-4 py-3 rounded-lg bg-linear-to-r from-brand-500 to-brand-400 hover:from-brand-400 hover:to-brand-300 text-neutral-950 font-semibold transition-all text-center shadow-lg shadow-brand-500/25">Get Started</a>
                </div>
            {/if}
        </div>
</div>
{/if}

<!-- Main Content Spacer for Mobile -->
<div class="lg:hidden h-16"></div>
