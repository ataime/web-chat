<template>
  <div id="app">
    <h1>群聊</h1>
    <div id="chatbox">
      <div v-for="msg in messages" :key="msg.time">
        <strong>{{ msg.username }}</strong>: {{ msg.message }}
      </div>
    </div>
    <input v-model="message" @keyup.enter="sendMessage" placeholder="输入消息...">
  </div>
</template>

<script>
export default {
  data() {
    return {
      ws: null,
      message: '',
      messages: []
    };
  },

  methods: {
    connectWebSocket() {
      this.ws = new WebSocket('ws://localhost:8000/ws');
      this.ws.onmessage = (e) => {
        this.messages.push(JSON.parse(e.data));
      };
      this.ws.onclose = () => {
        console.log('WebSocket disconnected. Attempting to reconnect...');
        setTimeout(() => {
          this.connectWebSocket();
        }, 5000); // 5 秒后重连
      };
    },

    sendMessage() {
      const msg = { username: '用户', message: this.message };
      this.ws.send(JSON.stringify(msg));
      this.message = '';
    }
  },

  created() {
    this.connectWebSocket();
  }

};
</script>

<style>
#app {
  text-align: center;
}

#chatbox {
  height: 300px;
  overflow-y: scroll;
  border: 1px solid #ddd;
  margin: auto;
  padding: 10px;
}

input {
  width: 100%;
  padding: 10px;
}
</style>
