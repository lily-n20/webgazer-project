<script lang="ts">
  import { QUIZ } from '$lib/studyText';
  import { QuizQuestion } from '$lib/components/quiz';
  import { submitCompleteSession } from '$lib/api';

  let answers: Record<string, number> = {};
  let submitted = false;
  let submitting = false;
  let submitError: string | null = null;
  let currentPage = 0;
  const questionsPerPage = 5;

  $: totalPages = Math.ceil(QUIZ.length / questionsPerPage);
  $: currentQuestions = QUIZ.slice(
    currentPage * questionsPerPage,
    (currentPage + 1) * questionsPerPage
  );

  function handleAnswerChange(questionId: string, answerIndex: number) {
    answers[questionId] = answerIndex;
    answers = { ...answers }; // trigger reactivity
  }

  function nextPage() {
    if (currentPage < totalPages - 1) {
      currentPage++;
    }
  }

  function previousPage() {
    if (currentPage > 0) {
      currentPage--;
    }
  }

  async function submit() {
    if (submitting) return;
    
    submitting = true;
    submitError = null;

    try {
      const success = await submitCompleteSession(answers);
      if (success) {
        submitted = true;
      } else {
        submitError = 'Failed to submit responses. Please try again.';
        submitting = false;
      }
    } catch (error) {
      console.error('Error submitting quiz:', error);
      submitError = error instanceof Error ? error.message : 'An error occurred while submitting.';
      submitting = false;
    }
  }
</script>

{#if submitted}
  <!-- Completion Screen -->
  <div class="min-h-screen bg-white flex items-center justify-center px-4">
    <div class="max-w-2xl w-full text-center space-y-6">
      <div class="space-y-4">
        <h1 class="text-5xl font-light text-gray-900 tracking-tight">Thank You!</h1>
        <p class="text-xl text-gray-600">You have completed the study.</p>
        <p class="text-gray-500">Your responses have been recorded.</p>
      </div>
    </div>
  </div>
{:else}
  <!-- Quiz Form -->
  <div class="min-h-screen bg-white px-4 py-10">
    <div class="max-w-3xl mx-auto space-y-6">
      <div class="text-center space-y-2">
        <h1 class="text-4xl font-light text-gray-900 tracking-tight">Comprehension Quiz</h1>
        <p class="text-gray-500">Answer the questions about the passage you just read.</p>
      </div>

      <form class="space-y-6" on:submit|preventDefault={submit}>
        {#if submitError}
          <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {submitError}
          </div>
        {/if}

        {#each currentQuestions as q (q.id)}
          <QuizQuestion
            question={q}
            answer={answers[q.id]}
            onAnswerChange={handleAnswerChange}
          />
        {/each}

        <div class="flex items-center justify-between pt-4">
          <button
            type="button"
            class="px-6 py-2 rounded-lg bg-gray-200 text-gray-800 hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
            on:click={previousPage}
            disabled={currentPage === 0}
          >Previous</button>

          <span class="text-gray-600">
            Page {currentPage + 1} of {totalPages}
          </span>

          {#if currentPage === totalPages - 1}
            <button
              class="px-6 py-2 rounded-lg bg-gray-900 text-white hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed"
              type="submit"
              disabled={submitting}
            >
              {submitting ? 'Submitting...' : 'Submit'}
            </button>
          {:else}
            <button
              type="button"
              class="px-6 py-2 rounded-lg bg-gray-900 text-white hover:bg-gray-800"
              on:click={nextPage}
            >Next</button>
          {/if}
        </div>
      </form>
    </div>
  </div>
{/if}
