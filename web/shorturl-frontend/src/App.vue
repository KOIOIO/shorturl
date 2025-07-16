<template>
  <div class="container">
    <h1 class="fade-in">🔗 <span class="highlight">短链接生成器</span></h1>
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

    <transition name="bounce">
      <div v-if="shortUrl" class="result fade-in">
        <p>生成成功：<a :href="shortUrl" target="_blank">{{ shortUrl }}</a></p>
      </div>
    </transition>
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
    let isRateLimited = false
    // 兼容后端 HTTP 500 且 code 字段不存在的情况
    if (error.response) {
      const data = error.response.data
      // code 字段优先判断
      if (data && String(data.code) === '5001') {
        alert('为防止数据库崩溃，请半个小时后再生成')
        isRateLimited = true
      } else if (error.response.status === 500 && error.response.statusText === 'Internal Server Error') {
        // 兼容后端直接返回500错误
        alert('为防止数据库崩溃，请半个小时后再生成')
        isRateLimited = true
      } else if (typeof data === 'string' && data.includes('rate limit')) {
        // 兼容后端直接返回字符串错误
        alert('为防止数据库崩溃，请半个小时后再生成')
        isRateLimited = true
      }
    }
    if (!isRateLimited) {
      alert('生成失败，请检查输入！')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@700&display=swap');

@keyframes fadeIn {
  0% { opacity: 0; transform: translateY(30px) scale(0.95); }
  60% { opacity: 1; transform: translateY(-8px) scale(1.03); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes bounceIn {
  0% { opacity: 0; transform: scale(0.7); }
  60% { opacity: 1; transform: scale(1.1); }
  100% { opacity: 1; transform: scale(1); }
}

.fade-in {
  animation: fadeIn 1s cubic-bezier(.68,-0.55,.27,1.55);
}

.bounce-enter-active {
  animation: bounceIn 0.7s cubic-bezier(.68,-0.55,.27,1.55);
}
.bounce-leave-active {
  animation: fadeIn 0.3s reverse;
}

.bg-animated {
  position: relative;
  min-height: 100vh;
  width: 100vw;
  overflow: hidden;
  background: linear-gradient(120deg, #f9d423 0%, #ff4e50 100%);
}

.container {
  max-width: 500px;
  margin: 60px auto;
  font-family: 'Montserrat', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  background: linear-gradient(135deg, #f9d423 0%, #ff4e50 100%);
  padding: 36px 32px 32px 32px;
  border-radius: 18px;
  box-shadow: 0 12px 32px rgba(255, 78, 80, 0.18), 0 2px 4px rgba(0,0,0,0.08);
  border: 2px solid #ff4e50;
  position: relative;
  z-index: 2;
}

h1 {
  text-align: center;
  color: #fff;
  margin-bottom: 32px;
  letter-spacing: 2px;
  font-size: 2.2rem;
  font-family: 'Montserrat', sans-serif;
}

.highlight {
  background: linear-gradient(90deg, #f9d423 30%, #ff4e50 70%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: bold;
}

.form-group {
  margin-bottom: 22px;
}

label {
  color: #ff4e50;
  font-weight: bold;
  margin-bottom: 6px;
  display: block;
  letter-spacing: 1px;
}

input, select {
  width: 100%;
  padding: 12px;
  border: 2px solid #f9d423;
  border-radius: 8px;
  box-sizing: border-box;
  font-size: 15px;
  background: #fffbe6;
  transition: border-color 0.3s, box-shadow 0.3s;
  outline: none;
}

input[type="text"] {
  background: #222;
  color: #fff;
  border: 2px solid #ff4e50;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

input[type="text"]:focus {
  border-color: #f9d423;
  box-shadow: 0 0 0 2px #ffb19955;
}

select {
  background: #222;
  color: #fff;
  border: 2px solid #ff4e50;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

select:focus {
  border-color: #f9d423;
  box-shadow: 0 0 0 2px #ffb19955;
}

button {
  background: linear-gradient(90deg, #ff4e50 0%, #f9d423 100%);
  color: #fff;
  border: none;
  padding: 14px;
  width: 100%;
  font-size: 18px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: bold;
  letter-spacing: 1px;
  box-shadow: 0 4px 16px rgba(255, 78, 80, 0.12);
  transition: background 0.3s, transform 0.2s;
}

button:disabled {
  background: #ffd6d6;
  color: #fff;
  cursor: not-allowed;
  opacity: 0.7;
}

button:hover:enabled {
  background: linear-gradient(90deg, #f9d423 0%, #ff4e50 100%);
  transform: translateY(-2px) scale(1.03);
}

.result {
  margin-top: 28px;
  background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
  padding: 18px;
  border-left: 6px solid #43e97b;
  color: #22543d;
  border-radius: 10px;
  font-size: 1.1rem;
  box-shadow: 0 2px 8px rgba(67, 233, 123, 0.08);
  animation: bounceIn 0.7s;
}

.result a {
  color: #ff4e50;
  font-weight: bold;
  word-break: break-all;
  text-decoration: underline;
  transition: color 0.2s;
}

.result a:hover {
  color: #f9d423;
}
</style>