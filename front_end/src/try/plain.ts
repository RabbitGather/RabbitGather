let step1_offer_sent = {
  type: "OFFER",
  payload: {
    sdp: {
      type: "offer",
      sdp: "v=0\r\no=- 9111744496538967905 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:xqg9\r\na=ice-pwd:p7BRT818terD2Q3uL8YFiS38\r\na=ice-options:trickle\r\na=fingerprint:sha-256 80:F7:EE:03:A2:42:93:40:48:85:50:73:7D:C0:C3:3A:42:D6:FE:95:2E:52:4A:F2:00:70:AC:36:07:2F:B0:E1\r\na=setup:actpass\r\na=mid:0\r\na=sctp-port:5000\r\na=max-message-size:262144\r\n",
    },
    type: "data",
    connectionId: "dc_li12ek0xez",
    browser: "chrome",
    label: "dc_li12ek0xez",
    reliable: true,
    serialization: "binary",
  },
  dst: "6ccaa070-0c8d-4d76-8af7-1fa9a2707eb3",
};
// +   src: "90b30418-54ba-4122-bdee-19aef3337a5a",

let step1_offer_recieve = {
  type: "OFFER",
  src: "90b30418-54ba-4122-bdee-19aef3337a5a",
  dst: "6ccaa070-0c8d-4d76-8af7-1fa9a2707eb3",
  payload: {
    sdp: {
      type: "offer",
      sdp: "v=0\r\no=- 9111744496538967905 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:xqg9\r\na=ice-pwd:p7BRT818terD2Q3uL8YFiS38\r\na=ice-options:trickle\r\na=fingerprint:sha-256 80:F7:EE:03:A2:42:93:40:48:85:50:73:7D:C0:C3:3A:42:D6:FE:95:2E:52:4A:F2:00:70:AC:36:07:2F:B0:E1\r\na=setup:actpass\r\na=mid:0\r\na=sctp-port:5000\r\na=max-message-size:262144\r\n",
    },
    type: "data",
    connectionId: "dc_li12ek0xez",
    browser: "chrome",
    label: "dc_li12ek0xez",
    reliable: true,
    serialization: "binary",
  },
};

let step2_answer_sent = {
  type: "ANSWER",
  payload: {
    sdp: {
      type: "answer",
      sdp: "v=0\r\no=- 7075592012677239613 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:uyOd\r\na=ice-pwd:ic1t1pTzsCyBqNfvNTOAJod7\r\na=ice-options:trickle\r\na=fingerprint:sha-256 D1:9D:30:2B:8D:79:0A:B6:B8:9C:8D:97:49:30:C4:5D:E0:AC:58:9A:6D:8E:C3:CA:45:28:FC:B8:74:18:6A:75\r\na=setup:active\r\na=mid:0\r\na=sctp-port:5000\r\na=max-message-size:262144\r\n",
    },
    type: "data",
    connectionId: "dc_li12ek0xez",
    browser: "chrome",
  },
  dst: "90b30418-54ba-4122-bdee-19aef3337a5a",
};

let step2_answer_recieve = {
  type: "ANSWER",
  payload: {
    sdp: {
      type: "answer",
      sdp: "v=0\r\no=- 7075592012677239613 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:uyOd\r\na=ice-pwd:ic1t1pTzsCyBqNfvNTOAJod7\r\na=ice-options:trickle\r\na=fingerprint:sha-256 D1:9D:30:2B:8D:79:0A:B6:B8:9C:8D:97:49:30:C4:5D:E0:AC:58:9A:6D:8E:C3:CA:45:28:FC:B8:74:18:6A:75\r\na=setup:active\r\na=mid:0\r\na=sctp-port:5000\r\na=max-message-size:262144\r\n",
    },
    type: "data",
    connectionId: "dc_li12ek0xez",
    browser: "chrome",
  },
  dst: "90b30418-54ba-4122-bdee-19aef3337a5a",
};

let sss = {
  type: "OFFER",
  payload: {
    sdp: {
      type: "offer",
      sdp: "v=0\r\no=- 6265712925003995093 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:GZ5O\r\na=ice-pwd:cJylSY5M10OPn54bejffQMOt\r\na=ice-options:trickle\r\na=fingerprint:sha-256 A4:FD:9C:69:53:A7:18:09:45:E0:97:27:BE:9B:E5:C5:D3:05:EA:E7:54:70:EF:71:9D:F5:51:76:D3:5D:47:E7\r\na=setup:actpass\r\na=mid:0\r\na=sctp-port:5000\r\na=max-message-size:262144\r\n",
    },
    type: "data",
    connectionId: "dc_ukde9blu45",
    browser: "chrome",
    label: "dc_ukde9blu45",
    reliable: true,
    serialization: "binary",
  },
  dst: "",
};
