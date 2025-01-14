import Alpine from "alpinejs";

type State = 'loading' | 'editing' | 'viewing'

Alpine.data('bookReview', () => ({
  state: 'viewing' as State,

  content: {
    ['x-show']() {
      return this.state !== 'editing'
    },

    [':aria-busy']() {
      return this.state === 'loading'
    }
  },

  islandTemplate: {
    ['x-if']() {
      return this.state !== 'viewing'
    },

    ['x-show']() {
      return this.state === 'editing'
    },
  },

  island: {
    ['@island-before-dispose']() {
      this.state = 'viewing';
    },
    
    ['@island-before-mount']() {
      this.state = 'editing';
    },

    ['@island-mount']() {
      this.state = 'editing';
    },

    ['@island-failed']() {
      this.state = 'viewing';
    }
  },

  editButton: {
    '@click'() {
      this.state = 'loading';
    },

    ':disabled'() {
      return this.state === 'loading'
    },
  },

  actions: {
    'x-show'() {
      return this.state !== 'editing'
    }
  }
}))