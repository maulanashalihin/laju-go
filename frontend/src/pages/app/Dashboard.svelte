<script lang="ts">
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
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="space-y-6">
        <div class="bg-white rounded-lg shadow p-6">
          <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
          <p class="mt-2 text-gray-600">Welcome back, {user?.name || 'User'}!</p>
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

        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-gray-900">Account Status</h3>
            <div class="mt-4 space-y-2">
              <div class="flex items-center justify-between">
                <span class="text-gray-600">Email</span>
                <span class="text-gray-900 font-medium">{user?.email}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600">Role</span>
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800 capitalize">
                  {user?.role}
                </span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600">Verified</span>
                {#if user?.email_verified}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                    Yes
                  </span>
                {:else}
                  <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                    No
                  </span>
                {/if}
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-gray-900">Quick Actions</h3>
            <div class="mt-4 space-y-3">
              <a href="/app/profile" class="block w-full text-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50">
                Edit Profile
              </a>
              <button class="block w-full text-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700">
                Upgrade Plan
              </button>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-semibold text-gray-900">System Info</h3>
            <div class="mt-4 space-y-2 text-sm text-gray-600">
              <p>Backend: Go Fiber</p>
              <p>Frontend: Svelte 5</p>
              <p>Database: SQLite</p>
            </div>
          </div>
        </div>
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
