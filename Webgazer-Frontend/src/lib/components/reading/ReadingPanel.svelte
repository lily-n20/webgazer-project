<script lang="ts">
  export let label: string;
  export let fontType: 'serif' | 'sans';
  export let text: string;
  export let done: boolean = false;
  export let onComplete: () => void;
  export let disabled: boolean = false;

  function fontClass(kind: 'serif' | 'sans') {
    return kind === 'serif' ? 'font-serif' : 'font-sans';
  }
</script>

<div class="border rounded-xl p-4 shadow-sm">
  <div class="flex items-center justify-between mb-2">
    <h3 class="text-lg font-medium">{label} — {fontType === 'serif' ? 'Serif' : 'Sans'}</h3>
  </div>
  <div class={`prose max-w-none leading-7 ${fontClass(fontType)}`}>
    <p class="whitespace-pre-wrap">{text}</p>
  </div>
</div>

