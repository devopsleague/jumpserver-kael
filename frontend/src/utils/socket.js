let ws= null; // 建立的连接
let lockReconnect= false; // 是否真正建立连接
let timeout= 10 * 1000; // 30秒一次心跳
let timeoutObj= null; // 心跳心跳倒计时
let serverTimeoutObj= null; // 心跳倒计时
let timeoutNum= null; // 断开 重连倒计时
let globalCallback = null; //监听服务端消息
let globalUri = null

// uri: 长链接地址
// callback: 服务端消息回调函数
export function createWebSocket(uri = globalUri, callback = globalCallback) {
  globalUri = uri
  globalCallback = callback
  ws = new WebSocket(uri)
  ws.onopen = ()=>{
    start()
  }
  ws.onmessage = onMessage
  ws.onerror = onError
  ws.onclose = onClose
  ws.onsend = onSend
}
 
// 发送消息
export function onSend(message){
  console.log(`发送消息`)
  if (typeof message !== 'string') {
    message = JSON.stringify(message)
  }
  ws.send(message)
}
 
// 接受服务端消息
export function onMessage(res){
  let msgData = res.data
  if (typeof msgData != 'object' && msgData != 'Connect success') {
    let data = msgData.replace(/\ufeff/g, '')
    let message = JSON.parse(data)
   // 服务端消息回掉
    globalCallback(message)
   // 重置心跳
    reset()
  }
}
 
// 连接失败
export function onError(){
  console.log('连接失败')
  reconnect()
}
 
// 连接关闭
export function onClose(){
  console.log('连接关闭')
  reconnect()
}

// 断开关闭
export function closeWs(){
  if(lockReconnect) {
    ws.close()
    ws = null
    lockReconnect = false
  }
}
 
// 发送心跳
export function start () {
  timeoutObj && clearTimeout(timeoutObj);
  serverTimeoutObj && clearTimeout(serverTimeoutObj);
  timeoutObj = setTimeout(function(){
    // 这里发送一个心跳，后端收到后，返回一个心跳消息
    if (ws.readyState == 1) {
      // 如果连接正常
      console.log('发送心跳')
      // ws.send('{key: heartbeat}')
    } else{
      // 否则重连
      reconnect()
    }
    serverTimeoutObj = setTimeout(function() {
      //超时关闭
      ws.close()
    }, timeout)
  }, timeout)
}

//重置心跳
export function reset(){
  clearTimeout(timeoutObj)
  clearTimeout(serverTimeoutObj)
  start()
}
 
// 重新连接
export function reconnect() {
  if(lockReconnect) {
    return
  }
  lockReconnect = true
  // 没连接上会一直重连，设置延迟避免请求过多
  timeoutNum && clearTimeout(timeoutNum)
  timeoutNum = setTimeout(function () {
    // 新连接
    createWebSocket()
    lockReconnect = false
  }, 10000)
}
