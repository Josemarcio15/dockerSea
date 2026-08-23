<script lang="ts">
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

<style>
  .dockersea-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 300px;
  }

  .dockersea-loading.fullScreen {
    min-height: 100vh;
    background: linear-gradient(180deg, #0b0f19 0%, #111827 40%, #0f172a 100%);
  }

  .sea-container {
    position: relative;
    width: 100%;
    max-width: 400px;
    height: 320px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-end;
    overflow: hidden;
  }

  .waves {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 400%;
    height: 120px;
    animation: waveMove 6s linear infinite;
  }

  .wave-1 {
    animation: wave1 4s ease-in-out infinite;
    opacity: 0.8;
  }
  .wave-2 {
    animation: wave2 5s ease-in-out infinite;
    opacity: 0.6;
    bottom: -5px;
  }
  .wave-3 {
    animation: wave3 3.5s ease-in-out infinite;
    opacity: 0.4;
    bottom: -10px;
  }

  @keyframes waveMove {
    0% {
      transform: translateX(0);
    }
    100% {
      transform: translateX(-25%);
    }
  }

  @keyframes wave1 {
    0%,
    100% {
      transform: translateY(0) scaleY(1);
    }
    50% {
      transform: translateY(-6px) scaleY(1.1);
    }
  }

  @keyframes wave2 {
    0%,
    100% {
      transform: translateY(0) scaleY(1);
    }
    50% {
      transform: translateY(4px) scaleY(0.9);
    }
  }

  @keyframes wave3 {
    0%,
    100% {
      transform: translateY(0) scaleY(1);
    }
    50% {
      transform: translateY(-3px) scaleY(1.05);
    }
  }

  /* ── Ship ── */
  .ship {
    position: absolute;
    bottom: 70px;
    left: 50%;
    transform: translateX(-50%);
    animation: shipBob 3s ease-in-out infinite;
  }

  @keyframes shipBob {
    0%,
    100% {
      transform: translateX(-50%) translateY(0) rotate(-1.5deg);
    }
    50% {
      transform: translateX(-50%) translateY(-6px) rotate(1.5deg);
    }
  }

  .ship-hull {
    position: relative;
    width: 120px;
    height: 30px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    border-radius: 4px 4px 12px 12px;
    box-shadow: 0 4px 20px rgba(99, 102, 241, 0.3);
  }

  .ship-deck {
    position: absolute;
    top: -4px;
    left: -4px;
    right: -4px;
    height: 8px;
    background: linear-gradient(135deg, #818cf8, #a78bfa);
    border-radius: 6px 6px 0 0;
  }

  .ship-stack {
    position: absolute;
    bottom: 100%;
    display: flex;
    gap: 2px;
  }

  .ship-stack.stack-2 {
    left: 20px;
  }

  .ship-stack:not(.stack-2) {
    right: 20px;
  }

  .container-box {
    width: 14px;
    border-radius: 2px;
    animation: containerGlow 2s ease-in-out infinite;
  }

  .container-box.c1 {
    height: 14px;
    background: #f472b6;
    animation-delay: 0s;
  }
  .container-box.c2 {
    height: 18px;
    background: #a78bfa;
    animation-delay: 0.3s;
  }
  .container-box.c3 {
    height: 12px;
    background: #34d399;
    animation-delay: 0.6s;
  }
  .container-box.c4 {
    height: 16px;
    background: #60a5fa;
    animation-delay: 0.9s;
  }
  .container-box.c5 {
    height: 20px;
    background: #fbbf24;
    animation-delay: 1.2s;
  }

  @keyframes containerGlow {
    0%,
    100% {
      opacity: 0.7;
    }
    50% {
      opacity: 1;
    }
  }

  .ship-chimney {
    position: absolute;
    top: -24px;
    right: 8px;
    width: 8px;
    height: 20px;
    background: linear-gradient(180deg, #4f46e5, #6366f1);
    border-radius: 2px 2px 0 0;
  }

  .smoke {
    position: absolute;
    width: 6px;
    height: 6px;
    background: rgba(148, 163, 184, 0.4);
    border-radius: 50%;
    left: 50%;
    transform: translateX(-50%);
  }

  .smoke.s1 {
    top: -6px;
    animation: smoke 2s ease-out infinite;
  }
  .smoke.s2 {
    top: -12px;
    animation: smoke 2s ease-out 0.5s infinite;
    width: 4px;
    height: 4px;
  }
  .smoke.s3 {
    top: -18px;
    animation: smoke 2s ease-out 1s infinite;
    width: 3px;
    height: 3px;
  }

  @keyframes smoke {
    0% {
      opacity: 0.6;
      transform: translateX(-50%) translateY(0) scale(1);
    }
    100% {
      opacity: 0;
      transform: translateX(-50%) translateY(-15px) scale(1.5);
    }
  }

  /* ── Whale ── */
  .whale {
    position: absolute;
    bottom: 45px;
    left: 10%;
    width: 60px;
    animation: whaleSwim 6s ease-in-out infinite;
  }

  .whale-svg {
    width: 100%;
    height: auto;
  }

  @keyframes whaleSwim {
    0%,
    100% {
      transform: translateX(0) scaleX(1);
    }
    25% {
      transform: translateX(10px) scaleX(1);
    }
    50% {
      transform: translateX(20px) scaleX(-1);
    }
    75% {
      transform: translateX(10px) scaleX(-1);
    }
  }

  /* ── Bubbles ── */
  .bubbles {
    position: absolute;
    bottom: 90px;
    left: 50%;
    transform: translateX(-50%);
    width: 200px;
    height: 100px;
  }

  .bubble {
    position: absolute;
    border-radius: 50%;
    background: radial-gradient(
      circle at 30% 30%,
      rgba(139, 92, 246, 0.3),
      transparent
    );
    animation: bubble 3s ease-in infinite;
  }

  .bubble.b1 {
    width: 8px;
    height: 8px;
    left: 20%;
    bottom: 0;
    animation-delay: 0s;
  }
  .bubble.b2 {
    width: 5px;
    height: 5px;
    left: 40%;
    bottom: 0;
    animation-delay: 0.8s;
  }
  .bubble.b3 {
    width: 6px;
    height: 6px;
    left: 60%;
    bottom: 0;
    animation-delay: 1.6s;
  }
  .bubble.b4 {
    width: 4px;
    height: 4px;
    left: 75%;
    bottom: 0;
    animation-delay: 2.2s;
  }
  .bubble.b5 {
    width: 7px;
    height: 7px;
    left: 30%;
    bottom: 0;
    animation-delay: 3s;
  }

  @keyframes bubble {
    0% {
      transform: translateY(0) scale(1);
      opacity: 0.6;
    }
    100% {
      transform: translateY(-60px) scale(0.5);
      opacity: 0;
    }
  }

  /* ── Loading Content ── */
  .loading-content {
    position: absolute;
    top: 20px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
  }

  .spinner-ring {
    width: 36px;
    height: 36px;
    border: 3px solid rgba(139, 92, 246, 0.15);
    border-top-color: #8b5cf6;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  .loading-message {
    font-size: 0.875rem;
    font-weight: 600;
    color: #8b5cf6;
    text-align: center;
    animation: pulse 2s ease-in-out infinite;
  }

  @keyframes pulse {
    0%,
    100% {
      opacity: 0.6;
    }
    50% {
      opacity: 1;
    }
  }

  .loading-dots {
    display: flex;
    gap: 4px;
  }

  .dot {
    width: 6px;
    height: 6px;
    background: #8b5cf6;
    border-radius: 50%;
    animation: dotBounce 1.4s ease-in-out infinite;
  }

  .dot:nth-child(2) {
    animation-delay: 0.2s;
  }
  .dot:nth-child(3) {
    animation-delay: 0.4s;
  }

  @keyframes dotBounce {
    0%,
    80%,
    100% {
      transform: scale(0.6);
      opacity: 0.4;
    }
    40% {
      transform: scale(1);
      opacity: 1;
    }
  }

  :global(.dark) .dockersea-loading.fullScreen {
    background: linear-gradient(180deg, #0b0f19 0%, #111827 40%, #0f172a 100%);
  }
</style>
