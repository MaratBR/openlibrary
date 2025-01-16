import Alpine from "alpinejs";

Alpine.data('bookRatingInput', () => ({
  value: 0,
  originalValue: 0,

  init() {
    
    const v = this.$root.getAttribute('data-value')
    if (v && !Number.isNaN(+v)) {
      const n = Math.floor(+v)
      if (n >= 0 && n <= 10) {
        this.value = n
        this.originalValue = n
      }
    }
  },

  trigger: {
    '@mousemove'(e: MouseEvent) {
      const rect = this.$el.getBoundingClientRect();
      const value = Math.max(2, Math.ceil((e.clientX - rect.x) / rect.width * 5) * 2);
      if (value !== this.value) {
        this.value = value;
      }
    },
    '@click'() {
      this.originalValue = this.value;
      this.$root.dispatchEvent(new CustomEvent('input', { detail: this.originalValue }))
    },
    '@mouseleave'() {
      this.value = this.originalValue
    }
  }
}))