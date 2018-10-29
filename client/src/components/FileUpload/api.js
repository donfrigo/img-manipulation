import openSocket from 'socket.io-client';
const  socket = openSocket('http://localhost:5000');

function socketSuccess(cb) {
    socket.on('finished', socketMsg => cb(null, socketMsg));
}

function socketError(cb) {
    socket.on('error', socketMsg => cb(null, socketMsg));
}

function getSocket() {
    return socket
}

export { socketSuccess, socketError, getSocket };