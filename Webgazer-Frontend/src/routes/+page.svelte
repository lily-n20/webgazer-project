<script lang="ts">
  import { goto } from '$app/navigation';
  let name = '';

  function start() {
    // not used in the flow yet, but persisted for later if needed
    localStorage.setItem('participant_name', name.trim());
    goto('/calibrate');
  }
</script>

<div class="min-h-screen bg-white flex items-center justify-center px-4">
  <div class="max-w-xl w-full space-y-8 text-center">
    <div class="space-y-4">
      <h1 class="text-5xl font-light text-gray-900 tracking-tight">Readability Study</h1>
    </div>

    <div class="space-y-3">
      <label for="name-input" class="block text-left text-sm text-gray-600">Name (optional for now)</label>
      <input
        id="name-input"
        class="w-full border rounded-lg px-3 py-2 focus:outline-none focus:ring focus:border-gray-900"
        bind:value={name}
        placeholder="Type your name"
      />
    </div>

    <div class="flex items-center justify-center gap-4">
      <button
        class="px-8 py-3 bg-gray-900 text-white rounded-lg font-medium hover:bg-gray-800 transition-colors shadow-sm"
        on:click={start}
      >
        Start Calibration
      </button>
    </div>
  </div>
</div>
