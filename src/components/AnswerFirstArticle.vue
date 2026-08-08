<template>
  <article class="answer-first">
    <header v-if="title || question" class="answer-first__head">
      <div v-if="eyebrow" class="dasm-caps">{{ eyebrow }}</div>
      <h1 v-if="title" class="dasm-title">{{ title }}</h1>
      <h2 v-else-if="question" class="answer-first__question">{{ question }}</h2>
    </header>

    <section class="answer-first__answer dasm-panel" aria-label="Direct answer">
      <p class="answer-first__lead">{{ answer }}</p>
    </section>

    <slot />

    <section v-if="faq?.length" class="answer-first__faq" aria-label="FAQ">
      <h2 class="answer-first__h2">FAQ</h2>
      <div v-for="(item, i) in faq" :key="i" class="faq-item dasm-panel">
        <h3 class="faq-item__q">{{ item.question }}</h3>
        <p class="faq-item__a">{{ item.answer }}</p>
      </div>
    </section>
  </article>
</template>

<script setup>
defineProps({
  eyebrow: String,
  title: String,
  question: String,
  answer: { type: String, required: true },
  faq: { type: Array, default: () => [] }
})
</script>

<style scoped>
.answer-first {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.answer-first__question {
  margin: 0.35rem 0 0;
  font-size: clamp(1.2rem, 2vw, 1.7rem);
  color: #1d2b36;
}

.answer-first__lead {
  margin: 0;
  font-size: 1.05rem;
  line-height: 1.55;
  color: #243746;
  font-weight: 500;
}

.answer-first__h2 {
  margin: 0.5rem 0;
  font-size: 1.15rem;
  color: #1d2b36;
}

.faq-item__q {
  margin: 0 0 0.4rem;
  font-size: 1rem;
  color: #1f6f62;
}

.faq-item__a {
  margin: 0;
  color: #4a5d6d;
  line-height: 1.55;
}
</style>
