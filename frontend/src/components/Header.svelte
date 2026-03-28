<script lang="ts">
    import { fly, fade } from "svelte/transition";
    import { page, router } from "@inertiajs/svelte";
    import { clickOutside } from "../lib/utils/helpers";
    import DarkModeToggle from "./DarkModeToggle.svelte";

    // Lucide icons
    import {
        LayoutDashboard,
        Settings,
        LogOut,
        Menu,
        X,
        User,
        Home,
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
    class="hidden lg:flex flex-col fixed left-0 top-0 h-full w-72 bg-white/95 dark:bg-slate-900/95 backdrop-blur-xl border-r border-slate-200 dark:border-slate-800 z-40 transition-all duration-300"
>
    <!-- Logo -->
    <div
        class="flex items-center gap-3 px-6 py-6 border-b border-slate-200 dark:border-slate-800/50"
    >
        <div class="flex items-center gap-2">
            <svg
                width="36"
                height="36"
                viewBox="0 0 100 100"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
            >
                <defs>
                    <linearGradient
                        id="grad1"
                        x1="0%"
                        y1="0%"
                        x2="100%"
                        y2="0%"
                    >
                        <stop
                            offset="0%"
                            style="stop-color:#f97316;stop-opacity:1"
                        />
                        <stop
                            offset="100%"
                            style="stop-color:#ea580c;stop-opacity:1"
                        />
                    </linearGradient>
                </defs>
                <path d="M30 10 H65 L55 50 H20 Z" fill="url(#grad1)" />
                <path d="M20 58 H85 L75 90 H10 Z" fill="url(#grad1)" />
                <rect
                    x="70"
                    y="58"
                    width="20"
                    height="32"
                    transform="skewX(-14)"
                    fill="white"
                    fill-opacity="0.1"
                />
            </svg>
            <div>
                <h1
                    class="text-xl font-black italic text-slate-900 dark:text-white"
                >
                    Laju<span class="text-orange-500">.dev</span>
                </h1>
                <p class="text-xs text-slate-500 dark:text-slate-400">
                    Dashboard
                </p>
            </div>
        </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-4 py-6 space-y-2">
        {#each menuLinks.filter((item) => item.show) as item}
            {@const Icon = item.icon}
            <a
                href={item.href}
                class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 group {item.group ===
                group
                    ? 'bg-orange-500/10 text-orange-400 border border-orange-500/20'
                    : 'text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-800/50 border border-transparent'}"
            >
                <Icon
                    size="20"
                    class={item.group === group
                        ? "text-orange-400"
                        : "text-slate-500 dark:text-slate-400 group-hover:text-slate-900 dark:group-hover:text-white"}
                    stroke-width="2"
                />
                {item.label}
                {#if item.group === group}
                    <div
                        class="ml-auto w-1.5 h-1.5 rounded-full bg-orange-500"
                    ></div>
                {/if}
            </a>
        {/each}
    </nav>

    <!-- Dark Mode Toggle -->
    <div class="px-4 py-4 border-t border-slate-200 dark:border-slate-800/50">
        <div
            class="flex items-center justify-between px-4 py-3 rounded-xl bg-slate-100/50 dark:bg-slate-800/50 border border-slate-200 dark:border-slate-700/50"
        >
            <div class="flex items-center gap-3">
                <div
                    class="w-8 h-8 rounded-lg bg-purple-500/10 flex items-center justify-center"
                >
                    <svg
                        class="w-4 h-4 text-purple-400"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 3v2.25m6.364.379l-1.591 1.591M20.25 10.5H18M7.757 14.743l-1.59 1.59M6 10.5H3.75m4.007-4.243l-1.59-1.59M15.75 18.75a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                        />
                    </svg>
                </div>
                <span
                    class="text-sm font-medium text-slate-700 dark:text-slate-300"
                    >Dark Mode</span
                >
            </div>
            <DarkModeToggle />
        </div>
    </div>

    <!-- User Section -->
    {#if user && user.id}
        <div class="p-4 border-t border-slate-200 dark:border-slate-800/50">
            <div
                class="bg-slate-100/50 dark:bg-slate-800/50 rounded-xl p-4 border border-slate-200 dark:border-slate-700/50"
            >
                <div class="flex items-center gap-3 mb-3">
                    <div
                        class="w-10 h-10 rounded-full bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center text-white font-bold text-sm ring-2 ring-slate-300 dark:ring-slate-700"
                    >
                        {user.name.charAt(0).toUpperCase()}
                    </div>
                    <div class="flex-1 min-w-0">
                        <p
                            class="text-sm font-semibold text-slate-900 dark:text-white truncate"
                        >
                            {user.name}
                        </p>
                        <p
                            class="text-xs text-slate-500 dark:text-slate-400 truncate"
                        >
                            {user.email || "Member"}
                        </p>
                    </div>
                </div>
                <button
                    onclick={handleLogout}
                    class="w-full flex items-center justify-center gap-2 px-3 py-2 rounded-lg bg-red-500/10 hover:bg-red-500/20 text-red-500 dark:text-red-400 text-sm font-medium transition-colors"
                >
                    <LogOut size="16" />
                    Logout
                </button>
            </div>
        </div>
    {:else}
        <div
            class="p-4 border-t border-slate-200 dark:border-slate-800/50 space-y-2"
        >
            <a
                href="/login"
                class="block w-full px-4 py-2.5 rounded-lg bg-slate-200 dark:bg-slate-800 hover:bg-slate-300 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300 text-sm font-medium transition-colors text-center"
                >Sign In</a
            >
            <a
                href="/login/register"
                class="block w-full px-4 py-2.5 rounded-lg bg-gradient-to-r from-orange-600 to-orange-500 hover:from-orange-500 hover:to-orange-400 text-white text-sm font-medium transition-all text-center shadow-lg shadow-orange-500/25"
                >Get Started</a
            >
        </div>
    {/if}
</aside>

<!-- Mobile Header -->
<header
    class="lg:hidden fixed top-0 left-0 right-0 z-50 bg-white dark:bg-slate-950 backdrop-blur-xl border-b border-slate-200 dark:border-slate-800"
>
    <div class="flex items-center justify-between px-4 h-16">
        <!-- Logo -->
        <a href="/app" class="flex items-center gap-2">
            <span
                class="text-lg font-black italic text-slate-900 dark:text-white"
                >Laju<span class="text-orange-500">.dev</span></span
            >
        </a>

        <!-- Right Actions -->
        <div class="flex items-center gap-2">
            {#if user && user.id}
                <div
                    class="relative"
                    use:clickOutside
                    onclick_outside={() => (isUserMenuOpen = false)}
                >
                    <button
                        onclick={() => (isUserMenuOpen = !isUserMenuOpen)}
                        class="w-9 h-9 rounded-full bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center text-white font-bold text-sm ring-2 ring-slate-300 dark:ring-slate-700"
                    >
                        {user.name.charAt(0).toUpperCase()}
                    </button>

                    {#if isUserMenuOpen}
                        <div
                            class="absolute right-0 mt-2 w-48 bg-white dark:bg-slate-900 rounded-xl shadow-xl border border-slate-200 dark:border-slate-800 overflow-hidden ring-1 ring-slate-900/10 dark:ring-white/10"
                            transition:fly={{ y: 10, duration: 200 }}
                        >
                            <div
                                class="p-3 border-b border-slate-200 dark:border-slate-800"
                            >
                                <p
                                    class="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase"
                                >
                                    Signed in as
                                </p>
                                <p
                                    class="text-sm font-semibold text-slate-900 dark:text-white truncate"
                                >
                                    {user.name}
                                </p>
                            </div>
                            <div class="p-2">
                                <a
                                    href="/app/profile"
                                    class="flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-white transition-colors"
                                >
                                    <User size="16" />
                                    Profile
                                </a>
                            </div>
                            <div
                                class="p-2 border-t border-slate-200 dark:border-slate-800"
                            >
                                <button
                                    onclick={handleLogout}
                                    class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-sm text-red-500 dark:text-red-400 hover:bg-red-500/10 transition-colors"
                                >
                                    <LogOut size="16" />
                                    Logout
                                </button>
                            </div>
                        </div>
                    {/if}
                </div>
            {:else}
                <a
                    href="/login"
                    class="px-4 py-2 rounded-lg bg-slate-200 dark:bg-slate-800 hover:bg-slate-300 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300 text-sm font-medium transition-colors"
                    >Sign In</a
                >
            {/if}

            <button
                onclick={() => (isMenuOpen = !isMenuOpen)}
                class="p-2 rounded-lg bg-slate-200 dark:bg-slate-800 text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors"
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
        <button
            class="absolute inset-0 w-full h-full bg-slate-900/50 backdrop-blur-sm"
            transition:fade={{ duration: 200 }}
            onclick={() => (isMenuOpen = false)}
            aria-label="Close menu"
        ></button>

        <div
            class="absolute right-0 top-0 h-full w-[85%] max-w-[320px] bg-white dark:bg-slate-900 shadow-2xl border-l border-slate-200 dark:border-slate-800 flex flex-col"
            transition:fly={{ x: 300, duration: 400, opacity: 1 }}
        >
            <!-- Header -->
            <div
                class="flex items-center justify-between p-4 border-b border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50"
            >
                <span class="text-base font-bold text-slate-900 dark:text-white"
                    >Menu</span
                >
                <button
                    onclick={() => (isMenuOpen = false)}
                    class="p-2 rounded-lg hover:bg-slate-200 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400 transition-colors"
                >
                    <X size="20" />
                </button>
            </div>

            <!-- Navigation -->
            <div
                class="flex-1 overflow-y-auto p-4 space-y-2 bg-white dark:bg-slate-900"
            >
                {#each menuLinks.filter((item) => item.show) as item}
                    {@const Icon = item.icon}
                    <a
                        href={item.href}
                        class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all {item.group ===
                        group
                            ? 'bg-orange-500/10 text-orange-400 border border-orange-500/20'
                            : 'text-slate-700 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-800/50 border border-transparent'}"
                    >
                        <Icon size="20" stroke-width="2" />
                        {item.label}
                    </a>
                {/each}
            </div>

            <!-- Dark Mode Toggle (Mobile) -->
            <div
                class="p-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50"
            >
                <div
                    class="flex items-center justify-between px-4 py-3 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700/50"
                >
                    <div class="flex items-center gap-3">
                        <div
                            class="w-8 h-8 rounded-lg bg-purple-500/10 flex items-center justify-center"
                        >
                            <svg
                                class="w-4 h-4 text-purple-400"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M12 3v2.25m6.364.379l-1.591 1.591M20.25 10.5H18M7.757 14.743l-1.59 1.59M6 10.5H3.75m4.007-4.243l-1.59-1.59M15.75 18.75a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                                />
                            </svg>
                        </div>
                        <span
                            class="text-sm font-medium text-slate-700 dark:text-slate-300"
                            >Dark Mode</span
                        >
                    </div>
                    <DarkModeToggle />
                </div>
            </div>

            <!-- Footer -->
            {#if user}
                <div
                    class="p-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50"
                >
                    <div
                        class="bg-slate-100/50 dark:bg-slate-800/50 rounded-xl p-4 border border-slate-200 dark:border-slate-700/50 mb-3"
                    >
                        <div class="flex items-center gap-3">
                            <div
                                class="w-10 h-10 rounded-full bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center text-white font-bold text-sm"
                            >
                                {user.name.charAt(0).toUpperCase()}
                            </div>
                            <div class="flex-1 min-w-0">
                                <p
                                    class="text-sm font-semibold text-slate-900 dark:text-white truncate"
                                >
                                    {user.name}
                                </p>
                                <p
                                    class="text-xs text-slate-500 dark:text-slate-400 truncate"
                                >
                                    {user.email || "Member"}
                                </p>
                            </div>
                        </div>
                    </div>
                    <button
                        onclick={handleLogout}
                        class="w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-red-500/10 hover:bg-red-500/20 text-red-500 dark:text-red-400 font-medium transition-colors"
                    >
                        <LogOut size="18" />
                        Logout
                    </button>
                </div>
            {:else}
                <div
                    class="p-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50 space-y-2"
                >
                    <a
                        href="/login"
                        class="block w-full px-4 py-3 rounded-lg bg-slate-200 dark:bg-slate-800 hover:bg-slate-300 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300 font-medium transition-colors text-center"
                        >Sign In</a
                    >
                    <a
                        href="/login/register"
                        class="block w-full px-4 py-3 rounded-lg bg-gradient-to-r from-orange-600 to-orange-500 hover:from-orange-500 hover:to-orange-400 text-white font-medium transition-all text-center shadow-lg shadow-orange-500/25"
                        >Get Started</a
                    >
                </div>
            {/if}
        </div>
    </div>
{/if}

<!-- Main Content Spacer for Mobile -->
<div class="lg:hidden h-16"></div>
