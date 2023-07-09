<script setup>
import { computed, onMounted, onUnmounted, onUpdated, ref, toRefs } from 'vue'
import MarkdownIt from 'markdown-it'
import mdKatex from '@traptitech/markdown-it-katex'
import mila from 'markdown-it-link-attributes'
import hljs from 'highlight.js'

const props = defineProps({
  message: Object,
  asRawText: Boolean,
  error: String
})

const { error } = toRefs(props)
const role = props.message?.role !== 'assistant'

const textRef = ref()

const mdi = new MarkdownIt({
  html: false,
  linkify: true,
  highlight(code, language) {
    const validLang = !!(language && hljs.getLanguage(language))
    if (validLang) {
      const lang = language ?? ''
      return highlightBlock(hljs.highlight(code, { language: lang }).value, lang)
    }
    return highlightBlock(hljs.highlightAuto(code).value, '')
  },
})

mdi.use(mila, { attrs: { target: '_blank', rel: 'noopener' } })
mdi.use(mdKatex, { blockClass: 'katexmath-block rounded-md p-[10px]', errorColor: ' #cc0000' })

const wrapClass = computed(() => {
  return [
    'text-wrap',
    'min-w-[20px]',
    'rounded-md',
    role ? 'bg-[#d2f9d1]' : 'bg-[#f4f6f8]',
    role ? 'dark:bg-[#a1dc95]' : 'dark:bg-[#1e1e20]',
    role ? 'message-request' : 'message-reply',
    error.value ? 'text-red-500' : ''
  ]
})

const text = computed(() => {
  const value = props.message?.content ?? ''
  console.log('value: ', value);
  if (props.message?.content) {
    return mdi.render(value)
  }
  return value
})

function highlightBlock(str, lang) {
  return `<pre class="code-block-wrapper"><div class="code-block-header"><span class="code-block-header__lang">${lang}</span><span class="code-block-header__copy">${'复制1'}</span></div><code class="hljs code-block-body ${lang}">${str}</code></pre>`
}

function addCopyEvents() {
  if (textRef.value) {
    const copyBtn = textRef.value.querySelectorAll('.code-block-header__copy')
    copyBtn.forEach((btn) => {
      btn.addEventListener('click', () => {
        const code = btn.parentElement?.nextElementSibling?.textContent
        if (code) {
        }
      })
    })
  }
}

function removeCopyEvents() {
  if (textRef.value) {
    const copyBtn = textRef.value.querySelectorAll('.code-block-header__copy')
    copyBtn.forEach((btn) => {
      btn.removeEventListener('click', () => {})
    })
  }
}

onMounted(() => {
  addCopyEvents()
})

onUpdated(() => {
  addCopyEvents()
})

onUnmounted(() => {
  removeCopyEvents()
})
</script>

<template>
  <div :class="wrapClass">
    <div ref="textRef" class="leading-relaxed break-words">
      <div class="markdown-body" v-html="text" />
    </div>
  </div>
</template>

<style lang="scss">
.message {
  & > div {
    display: inline-block;
    padding: 6px 10px;
    border-radius: 6px;
  }
}
</style>
