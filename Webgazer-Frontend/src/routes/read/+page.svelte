<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { get } from 'svelte/store';
  import { fetchStudyText, submitGazePoint, type Passage } from '$lib/api';
  import { WebGazerManager } from '$lib/components';
  import { ReadingPanel } from '$lib/components/reading';
  import { webgazerStore } from '$lib/stores/webgazer';

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

  // Gaze data collection
  let gazeCollectionInterval: ReturnType<typeof setInterval> | null = null;
  let sessionDbId: number | null = null;
  let gazeBuffer: Array<{ x: number; y: number; panel: string; phase: string; timestamp: number }> = [];
  const GAZE_COLLECTION_INTERVAL = 100; // Collect gaze every 100ms
  const GAZE_BATCH_SIZE = 10; // Submit in batches of 10 points

  // Gaze indicator (red dot) - set to false for production deployment
  const SHOW_GAZE_INDICATOR = true;
  let currentGaze: { x: number; y: number } | null = null;
  let hasGaze = false;
  let gazeUnsubscribe: (() => void) | null = null;

  // Fetch study text on mount
  onMount(async () => {
    // Get session ID from sessionStorage
    const sessionIdStr = sessionStorage.getItem('session_db_id');
    if (sessionIdStr) {
      sessionDbId = parseInt(sessionIdStr, 10);
    }

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

    // Start gaze collection
    startGazeCollection();

    // Subscribe to gaze store for indicator
    if (SHOW_GAZE_INDICATOR) {
      gazeUnsubscribe = webgazerStore.subscribe((state) => {
        currentGaze = state.currentGaze;
        hasGaze = state.hasGaze;
      });
    }
  });

  onDestroy(() => {
    // Stop gaze collection
    stopGazeCollection();
    // Submit any remaining buffered gaze points
    submitBufferedGazePoints();
    // Unsubscribe from gaze store
    if (gazeUnsubscribe) {
      gazeUnsubscribe();
    }
  });

  function loadPassage(index: number) {
    if (index >= passages.length) {
      // All passages completed, go to quiz
      // Submit any remaining gaze data before leaving
      submitBufferedGazePoints();
      saveAllPreferences();
      setTimeout(() => {
        goto('/quiz');
      }, 500);
      return;
    }
    
    // Submit gaze data from previous passage before loading new one
    submitBufferedGazePoints();
    
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

    // Auto-start reading when passage loads (this enables gaze collection)
    setTimeout(() => {
      start();
    }, 100);
  }

  function handleWebGazerInitialized(instance: any) {
    wgInstance = instance;
    // Hide video & overlays during reading
    instance.showVideo(false)
      .showFaceOverlay(false)
      .showFaceFeedbackBox(false)
      .showPredictionPoints(false);
  }

  // Determine which panel the gaze is on based on x coordinate
  function getPanelFromGaze(x: number): string {
    const screenWidth = window.innerWidth;
    const midpoint = screenWidth / 2;
    return x < midpoint ? 'A' : 'B';
  }

  // Get current reading phase based on which panel user is looking at
  function getCurrentPhase(panel: string): string {
    if (!started) return 'waiting';
    
    // Phase is determined by which panel the user is currently looking at
    // This gives us more accurate phase tracking based on actual gaze behavior
    if (panel === 'A') {
      return 'reading_A';
    } else if (panel === 'B') {
      return 'reading_B';
    }
    
    // Fallback: if we can't determine panel, use time-based logic
    if (!doneA) return 'reading_A';
    if (!doneB) return 'reading_B';
    return 'completed';
  }

  // Collect gaze data periodically
  function startGazeCollection() {
    if (gazeCollectionInterval) return; // Already started

    gazeCollectionInterval = setInterval(() => {
      // Check for session ID update (in case it was set after page load)
      if (!sessionDbId) {
        const sessionIdStr = sessionStorage.getItem('session_db_id');
        if (sessionIdStr) {
          sessionDbId = parseInt(sessionIdStr, 10);
          console.log('Gaze collection: Session ID found:', sessionDbId);
        } else {
          console.log('Gaze collection: No session ID yet, skipping...');
          return; // Still no session ID, skip collection
        }
      }

      if (!started) {
        // Don't log every interval, just occasionally
        if (Math.random() < 0.01) { // ~1% of the time
          console.log('Gaze collection: Waiting for reading to start...');
        }
        return; // Only collect when reading has started
      }

      const gazeState = get(webgazerStore);
      if (gazeState.currentGaze && gazeState.hasGaze) {
        const panel = getPanelFromGaze(gazeState.currentGaze.x);
        const phase = getCurrentPhase(panel);

        // Add to buffer
        gazeBuffer.push({
          x: gazeState.currentGaze.x,
          y: gazeState.currentGaze.y,
          panel: panel,
          phase: phase,
          timestamp: Date.now()
        });

        // Submit in batches
        if (gazeBuffer.length >= GAZE_BATCH_SIZE) {
          console.log(`Submitting batch of ${gazeBuffer.length} gaze points`);
          submitBufferedGazePoints();
        }
      }
    }, GAZE_COLLECTION_INTERVAL);
  }

  function stopGazeCollection() {
    if (gazeCollectionInterval) {
      clearInterval(gazeCollectionInterval);
      gazeCollectionInterval = null;
    }
  }

  async function submitBufferedGazePoints() {
    if (gazeBuffer.length === 0 || !sessionDbId) {
      if (gazeBuffer.length > 0 && !sessionDbId) {
        console.warn('Cannot submit gaze points: No session ID');
      }
      return;
    }

    console.log(`Submitting ${gazeBuffer.length} gaze points to session ${sessionDbId}`);

    // Submit all buffered points
    const promises = gazeBuffer.map(point =>
      submitGazePoint({
        session_id: sessionDbId!,
        x: point.x,
        y: point.y,
        panel: point.panel,
        phase: point.phase
      }).catch((error) => {
        console.error('Failed to submit gaze point:', error, point);
        return false;
      })
    );

    // Wait for all submissions (but don't block on failures)
    const results = await Promise.allSettled(promises);
    const successCount = results.filter(r => r.status === 'fulfilled' && r.value === true).length;
    console.log(`Successfully submitted ${successCount}/${gazeBuffer.length} gaze points`);

    // Clear buffer
    gazeBuffer = [];
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

  async function selectFontPreference(preference: 'A' | 'B') {
    if (!currentPassage) return;
    
    // Submit any remaining gaze data before moving to next passage
    await submitBufferedGazePoints();
    
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

<!-- Gaze indicator (red dot) - only shown when SHOW_GAZE_INDICATOR is true -->
{#if SHOW_GAZE_INDICATOR && hasGaze && currentGaze}
  <div
    class="fixed w-4 h-4 rounded-full ring-2 ring-white pointer-events-none z-50 transition-opacity duration-100"
    style={`left:${currentGaze.x - 8}px; top:${currentGaze.y - 8}px; background: rgba(239, 68, 68, 0.9); box-shadow: 0 0 8px rgba(239, 68, 68, 0.5);`}
    aria-label="Current gaze position"
    aria-hidden="true"
  ></div>
{/if}
