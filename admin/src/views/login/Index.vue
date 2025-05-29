<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-title">酒馆后台管理系统</div>
      
      <el-form 
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        auto-complete="on"
        label-position="left"
      >
        <div class="title-container">
          <h3 class="title">管理员登录</h3>
        </div>
        
        <el-form-item prop="username">
          <el-input
            ref="usernameRef"
            v-model="loginForm.username"
            placeholder="用户名"
            name="username"
            type="text"
            auto-complete="on"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            ref="passwordRef"
            v-model="loginForm.password"
            :type="passwordVisible ? 'text' : 'password'"
            placeholder="密码"
            name="password"
            auto-complete="on"
            prefix-icon="Lock"
            @keyup.enter="handleLogin"
          >
            <template #suffix>
              <el-icon class="password-icon" @click="togglePasswordVisible">
                <component :is="passwordVisible ? 'View' : 'Hide'" />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-button 
          :loading="loading" 
          type="primary" 
          class="login-button" 
          @click="handleLogin"
        >
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { login } from '../../api/user';

const router = useRouter();
const route = useRoute();

// 表单引用
const loginFormRef = ref(null);
const usernameRef = ref(null);
const passwordRef = ref(null);

// 密码可见性
const passwordVisible = ref(false);

// 加载状态
const loading = ref(false);

// 登录表单数据
const loginForm = reactive({
  username: 'admin',
  password: ''
});

// 表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
  ]
};

// 切换密码可见性
const togglePasswordVisible = () => {
  passwordVisible.value = !passwordVisible.value;
};

// 处理登录
const handleLogin = () => {
  loginFormRef.value.validate(async (valid) => {
    if (!valid) return;
    
    loading.value = true;
    
    try {
      const response = await login({
        username: loginForm.username,
        password: loginForm.password
      });
      
      // 保存token
      localStorage.setItem('token', response.data.token);
      
      // 登录成功提示
      ElMessage({
        message: '登录成功',
        type: 'success',
        duration: 2000
      });
      
      // 重定向到之前的页面或首页
      const redirect = route.query.redirect || '/dashboard';
      router.push(redirect);
    } catch (error) {
      console.error('Login failed:', error);
    } finally {
      loading.value = false;
    }
  });
};

// 挂载后自动聚焦用户名输入框
onMounted(() => {
  if (usernameRef.value) {
    usernameRef.value.focus();
  }
});
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #2d3a4b;
  background-image: linear-gradient(135deg, #2d3a4b, #3f4e63);
}

.login-box {
  width: 420px;
  padding: 30px;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.login-title {
  font-size: 26px;
  font-weight: bold;
  text-align: center;
  color: #304156;
  margin-bottom: 30px;
}

.title-container {
  margin-bottom: 20px;
}

.title {
  font-size: 20px;
  color: #5a5e66;
  text-align: center;
  font-weight: normal;
}

.login-form {
  width: 100%;
}

.login-button {
  width: 100%;
  margin-top: 20px;
  border: none;
  height: 40px;
  font-size: 16px;
}

.password-icon {
  cursor: pointer;
  font-size: 16px;
}
</style> 