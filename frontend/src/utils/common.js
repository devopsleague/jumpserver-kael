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
