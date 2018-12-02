const SockJS = require("sockjs-client");

// Should change when you reconnect
const terminalId = '1d29e2561dac917f8616a6325f0cb664';

const sock = new SockJS(`http://localhost:9090/api/sockjs?${terminalId}`);

sock.onopen = function () {
    console.log('websocket open');
    sock.send(JSON.stringify({'Op': 'bind', 'SessionID': terminalId}));
};

sock.onmessage = function (evt) {
    let msg = JSON.parse(evt.data);
    switch (msg['Op']) {
        case 'stdout':
            console.log(msg['Data']);
            break;
        case 'oob':
            console.log(msg['Data']);
            break;
        default:
        // console.error('Unexpected message type:', msg);
    }
};

sock.onclose = function (evt) {
    if (evt.reason !== '' && evt.code < 1000) {
        console.log(evt.reason, null);
    } else {
        console.log('Connnection closed', null);
    }
    sock.close();
};