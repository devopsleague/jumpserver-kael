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
    dom?.scrollIntoView(scrollOption)
  }, 0)
}

export function getUrlParams(url = window.location.search) {
  if (url.indexOf('?') === -1) {
    return {}
  }
  const urlStr = url.split('?')[1]
  const obj = {}
  const paramsArr = urlStr.split('&')
  for (let i = 0, len = paramsArr.length; i < len; i++) {
    const arr = paramsArr[i].split('=')
    obj[arr[0]] = arr[1]
  }
  return obj
}

export function isMobile() {
  const flag = navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Windows Phone)/i)
  return flag;
}
