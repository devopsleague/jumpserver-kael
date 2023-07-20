<template>
  <svg 
    :class="svgClass" 
    aria-hidden="true"
    :style="{ width: size + 'px', height: size + 'px', color: color }"
  >
    <use class="svg-use" :href="symbolId" :xlink:href="symbolId" />
  </svg>
</template>
 
<script>
  import { defineComponent, computed } from 'vue'
 
  export default defineComponent({
    name: 'SvgIcon',
    props: {
      prefix: {
        type: String,
        default: 'icon'
      },
      name: {
        type: String,
        required: true
      },
      className: {
        type: String,
        default: ''
      },
      color: {
        type: String,
        default: "#C5C5D2"
      },
      size: {
        type: [String, Number],
        default: 16
      }
    },
    setup(props) {
      const symbolId = computed(() => `#${props.name}`)
      const svgClass = computed(() => {
        if (props.className) {
          return `svg-icon ${props.className}`
        }
        return 'svg-icon'
      })
      return { symbolId, svgClass }
    }
  })
</script>
<style scoped>
  .svg-icon {
    display: inline-block;
    vertical-align: middle !important;
    vertical-align: -0.1em;
    fill: currentColor;
    overflow: hidden;
  } 
</style>
