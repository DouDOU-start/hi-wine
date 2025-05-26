const fs = require('fs');
const path = require('path');

const icons = [
  {
    name: 'tab-home',
    normal: `<svg width="81" height="81" viewBox="0 0 81 81" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="81" height="81" rx="20" fill="#fff"/>
  <path d="M20 40L40.5 22L61 40" stroke="#f7cac9" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
  <rect x="27" y="40" width="27" height="19" rx="5" stroke="#92a8d1" stroke-width="4"/>
</svg>`,
    selected: `<svg width="81" height="81" viewBox="0 0 81 81" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="81" height="81" rx="20" fill="url(#paint0_linear)"/>
  <path d="M20 40L40.5 22L61 40" stroke="#fff" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
  <rect x="27" y="40" width="27" height="19" rx="5" stroke="#fff" stroke-width="4"/>
  <defs>
    <linearGradient id="paint0_linear" x1="0" y1="0" x2="81" y2="81" gradientUnits="userSpaceOnUse">
      <stop stop-color="#f7cac9"/>
      <stop offset="1" stop-color="#92a8d1"/>
    </linearGradient>
  </defs>
</svg>`
  },
  {
    name: 'tab-cart',
    normal: `<svg width="81" height="81" viewBox="0 0 81 81" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="81" height="81" rx="20" fill="#fff"/>
  <circle cx="30" cy="62" r="4" fill="#f7cac9"/>
  <circle cx="54" cy="62" r="4" fill="#92a8d1"/>
  <path d="M22 25H27L32 54H54L59 34H29" stroke="#f7cac9" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`,
    selected: `<svg width="81" height="81" viewBox="0 0 81 81" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="81" height="81" rx="20" fill="url(#paint0_linear)"/>
  <circle cx="30" cy="62" r="4" fill="#fff"/>
  <circle cx="54" cy="62" r="4" fill="#fff"/>
  <path d="M22 25H27L32 54H54L59 34H29" stroke="#fff" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
  <defs>
    <linearGradient id="paint0_linear" x1="0" y1="0" x2="81" y2="81" gradientUnits="userSpaceOnUse">
      <stop stop-color="#f7cac9"/>
      <stop offset="1" stop-color="#92a8d1"/>
    </linearGradient>
  </defs>
</svg>`
  },
  {
    name: 'tab-user',
    normal: `<svg width="81" height="81" viewBox="0 0 81 81" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="81" height="81" rx="20" fill="#fff"/>
  <circle cx="40.5" cy="36" r="10" stroke="#92a8d1" stroke-width="4"/>
  <path d="M24 59c0-8.284 7.716-15 16.5-15s16.5 6.716 16.5 15" stroke="#f7cac9" stroke-width="4" stroke-linecap="round"/>
</svg>`,
    selected: `<svg width="81" height="81" viewBox="0 0 81 81" fill="none" xmlns="http://www.w3.org/2000/svg">
  <rect width="81" height="81" rx="20" fill="url(#paint0_linear)"/>
  <circle cx="40.5" cy="36" r="10" stroke="#fff" stroke-width="4"/>
  <path d="M24 59c0-8.284 7.716-15 16.5-15s16.5 6.716 16.5 15" stroke="#fff" stroke-width="4" stroke-linecap="round"/>
  <defs>
    <linearGradient id="paint0_linear" x1="0" y1="0" x2="81" y2="81" gradientUnits="userSpaceOnUse">
      <stop stop-color="#f7cac9"/>
      <stop offset="1" stop-color="#92a8d1"/>
    </linearGradient>
  </defs>
</svg>`
  }
];

const staticDir = path.join(__dirname, 'static');
if (!fs.existsSync(staticDir)) fs.mkdirSync(staticDir);

icons.forEach(icon => {
  fs.writeFileSync(path.join(staticDir, `${icon.name}.svg`), icon.normal);
  fs.writeFileSync(path.join(staticDir, `${icon.name}-selected.svg`), icon.selected);
});

console.log('INS风tabBar SVG图标已生成到 static 目录'); 