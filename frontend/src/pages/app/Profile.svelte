<script lang="ts">
  import { router } from '@inertiajs/svelte';

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

  let name = $state('');
  let avatar = $state('');
  let isSubmitting = $state(false);

  $effect(() => {
    if (user) {
      name = user.name || '';
      avatar = user.avatar || '';
    }
  });

  async function handleSubmit(e: Event) {
    e.preventDefault();
    isSubmitting = true;

    try {
      router.put('/app/profile', { name, avatar });
    } catch (err) {
      alert('An error occurred');
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="min-h-screen bg-gray-50">
  <!-- Navigation -->
  <nav class="bg-white shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <a href="/app" class="text-xl font-bold text-indigo-600">VeloStack</a>
        </div>
        <div class="flex items-center space-x-4">
          <a href="/app" class="text-gray-600 hover:text-gray-900">Dashboard</a>
          <a href="/app/profile" class="text-gray-600 hover:text-gray-900">Profile</a>
          <form action="/logout" method="POST">
            <button type="submit" class="text-gray-600 hover:text-gray-900">Logout</button>
          </form>
        </div>
      </div>
    </div>
  </nav>

  <!-- Main Content -->
  <main class="py-8">
    <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="space-y-6">
        <div class="bg-white rounded-lg shadow p-6">
          <h1 class="text-2xl font-bold text-gray-900">Profile Settings</h1>
          <p class="mt-2 text-gray-600">Manage your account information</p>
        </div>

        {#if success}
          <div class="bg-green-50 border border-green-200 text-green-800 px-4 py-3 rounded-lg">
            {success}
          </div>
        {/if}

        {#if error}
          <div class="bg-red-50 border border-red-200 text-red-800 px-4 py-3 rounded-lg">
            {error}
          </div>
        {/if}

        <form onsubmit={handleSubmit} class="bg-white rounded-lg shadow p-6 space-y-6">
          <div>
            <label for="email" class="form-label">Email</label>
            <input
              type="email"
              id="email"
              value={user?.email}
              disabled
              class="form-input bg-gray-100 cursor-not-allowed"
            />
            <p class="mt-1 text-sm text-gray-500">Email cannot be changed</p>
          </div>

          <div>
            <label for="name" class="form-label">Full Name</label>
            <input
              type="text"
              id="name"
              bind:value={name}
              class="form-input"
              required
            />
          </div>

          <div>
            <label for="avatar" class="form-label">Avatar URL</label>
            <input
              type="url"
              id="avatar"
              bind:value={avatar}
              class="form-input"
              placeholder="https://example.com/avatar.jpg"
            />
          </div>

          <div class="flex items-center justify-between">
            <a href="/app" class="text-gray-600 hover:text-gray-900">
              ← Back to Dashboard
            </a>
            <button
              type="submit"
              disabled={isSubmitting}
              class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  </main>

  <!-- Footer -->
  <footer class="bg-white border-t mt-auto">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <p class="text-center text-gray-500 text-sm">
        &copy; 2026 VeloStack. Built with Go Fiber + Svelte 5.
      </p>
    </div>
  </footer>
</div>
