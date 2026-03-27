<script lang="ts">
    import { router } from "@inertiajs/svelte";
    import Input from "@components/Input.svelte";
    import Button from "@components/Button.svelte";

    let email = "";
    let password = "";
    let error = "";
    let loading = false;

    function handleSubmit() {
        error = "";
        loading = true;

        if (!email || !password) {
            error = "Email and password are required";
            loading = false;
            return;
        }

        router.post(
            "/login/login",
            { email, password },
            {
                onSuccess: (page) => {
                    loading = false;
                },
                onError: (errors) => {
                    loading = false;
                    error = errors.error || "Login failed";
                },
            },
        );
    }
</script>

<svelte:head>
    <title>Login - VeloStack</title>
</svelte:head>

<div
    class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
>
    <div class="max-w-md w-full space-y-8">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Sign in to your account
            </h2>
            <p class="mt-2 text-center text-sm text-gray-600">
                Or <a
                    href="/register"
                    class="font-medium text-indigo-600 hover:text-indigo-500"
                    >create a new account</a
                >
            </p>
        </div>

        {#if error}
            <div class="rounded-md bg-red-50 p-4">
                <div class="text-sm text-red-700">{error}</div>
            </div>
        {/if}

        <form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
            <div class="rounded-md shadow-sm -space-y-px">
                <div>
                    <Input
                        type="email"
                        id="email"
                        name="email"
                        placeholder="Email address"
                        bind:value={email}
                        required={true}
                    />
                </div>
                <div>
                    <Input
                        type="password"
                        id="password"
                        name="password"
                        placeholder="Password"
                        bind:value={password}
                        required={true}
                    />
                </div>
            </div>

            <div class="flex items-center justify-between">
                <div class="flex items-center">
                    <input
                        id="remember-me"
                        name="remember-me"
                        type="checkbox"
                        class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                    <label
                        for="remember-me"
                        class="ml-2 block text-sm text-gray-900"
                    >
                        Remember me
                    </label>
                </div>
            </div>

            <div>
                <Button type="submit" disabled={loading}>
                    {#if loading}Signing in...{:else}Sign in{/if}
                </Button>
            </div>

            <div class="relative">
                <div class="absolute inset-0 flex items-center">
                    <div class="w-full border-t border-gray-300"></div>
                </div>
                <div class="relative flex justify-center text-sm">
                    <span class="px-2 bg-gray-50 text-gray-500"
                        >Or continue with</span
                    >
                </div>
            </div>

            <div>
                <a
                    href="/auth/google"
                    class="w-full flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50"
                >
                    <svg class="w-5 h-5 mr-2" viewBox="0 0 24 24">
                        <path
                            fill="#4285F4"
                            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                        />
                        <path
                            fill="#34A853"
                            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                        />
                        <path
                            fill="#FBBC05"
                            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                        />
                        <path
                            fill="#EA4335"
                            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                        />
                    </svg>
                    Google
                </a>
            </div>
        </form>
    </div>
</div>
