<script lang="ts">
    import { router } from "@inertiajs/svelte";
    import { fly } from "svelte/transition";
    import Header from "../../components/Header.svelte";
    import DarkModeToggle from "../../components/DarkModeToggle.svelte";
    import { Toast } from "../../lib/utils/helpers";
    import { Upload, Lock, User, Mail } from "lucide-svelte";

    interface User {
        id: number;
        email: string;
        name: string;
        avatar: string;
        role: string;
        email_verified: boolean;
    }

    interface Props {
        user?: User;
        success?: string;
        error?: string;
    }

    let { user, success, error }: Props = $props();

    let profileForm = $state({
        name: "",
        email: "",
        avatar: "",
    });

    let passwordForm = $state({
        current_password: "",
        new_password: "",
        confirm_password: "",
    });

    let isProfileLoading = $state(false);
    let isPasswordLoading = $state(false);
    let previewUrl = $state<string | null>(null);
    let showPassword = $state(false);

    $effect(() => {
        if (user) {
            profileForm.name = user.name || "";
            profileForm.email = user.email || "";
            profileForm.avatar = user.avatar || "";
            previewUrl = user.avatar ? `/api/avatar/${user.id}?v=${Date.now()}` : null;
        }
    });

    function handleAvatarChange(event: Event) {
        const target = event.target as HTMLInputElement;
        const file = target.files?.[0];
        if (file) {
            const formData = new FormData();
            formData.append("file", file);
            isProfileLoading = true;
            fetch("/app/upload", {
                method: "POST",
                body: formData,
            })
                .then((response) => response.json())
                .then((data) => {
                    if (data.success && data.url) {
                        // Auto-save avatar URL to database via router
                        router.put("/app/profile", {
                            avatar: data.url,
                        }, {
                            onError: (error) => {
                                isProfileLoading = false;
                                Toast("Failed to save avatar: " + (error as any).message, "error");
                            },
                            onFinish: () => {
                                setTimeout(() => {
                                    isProfileLoading = false;
                                    // Reload page to get fresh data from server
                                    window.location.reload();
                                }, 500);
                            },
                        });
                    } else {
                        isProfileLoading = false;
                        Toast(data.error || "Failed to upload avatar", "error");
                    }
                })
                .catch((error) => {
                    isProfileLoading = false;
                    Toast("Failed to upload avatar", "error");
                    console.error("Upload error:", error);
                });
        }
    }

    function handleProfileSubmit(e: Event) {
        e.preventDefault();
        isProfileLoading = true;
        router.put("/app/profile", profileForm, {
            onFinish: () => {
                isProfileLoading = false;
            },
        });
    }

    function handlePasswordSubmit(e: Event) {
        e.preventDefault();

        if (passwordForm.new_password !== passwordForm.confirm_password) {
            Toast("Passwords don't match", "error");
            return;
        }

        if (!passwordForm.current_password || !passwordForm.new_password || !passwordForm.confirm_password) {
            Toast("Please fill all fields", "error");
            return;
        }

        if (passwordForm.new_password.length < 8) {
            Toast("Password must be at least 8 characters", "error");
            return;
        }

        isPasswordLoading = true;
        router.put("/app/profile/password", passwordForm, {
            onFinish: () => {
                isPasswordLoading = false;
                passwordForm.current_password = "";
                passwordForm.new_password = "";
                passwordForm.confirm_password = "";
            },
        });
    }
</script>

<Header group="profile" />

<!-- Main Content -->
<div class="relative min-h-screen bg-white dark:bg-slate-950">
    <!-- Desktop Sidebar Spacer -->
    <div
        class="hidden lg:block w-72 fixed inset-y-0 left-0 pointer-events-none"
    ></div>

    <!-- Background Effects -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none">
        <div
            class="absolute top-0 -left-4 w-96 h-96 bg-orange-500/10 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob"
        ></div>
        <div
            class="absolute top-0 -right-4 w-96 h-96 bg-purple-500/10 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-blob animation-delay-2000"
        ></div>
    </div>

    <!-- Page Header -->
    <div
        class="relative pt-8 pb-12 px-6 border-b border-slate-200 dark:border-slate-800/50 bg-white/50 dark:bg-slate-950/50 backdrop-blur-xl"
    >
        <div class="max-w-5xl mx-auto">
            <div
                class="flex items-center gap-2 text-sm text-slate-500 dark:text-slate-400 mb-4"
            >
                <a
                    href="/app"
                    class="hover:text-orange-400 transition-colors"
                >Dashboard</a>
                <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 5l7 7-7 7"
                    />
                </svg>
                <span class="text-slate-700 dark:text-slate-300">Settings</span>
            </div>
            <h1 class="text-3xl font-bold text-slate-900 dark:text-white mb-2">
                Account Settings
            </h1>
            <p class="text-slate-600 dark:text-slate-400">
                Manage your profile and security preferences
            </p>
        </div>
    </div>

    <!-- Content Area -->
    <div class="relative max-w-5xl mx-auto px-6 py-12">
        <!-- Flash Messages -->
        {#if success}
            <div
                class="mb-6 bg-green-500/10 border border-green-500/20 backdrop-blur-xl text-green-600 dark:text-green-400 rounded-2xl shadow-2xl p-4 flex items-center gap-3 animate-in slide-in-from-top-2 duration-300"
                in:fly={{ y: 20, duration: 300 }}
            >
                <div
                    class="w-8 h-8 rounded-full bg-green-500/20 flex items-center justify-center flex-shrink-0"
                >
                    <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M5 13l4 4L19 7"
                        />
                    </svg>
                </div>
                <p class="text-sm font-medium">{success}</p>
            </div>
        {/if}

        {#if error}
            <div
                class="mb-6 bg-red-500/10 border border-red-500/20 backdrop-blur-xl text-red-600 dark:text-red-400 rounded-2xl shadow-2xl p-4 flex items-center gap-3 animate-in slide-in-from-top-2 duration-300"
                in:fly={{ y: 20, duration: 300 }}
            >
                <div
                    class="w-8 h-8 rounded-full bg-red-500/20 flex items-center justify-center flex-shrink-0"
                >
                    <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                    </svg>
                </div>
                <p class="text-sm font-medium">{error}</p>
            </div>
        {/if}

        <!-- Profile Overview Card -->
        <div
            class="bg-gradient-to-br from-white to-slate-50 dark:from-slate-900 dark:to-slate-900/50 rounded-3xl border border-slate-200 dark:border-slate-800 p-8 mb-8"
            in:fly={{ y: 20, duration: 600 }}
        >
            <div
                class="flex flex-col sm:flex-row items-center sm:items-start gap-6"
            >
                <!-- Avatar -->
                <div class="relative group">
                    <div
                        class="w-28 h-28 rounded-2xl bg-gradient-to-br from-orange-500 to-purple-500 p-1 shadow-2xl shadow-orange-500/20"
                    >
                        <div
                            class="w-full h-full rounded-xl bg-white dark:bg-slate-900 overflow-hidden"
                        >
                            {#if previewUrl}
                                <img
                                    src={previewUrl}
                                    alt="Profile"
                                    class="w-full h-full object-cover"
                                />
                            {:else}
                                <div
                                    class="w-full h-full flex items-center justify-center"
                                >
                                    <span
                                        class="text-4xl font-bold text-orange-500"
                                    >{user?.name?.charAt(0)?.toUpperCase() || ""}</span>
                                </div>
                            {/if}
                        </div>
                    </div>
                    <label
                        class="absolute bottom-0 right-0 w-10 h-10 bg-orange-500 hover:bg-orange-600 text-white rounded-xl flex items-center justify-center cursor-pointer transition-all shadow-lg shadow-orange-500/30 group-hover:scale-110"
                    >
                        <Upload class="w-5 h-5" />
                        <input
                            type="file"
                            accept="image/*"
                            onchange={handleAvatarChange}
                            class="hidden"
                        />
                    </label>
                </div>

                <!-- User Info -->
                <div class="flex-1 text-center sm:text-left">
                    <h2
                        class="text-2xl font-bold text-slate-900 dark:text-white mb-1"
                    >
                        {user?.name || ""}
                    </h2>
                    <p class="text-slate-600 dark:text-slate-400 mb-4">
                        {user?.email || ""}
                    </p>
                    <div
                        class="flex flex-wrap justify-center sm:justify-start gap-2"
                    >
                        <span
                            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium bg-orange-500/10 text-orange-400 border border-orange-500/20"
                        >
                            <div
                                class="w-1.5 h-1.5 rounded-full bg-orange-500"
                            ></div>
                            Active Member
                        </span>
                        <span
                            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium bg-slate-200 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-300 dark:border-slate-700"
                        >
                            <svg
                                class="w-3 h-3"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                                />
                            </svg>
                            Verified
                        </span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Appearance Settings -->
        <div
            class="bg-white/50 dark:bg-slate-900/50 rounded-2xl border border-slate-200 dark:border-slate-800 p-6 mb-8"
            in:fly={{ y: 20, duration: 600, delay: 50 }}
        >
            <div class="flex items-center gap-3 mb-6">
                <div
                    class="w-10 h-10 rounded-xl bg-purple-500/10 flex items-center justify-center"
                >
                    <svg
                        class="w-5 h-5 text-purple-400"
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
                <div>
                    <h3
                        class="text-lg font-semibold text-slate-900 dark:text-white"
                    >
                        Appearance
                    </h3>
                    <p class="text-sm text-slate-600 dark:text-slate-500">
                        Customize how VeloStack looks on your device
                    </p>
                </div>
            </div>

            <div class="flex items-center justify-between">
                <div>
                    <p
                        class="text-sm font-medium text-slate-900 dark:text-white"
                    >
                        Dark Mode
                    </p>
                    <p class="text-xs text-slate-500 dark:text-slate-400 mt-1">
                        Switch between light and dark themes
                    </p>
                </div>
                <div class="scale-110 origin-left">
                    <DarkModeToggle />
                </div>
            </div>
        </div>

        <!-- Settings Grid -->
        <div class="grid md:grid-cols-2 gap-6">
            <!-- Personal Information -->
            <div
                class="bg-white/50 dark:bg-slate-900/50 rounded-2xl border border-slate-200 dark:border-slate-800 p-6"
                in:fly={{ y: 20, duration: 600, delay: 100 }}
            >
                <div class="flex items-center gap-3 mb-6">
                    <div
                        class="w-10 h-10 rounded-xl bg-blue-500/10 flex items-center justify-center"
                    >
                        <User class="w-5 h-5 text-blue-400" />
                    </div>
                    <div>
                        <h3
                            class="text-lg font-semibold text-slate-900 dark:text-white"
                        >
                            Personal Information
                        </h3>
                        <p class="text-sm text-slate-600 dark:text-slate-500">
                            Update your personal details
                        </p>
                    </div>
                </div>

                <form
                    onsubmit={handleProfileSubmit}
                    class="space-y-5"
                >
                    <div>
                        <label
                            for="name"
                            class="block text-sm font-medium text-slate-700 dark:text-slate-400 mb-2"
                        >Full Name</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <User class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={profileForm.name}
                                type="text"
                                id="name"
                                class="w-full pl-12 pr-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800/50 border border-slate-300 dark:border-slate-700 focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 text-slate-900 dark:text-white placeholder-slate-500 transition-all outline-none"
                                placeholder="Your full name"
                            />
                        </div>
                    </div>

                    <div>
                        <label
                            for="email"
                            class="block text-sm font-medium text-slate-700 dark:text-slate-400 mb-2"
                        >Email Address</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Mail class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={profileForm.email}
                                type="email"
                                id="email"
                                class="w-full pl-12 pr-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800/50 border border-slate-300 dark:border-slate-700 focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 text-slate-900 dark:text-white placeholder-slate-500 transition-all outline-none"
                                placeholder="you@example.com"
                            />
                        </div>
                    </div>

                    <div class="pt-4">
                        <button
                            type="submit"
                            disabled={isProfileLoading}
                            class="w-full px-6 py-3 rounded-xl bg-gradient-to-r from-orange-600 to-orange-500 hover:from-orange-500 hover:to-orange-400 text-white font-semibold transition-all shadow-lg shadow-orange-500/25 hover:shadow-orange-500/40 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                        >
                            {#if isProfileLoading}
                                <svg
                                    class="animate-spin h-5 w-5"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                >
                                    <circle
                                        class="opacity-25"
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        stroke-width="4"
                                    ></circle>
                                    <path
                                        class="opacity-75"
                                        fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                    ></path>
                                </svg>
                                Saving...
                            {:else}
                                <svg
                                    class="w-5 h-5"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M5 13l4 4L19 7"
                                    />
                                </svg>
                                Save Changes
                            {/if}
                        </button>
                    </div>
                </form>
            </div>

            <!-- Change Password -->
            <div
                class="bg-white/50 dark:bg-slate-900/50 rounded-2xl border border-slate-200 dark:border-slate-800 p-6"
                in:fly={{ y: 20, duration: 600, delay: 200 }}
            >
                <div class="flex items-center gap-3 mb-6">
                    <div
                        class="w-10 h-10 rounded-xl bg-red-500/10 flex items-center justify-center"
                    >
                        <Lock class="w-5 h-5 text-red-400" />
                    </div>
                    <div>
                        <h3
                            class="text-lg font-semibold text-slate-900 dark:text-white"
                        >
                            Security
                        </h3>
                        <p class="text-sm text-slate-600 dark:text-slate-500">
                            Update your password
                        </p>
                    </div>
                </div>

                <form onsubmit={handlePasswordSubmit} class="space-y-5">
                    <div>
                        <label
                            for="current_password"
                            class="block text-sm font-medium text-slate-700 dark:text-slate-400 mb-2"
                        >Current Password</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Lock class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={passwordForm.current_password}
                                type={showPassword ? "text" : "password"}
                                id="current_password"
                                class="w-full pl-12 pr-12 py-3 rounded-xl bg-slate-100 dark:bg-slate-800/50 border border-slate-300 dark:border-slate-700 focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 text-slate-900 dark:text-white placeholder-slate-500 transition-all outline-none"
                                placeholder="••••••••"
                            />
                        </div>
                    </div>

                    <div>
                        <label
                            for="new_password"
                            class="block text-sm font-medium text-slate-700 dark:text-slate-400 mb-2"
                        >New Password</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Lock class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={passwordForm.new_password}
                                type={showPassword ? "text" : "password"}
                                id="new_password"
                                class="w-full pl-12 pr-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800/50 border border-slate-300 dark:border-slate-700 focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 text-slate-900 dark:text-white placeholder-slate-500 transition-all outline-none"
                                placeholder="••••••••"
                                minlength="8"
                            />
                        </div>
                    </div>

                    <div>
                        <label
                            for="confirm_password"
                            class="block text-sm font-medium text-slate-700 dark:text-slate-400 mb-2"
                        >Confirm New Password</label>
                        <div class="relative">
                            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                <Lock class="w-5 h-5 text-slate-500" />
                            </div>
                            <input
                                bind:value={passwordForm.confirm_password}
                                type={showPassword ? "text" : "password"}
                                id="confirm_password"
                                class="w-full pl-12 pr-4 py-3 rounded-xl bg-slate-100 dark:bg-slate-800/50 border border-slate-300 dark:border-slate-700 focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 text-slate-900 dark:text-white placeholder-slate-500 transition-all outline-none"
                                placeholder="••••••••"
                            />
                        </div>
                    </div>

                    <div class="flex items-center gap-2">
                        <input
                            type="checkbox"
                            id="show_password"
                            bind:checked={showPassword}
                            class="w-4 h-4 rounded border-slate-300 text-orange-500 focus:ring-orange-500"
                        />
                        <label for="show_password" class="text-sm text-slate-600 dark:text-slate-400">
                            Show passwords
                        </label>
                    </div>

                    <div class="pt-4">
                        <button
                            type="submit"
                            disabled={isPasswordLoading}
                            class="w-full px-6 py-3 rounded-xl bg-gradient-to-r from-red-600 to-red-500 hover:from-red-500 hover:to-red-400 text-white font-semibold transition-all shadow-lg shadow-red-500/25 hover:shadow-red-500/40 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                        >
                            {#if isPasswordLoading}
                                <svg
                                    class="animate-spin h-5 w-5"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                >
                                    <circle
                                        class="opacity-25"
                                        cx="12"
                                        cy="12"
                                        r="10"
                                        stroke="currentColor"
                                        stroke-width="4"
                                    ></circle>
                                    <path
                                        class="opacity-75"
                                        fill="currentColor"
                                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                    ></path>
                                </svg>
                                Changing...
                            {:else}
                                <svg
                                    class="w-5 h-5"
                                    fill="none"
                                    stroke="currentColor"
                                    viewBox="0 0 24 24"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                                    />
                                </svg>
                                Change Password
                            {/if}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
