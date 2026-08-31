<script lang="ts">
  import "./DockerseaLoading.css";
  interface Props {
    message?: string;
    fullScreen?: boolean;
  }
  let {
    message = "Navegando pelos mares do Docker...",
    fullScreen = false,
  }: Props = $props();
</script>

<div
  class="dockersea-loading {fullScreen ? 'fixed inset-0 z-50' : 'w-full'}"
  class:fullScreen
>
  <!-- SVG Waves Background -->
  <div class="sea-container">
    <svg
      class="waves"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 150 28"
      preserveAspectRatio="none"
    >
      <defs>
        <linearGradient id="wave-gradient" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" stop-color="rgba(14, 165, 233, 0.45)" />
          <stop offset="50%" stop-color="rgba(59, 130, 246, 0.65)" />
          <stop offset="100%" stop-color="rgba(6, 182, 212, 0.45)" />
        </linearGradient>
        <linearGradient id="wave-gradient-2" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" stop-color="rgba(59, 130, 246, 0.35)" />
          <stop offset="50%" stop-color="rgba(14, 165, 233, 0.5)" />
          <stop offset="100%" stop-color="rgba(59, 130, 246, 0.35)" />
        </linearGradient>
        <linearGradient id="wave-gradient-3" x1="0%" y1="0%" x2="100%" y2="0%">
          <stop offset="0%" stop-color="rgba(6, 182, 212, 0.25)" />
          <stop offset="50%" stop-color="rgba(37, 99, 235, 0.4)" />
          <stop offset="100%" stop-color="rgba(6, 182, 212, 0.25)" />
        </linearGradient>
      </defs>
      <path
        class="wave wave-1"
        fill="url(#wave-gradient)"
        d="M0,10 C30,20 60,0 90,10 C120,20 150,10 150,10 L150,28 L0,28 Z"
      />
      <path
        class="wave wave-2"
        fill="url(#wave-gradient-2)"
        d="M0,15 C30,5 60,25 90,15 C120,5 150,15 150,15 L150,28 L0,28 Z"
      />
      <path
        class="wave wave-3"
        fill="url(#wave-gradient-3)"
        d="M0,20 C30,30 60,10 90,20 C120,30 150,20 150,20 L150,28 L0,28 Z"
      />
    </svg>

    <!-- Docker Whale / Container Ship -->
    <div class="ship">
      <div class="ship-hull">
        <div class="ship-deck"></div>
        <div class="ship-stack">
          <div class="container-box c1"></div>
          <div class="container-box c2"></div>
          <div class="container-box c3"></div>
        </div>
        <div class="ship-stack stack-2">
          <div class="container-box c4"></div>
          <div class="container-box c5"></div>
        </div>
        <div class="ship-chimney">
          <div class="smoke s1"></div>
          <div class="smoke s2"></div>
          <div class="smoke s3"></div>
        </div>
      </div>
    </div>

    <!-- Whale silhouette -->
    <div class="whale">
      <svg viewBox="0 0 100 50" class="whale-svg">
        <path
          d="M10,35 Q5,30 8,25 Q10,22 15,23 Q20,20 30,18 Q40,16 50,17
					Q60,18 70,20 Q80,22 88,25 Q92,27 95,30 Q98,33 97,37
					Q96,40 93,42 Q90,43 85,41 Q80,38 75,37
					Q70,36 60,35 L50,35 L40,35
					Q30,35 25,36 Q20,37 15,38 Q12,39 10,38 Q8,37 10,35Z"
          fill="rgba(139,92,246,0.3)"
        />
        <circle cx="82" cy="30" r="2" fill="rgba(255,255,255,0.4)" />
      </svg>
    </div>

    <!-- Dots / Bubbles -->
    <div class="bubbles">
      <div class="bubble b1"></div>
      <div class="bubble b2"></div>
      <div class="bubble b3"></div>
      <div class="bubble b4"></div>
      <div class="bubble b5"></div>
    </div>

    <!-- Loading text -->
    <div class="loading-content">
      <div class="spinner-ring"></div>
      <p class="loading-message">{message}</p>
      <div class="loading-dots">
        <span class="dot"></span>
        <span class="dot"></span>
        <span class="dot"></span>
      </div>
    </div>
  </div>
</div>
