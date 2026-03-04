/**
 * Initializes the mouse-following glow effect on elements with the 'glow-effect' class.
 * It updates CSS variables --x and --y with the mouse position relative to each element.
 * This runs globally to allow "proximity" lighting effects (glow before hover).
 */
export function initGlowEffect(): void {
  const updateMousePosition = (e: MouseEvent) => {
    // Optimization: Use requestAnimationFrame to prevent layout thrashing
    requestAnimationFrame(() => {
      const elements = document.querySelectorAll('.glow-effect');

      elements.forEach((el) => {
        const rect = el.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;

        (el as HTMLElement).style.setProperty('--x', `${x}px`);
        (el as HTMLElement).style.setProperty('--y', `${y}px`);
      });
    });
  };

  document.addEventListener('mousemove', updateMousePosition);
}
