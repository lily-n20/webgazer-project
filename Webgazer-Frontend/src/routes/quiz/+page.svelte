<script lang="ts">
  import { QUIZ } from '$lib/studyText';
  import { QuizQuestion } from '$lib/components/quiz';

  let answers: Record<string, number> = {};
  let submitted = false;

  function handleAnswerChange(questionId: string, answerIndex: number) {
    answers[questionId] = answerIndex;
    answers = { ...answers }; // trigger reactivity
  }

  function submit() {
    submitted = true;
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
        {#each QUIZ as q (q.id)}
          <QuizQuestion
            question={q}
            answer={answers[q.id]}
            onAnswerChange={handleAnswerChange}
          />
        {/each}

        <div class="flex items-center justify-center">
          <button
            class="px-6 py-2 rounded-lg bg-gray-900 text-white hover:bg-gray-800"
            type="submit"
          >Submit</button>
        </div>
      </form>
    </div>
  </div>
{/if}
