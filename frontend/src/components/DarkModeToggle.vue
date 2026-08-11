<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { Sun, Moon } from "lucide-vue-next";

const darkMode = ref(false);
const mounted = ref(false);

function applyDarkMode(isDark: boolean) {
  if (isDark) {
    document.documentElement.classList.add("dark");
  } else {
    document.documentElement.classList.remove("dark");
  }
}

function toggleDarkMode() {
  darkMode.value = !darkMode.value;
  applyDarkMode(darkMode.value);
  localStorage.setItem("darkMode", darkMode.value.toString());
}

let mediaQuery: MediaQueryList | null = null;

function handleChange(e: MediaQueryListEvent) {
  if (localStorage.getItem("darkMode") === null) {
    darkMode.value = e.matches;
    applyDarkMode(darkMode.value);
  }
}

// Initialize dark mode immediately (before component mounts)
(function initDarkMode() {
  const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const savedMode = localStorage.getItem("darkMode");
  darkMode.value = savedMode === null ? systemPrefersDark : savedMode === "true";

  // Apply saved preference immediately
  applyDarkMode(darkMode.value);

  // Mark as mounted after applying
  mounted.value = true;

  // Add transition class after initial load to prevent flash
  setTimeout(() => {
    document.documentElement.classList.add("transition-colors");
  }, 100);
})();

onMounted(() => {
  mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  mediaQuery.addEventListener("change", handleChange);
});

onUnmounted(() => {
  if (mediaQuery) {
    mediaQuery.removeEventListener("change", handleChange);
  }
});
</script>

<template>
  <button
    @click="toggleDarkMode"
    class="p-2 rounded-lg hover:bg-neutral-200/80 dark:hover:bg-neutral-800 focus:outline-none focus:ring-2 focus:ring-neutral-300 dark:focus:ring-neutral-700"
    aria-label="Toggle dark mode"
  >
    <template v-if="mounted">
      <Sun v-if="darkMode" class="w-5 h-5 text-neutral-800 dark:text-neutral-200" />
      <Moon v-else class="w-5 h-5 text-neutral-800 dark:text-neutral-200" />
    </template>
    <template v-else>
      <!-- Placeholder to prevent layout shift -->
      <div class="w-5 h-5"></div>
    </template>
  </button>
</template>
