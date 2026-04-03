<template>
  <div class="svg-badge" :class="`badge-${type}`" :style="badgeStyle" :title="title">
    <svg v-if="type === 'newcomer'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="newcomer-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#86efac;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#22c55e;stop-opacity:1" />
        </linearGradient>
        <filter id="newcomer-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#166534" flood-opacity="0.3"/>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#newcomer-grad)" filter="url(#newcomer-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#f0fdf4" stroke="#22c55e" stroke-width="2"/>
      <g transform="translate(50, 50)">
        <circle cx="0" cy="-8" r="12" fill="#22c55e"/>
        <circle cx="0" cy="-8" r="8" fill="#86efac"/>
        <path d="M-12 10 Q0 0 12 10" stroke="#22c55e" stroke-width="3" fill="none" stroke-linecap="round"/>
        <path d="M-18 18 Q0 35 18 18" stroke="#22c55e" stroke-width="2" fill="none" stroke-linecap="round"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'first-post'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="first-post-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#93c5fd;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#3b82f6;stop-opacity:1" />
        </linearGradient>
        <filter id="first-post-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#1e40af" flood-opacity="0.3"/>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#first-post-grad)" filter="url(#first-post-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#eff6ff" stroke="#3b82f6" stroke-width="2"/>
      <g transform="translate(50, 50)">
        <rect x="-20" y="-18" width="40" height="36" rx="3" fill="#3b82f6"/>
        <rect x="-17" y="-15" width="34" height="30" rx="2" fill="#eff6ff"/>
        <line x1="-12" y1="-8" x2="12" y2="-8" stroke="#3b82f6" stroke-width="2" stroke-linecap="round"/>
        <line x1="-12" y1="-2" x2="12" y2="-2" stroke="#3b82f6" stroke-width="2" stroke-linecap="round"/>
        <line x1="-12" y1="4" x2="8" y2="4" stroke="#3b82f6" stroke-width="2" stroke-linecap="round"/>
        <path d="M18 -15 L28 -5 L22 1" stroke="#3b82f6" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'first-comment'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="first-comment-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#fde047;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#eab308;stop-opacity:1" />
        </linearGradient>
        <filter id="first-comment-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#854d0e" flood-opacity="0.3"/>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#first-comment-grad)" filter="url(#first-comment-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#fefce8" stroke="#eab308" stroke-width="2"/>
      <g transform="translate(50, 50)">
        <path d="M-20 -15 Q-25 -15 -25 -10 L-25 12 Q-25 17 -20 17 L20 17 Q25 17 25 12 L25 -10 Q25 -15 20 -15 Z" fill="#eab308"/>
        <path d="M-17 -12 Q-21 -12 -21 -8 L-21 10 Q-21 14 -17 14 L17 14 Q21 14 21 10 L21 -8 Q21 -12 17 -12 Z" fill="#fefce8"/>
        <circle cx="-10" cy="0" r="3" fill="#eab308"/>
        <circle cx="5" cy="0" r="3" fill="#eab308"/>
        <circle cx="15" cy="-6" r="3" fill="#eab308"/>
        <path d="M-18 20 L-12 14 L-6 20" fill="#eab308"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'writer'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="writer-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#fdba74;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#ea580c;stop-opacity:1" />
        </linearGradient>
        <filter id="writer-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#9a3412" flood-opacity="0.3"/>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#writer-grad)" filter="url(#writer-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#fff7ed" stroke="#ea580c" stroke-width="2"/>
      <g transform="translate(50, 50)">
        <path d="M-5 20 L-2 8 L18 -12 L24 -6 L4 14 L-2 17 Z" fill="#ea580c"/>
        <path d="M-2 8 L6 0 L10 4 L2 12 Z" fill="#c2410c"/>
        <circle cx="20" cy="-10" r="5" fill="#fdba74" stroke="#ea580c" stroke-width="1.5"/>
        <line x1="-6" y1="22" x2="6" y2="22" stroke="#c2410c" stroke-width="3" stroke-linecap="round"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'community-star'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="community-star-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#c4b5fd;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#7c3aed;stop-opacity:1" />
        </linearGradient>
        <filter id="community-star-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#581c87" flood-opacity="0.3"/>
        </filter>
        <filter id="glow">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#community-star-grad)" filter="url(#community-star-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#faf5ff" stroke="#7c3aed" stroke-width="2"/>
      <g transform="translate(50, 50)" filter="url(#glow)">
        <path d="M0 -22 L5 -8 L20 -8 L8 2 L13 18 L0 8 L-13 18 L-8 2 L-20 -8 L-5 -8 Z" fill="#7c3aed"/>
        <path d="M0 -15 L3 -5 L12 -5 L5 0 L8 11 L0 6 L-8 11 L-5 0 L-12 -5 L-3 -5 Z" fill="#c4b5fd"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'popular'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="popular-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#fca5a5;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#dc2626;stop-opacity:1" />
        </linearGradient>
        <filter id="popular-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#7f1d1d" flood-opacity="0.3"/>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#popular-grad)" filter="url(#popular-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#fef2f2" stroke="#dc2626" stroke-width="2"/>
      <g transform="translate(50, 50)">
        <path d="M0 20 C-18 10 -25 -2 -25 -10 C-25 -18 -18 -24 -12 -24 C-6 -24 0 -18 0 -12 C0 -18 6 -24 12 -24 C18 -24 25 -18 25 -10 C25 -2 18 10 0 20 Z" fill="#dc2626"/>
        <path d="M0 15 C-14 7 -20 -2 -20 -8 C-20 -14 -14 -19 -10 -19 C-6 -19 0 -14 0 -10 C0 -14 6 -19 10 -19 C14 -19 20 -14 20 -8 C20 -2 14 7 0 15 Z" fill="#fef2f2"/>
        <circle cx="-10" cy="-5" r="2" fill="#dc2626"/>
        <circle cx="10" cy="-5" r="2" fill="#dc2626"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'gold-comment'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="gold-comment-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#fef08a;stop-opacity:1" />
          <stop offset="50%" style="stop-color:#facc15;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#ca8a04;stop-opacity:1" />
        </linearGradient>
        <filter id="gold-comment-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="2" stdDeviation="2" flood-color="#713f12" flood-opacity="0.3"/>
        </filter>
        <filter id="gold-shine">
          <feGaussianBlur stdDeviation="1" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#gold-comment-grad)" filter="url(#gold-comment-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#fefce8" stroke="#ca8a04" stroke-width="2"/>
      <g transform="translate(50, 50)" filter="url(#gold-shine)">
        <circle cx="0" cy="0" r="18" fill="#ca8a04"/>
        <circle cx="0" cy="0" r="14" fill="#facc15"/>
        <path d="M0 -10 L-4 -2 L-12 -2 L-6 4 L-8 12 L0 7 L8 12 L6 4 L12 -2 L4 -2 Z" fill="#fefce8"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'opinion-leader'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="opinion-leader-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#fef08a;stop-opacity:1" />
          <stop offset="30%" style="stop-color:#facc15;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#854d0e;stop-opacity:1" />
        </linearGradient>
        <filter id="opinion-leader-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="3" stdDeviation="3" flood-color="#422006" flood-opacity="0.4"/>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#opinion-leader-grad)" filter="url(#opinion-leader-shadow)"/>
      <circle cx="50" cy="50" r="38" fill="#fefce8" stroke="#854d0e" stroke-width="2.5"/>
      <g transform="translate(50, 50)">
        <path d="M-20 -15 L-15 -25 L-7 -18 L0 -28 L7 -18 L15 -25 L20 -15 L12 -12 L0 -18 L-12 -12 Z" fill="#854d0e"/>
        <circle cx="0" cy="2" r="12" fill="#854d0e"/>
        <circle cx="0" cy="2" r="9" fill="#facc15"/>
        <circle cx="-3" cy="-1" r="2.5" fill="#fefce8"/>
        <circle cx="3" cy="-1" r="2.5" fill="#fefce8"/>
        <path d="M-4 6 Q0 9 4 6" stroke="#854d0e" stroke-width="1.5" fill="none" stroke-linecap="round"/>
      </g>
    </svg>
    
    <svg v-else-if="type === 'legend'" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <linearGradient id="legend-grad1" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" style="stop-color:#ef4444;stop-opacity:1" />
          <stop offset="25%" style="stop-color:#f59e0b;stop-opacity:1" />
          <stop offset="50%" style="stop-color:#10b981;stop-opacity:1" />
          <stop offset="75%" style="stop-color:#3b82f6;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#8b5cf6;stop-opacity:1" />
        </linearGradient>
        <linearGradient id="legend-grad2" x1="100%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" style="stop-color:#ef4444;stop-opacity:1" />
          <stop offset="25%" style="stop-color:#f59e0b;stop-opacity:1" />
          <stop offset="50%" style="stop-color:#10b981;stop-opacity:1" />
          <stop offset="75%" style="stop-color:#3b82f6;stop-opacity:1" />
          <stop offset="100%" style="stop-color:#8b5cf6;stop-opacity:1" />
        </linearGradient>
        <filter id="legend-shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="3" stdDeviation="4" flood-color="#1f2937" flood-opacity="0.5"/>
        </filter>
        <filter id="rainbow-glow">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
      </defs>
      <circle cx="50" cy="50" r="45" fill="url(#legend-grad1)" filter="url(#legend-shadow)">
        <animateTransform attributeName="transform" type="rotate" from="0 50 50" to="360 50 50" dur="10s" repeatCount="indefinite"/>
      </circle>
      <circle cx="50" cy="50" r="38" fill="white" stroke="url(#legend-grad2)" stroke-width="3"/>
      <g transform="translate(50, 50)" filter="url(#rainbow-glow)">
        <path d="M-18 5 L-15 -22 L0 -10 L15 -22 L18 5 L25 20 L0 12 L-25 20 Z" fill="url(#legend-grad1)"/>
        <path d="M-12 2 L-10 -16 L0 -6 L10 -16 L12 2 L17 14 L0 8 L-17 14 Z" fill="white"/>
        <polygon points="0,-18 5,-10 14,-10 7,-4 10,5 0,0 -10,5 -7,-4 -14,-10 -5,-10" fill="url(#legend-grad2)"/>
      </g>
    </svg>
  </div>
</template>

<script setup>
const props = defineProps({
  type: {
    type: String,
    required: true
  },
  size: {
    type: Number,
    default: 48
  },
  title: {
    type: String,
    default: ''
  }
})

const badgeStyle = {
  width: props.size + 'px',
  height: props.size + 'px'
}
</script>

<style scoped>
.svg-badge {
  display: inline-block;
  transition: all 0.3s ease;
}

.svg-badge:hover {
  transform: scale(1.15);
}

.svg-badge svg {
  width: 100%;
  height: 100%;
  display: block;
}
</style>
