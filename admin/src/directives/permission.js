import { useUserStore } from '@/stores/user';

export default {
  // 在绑定元素的父组件挂载之前调用
  beforeMount(el, binding) {
    const { value } = binding;
    const userStore = useUserStore();
    
    if (value && !userStore.hasPermission(value)) {
      // 如果没有权限，移除元素
      el.parentNode && el.parentNode.removeChild(el);
    }
  }
}; 