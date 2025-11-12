<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { SAMPLE_TEXT } from '$lib/studyText';
  import { WebGazerManager } from '$lib/components';
  import { ReadingPanel } from '$lib/components/reading';

  let wgInstance: any = null;
  let started = false;
  let doneA = false;
  let doneB = false;
  let t0 = 0;
  let timeA = 0;
  let timeB = 0;
  let fontPreference: 'A' | 'B' | null = null;

  // randomly assign which side is Serif vs Sans
  const fonts: { left: 'serif' | 'sans', right: 'serif' | 'sans' } = Math.random() < 0.5
    ? { left: 'serif', right: 'sans' }
    : { left: 'sans', right: 'serif' };

  function handleWebGazerInitialized(instance: any) {
    wgInstance = instance;
    // Hide video & overlays during reading
    instance.showVideo(false)
      .showFaceOverlay(false)
      .showFaceFeedbackBox(false)
      .showPredictionPoints(false);
  }

  function start() {
    if (started) return;
    started = true;
    t0 = performance.now();
  }

  function completeA() {
    if (!started || doneA) return;
    timeA = performance.now() - t0;
    doneA = true;
    // restart timer for B
    t0 = performance.now();
  }

  function completeB() {
    if (!started || doneB) return;
    timeB = performance.now() - t0;
    doneB = true;
  }

  function selectFontPreference(preference: 'A' | 'B') {
    fontPreference = preference;
    // Store preference
    sessionStorage.setItem('font_preference', preference);
    sessionStorage.setItem('font_preferred_type', preference === 'A' ? fonts.left : fonts.right);
    
    // Stash results and continue to quiz
    sessionStorage.setItem('font_left', fonts.left);
    sessionStorage.setItem('font_right', fonts.right);
    sessionStorage.setItem('time_left_ms', String(fonts.left === 'serif' ? timeA : timeB));
    sessionStorage.setItem('time_right_ms', String(fonts.right === 'serif' ? timeB : timeA));
    sessionStorage.setItem('timeA_ms', String(timeA));
    sessionStorage.setItem('timeB_ms', String(timeB));
    
    // Navigate to quiz after a brief delay
    setTimeout(() => {
      goto('/quiz');
    }, 500);
  }

</script>

<WebGazerManager
  showVideo={false}
  showFaceOverlay={false}
  showFaceFeedbackBox={false}
  showPredictionPoints={false}
  onInitialized={handleWebGazerInitialized}
/>

<div class="min-h-screen bg-white flex flex-col">
  <div class="flex-1 flex flex-col items-center justify-center px-8 py-10">
    <div class="flex-1 w-full flex flex-col items-center justify-center gap-8 px-8">
      <div class="text-center mb-6">
        <h1 class="text-4xl font-light text-gray-900 tracking-tight">Which font did you prefer?</h1>
      </div>
      
      <div class="w-full flex items-center justify-center gap-70">
        <div class="flex-1 max-w-xl flex flex-col items-center gap-8">
          <ReadingPanel
            label="Box A"
            fontType={fonts.left}
            text={SAMPLE_TEXT}
          />
          <button
            class="px-10 py-3 mt-5 rounded-lg border-2 border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 transition-colors disabled:opacity-50 disabled:cursor-not-allowed
                   {fontPreference === 'A' ? 'bg-gray-100 border-gray-500' : ''}"
            on:click={() => selectFontPreference('A')}
            disabled={fontPreference !== null}
          >
            Box A
          </button>
        </div>

        <div class="flex-1 max-w-xl flex flex-col items-center gap-8">
          <ReadingPanel
            label="Box B"
            fontType={fonts.right}
            text={SAMPLE_TEXT}
          />
          <button
            class="px-10 py-3 mt-5 rounded-lg border-2 border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 transition-colors disabled:opacity-50 disabled:cursor-not-allowed
                   {fontPreference === 'B' ? 'bg-gray-100 border-gray-500' : ''}"
            on:click={() => selectFontPreference('B')}
            disabled={fontPreference !== null}
          >
            Box B
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
