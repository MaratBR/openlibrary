import { connectUserStream, disconnectUserStream, resumeUserStream } from './api'

window.addEventListener('pagehide', disconnectUserStream)
window.addEventListener('pageshow', resumeUserStream)
void connectUserStream()
