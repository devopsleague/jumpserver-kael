<template>
  <div>
    <div>
      <strong>Test</strong>: <span>{{ show_message }}</span>
    </div>
    <div>
      <input v-model="inputMessage" type="text"/>
      <button @click="sendMessage">Send</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      show_message: "",
      inputMessage: "",
      websocket: null,
    };
  },
  mounted() {
    this.websocket = new WebSocket("ws://localhost:8080/chat");
    this.websocket.onmessage = this.handleMessage;
  },
  methods: {
    handleMessage(event) {
      const data = JSON.parse(event.data);
      console.log(data.message.content);
      this.show_message = data.message.content; // 给show_message赋值

      // 其他逻辑
    },
    sendMessage() {
      const message = {
        content: this.inputMessage,
        sender: "user",
        new_conversation: true,
        model: 'gpt_3_5',
      };
      this.websocket.send(JSON.stringify(message));
      this.inputMessage = "";
    },
  },
};
</script>
