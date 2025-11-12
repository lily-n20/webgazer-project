<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { fetchStudyText, type Passage } from '$lib/api';
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
  let loading = true;
  
  // Passage management
  let passages: Passage[] = [];
  let currentPassageIndex = 0;
  let currentPassage: Passage | null = null;
  let fonts: { left: 'serif' | 'sans', right: 'serif' | 'sans' } = { left: 'serif', right: 'sans' };
  let defaultFonts: { left: 'serif' | 'sans', right: 'serif' | 'sans' } = { left: 'serif', right: 'sans' };
  
  // Store preferences for each passage
  let passagePreferences: Array<{
    passageId: number;
    preference: 'A' | 'B';
    fontType: string;
    timeA: number;
    timeB: number;
  }> = [];

  // Fetch study text on mount
  onMount(async () => {
    const textData = await fetchStudyText();
    if (textData) {
      sessionStorage.setItem('study_text_id', String(textData.id));
      
      // Store default fonts from study text
      if (textData.font_left && textData.font_right) {
        defaultFonts = {
          left: textData.font_left as 'serif' | 'sans',
          right: textData.font_right as 'serif' | 'sans'
        };
      }
      
      // Handle multiple passages
      if (textData.passages && textData.passages.length > 0) {
        passages = textData.passages.sort((a: Passage, b: Passage) => a.order - b.order);
        loadPassage(0);
      } else if (textData.content) {
        // Legacy: use single content field
        passages = [{
          id: 0,
          study_text_id: textData.id,
          order: 0,
          content: textData.content,
          font_left: textData.font_left,
          font_right: textData.font_right
        }];
        loadPassage(0);
      } else {
        loading = false;
        return;
      }
    } else {
      loading = false;
      return;
    }
    loading = false;
  });

  function loadPassage(index: number) {
    if (index >= passages.length) {
      // All passages completed, go to quiz
      saveAllPreferences();
      setTimeout(() => {
        goto('/quiz');
      }, 500);
      return;
    }
    
    currentPassageIndex = index;
    currentPassage = passages[index];
    
    // Reset state for new passage
    started = false;
    doneA = false;
    doneB = false;
    fontPreference = null;
    timeA = 0;
    timeB = 0;
    t0 = 0;
    
    // Set fonts for this passage (use passage fonts if available, otherwise use default)
    if (currentPassage.font_left && currentPassage.font_right) {
      fonts = {
        left: currentPassage.font_left as 'serif' | 'sans',
        right: currentPassage.font_right as 'serif' | 'sans'
      };
    } else {
      fonts = { ...defaultFonts };
    }
  }

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
    if (!currentPassage) return;
    
    fontPreference = preference;
    const preferredFontType = preference === 'A' ? fonts.left : fonts.right;
    
    // Store preference for this passage
    passagePreferences.push({
      passageId: currentPassage.id,
      preference: preference,
      fontType: preferredFontType,
      timeA: timeA,
      timeB: timeB
    });
    
    // Store in sessionStorage for this passage
    const passageKey = `passage_${currentPassage.id}`;
    sessionStorage.setItem(`${passageKey}_preference`, preference);
    sessionStorage.setItem(`${passageKey}_font_type`, preferredFontType);
    sessionStorage.setItem(`${passageKey}_timeA`, String(timeA));
    sessionStorage.setItem(`${passageKey}_timeB`, String(timeB));
    sessionStorage.setItem(`${passageKey}_font_left`, fonts.left);
    sessionStorage.setItem(`${passageKey}_font_right`, fonts.right);
    
    // Move to next passage after a brief delay
    setTimeout(() => {
      loadPassage(currentPassageIndex + 1);
    }, 500);
  }

  function saveAllPreferences() {
    // Save summary data to sessionStorage for final submission
    const allPreferences = passagePreferences.map((p, idx) => ({
      passage_index: idx,
      passage_id: p.passageId,
      preference: p.preference,
      font_type: p.fontType,
      timeA: p.timeA,
      timeB: p.timeB
    }));
    
    sessionStorage.setItem('all_passage_preferences', JSON.stringify(allPreferences));
    
    // Also save legacy format for backward compatibility (use last passage)
    if (passagePreferences.length > 0) {
      const last = passagePreferences[passagePreferences.length - 1];
      sessionStorage.setItem('font_preference', last.preference);
      sessionStorage.setItem('font_preferred_type', last.fontType);
      sessionStorage.setItem('timeA_ms', String(last.timeA));
      sessionStorage.setItem('timeB_ms', String(last.timeB));
    }
  }

  $: passageProgress = passages.length > 0 
    ? `${currentPassageIndex + 1} / ${passages.length}`
    : '';

</script>

<WebGazerManager
  showVideo={false}
  showFaceOverlay={false}
  showFaceFeedbackBox={false}
  showPredictionPoints={false}
  onInitialized={handleWebGazerInitialized}
/>

<div class="min-h-screen bg-white flex flex-col">
  {#if loading}
    <div class="flex-1 flex items-center justify-center">
      <p class="text-gray-500">Loading study text...</p>
    </div>
  {:else if currentPassage}
    <div class="flex-1 flex flex-col items-center justify-center px-8 py-10">
      <div class="flex-1 w-full flex flex-col items-center justify-center gap-8 px-8">
        <div class="text-center mb-6">
          <h1 class="text-4xl font-light text-gray-900 tracking-tight">Which font did you prefer?</h1>
          <!-- {#if passages.length > 1}
            <p class="text-lg text-gray-500 mt-2">Passage {passageProgress}</p>
          {/if} -->
        </div>
        
        <div class="w-full flex items-center justify-center gap-70">
        <div class="flex-1 max-w-xl flex flex-col items-center gap-8">
          <ReadingPanel
            label="Box A"
            fontType={fonts.left}
            text={currentPassage.content}
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
            text={currentPassage.content}
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
  {:else}
    <div class="flex-1 flex items-center justify-center">
      <p class="text-gray-500">No passages available. Please refresh the page.</p>
    </div>
  {/if}
</div>
