import _ from 'lodash'

export const copy = _.throttle(function(value) {
  const inputDom = document.createElement('input')
  inputDom.id = 'createInputDom'
  inputDom.value = value
  document.body.appendChild(inputDom)
  inputDom.select()
  document?.execCommand('copy')
  document.body.removeChild(inputDom)
}, 1400)


export const pageScroll = (
  el,
  scrollOption = {
    behavior: 'smooth',
    block: 'end'
  }
) => {
  setTimeout(() => {
    const dom = document.getElementById(el)
    dom.scrollIntoView(scrollOption)
  }, 0)
}
