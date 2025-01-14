declare global {
  interface OLLogger {
    debug: typeof console.debug;
    info: typeof console.log;
    warn: typeof console.warn;
    error: typeof console.error;
  }

  interface Window {
    createLogger(name: string): OLLogger;
  }
}

function createLogger(name: string): OLLogger {
  if (name === '') {
    return console
  }

  const prefix = [`%c[${name}]`, 'color:#A676FF']

  return {
    debug: console.debug.bind(console, ...prefix),
    info: console.log.bind(console, ...prefix),
    warn: console.warn.bind(console, ...prefix),
    error: console.error.bind(console, ...prefix),
  }
}


window.createLogger = createLogger;