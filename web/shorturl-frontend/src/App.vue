<template>
  <div class="container">
    <h1 class="fade-in">🔗 短链接生成器</h1>
    <form @submit.prevent="generateShortURL" class="form fade-in">
      <div class="form-group">
        <label>原始链接:</label>
        <input v-model="url" type="text" required placeholder="请输入原始链接..." />
      </div>
      <div class="form-group">
        <label>过期时间:</label>
        <select v-model="expiration">
          <option value="30m">30 分钟</option>
          <option value="1h">1 小时</option>
          <option value="1d">1 天</option>
        </select>
      </div>
      <button type="submit" :disabled="loading">
        {{ loading ? '生成中...' : '✨ 生成短链接' }}
      </button>
    </form>

    <div v-if="shortUrl" class="result fade-in">
      <p>生成成功：<a :href="shortUrl" target="_blank">{{ shortUrl }}</a></p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const url = ref('')
const expiration = ref('1h')
const shortUrl = ref('')
const loading = ref(false)

const generateShortURL = async () => {
  loading.value = true
  try {
    const response = await axios.post('http://localhost:8080/generate', {
      url: url.value,
      expiration: expiration.value,
    }, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      transformRequest: [(data) => {
        return Object.entries(data).map(([key, val]) => `${encodeURIComponent(key)}=${encodeURIComponent(val)}`).join('&')
      }]
    })

    shortUrl.value = `http://localhost:8080/${response.data.short_url}`
  } catch (error) {
    alert('生成失败，请检查输入！')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@keyframes fadeIn {
  0% { opacity: 0; transform: translateY(10px); }
  100% { opacity: 1; transform: translateY(0); }
}

.fade-in {
  animation: fadeIn 0.8s ease-in;
}

.container {
  max-width: 500px;
  margin: 50px auto;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: #f9f9f9;
  padding: 30px;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

h1 {
  text-align: center;
  color: #2f855a;
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
}

input, select {
  width: 100%;
  padding: 10px;
  border: 1px solid #cbd5e0;
  border-radius: 6px;
  box-sizing: border-box;
  font-size: 14px;
}

button {
  background-color: #38a169;
  color: white;
  border: none;
  padding: 12px;
  width: 100%;
  font-size: 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:disabled {
  background-color: #a0aec0;
  cursor: not-allowed;
}

button:hover:enabled {
  background-color: #2f855a;
}

.result {
  margin-top: 25px;
  background: #e6fffa;
  padding: 15px;
  border-left: 4px solid #38a169;
  color: #22543d;
  border-radius: 8px;
}
</style>
