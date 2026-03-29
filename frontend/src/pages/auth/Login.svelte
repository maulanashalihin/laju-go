<script lang="ts">
    import { router } from "@inertiajs/svelte";
    import { Lock, Mail, ArrowRight, Eye, EyeOff } from "lucide-svelte";

    let form = $state({
        email: "",
        password: "",
    });

    let isLoading = $state(false);
    let showPassword = $state(false);

    interface Props {
        flash?: {
            error?: string;
        };
    }

    let { flash }: Props = $props();

    function submitForm(e: Event) {
        e.preventDefault();
        isLoading = true;
        router.post("/login", { email: form.email, password: form.password }, {
            onFinish: () => {
                setTimeout(() => {
                    isLoading = false;
                }, 500);
            },
        });
    }
</script>

<svelte:head>
    <title>Login - VeloStack</title>
</svelte:head>

<section class="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex">
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute -top-40 -right-40 w-80 h-80 bg-orange-500/20 rounded-full blur-3xl"></div>
        <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-purple-500/20 rounded-full blur-3xl"></div>
        <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-orange-500/10 rounded-full blur-3xl"></div>
    </div>

    <!-- Left Side - Branding (Desktop) -->
    <div class="hidden lg:flex lg:w-1/2 relative items-center justify-center p-12">
        <div class="relative z-10 max-w-lg">
            <div class="mb-8">
                <svg
                    width="80"
                    height="80"
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
            </div>
            <h1 class="text-4xl font-bold text-white mb-4">
                Welcome back
            </h1>
            <p class="text-slate-400 text-lg leading-relaxed">
                Sign in to continue building amazing applications with the high-performance Go + Svelte framework.
            </p>
            <div class="mt-12 grid grid-cols-3 gap-6">
                <div class="text-center">
                    <div class="text-3xl font-bold text-orange-400">11x</div>
                    <div class="text-sm text-slate-500 mt-1">Faster</div>
                </div>
                <div class="text-center">
                    <div class="text-3xl font-bold text-orange-400">1.5ms</div>
                    <div class="text-sm text-slate-500 mt-1">Latency</div>
                </div>
                <div class="text-center">
                    <div class="text-3xl font-bold text-orange-400">99.9%</div>
                    <div class="text-sm text-slate-500 mt-1">Uptime</div>
                </div>
            </div>
        </div>
    </div>

    <!-- Right Side - Login Form -->
    <div class="w-full lg:w-1/2 flex items-center justify-center p-6 lg:p-12">
        <div class="w-full max-w-md">
            <!-- Mobile Logo -->
            <div class="lg:hidden mb-8 flex justify-center">
                <svg
                    width="64"
                    height="64"
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
            </div>

            <div class="bg-slate-800/50 backdrop-blur-xl rounded-2xl border border-slate-700/50 p-8 shadow-2xl">
                <div class="text-center mb-8">
                    <h2 class="text-2xl font-bold text-white">Sign in</h2>
                    <p class="text-slate-400 mt-2">Enter your credentials to continue</p>
                </div>

                {#if flash?.error}
                    <div class="mb-6 p-4 rounded-xl bg-red-500/10 border border-red-500/20 flex items-start gap-3">
                        <svg class="w-5 h-5 text-red-400 shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <span class="text-red-400 text-sm">{flash.error}</span>
                    </div>
                {/if}

                <a href="/auth/google"
                   class="w-full flex items-center justify-center gap-3 px-4 py-3 rounded-xl border border-slate-600 bg-slate-700/50 text-white font-medium hover:bg-slate-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-orange-500/50 focus:ring-offset-2 focus:ring-offset-slate-800">
                    <svg class="h-5 w-5" viewBox="0 0 24 24">
                        <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                        <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                        <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                        <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
                    </svg>
                    Continue with Google
                </a>

                <div class="relative my-8">
                    <div class="absolute inset-0 flex items-center">
                        <div class="w-full border-t border-slate-700"></div>
                    </div>
                    <div class="relative flex justify-center">
                        <span class="px-4 text-sm text-slate-500 bg-slate-800/50">or continue with email</span>
                    </div>
                </div>

                <form class="space-y-5" onsubmit={submitForm}>
                    <div class="space-y-2">
                        <label for="email" class="block text-sm font-medium text-slate-300">Email</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Mail class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={form.email}
                                required
                                type="email"
                                name="email"
                                id="email"
                                class="w-full pl-12 pr-4 py-3 rounded-xl bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:border-orange-500 focus:ring-2 focus:ring-orange-500/20 transition-colors duration-200"
                                placeholder="you@example.com"
                            />
                        </div>
                    </div>

                    <div class="space-y-2">
                        <label for="password" class="block text-sm font-medium text-slate-300">Password</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Lock class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={form.password}
                                required
                                type={showPassword ? 'text' : 'password'}
                                name="password"
                                id="password"
                                placeholder="••••••••"
                                class="w-full pl-12 pr-12 py-3 rounded-xl bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:border-orange-500 focus:ring-2 focus:ring-orange-500/20 transition-colors duration-200"
                            />
                            <button
                                type="button"
                                onclick={() => showPassword = !showPassword}
                                class="absolute inset-y-0 right-0 pr-4 flex items-center text-slate-500 hover:text-slate-300 transition-colors"
                            >
                                {#if showPassword}
                                    <EyeOff class="w-5 h-5" />
                                {:else}
                                    <Eye class="w-5 h-5" />
                                {/if}
                            </button>
                        </div>
                    </div>

                    <div class="flex justify-end">
                        <a href="/forgot-password" class="text-sm text-orange-400 hover:text-orange-300 transition-colors">
                            Forgot password?
                        </a>
                    </div>

                    <button
                        type="submit"
                        disabled={isLoading}
                        class="w-full py-3 px-4 rounded-xl bg-gradient-to-r from-orange-500 to-orange-600 text-white font-semibold hover:from-orange-600 hover:to-orange-700 focus:outline-none focus:ring-2 focus:ring-orange-500/50 focus:ring-offset-2 focus:ring-offset-slate-800 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                        {#if isLoading}
                            <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Signing in...
                        {:else}
                            Sign in
                            <ArrowRight class="w-5 h-5" />
                        {/if}
                    </button>
                </form>

                <p class="mt-8 text-center text-sm text-slate-400">
                    Don't have an account?
                    <a href="/register" class="text-orange-400 hover:text-orange-300 font-medium transition-colors">
                        Create one
                    </a>
                </p>
            </div>
        </div>
    </div>
</section>
