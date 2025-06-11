import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import 'element-plus/dist/index.css';
import './assets/css/main.css';
import './assets/css/theme.css';
import { registerGlobalComponents } from './components';
import { registerDirectives } from './directives';
import { pinia } from './stores';

const app = createApp(App);

app.use(router);
app.use(pinia);

// 注册全局组件
registerGlobalComponents(app);

// 注册自定义指令
registerDirectives(app);

app.mount('#app'); 